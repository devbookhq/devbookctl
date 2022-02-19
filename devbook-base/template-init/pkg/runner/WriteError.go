package runner

type WriteError struct {
	msg        string
	SocketPath string
	// Number of bytes written before the error.
	Offset int
	// Byte at `Offset`.
	ErrByte byte
}

func (err *WriteError) Error() string { return err.msg }
