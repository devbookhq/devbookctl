package env

import (
	"bytes"
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func PullBaseEnv(ctx context.Context, client *docker.Client) error {
	outputbuf := bytes.NewBuffer(nil)

	// Pull base image path from registry
	err := client.PullImage(docker.PullImageOptions{
		Repository:   baseImagePath,
		Context:      ctx,
		OutputStream: outputbuf,
	}, docker.AuthConfiguration{})

	// fmt.Println(outputbuf)

	if err != nil {
		return fmt.Errorf("cannot pull Devbook base Docker image: %v", err)
	}

	// Retag base image from registry as "devbook" image
	err = client.TagImage(baseImagePath, docker.TagImageOptions{
		Repo:    baseImageAlias,
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("cannot alias Devbook base Docker image: %v", err)
	}

	// fmt.Println()
	return nil
}
