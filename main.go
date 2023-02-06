package main

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

func main() {
	_, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	var rootCmd = &cobra.Command{
		Use:   "whale",
		Short: "A Docker cli",
	}

	var ctnCmd = &cobra.Command{
		Use:   "ctn",
		Short: "List all containers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("listing containers..")
		},
	}

	var imgCmd = &cobra.Command{
		Use:   "img",
		Short: "List all images",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("listing images..")
		},
	}

	rootCmd.AddCommand(ctnCmd, imgCmd)
	rootCmd.Execute()
}
