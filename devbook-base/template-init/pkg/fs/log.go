package fs

import (
	"github.com/DevbookHQ/template-init/pkg/log"
	"github.com/sirupsen/logrus"
)

func WithEvent(ev Event) log.Opt {
	return func(entry *logrus.Entry) *logrus.Entry {
		return entry.WithFields(logrus.Fields{
			"op":   ev.Operation(),
			"path": ev.Path(),
			"type": ev.Type(),
		})
	}
}
