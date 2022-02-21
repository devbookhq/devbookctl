package environment

import (
	"fmt"
	"os"
)

// All the values for these variables should be already set as environment variables during template dockerization.
type Environment struct {
	// `RUNNER_SOCKET_PATH` is a path to a unix socket that's used for communication with Runner.
	RUNNER_SOCKET_PATH string

	// Template specific variables.
	ROOT_DIR       string // Root dir for a template.
	CODE_CELLS_DIR string // A directory where we should write code cells.
	START_CMD      string // A shell command that starts the template, e.g. `npm run dev` for NextJS
}

func New() (*Environment, error) {
	env := &Environment{}

	env.RUNNER_SOCKET_PATH = os.Getenv("runner_socket_path")
	if env.RUNNER_SOCKET_PATH == "" {
		return nil, fmt.Errorf("env 'RUNNER_SOCKET_PATH' is empty")
	}

	env.ROOT_DIR = os.Getenv("root_dir")
	if env.ROOT_DIR == "" {
		return nil, fmt.Errorf("cannot create template, env 'root_dir' is empty")
	}
	env.CODE_CELLS_DIR = env.ROOT_DIR

	env.START_CMD = os.Getenv("start_cmd")
	if env.START_CMD == "" {
		return nil, fmt.Errorf("cannot create template, env 'start_cmd' is empty")
	}

	return env, nil
}
