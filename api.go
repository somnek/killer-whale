package main

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func stopContainer(c *docker.Client, id string) {
	if err := c.StopContainer(id, 5); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Container stopped")
}

func removeContainer(c *docker.Client, id string) {
	opts := docker.RemoveContainerOptions{
		ID: id,
	}
	if err := c.RemoveContainer(opts); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Container removed")
}

func listContainers(c *docker.Client) {
	containers, err := c.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		fmt.Println("Name", container.ID)
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
