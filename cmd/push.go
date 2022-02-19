package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"

	"github.com/devbookhq/devbookctl/internal/env"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Build and push a VM environment",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		dir, err := os.Getwd()
		if err != nil {
			log.Fatalln("Error determining current working dir:", err)
		}

		fmt.Print("Parsing config...")
		confPath := filepath.Join(dir, "dbk.toml")
		conf, err := env.ParseConfig(confPath)
		if err != nil {
			log.Fatalln("\nError parsing config:", err)
		}
		fmt.Println("done")

		fmt.Print("Initializing Docker...")
		client, err := docker.NewClientFromEnv()
		if err != nil {
			log.Fatalln("\nError initializing Docker client:", err)
		}
		fmt.Println("done")

		fmt.Print("Updating Devbook base image...")
		if err = env.PullBaseEnv(ctx, client); err != nil {
			log.Fatalln("\nError pulling base env:", err)
		}
		fmt.Println("done")

		fmt.Print("Building custom Devbook env...")
		imageName, err := env.BuildEnv(ctx, client, conf, dir)
		if err != nil {
			log.Fatalln("\nError building custom env:", err)
		}
		fmt.Println("done")

		fmt.Print("Pushing custom Devbook env...")
		if err = env.PushEnv(ctx, client, conf, imageName); err != nil {
			log.Fatalln("\nError pushing custom env:", err)
		}
		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
