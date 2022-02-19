package template

import (
	"sync"
)

type ExecutionID = string

type runningCommands struct {
	lock sync.RWMutex
	cmap map[ExecutionID]*Command
}

func newRunningCommands() *runningCommands {
	return &runningCommands{
		cmap: make(map[ExecutionID]*Command),
	}
}

func (rc *runningCommands) add(id ExecutionID, cmd *Command) {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	rc.cmap[id] = cmd
}

func (rc *runningCommands) get(id ExecutionID) *Command {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return rc.cmap[id]
}

func (rc *runningCommands) getAll() []*Command {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	cmds := make([]*Command, 0, len(rc.cmap))
	for _, cmd := range rc.cmap {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (rc *runningCommands) remove(id ExecutionID) {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	delete(rc.cmap, id)
}
