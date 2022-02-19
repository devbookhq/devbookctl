package env

import (
	"fmt"
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
)

type EnvConfig struct {
	ID       string `toml:"id"`
	RootDir  string `toml:"root_dir"`
	StartCmd string `toml:"start_cmd"`
}

var registryServer = "us-central1-docker.pkg.dev"
var registryPath = registryServer + "/devbookhq/devbook-runner-templates"
var baseImageName = "base"
var baseImagePath = registryPath + "/" + baseImageName
var baseImageAlias = "devbook"

func ParseConfig(confPath string) (*EnvConfig, error) {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %v", err)
	}

	var conf EnvConfig
	_, err = toml.Decode(string(data), &conf)
	if err != nil {
		return nil, fmt.Errorf("cannot decode config content: %v", err)
	}

	match, err := regexp.MatchString("^[a-z][a-z0-9-]*$", conf.ID)
	if err != nil {
		return nil, fmt.Errorf("cannot check config ID field: %v", err)
	}

	if !match {
		return nil, fmt.Errorf("id in config must start with a letter and contain only lowercase letters, numbers, or a dash '-'")
	}

	return &conf, nil
}
