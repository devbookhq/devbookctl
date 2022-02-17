package fs

import (
	"github.com/DevbookHQ/template-init/pkg/fs/item"
)

// `Event` is both emitted and accepted by `Filesystem`.
// It's accepted as an operation that the filesystem should do.
// It's emitted as an operation that has happened in the filesystem.
type Event interface {
	Operation() Op
	Path() string
	Type() item.Type
}

type CreateEvent struct {
	path     string
	itemType item.Type
}

func (ce *CreateEvent) Operation() Op   { return CREATE_OP }
func (ce *CreateEvent) Path() string    { return ce.path }
func (ce *CreateEvent) Type() item.Type { return ce.itemType }

type WriteEvent struct {
	path string
}

func (we *WriteEvent) Operation() Op   { return WRITE_OP }
func (we *WriteEvent) Path() string    { return we.path }
func (we *WriteEvent) Type() item.Type { return item.FILE }

type RemoveEvent struct {
	path     string
	itemType item.Type
}

func (re *RemoveEvent) Operation() Op   { return REMOVE_OP }
func (re *RemoveEvent) Path() string    { return re.path }
func (re *RemoveEvent) Type() item.Type { return re.itemType }
