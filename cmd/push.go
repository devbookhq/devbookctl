package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"

	"github.com/devbookhq/devbookctl/internal/env"
)

const (
	configName     = "dbk.toml"
	dockerfileName = "Dockerfile.dbk"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Build and push a custom environment for Devbook VM",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error determining current working dir\n> %s\n", err)
		}

		confPath := filepath.Join(dir, configName)
		conf, err := env.ParseConfig(confPath)
		if err != nil {
			log.Fatalf("\nError with config (dbk.toml)\n> %s\n", err)
		}

		dockerfilePath := filepath.Join(dir, dockerfileName)
		if _, err := os.Stat(dockerfilePath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("\nError with Dockerfile.dbk\n> file %s is missing", dockerfilePath)
		}

		fmt.Printf("\nBuilding and pushing env with ID: \"%s\"\n", conf.ID)

		fmt.Print("- Initializing Docker ")
		client, err := docker.NewClientFromEnv()
		if err != nil {
			log.Fatalf("\n\nError initializing Docker\n> %s\n", err)
		}

		fmt.Println("(done)")

		fmt.Print("- Updating Devbook base image ")
		if err = env.PullBaseEnv(ctx, client); err != nil {
			log.Fatalf("\n\nError updating base image\n> %s\n", err)
		}
		fmt.Println("(done)")

		fmt.Printf("- Building custom Devbook env \"%s\" from Devbook.dbk\n\n", conf.ID)
		imageName, err := env.BuildEnv(ctx, client, conf, dir, dockerfileName)
		if err != nil {
			log.Fatalf("\n\nError building custom env\n> %s\n", err)
		}
		fmt.Println("")

		fmt.Printf("- Pushing custom Devbook env \"%s\"\n\n", conf.ID)
		if err = env.PushEnv(ctx, client, conf, imageName); err != nil {
			log.Fatalf("\n\nError pushing custom env\n> %s\n", err)
		}
		fmt.Printf("\nCreated and pushed custom env with id \"%s\"\n", conf.ID)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
