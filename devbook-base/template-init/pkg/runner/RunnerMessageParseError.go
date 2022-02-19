package runner

type RunnerMessageParseError struct {
	msg string
	// Number of bytes parsed before the error byte.
	Offset int64
	// Byte at `Offset`.
	ErrByte byte
}

func (err *RunnerMessageParseError) Error() string { return err.msg }
