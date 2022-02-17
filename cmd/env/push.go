package env

import (
	"bytes"
	"fmt"

	"github.com/devbookhq/devbookctl/cmd/utils"
	docker "github.com/fsouza/go-dockerclient"
)

func PushEnv(client *docker.Client, conf *EnvConfig) {

	var buf bytes.Buffer

	opts := docker.PushImageOptions{
		Name:         registryPath + conf.Id,
		OutputStream: &buf,
	}

	auth := docker.AuthConfiguration{}

	err := client.PushImage(opts, auth)
	utils.Check(err)

	fmt.Println(buf)
}
