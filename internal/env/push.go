package env

import (
	"bytes"
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func PushEnv(ctx context.Context, client *docker.Client, conf *EnvConfig, imageName string) error {
	outputbuf := bytes.NewBuffer(nil)

	err := client.PushImage(docker.PushImageOptions{
		Name:         imageName,
		OutputStream: outputbuf,
		Tag:          "latest",
		Context:      ctx,
	}, docker.AuthConfiguration{
		// Docker push requires that the `X-Registry-Auth` header is present - it can be even an empty string.
		RegistryToken: "",
		ServerAddress: registryServer,
	})

	// fmt.Println(outputbuf)

	if err != nil {
		return fmt.Errorf("cannot push custom env Docker image: %v", err)
	}

	return nil
}
