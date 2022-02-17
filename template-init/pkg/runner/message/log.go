package message

import (
	"github.com/DevbookHQ/template-init/pkg/log"
	"github.com/sirupsen/logrus"
)

func WithMessageData(data MessagePayload) log.Opt {
	f := logrus.Fields{}
	f["type"] = data.Type

	if p, ok := data.Payload.(MessagePayloadCommand); ok {
		f["command"] = p.Command
	}

	if p, ok := data.Payload.(MessagePayloadCodeCells); ok {
		f["codeCells"] = p.CodeCells
	}

	if p, ok := data.Payload.(MessagePayloadInstallPkgs); ok {
		f["packages"] = p.Packages
	}

	return func(entry *logrus.Entry) *logrus.Entry {
		return entry.WithFields(f)
	}
}
