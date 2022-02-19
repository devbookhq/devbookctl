package template

import (
	"github.com/DevbookHQ/template-init/pkg/log"
	"github.com/sirupsen/logrus"
)

func WithTemplate(t *Template) log.Opt {
	return func(entry *logrus.Entry) *logrus.Entry {
		return entry.WithFields(logrus.Fields{
			"startCommand": t.startCommand,
			"rootDir":      t.RootDir,
			"codeCellsDir": t.CodeCellsDir,
			"state":        t.State,
		})
	}
}
