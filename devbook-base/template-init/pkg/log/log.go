package log

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Entry = logrus.Entry
type Fields = logrus.Fields

var (
	logFile *os.File
	l       = logrus.New()
)

func init() {
	path := os.Getenv("TINIT_LOG_FILE")
	if path == "" {
		l.Warn("The 'TINIT_LOG_FILE' env var is empty, logs aren't written to a file")
	} else {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		logFile = f
		l.SetOutput(logFile)
	}
}

func Cleanup() error {
	if logFile == nil {
		return nil
	}

	err := logFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close the log file: %s", err)
	}
	return nil
}

func SetOutput(out io.Writer) {
	l.SetOutput(out)
}

func Info(args ...interface{}) {
	l.Info(args...)
}

func Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func Fatal(args ...interface{}) {
	l.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	l.Fatalf(format, args...)
}

func Log(opts ...Opt) *Entry {
	// PERF: By instantiating a new entry everytime we call `Log` we will probably
	// loose significant performance of Logrus.
	contextLogger := logrus.NewEntry(l)

	for _, opt := range opts {
		contextLogger = opt(contextLogger)
	}

	return contextLogger
}
