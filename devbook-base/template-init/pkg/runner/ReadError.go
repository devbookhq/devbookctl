package runner

type ReadError struct {
	// `ReadError` has no offset field because we always read only single byte.

	msg        string
	SocketPath string
	ErrByte    byte
}
type ReadEOF ReadError

func (err *ReadError) Error() string { return err.msg }
func (err *ReadEOF) Error() string   { return err.msg }
