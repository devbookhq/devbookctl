package env

import (
	"bytes"
	"fmt"

	"github.com/devbookhq/devbookctl/cmd/err"
	docker "github.com/fsouza/go-dockerclient"
)

func PushEnv(client *docker.Client, conf *EnvConfig, imageName string) {

	var buf bytes.Buffer

	opts := docker.PushImageOptions{
		Name:         imageName,
		OutputStream: &buf,
	}

	auth := docker.AuthConfiguration{}

	pushErr := client.PushImage(opts, auth)
	err.Check(pushErr)

	fmt.Println(buf)
}
