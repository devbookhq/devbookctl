package env

import (
	"fmt"
	"os"

	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/devbookhq/devbookctl/cmd/err"
)

type EnvConfig struct {
	Id       string `toml:"id"`
	RootDir  string `toml:"root_dir"`
	StartCmd string `toml:"start_cmd"`
}

func ParseConfig(confPath string) EnvConfig {
	data, readErr := os.ReadFile(confPath)
	err.Check(readErr)

	var conf EnvConfig
	_, decodeErr := toml.Decode(string(data), &conf)
	err.Check(decodeErr)

	match, matchErr := regexp.MatchString("^[a-z][a-z0-9-]*$", conf.Id)
	err.Check(matchErr)

	if !match {
		fmt.Println("env id must start with a letter and contain only lowercase letters, numbers, or a dash '-'")
		os.Exit(1)
	}

	return conf
}
