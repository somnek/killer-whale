package main

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func listContainers(c *docker.Client) {
	imgs, err := c.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range imgs {
		fmt.Println("ID", img.ID)
	}
}

func listImages(c *docker.Client) {
	imgs, err := c.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range imgs {
		fmt.Println("ID", img.ID)
	}
}
