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
	_ = c.RemoveContainer(opts)
}

func restartContainer(c *docker.Client, id string) {
	_ = c.RestartContainer(id, 5)
}

func unPauseContainer(c *docker.Client, id string) {
	_ = c.UnpauseContainer(id)
}

func pauseContainer(c *docker.Client, id string) {
	_ = c.PauseContainer(id)
}

func killContainer(c *docker.Client, id string) {
	opts := docker.KillContainerOptions{
		ID: id,
	}
	_ = c.KillContainer(opts)
}

func startContainer(c *docker.Client, id string) {
	_ = c.StartContainer(id, nil)
}

func stopContainer(c *docker.Client, id string) {
	_ = c.StopContainer(id, 5)
}

func listContainers(c *docker.Client, showAll bool) []docker.APIContainers {
	containers, err := c.ListContainers(docker.ListContainersOptions{All: showAll})
	if err != nil {
		log.Fatal(err)
	}
	return containers
}
