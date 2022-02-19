package log

import (
	"github.com/sirupsen/logrus"
)

type Opt func(*logrus.Entry) *logrus.Entry

func WithError(err error) Opt {
	return func(entry *logrus.Entry) *logrus.Entry {
		return entry.WithError(err)
	}
}

func WithField(key string, value interface{}) Opt {
	return func(entry *logrus.Entry) *logrus.Entry {
		return entry.WithField(key, value)
	}
}

func WithFields(fields logrus.Fields) Opt {
	return func(entry *logrus.Entry) *logrus.Entry {
		return entry.WithFields(fields)
	}
}
