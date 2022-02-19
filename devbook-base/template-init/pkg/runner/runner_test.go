package runner

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"testing"

	msg "github.com/DevbookHQ/template-init/pkg/runner/message"
	"github.com/DevbookHQ/template-init/pkg/template"
)

const (
	sock      = "./runner_test.socket"
	delimiter = "\f"
)

type mockServer struct {
	c net.Conn
}

func TestParsingInstallPkgsMessage(t *testing.T) {
	cleanup()

	pkgString, pkgNames := longPackages()

	tests := []struct {
		json     string
		expected msg.Message
	}{
		{
			json: `
			{
				"type": "message",
				"data": {
					"type": "InstallPackages",
					"payload": {
						"packages": [` + pkgString + `]
					}
				}
			}
			`,
			expected: msg.Message{
				Type: "message",
				Data: msg.MessagePayload{
					Type: msg.InMessageInstallPkgs,
					Payload: msg.MessagePayloadInstallPkgs{
						Packages: pkgNames,
					},
				},
			},
		},
	}

	server, rconn := createRunnerConn(t)

	packages := make(chan msg.Message)
	rconn.Subscribe(msg.InMessageInstallPkgs, packages)

	server.write(t, []string{tests[0].json})

	select {
	case p := <-packages:
		testInstallPkgsMessage(t, p, tests[0].expected)
	case err := <-rconn.Err:
		// Ignore socket EOF error.
		if _, ok := err.(*ReadEOF); !ok {
			t.Fatal(err)
		}
	}
}

func TestParsingCodeCellsMessage(t *testing.T) {
	cleanup()

	// Create a long string to test whether we can handle arbitrary long messages.
	ccCode := longCodeCellCode()

	tests := []struct {
		json     string
		expected msg.Message
	}{
		{
			json: `
			{
				"type": "message",
				"data": {
					"type": "CodeCells",
					"payload": {
						"codeCells": [
							{
								"name": "cc1",
								"code": "` + ccCode + `"` +
				`},
							{
								"name": "cc2",
								"code": "` + ccCode + `"` +
				`}
						]
					}
				}
			}
			`,
			expected: msg.Message{
				Type: "message",
				Data: msg.MessagePayload{
					Type: msg.InMessageCodeCells,
					Payload: msg.MessagePayloadCodeCells{
						CodeCells: []template.CodeCell{
							{
								Name: "cc1",
								Code: ccCode,
							},
							{
								Name: "cc2",
								Code: ccCode,
							},
						},
					},
				},
			},
		},
	}

	server, rconn := createRunnerConn(t)

	codeCells := make(chan msg.Message)
	rconn.Subscribe(msg.InMessageCodeCells, codeCells)

	server.write(t, []string{tests[0].json})

	select {
	case cc := <-codeCells:
		testCodeCellsMessage(t, cc, tests[0].expected)
	case err := <-rconn.Err:
		// Ignore socket EOF error.
		if _, ok := err.(*ReadEOF); !ok {
			t.Fatal(err)
		}
	}
}

func TestParsingCommandMessage(t *testing.T) {
	cleanup()

	// Create a long string to test whether we can handle arbitrary long messages.
	commandStr := longCommand()

	tests := []struct {
		json     string
		expected msg.Message
	}{
		{
			json: `
				{
					"type": "message",
					"data": {
						"type": "Command",
						"payload": {
							"command": "` + commandStr + `"` +
				`}
					}
				}
			`,
			expected: msg.Message{
				Type: "message",
				Data: msg.MessagePayload{
					Type: msg.InMessageCommand,
					Payload: msg.MessagePayloadCommand{
						Command: commandStr,
					},
				},
			},
		},
	}

	server, rconn := createRunnerConn(t)

	commands := make(chan msg.Message)
	rconn.Subscribe(msg.InMessageCommand, commands)

	server.write(t, []string{tests[0].json})

	select {
	case cmd := <-commands:
		testCommandMessage(t, cmd, tests[0].expected)
	case err := <-rconn.Err:
		// Ignore socket EOF error.
		if _, ok := err.(*ReadEOF); !ok {
			t.Fatal(err)
		}
	}
}

// Test whether we are able to correctly parse multiple JSON objects in a single byte stream.
func TestParsingMultipleMessages(t *testing.T) {
	// Try to cleanup from a previous run.
	cleanup()

	ccCode := longCodeCellCode()
	ccMsgs := []string{
		`
			{
				"type": "message",
				"data": {
					"type": "CodeCells",
					"payload": {
						"codeCells": [
							{
								"name": "cc1",
								"code": "` + ccCode + `"` +
			`},
							{
								"name": "cc2",
								"code": "` + ccCode + `"` +
			`}
						]
					}
				}
			}
`,
		`
			{
				"type": "message",
				"data": {
					"type": "CodeCells",
					"payload": {
						"codeCells": [
							{
								"name": "cc1",
								"code": "` + ccCode + `"` +
			`},
							{
								"name": "cc2",
								"code": "` + ccCode + `"` +
			`}
						]
					}
				}
			}
`,
	}

	commandStr := longCommand()
	cmdMsgs := []string{
		`
				{
					"type": "message",
					"data": {
						"type": "Command",
						"payload": {
							"command": "` + commandStr + `"` +
			`}
					}
				}
`,
		`
				{
					"type": "message",
					"data": {
						"type": "Command",
						"payload": {
							"command": "` + commandStr + `"` +
			`}
					}
				}
`,
		`
				{
					"type": "message",
					"data": {
						"type": "Command",
						"payload": {
							"command": "` + commandStr + `"` +
			`}
					}
				}
`,
		`
				{
					"type": "message",
					"data": {
						"type": "Command",
						"payload": {
							"command": "` + commandStr + `"` +
			`}
					}
				}
`,
		`
				{
					"type": "message",
					"data": {
						"type": "Command",
						"payload": {
							"command": "` + commandStr + `"` +
			`}
					}
				}
`,
	}

	pkgString, _ := longPackages()
	installPkgsMsgs := []string{
		`
			{
				"type": "message",
				"data": {
					"type": "InstallPackages",
					"payload": {
						"packages": [` + pkgString + `]
					}
				}
			}
		`,
		`
			{
				"type": "message",
				"data": {
					"type": "InstallPackages",
					"payload": {
						"packages": [` + pkgString + `]
					}
				}
			}
		`,
		`
			{
				"type": "message",
				"data": {
					"type": "InstallPackages",
					"payload": {
						"packages": [` + pkgString + `]
					}
				}
			}
		`,
	}

	server, rconn := createRunnerConn(t)

	codeCells := make(chan msg.Message)
	rconn.Subscribe(msg.InMessageCodeCells, codeCells)

	commands := make(chan msg.Message)
	rconn.Subscribe(msg.InMessageCommand, commands)

	packages := make(chan msg.Message)
	rconn.Subscribe(msg.InMessageInstallPkgs, packages)

	all := append(ccMsgs, cmdMsgs...)
	all = append(all, installPkgsMsgs...)
	server.write(t, all)

	ccCount := 0
	cmdCount := 0
	pkgCount := 0
	for {
		select {
		case <-codeCells:
			ccCount += 1
			t.Log("Received CodeCellMessage")
		case <-commands:
			cmdCount += 1
			t.Log("Received CommandMessage")
		case <-packages:
			pkgCount += 1
			t.Log("Received InstallPackageMessage")
		case err := <-rconn.Err:
			// Ignore socket EOF error.
			if _, ok := err.(*ReadEOF); !ok {
				t.Fatal(err)
			}
		}

		if ccCount == len(ccMsgs) &&
			cmdCount == len(cmdMsgs) &&
			pkgCount == len(installPkgsMsgs) {
			break
		}
	}
}

func testInstallPkgsMessage(t *testing.T, parsed, expected msg.Message) {
	testBaseMessage(t, parsed, expected)

	parsedPayload, ok := parsed.Data.Payload.(msg.MessagePayloadInstallPkgs)
	if !ok {
		t.Fatalf("parsed.Data.Payload not `MessagePayloadInstallPkgs`. got=%T", parsed.Data.Payload)
	}

	expectedPayload, ok := expected.Data.Payload.(msg.MessagePayloadInstallPkgs)
	if !ok {
		t.Fatalf("expected.Data.Payload not `MessagePayloadInstallPkgs`. got=%T", expected.Data.Payload)
	}

	if len(expectedPayload.Packages) != len(parsedPayload.Packages) {
		t.Fatalf(
			"Lengths of parsedPayload.Packages and expectedPayload.Packages don't match. expected=%d got=%d",
			len(expectedPayload.Packages),
			len(parsedPayload.Packages),
		)
	}

	sort.Strings(expectedPayload.Packages)
	sort.Strings(parsedPayload.Packages)

	for idx, expected := range expectedPayload.Packages {
		got := parsedPayload.Packages[idx]
		if expected != got {
			t.Errorf("expected package name and parsed package name don't match. expected=%s got=%s", expected, got)
		}
	}
}

func testCodeCellsMessage(t *testing.T, parsed, expected msg.Message) {
	testBaseMessage(t, parsed, expected)

	parsedPayload, ok := parsed.Data.Payload.(msg.MessagePayloadCodeCells)
	if !ok {
		t.Fatalf("parsed.Data.Payload not `MessagePayloadCodeCell`. got=%T", parsed.Data.Payload)
	}

	expectedPayload, ok := expected.Data.Payload.(msg.MessagePayloadCodeCells)
	if !ok {
		t.Fatalf("expected.Data.Payload not `MessagePayloadCodeCell`. got=%T", parsed.Data.Payload)
	}

	if len(parsedPayload.CodeCells) != len(expectedPayload.CodeCells) {
		t.Fatalf(
			"parsedPayload.CodeCells does not contain %d code cells. got=%d",
			len(expectedPayload.CodeCells),
			len(parsedPayload.CodeCells),
		)
	}

	for idx, cc := range parsedPayload.CodeCells {
		expectedCC := expectedPayload.CodeCells[idx]
		testCodeCell(t, &cc, &expectedCC)
	}
}

func testCodeCell(t *testing.T, expectedCC, parsedCC *template.CodeCell) {
	if expectedCC.Name != parsedCC.Name {
		t.Fatalf("expectedCC.Name not %s. got=%s", parsedCC.Name, expectedCC.Name)
	}

	if expectedCC.Code != parsedCC.Code {
		t.Fatalf("expectedCC.Code not %s. got=%s", parsedCC.Code, expectedCC.Code)
	}
}

func testCommandMessage(t *testing.T, parsed, expected msg.Message) {
	testBaseMessage(t, parsed, expected)

	parsedPayload, ok := parsed.Data.Payload.(msg.MessagePayloadCommand)
	if !ok {
		t.Fatalf("parsed.Data.Payload not `MessagePayloadCommand`. got=%T", parsed.Data.Payload)
	}

	expectedPayload, ok := expected.Data.Payload.(msg.MessagePayloadCommand)
	if !ok {
		t.Fatalf("expected.Data.Payload not `MessagePayloadCommand`. got=%T", parsed.Data.Payload)
	}

	if parsedPayload.Command != expectedPayload.Command {
		t.Fatalf("parsedPayload.Command has unexpected value. expected=%s, got=%s", expectedPayload.Command, parsedPayload.Command)
	}
}

func testBaseMessage(t *testing.T, parsed, expected msg.Message) {
	if parsed.Type != expected.Type {
		t.Fatalf("parsed.Type has unexpected value. expected=%s, got=%s", expected.Type, parsed.Type)
	}

	if parsed.Data.Type != expected.Data.Type {
		t.Fatalf("parsed.Data.Type has unexpected value. expected=%s, got=%s", expected.Data.Type, parsed.Data.Type)
	}
}

func createRunnerConn(t *testing.T) (mockServer, Connection) {
	// Start a mock server that will write to the socket file.
	conn := make(chan net.Conn)
	s := mockServer{}
	s.start(t, conn)

	// Create new `Connection` that connects to a mock server.
	rconn := NewConnection(sock)
	rconn.disableLogging()

	go rconn.DialAndRead()
	// Wait for rconn to connect to the mock server.
	<-rconn.Ready

	// Wait for the connection to establish.
	c := <-conn
	s.c = c

	return s, rconn
}

// `start` starts a mock socket server.
func (s *mockServer) start(t *testing.T, connChan chan<- net.Conn) {
	ln, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatalf("Failed to start socket server: %s\n", err)
	}

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				t.Fatalf("Failed to accept incoming connection to server: %s\n", err)
			}
			t.Log("Accepted connection")
			connChan <- c
		}
	}()
}

// `write` writes bytes into an active connection via the socket file.
func (s *mockServer) write(t *testing.T, messages []string) {
	if s.c == nil {
		t.Fatal("Cannot write, no connection to server")
	}

	// Write all messages to a socket as a single string with delimiters.
	allMsgs := strings.Join(messages, delimiter)
	// `strings.Join()` puts a delimiter between each two elements sowe must append the delimiter at the end of all messages.
	allMsgs += delimiter

	_, err := s.c.Write([]byte(allMsgs))
	if err != nil {
		t.Fatalf("Failed to write: %s\n", err)
	}
}

// `cleanup` removes a socket file from the previous run.
func cleanup() {
	err := os.Remove(sock)
	if os.IsNotExist(err) {
		// Ignore if the socket file doesn't exist.
		return
	}
	if err != nil {
		panic(err)
	}
}

func longCodeCellCode() string {
	return strings.Repeat("LoremIpsum", 1024)
}

func longCommand() string {
	return strings.Repeat("echo ", 1024)
}

func longPackages() (string, []string) {
	l := 1024

	pkgNames := make([]string, 0, l)
	pkgNames = append(pkgNames, "packageName-0")
	pkgString := "\"packageName-0\""

	for i := 0; i < l; i++ {
		n := fmt.Sprintf("packageName-%d", i+1)
		pkgNames = append(pkgNames, fmt.Sprintf(n))
		pkgString = pkgString + fmt.Sprintf(", \"%s\"", n)
	}

	return pkgString, pkgNames
}
