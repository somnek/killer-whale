package main

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func listContainers(c *docker.Client) {
	containers, err := c.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		fmt.Println("Name", container.Names)
	}
}

func listImages(c *docker.Client) {
	images, err := c.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}

	for _, image := range images {
		fmt.Println("ID", image.ID)
	}
}
