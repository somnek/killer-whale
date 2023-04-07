package main

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func unPauseContainer(c *docker.Client, id string) {
	if err := c.UnpauseContainer(id); err != nil {
		log.Fatal(err)
	}
}

func pauseContainer(c *docker.Client, id string) {
	if err := c.PauseContainer(id); err != nil {
		log.Fatal(err)
	}
}

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

func listContainers(c *docker.Client, showAll bool) []docker.APIContainers {
	containers, err := c.ListContainers(docker.ListContainersOptions{All: showAll})
	if err != nil {
		log.Fatal(err)
	}
	return containers
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
