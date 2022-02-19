package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	fsitem "github.com/DevbookHQ/template-init/pkg/fs/item"
	"github.com/DevbookHQ/template-init/pkg/log"
	"github.com/fsnotify/fsnotify"
)

type Op string

const (
	LIST_OP   Op = "List"
	CREATE_OP Op = "Create"
	WRITE_OP  Op = "Write"
	REMOVE_OP Op = "Remove"
	RENAME_OP Op = "Rename"
)

type Filesystem struct {
	RootDir string
	// TODO: Add "Watcher" prefix to both?
	WatcherEvents chan Event
	WatcherErrors chan error
	fswatcher     *fsnotify.Watcher
}

func New(rootDir string) (*Filesystem, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create new filesytem watcher: %s", err)
	}

	fs := &Filesystem{
		RootDir:       rootDir,
		WatcherEvents: make(chan Event),
		WatcherErrors: make(chan error),
		fswatcher:     w,
	}

	go fs.handleWatcherErrors()
	go fs.handleWatcherEvents()

	log.Log(
		log.WithField("dir", rootDir),
	).Info("Adding filesystem watcher for dir")
	err = fs.fswatcher.Add(rootDir)
	if err != nil {
		log.Log(
			log.WithField("dir", rootDir),
		).Error("Failed to add filesystem watcher for dir")
		return nil, fmt.Errorf("failed to add filesystem watcher: %s", err)
	}
	log.Log(
		log.WithField("dir", rootDir),
	).Info("Added filesystem watcher for dir")

	return fs, nil
}

func (fs *Filesystem) Close() {
	fs.fswatcher.Close()
	close(fs.WatcherErrors)
	close(fs.WatcherEvents)
}

// `ListDir` returns a content of the dir at `path`.
// The `path` is expected to start with the "/" character as a root path of this Filesystem.
func (fs *Filesystem) ListDir(path string) ([]fsitem.Item, error) {
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("dir", path),
	).Info("Listing dir")

	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("Cannot ListDir, path '%s' doesn't start with '/'", path)
	}

	p := filepath.Join(fs.RootDir, path)
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to read directory '%s': %s", p, err)
	}

	items := make([]fsitem.Item, len(files))
	for idx, f := range files {
		// Since we checked that the requested `path` starts wih "/" (= `Filesystem.RootDir`)
		// we create the final path for a file by prefixing file's name with the `path` param.
		p := filepath.Join(path, f.Name())

		if f.IsDir() {
			items[idx] = &DirItem{p}
		} else {
			items[idx] = &FileItem{p}
		}
	}

	// Start watching the dir for changes.
	// No need to check if we are already watching the dir, `fswatcher` handles that for us.
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("path", p),
	).Info("Adding path to watcher")
	if err := fs.fswatcher.Add(p); err != nil {
		return nil, fmt.Errorf("failed to add path '%s' to watcher: %s", p, err)
	}

	return items, nil
}

func (fs *Filesystem) CreateDir(path string) error {
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("dir", path),
	).Info("Creating dir")

	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("Cannot CreateDir, path '%s' doesn't start with '/'", path)
	}

	p := filepath.Join(fs.RootDir, path)
	// TODO 0644 move to const
	if err := os.MkdirAll(p, 0644); err != nil {
		return fmt.Errorf("Failed to MkdirAll for path '%s': %s", p, err)
	}

	// Start watching the new dir.
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("path", p),
	).Info("Adding path to watcher")
	if err := fs.fswatcher.Add(p); err != nil {
		return fmt.Errorf("Failed to add path '%s' to watcher: %s", p, err)
	}

	return nil
}

func (fs *Filesystem) GetFile(path string) ([]byte, error) {
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("path", path),
	).Info("Getting file")

	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("Cannot GetFile, path '%s' doesn't start with '/'", path)
	}

	p := filepath.Join(fs.RootDir, path)
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file '%s': %s", p, err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file '%s': %s", p, err)
	}

	return b, nil
}

func (fs *Filesystem) WriteFile(path, content string) error {
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("file", path),
	).Info("Writing file")

	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("Cannot WriteFile, path '%s' doesn't start with '/'", path)
	}

	filePath := filepath.Join(fs.RootDir, path)

	// Create all parent dirs, if needed.
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0644); err != nil {
		return fmt.Errorf("Failed to MkdirAll for '%s': %s", dirPath, err)
	}

	// Start watching the new dir.
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("path", dirPath),
	).Info("Adding path to watcher")
	if err := fs.fswatcher.Add(dirPath); err != nil {
		return fmt.Errorf("Failed to add path '%s' to watcher: %s", dirPath, err)
	}

	// Finally, create the file.
	// TODO 0644 move to const
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("Failed to WriteFile '%s': %s", filePath, err)
	}

	return nil
}

func (fs *Filesystem) RemoveFileOrDir(path string) error {
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("dir", path),
	).Info("Removing file or dir")

	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("Cannot RemoveFileOrDir, path '%s' doesn't start with '/'", path)
	}

	p := filepath.Join(fs.RootDir, path)
	if err := os.RemoveAll(p); err != nil {
		return fmt.Errorf("Failed to remove path '%s': %s", p, err)
	}

	// TODO: If we are deleting a dir, we should remove it from the watcher.
	log.Log(
		log.WithField("fsRoot", fs.RootDir),
		log.WithField("path", p),
	).Info("Removing path from watcher")
	if err := fs.fswatcher.Remove(p); err != nil {
		return fmt.Errorf("Failed to remove path '%s' from watcher: %s", p, err)
	}

	return nil
}

func (fs *Filesystem) handleWatcherErrors() {
	for {
		err, ok := <-fs.fswatcher.Errors
		if !ok {
			return
		}

		if err != nil {
			fs.WatcherErrors <- fmt.Errorf("Filesystem watcher error: %s", err)
		}
	}
}

func (fs *Filesystem) handleWatcherEvents() {
	for {
		event, ok := <-fs.fswatcher.Events
		if !ok {
			return
		}

		// Create.
		if event.Op&fsnotify.Create == fsnotify.Create {
			itemType, err := fs.itemType(event.Name)
			if err != nil {
				fs.WatcherErrors <- fmt.Errorf("CreateEvent error when finding out item type: %s", err)
				continue
			}
			fsPath, err := fs.convertSystemPathToFSPath(event.Name)
			if err != nil {
				fs.WatcherErrors <- fmt.Errorf("CreateEvent error when converting path: %s", err)
				continue
			}
			fs.WatcherEvents <- &CreateEvent{path: fsPath, itemType: itemType}
		}

		// Write.
		if event.Op&fsnotify.Write == fsnotify.Write {
			fsPath, err := fs.convertSystemPathToFSPath(event.Name)
			if err != nil {
				fs.WatcherErrors <- fmt.Errorf("CreateEvent error when converting path: %s", err)
				continue
			}
			fs.WatcherEvents <- &WriteEvent{path: fsPath}
		}

		// Remove.
		if event.Op&fsnotify.Remove == fsnotify.Remove {
			itemType, err := fs.itemType(event.Name)
			if err != nil {
				fs.WatcherErrors <- fmt.Errorf("RemoveEvent error when finding out item type: %s", err)
				continue
			}
			fsPath, err := fs.convertSystemPathToFSPath(event.Name)
			if err != nil {
				fs.WatcherErrors <- fmt.Errorf("RemoveEvent error when converting path: %s", err)
				continue
			}

			fs.WatcherEvents <- &RemoveEvent{path: fsPath, itemType: itemType}
		}
	}
}

// `itemType` returns `fsitem.DIR` if the path is pointing to a dir or `fsitem.FILE` otherwise.
// The `path` parameter is the absolute path in the OS filesystem.
func (fs *Filesystem) itemType(path string) (fsitem.Type, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("Failed to stat '%s': %s", path, err)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return fsitem.DIR, nil
	case mode.IsRegular():
		return fsitem.FILE, nil
	default:
		return "", fmt.Errorf("Unexpected FileMode: %d", mode)
	}
}

// `convertSystemPathToFSPath` converts an absolute path return by OS
// to a path that is absolute in terms of this Filesystem and its `RootDir`.
// Example:
// This Filesystem has a root dir "/home/runner".
// OS returns a path to file that looks like "/home/runner/src/hello.txt".
// `convertSystemPathToFSPath` removes the "/home/runner" prefix and returns "/src/hello.txt".
func (fs *Filesystem) convertSystemPathToFSPath(sysPath string) (string, error) {
	if strings.HasPrefix(sysPath, fs.RootDir) {
		return strings.TrimPrefix(sysPath, fs.RootDir), nil
	} else {
		return "", fmt.Errorf(
			"Cannot convert convert path system path '%s' to filesystem path - doesn't have the '%s' prefix",
			sysPath,
			fs.RootDir,
		)
	}
}
