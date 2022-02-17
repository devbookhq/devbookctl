package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/DevbookHQ/template-init/pkg/environment"
	"github.com/DevbookHQ/template-init/pkg/fs"
)

type TemplateState uint

const (
	TemplateStateDone TemplateState = iota
	TemplateStateWaiting
	TemplateStateRunning
)

type CodeCell struct {
	// Name *must* contain the full code cell's name with extension.
	Name string `json:"name"`
	Code string `json:"code"`
}

// `Template` represents a template that is currently running in this environment.
// For example "nextjs-v11-components" or "nodejs-v16" template.
type Template struct {
	// `RootDir` is the root directory of a template. All files required by the template must be placed here at runtime.
	RootDir string
	// `CodeCellsDir` is a directrory where code cells are stored.
	CodeCellsDir string
	// `startCommand` is a shell command that should be executed to start the template.
	startCommand *exec.Cmd
	State        TemplateState
	once         sync.Once
	// `Done` channel sends exactly once either an error or nil when the template's process exits or is stopped via `Stop()`.
	Done        chan error
	Filesystem  *fs.Filesystem
	runningCmds *runningCommands
}

// `FromEnvironment` returns a template based on the current environment.
func FromEnvironment(env *environment.Environment) (*Template, error) {
	// "npm run dev" -> ["npm", "run", "dev"]
	splits := strings.Split(env.START_CMD, " ")
	cmd := exec.Command(splits[0], splits[1:]...)

	cmd.Dir = env.ROOT_DIR
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	templateFS, err := fs.New(env.ROOT_DIR)
	if err != nil {
		return nil, fmt.Errorf("failed to create new filesystem for template: %s", err)
	}

	return &Template{
		RootDir:      env.ROOT_DIR,
		CodeCellsDir: env.CODE_CELLS_DIR,
		startCommand: cmd,
		State:        TemplateStateWaiting,
		Done:         make(chan error),
		once:         sync.Once{},
		Filesystem:   templateFS,
		runningCmds:  newRunningCommands(),
	}, nil
}

// `Start` starts the template and waits for it to finish.
func (t *Template) Start() {
	// Already running.
	if t.State == TemplateStateRunning {
		return
	}

	if err := t.startCommand.Start(); err != nil {
		err := fmt.Errorf(
			"failed to startCommand.Start():\ncommand=%+v\nerr=%s",
			t.startCommand,
			err,
		)
		t.markDone(err)
		return
	}
	t.State = TemplateStateRunning

	go func() {
		if err := t.startCommand.Wait(); err != nil {
			err := fmt.Errorf(
				"failed to startCommand.Wait():\ncommand=%+v\nerr=%s",
				t.startCommand,
				err,
			)
			t.markDone(err)
		}
		t.markDone(nil)
	}()
}

func (t *Template) markDone(err error) {
	t.once.Do(func() {
		t.Done <- err
		close(t.Done)
		t.State = TemplateStateDone
	})
}

// `Stop` stops the running running template.
func (t *Template) Stop() {
	// Already stopped.
	if t.State == TemplateStateDone {
		return
	}

	err := t.startCommand.Process.Kill()
	if err != nil {
		err = fmt.Errorf(
			"failed to startCommand.Process.Kill():\ncommand=%s\nerr=%s",
			t.startCommand,
			err,
		)
	}
	t.markDone(err)
}

// `UpdateCodeCells` writes code cells inside `template.CodeCellsDir`.
func (t *Template) UpdateCodeCells(ccs []CodeCell) error {
	for _, cc := range ccs {
		p := t.CodeCellsDir + "/" + cc.Name
		if err := os.WriteFile(p, []byte(cc.Code), 0644); err != nil {
			return fmt.Errorf(
				"failed to os.WriteFile:\ncodeCell=%s\npath=%s\ncontent=%s\nerr=%s",
				cc.Name,
				p,
				cc.Code,
				err,
			)
		}
	}
	return nil
}

// `RunningCommands` returns a slice of all running commands.
func (t *Template) RunningCommands() []*Command {
	return t.runningCmds.getAll()
}

// `KillCommand` kills the running command.
func (t *Template) KillCommand(executionID string) {
	cmd := t.runningCmds.get(executionID)
	if cmd != nil {
		cmd.kill()
		t.runningCmds.remove(executionID)
	}
}

// `ExecCommand` starts a command inside the template's root directory and waits for it to finish.
// If a command with the same `executionID` is already running, it will kill it before starts the new one.
func (t *Template) ExecCommand(executionID, command string, stdout, stderr chan<- string, done chan<- error) {
	cmd := newCommand(executionID, t.RootDir, command)
	cmd.stdout = stdout
	cmd.stderr = stderr
	cmd.done = done

	// If a client sent a command with the same executionID, kill the old one.
	existingCmd := t.runningCmds.get(executionID)
	if existingCmd != nil {
		// TODO: This isn't the best solution because `kill()` sends to the `done` channel.
		//existingCmd.kill()
	}

	t.runningCmds.add(executionID, cmd)

	cmd.startAndWait()

	t.runningCmds.remove(executionID)
}

func (t *Template) InstallPackages(pkgs []string) ([]byte, error) {
	err, installed := t.installedPackages()
	if err != nil {
		return nil, err
	}

	// Filter out packages that are already installed.
	diff := t.diff(installed, pkgs)

	if len(diff) > 0 {
		args := append([]string{"install"}, diff...)
		cmd := exec.Command("npm", args...)
		cmd.Dir = t.RootDir
		return cmd.CombinedOutput()
	}

	return nil, nil
}

// `diff` returns a difference of two string arrays: diff = a - b.
func (t *Template) diff(a, b []string) []string {
	m := make(map[string]struct{})
	for _, el := range a {
		m[el] = struct{}{}
	}

	var diff []string
	for _, el := range b {
		if _, found := m[el]; !found {
			diff = append(diff, el)
		}
	}

	return diff
}

// `installedPackages` reads currently installed packages based on the field "dependencies"
// in the "package.json".
func (t *Template) installedPackages() (error, []string) {
	// We only have NPM based templates right now so we can just read template's "package.json".
	var pkgJSON struct {
		Deps map[string]string `json:"dependencies"`
	}

	p := t.RootDir + "/" + "package.json"
	f, err := os.Open(p)
	if err != nil {
		return err, nil
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err, nil
	}

	if err := json.Unmarshal(b, &pkgJSON); err != nil {
		return err, nil
	}

	pkgs := make([]string, 0, len(pkgJSON.Deps))
	for k := range pkgJSON.Deps {
		pkgs = append(pkgs, k)
	}

	return nil, pkgs
}
