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

	var listCmd = &cobra.Command{
		Use:   "container",
		Short: "List all containers",
		Run: func(cmd *cobra.Command, args []string) {
			listContainers(client)
		},
	}

	var xCmd = &cobra.Command{
		Use:   "x",
		Short: "...",
		Run: func(cmd *cobra.Command, args []string) {
			id := "d1f02035c4a1366b01f21d078be53a09d3899669152b541a8203aef297e7646c"
			stopContainer(client, id)
			removeContainer(client, id)
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "...",
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

	rootCmd.AddCommand(listCmd, stopCmd, imgCmd, xCmd)
	rootCmd.Execute()
}
