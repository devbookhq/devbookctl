package env

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func PullBaseEnv(ctx context.Context, client *docker.Client) error {
	// Pull base image path from registry
	if err := client.PullImage(docker.PullImageOptions{
		Repository: baseImagePath,
		Context:    ctx,
	}, docker.AuthConfiguration{}); err != nil {
		return fmt.Errorf("cannot pull Devbook base Docker image: %v", err)
	}

	// Retag base image from registry as "devbook" image
	if err := client.TagImage(baseImagePath, docker.TagImageOptions{
		Repo:    baseImageAlias,
		Context: ctx,
	}); err != nil {
		return fmt.Errorf("cannot alias Devbook base Docker image: %v", err)
	}

	return nil
}
