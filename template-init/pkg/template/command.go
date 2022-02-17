package template

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type Command struct {
	cmd         *exec.Cmd
	executionID string
	startedAt   int64

	outScanner *bufio.Scanner
	errScanner *bufio.Scanner

	stdout chan<- string
	stderr chan<- string
	done   chan<- error
}

func newCommand(eid ExecutionID, rootDir, cmdString string) *Command {
	cmd := exec.Command("sh", "-c", cmdString)
	cmd.Dir = rootDir
	return &Command{
		cmd:         cmd,
		executionID: eid,
		startedAt:   time.Now().UnixMilli(),
	}
}

func (c *Command) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExecutionID string `json:"executionID"`
		StartedAt   int64  `json:"startedAt"`
	}{
		ExecutionID: c.executionID,
		StartedAt:   c.startedAt,
	})
}

func (c *Command) String() string {
	return fmt.Sprintf("command=%s", c.cmd)
}

func (c *Command) scanStdout(stdoutPipe io.ReadCloser) {
	outScanner := bufio.NewScanner(stdoutPipe)
	outScanner.Split(bufio.ScanLines)

	for outScanner.Scan() {
		c.stdout <- outScanner.Text()
	}
}

func (c *Command) scanStderr(stderrPipe io.ReadCloser) {
	errScanner := bufio.NewScanner(stderrPipe)
	errScanner.Split(bufio.ScanLines)

	for errScanner.Scan() {
		c.stderr <- errScanner.Text()
	}
}

func (c *Command) startAndWait() {
	// Stream stdout to the stdout channel.
	stdoutPipe, err := c.cmd.StdoutPipe()
	if err != nil {
		c.done <- fmt.Errorf(
			"failed to cmd.StdoutPipe():\ncommand=%s\nerr=%s",
			c.cmd,
			err,
		)
		return
	}
	go c.scanStdout(stdoutPipe)

	// Stream stderr to the stderr channel.
	stderrPipe, err := c.cmd.StderrPipe()
	if err != nil {
		c.done <- fmt.Errorf(
			"failed to cmd.StderrPipe():\ncommand=%s\nerr=%s",
			c.cmd,
			err,
		)
		return
	}
	go c.scanStderr(stderrPipe)

	// Start the command.
	if err := c.cmd.Start(); err != nil {
		c.done <- fmt.Errorf(
			"failed to cmd.Start():\ncommand=%s\nerr=%s",
			c.cmd,
			err,
		)
		return
	}

	// Wait for the command to finish.
	if err := c.cmd.Wait(); err != nil {
		c.done <- fmt.Errorf(
			"failed to cmd.Wait():\ncommand=%s\nerr=%s",
			c.cmd,
			err,
		)
		return
	}

	c.done <- nil
}

func (c *Command) kill() {
	if err := c.cmd.Process.Kill(); err != nil {
		c.done <- fmt.Errorf(
			"failed to cmd.Process.Kill():\ncommand=%s\nerr=%s",
			c.cmd,
			err,
		)
		return
	}
	c.done <- nil
}
