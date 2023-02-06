package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

func main() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// cli
	var rootCmd = &cobra.Command{
		Use:   "whale",
		Short: "A Docker cli",
	}

	var ctnCmd = &cobra.Command{
		Use:   "container",
		Short: "List all containers",
		Run: func(cmd *cobra.Command, args []string) {
			listContainers(client)
		},
	}

	var imgCmd = &cobra.Command{
		Use:   "image",
		Short: "List all images",
		Run: func(cmd *cobra.Command, args []string) {
			listImages(client)
		},
	}

	rootCmd.AddCommand(ctnCmd, imgCmd)
	rootCmd.Execute()
}
