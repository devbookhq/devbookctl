package env

import (
	"bytes"
	"fmt"
	
	"github.com/devbookhq/devbookctl/cmd/utils"

	docker "github.com/fsouza/go-dockerclient"
)

var registryPath = "us-central1-docker.pkg.dev/devbookhq/devbook-runner-templates/"

func BuildEnv(client *docker.Client, conf *EnvConfig) {
	var buf bytes.Buffer

	opts := docker.BuildImageOptions{
		Name:         registryPath + conf.Id,
		OutputStream: &buf,
	}

	err := client.BuildImage(opts)
	utils.Check(err)

	fmt.Println(buf)
}
