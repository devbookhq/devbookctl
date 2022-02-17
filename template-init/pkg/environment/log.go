package environment

import (
	"github.com/DevbookHQ/template-init/pkg/log"
	"github.com/sirupsen/logrus"
)

func WithEnv(env *Environment) log.Opt {
	return func(entry *logrus.Entry) *logrus.Entry {
		return entry.WithFields(logrus.Fields{
			"runnerSocketPath": env.RUNNER_SOCKET_PATH,
			"rootDir":          env.ROOT_DIR,
			"codeCellsDir":     env.CODE_CELLS_DIR,
			"startCMD":         env.START_CMD,
		})
	}
}
