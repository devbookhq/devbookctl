package env

import (
	"os"

	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/devbookhq/devbookctl/cmd/utils"
)

type EnvConfig struct {
	Id string `toml:"id"`

	FilesDir     string `toml:"files_dir"`
	CodeCellsDir string `toml:"code_cells_dir"`

	SetupCmd string `toml:"setup_cmd"`
	StartCmd string `toml:"start_cmd"`
}

func ParseConfig(confPath string) EnvConfig {
	data, readErr := os.ReadFile(confPath)
	utils.Check(readErr)

	var conf EnvConfig
	_, decodeErr := toml.Decode(string(data), &conf)
	utils.Check(decodeErr)

	// Validate id so it won't contain special chars -> only a-Z,0-9,- allowed
	if regexp.MatchString(, Id) {

	}

	return conf
}
