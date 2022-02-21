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
			log.Fatalf("Error determining current working dir\n> %v\n", err)
		}

		confPath := filepath.Join(dir, configName)
		conf, err := env.ParseConfig(confPath)
		if err != nil {
			log.Fatalf("\nError with config (dbk.toml)\n> %v\n", err)
		}

		dockerfilePath := filepath.Join(dir, dockerfileName)
		if _, err := os.Stat(dockerfilePath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("\nError with Dockerfile.dbk\n> file %v is missing", dockerfilePath)
		}

		fmt.Printf("\nBuilding and pushing **%v**\n", conf.ID)

		fmt.Print("- Initializing Docker ")
		client, err := docker.NewClientFromEnv()
		if err != nil {
			log.Fatalf("\n\nError initializing Docker\n> %v\n", err)
		}
		fmt.Println("(done)")

		fmt.Print("- Updating Devbook base image ")
		if err = env.PullBaseEnv(ctx, client); err != nil {
			log.Fatalf("\n\nError updating base image\n> %v\n", err)
		}
		fmt.Println("(done)")

		fmt.Print("- Building custom Devbook env ")
		imageName, err := env.BuildEnv(ctx, client, conf, dir, dockerfileName)
		if err != nil {
			log.Fatalf("\n\nError building custom env\n> %v\n", err)
		}
		fmt.Println("(done)")

		fmt.Print("- Pushing custom Devbook env ")
		if err = env.PushEnv(ctx, client, conf, imageName); err != nil {
			log.Fatalf("\n\nError pushing custom env\n> %v\n", err)
		}
		fmt.Println("(done)")
		fmt.Printf("\nPushed custom env with id **%v**\n", conf.ID)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
