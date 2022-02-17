package runner

type ConnError struct {
	msg        string
	SocketPath string
}

func (err *ConnError) Error() string { return err.msg }
