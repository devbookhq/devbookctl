package cmd

import (
	"context"
	"fmt"
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
			fmt.Fprintf(os.Stderr, "\nError determining current working dir: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Parsing config...")
		confPath := filepath.Join(dir, "dbk.toml")
		conf, err := env.ParseConfig(confPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError parsing config: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "done\n")

		fmt.Fprintf(os.Stdout, "Initializing Docker...")
		client, err := docker.NewClientFromEnv()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError initializing Docker client: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "done\n")

		fmt.Fprintf(os.Stdout, "Updating Devbook base image...")
		err = env.PullBaseEnv(ctx, client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError pulling base env: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "done\n")
		
		fmt.Fprintf(os.Stdout, "Building custom Devbook env...")
		imageName, err := env.BuildEnv(ctx, client, &conf, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError building custom env: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "done\n")

		fmt.Fprintf(os.Stdout, "Pushing custom Devbook env...")
		err = env.PushEnv(ctx, client, &conf, imageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError pushing custom env: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "done\n")
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
