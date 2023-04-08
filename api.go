package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func removeContainer(c *docker.Client, id string) {
	opts := docker.RemoveContainerOptions{
		ID:    id,
		Force: true,
	}
	if err := c.RemoveContainer(opts); err != nil {
		log.Fatal(err)
	}
}

func restartContainer(c *docker.Client, id string) {
	if err := c.RestartContainer(id, 5); err != nil {
		log.Fatal(err)
	}
}

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

func killContainer(c *docker.Client, id string) {
	opts := docker.KillContainerOptions{
		ID: id,
	}
	if err := c.KillContainer(opts); err != nil {
		log.Fatal(err)
	}
}

func startContainer(c *docker.Client, id string) {
	if err := c.StartContainer(id, nil); err != nil {
		log.Fatal(err)
	}
}

func stopContainer(c *docker.Client, id string) {
	if err := c.StopContainer(id, 5); err != nil {
		log.Fatal(err)
	}
}

func listContainers(c *docker.Client, showAll bool) []docker.APIContainers {
	containers, err := c.ListContainers(docker.ListContainersOptions{All: showAll})
	if err != nil {
		log.Fatal(err)
	}
	return containers
}
