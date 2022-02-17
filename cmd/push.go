package cmd

import (
	"path/filepath"
	"os"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"

	"github.com/devbookhq/devbookctl/cmd/env"
	"github.com/devbookhq/devbookctl/cmd/utils"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Build and push a VM environment",
	Run: func(cmd *cobra.Command, args []string) {
		dir, dirErr := os.Getwd()
		utils.Check(dirErr)

		confPath:= filepath.Join(dir, "dbk.toml")
		conf := env.ParseConfig(confPath)

		client, dockerErr := docker.NewClientFromEnv()
		utils.Check(dockerErr)

		env.BuildEnv(client, &conf)
		env.PushEnv(client, &conf)
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
