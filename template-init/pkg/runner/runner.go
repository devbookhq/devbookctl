package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"

	"github.com/DevbookHQ/template-init/pkg/log"
	msg "github.com/DevbookHQ/template-init/pkg/runner/message"
)

type EnvironmentStatus int

const (
	EnvReady EnvironmentStatus = iota
)

// TODO: Improve
func (s EnvironmentStatus) String() string {
	switch s {
	case EnvReady:
		return "Ready"
	}
	return "unknown"
}

// `Connection` represents connection to Runner.
type Connection struct {
	socketPath  string
	subscribers map[msg.MessageType][]chan msg.Message
	c           net.Conn
	rm          sync.RWMutex
	// `delimiter` splits incoming JSON messages.
	// Data coming to the socket is just a stream of bytes so we need a way to tell when a message ends.
	delimiter string
	Err       chan error
	Ready     chan struct{}
}

// `NewConnection` creates a new connection to Runner via a unix socket.
func NewConnection(socket string) Connection {
	return Connection{
		socketPath:  socket,
		subscribers: map[msg.MessageType][]chan msg.Message{},
		// Each JSON message is delimited by the '\f' character.
		delimiter: "\f",
		Err:       make(chan error),
		Ready:     make(chan struct{}),
	}
}

// For testing.
func (rc *Connection) disableLogging() {
	log.SetOutput(ioutil.Discard)
}

// `read` continuously reads messages from the connection.
// It ignores any errors that might happen during reading & unmarshalling and only logs them to stderr.
func (rc *Connection) read(messages chan []byte) {

	reader := bufio.NewReader(rc.c)
	var currentMsg []byte

	for {
		b, err := reader.ReadByte()
		if err != nil && err != io.EOF {
			rerr := &ReadError{
				msg:        fmt.Sprintf("failed to read single byte from socket: %s", err),
				SocketPath: rc.socketPath,
				ErrByte:    b,
			}
			rc.Err <- rerr
			continue
		}

		if err == io.EOF {
			rerr := &ReadEOF{
				msg:        fmt.Sprintf("EOF byte when reading from socket"),
				SocketPath: rc.socketPath,
			}
			rc.Err <- rerr
			continue
		}

		if string(b) == rc.delimiter {
			// A message is complete.
			messages <- currentMsg
			currentMsg = nil
			continue
		}

		currentMsg = append(currentMsg, b)
	}
}

// `Dial` connects to the Runner server and continuously keeps reading messages from the server. It keeps reading even if it encounters an error.
func (rc *Connection) DialAndRead() {
	// Connect to the server.
	c, err := net.Dial("unix", rc.socketPath)
	if err != nil {
		err = &ConnError{
			msg:        fmt.Sprintf("failed to dial Runner socket: %s", err),
			SocketPath: rc.socketPath,
		}
		rc.Err <- err
		return
	}
	rc.c = c

	messages := make(chan []byte)
	go func() {
		for {
			m := <-messages
			rc.parseMessage(m)
		}
	}()

	rc.Ready <- struct{}{}
	// Start continuously reading messages from Runner.
	rc.read(messages)
}

func (rc *Connection) parseMessage(m []byte) {
	msg := msg.Message{}
	if err := json.Unmarshal(m, &msg); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			parseErr := &RunnerMessageParseError{
				msg:     fmt.Sprintf("JSON syntax error when unmarshalling message in parseMessage(): %s", syntaxErr),
				Offset:  syntaxErr.Offset,
				ErrByte: m[syntaxErr.Offset],
			}
			rc.Err <- parseErr
		} else {
			rc.Err <- fmt.Errorf("failed to unmarshall incoming Runner message: %s", err)
		}
		return
	}
	rc.publish(msg)
}

// `Close` closes the connection.
func (rc *Connection) Close() {
	if rc.c != nil {
		if err := rc.c.Close(); err != nil {
			err = &ConnError{
				msg:        fmt.Sprintf("failed to close socket connection: %s", err),
				SocketPath: rc.socketPath,
			}
			rc.Err <- err
		}
	}
}

// `Subscribe` subscribes channel `ch` to Runner messages of type `t`.
func (rc *Connection) Subscribe(t msg.MessageType, ch chan msg.Message) {
	rc.rm.Lock()
	defer rc.rm.Unlock()
	// No need to check whether `rc.subscribers[t]` is nil
	// because in Go one can append to a nil slice.
	rc.subscribers[t] = append(rc.subscribers[t], ch)
}

// `publish` publishes `msg` to all subscribers.
func (rc *Connection) publish(m msg.Message) {
	rc.rm.RLock()
	defer rc.rm.RUnlock()

	if channels, found := rc.subscribers[m.Data.Type]; found {
		// We create a new slice of channels for the `topic` so we stay thread safe.
		// Because when passing a slice to a function the passed value refers to the
		// same underlying array.
		// See for example https://stackoverflow.com/questions/39993688/are-slices-passed-by-value
		slice := append([]chan msg.Message{}, channels...)
		go func(m msg.Message, chans []chan msg.Message) {
			for _, ch := range chans {
				ch <- m
			}
		}(m, slice)
	}
}

// `Send` sends a message of type `t` with `data` to a remote Runner.
func (rc *Connection) Send(payload msg.MessagePayload) {
	msg := msg.Message{
		Type: msg.IPCEventMessage,
		Data: payload,
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		rc.Err <- fmt.Errorf("failed to marshal JSON in Send(): %s", err)
		return
	}

	// The delimiter character works as a separator between messages that the server expects. Otherwise, the server wouldn't
	// know when a message ends because we are just writing streams of data to a socket file.
	bs := append(b, []byte(rc.delimiter)...)
	n, err := rc.c.Write(bs)
	if err != nil {
		err = &WriteError{
			msg:        fmt.Sprintf("failed to write to connection in Send(): %s", err),
			SocketPath: rc.socketPath,
			Offset:     n,
			ErrByte:    bs[n],
		}
		rc.Err <- err
	}
}

// `SendReady` sends the `EnvReady` message to a remote Runner signaling that the environment is ready to handle incoming messages.
func (rc *Connection) SendReady() {
	msgPayload := msg.MessagePayload{
		Type: msg.OutMessageStatus,
		Payload: msg.MessagePayloadStatus{
			Status: EnvReady.String(),
		},
	}
	rc.Send(msgPayload)
}
