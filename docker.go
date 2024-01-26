package main

import (
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

////////////////////////////////////////////
// TODO: handle errors return from docker
// client instead of assigning to _
////////////////////////////////////////////

// ---------------- Volume ----------------
func removeVolume(c *docker.Client, id string) {
	opts := docker.RemoveVolumeOptions{
		Name:  id,
		Force: true,
	}
	_ = c.RemoveVolumeWithOptions(opts)
}

func listVolumes(c *docker.Client) []docker.Volume {
	opts := docker.ListVolumesOptions{}
	volumes, err := c.ListVolumes(opts)
	if err != nil {
		log.Fatal(err)
	}
	return volumes
}

// ---------------- Image ----------------
func removeImage(c *docker.Client, id string) {
	opts := docker.RemoveImageOptions{
		Force: true,
	}
	// just tell em to remove the container that use this image first
	_ = c.RemoveImageExtended(id, opts)
}

func listImages(c *docker.Client, showAll bool) []docker.APIImages {
	opts := docker.ListImagesOptions{
		All: showAll,
	}
	images, err := c.ListImages(opts)
	if err != nil {
		log.Fatal(err)
	}
	return images
}

// ---------------- Container ----------------
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

func unpauseContainer(c *docker.Client, id string) {
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
