package main

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func unpauseAndWriteLog(container Container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	state := container.state
	id := container.id

	var logs string
	if state == "paused" {
		unpauseContainer(client, id)
		logs = "✅ Unpaused " + container.name + "\n"
	} else {
		logs = "🚧  " + container.name + " is not running\n"
	}
	return logs
}

func pauseAndWriteLog(container Container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	state := container.state
	id := container.id

	var logs string
	if state == "running" {
		pauseContainer(client, id)
		logs = "⏳ Paused " + container.name + "\n"
	} else {
		logs = "🚧  " + container.name + " is not running\n"
	}
	return logs
}

func startAndWriteLog(container Container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	state := container.state
	id := container.id

	var logs string
	if state == "exited" || state == "created" {
		go startContainer(client, id)
		if err != nil {
			logs = fmt.Sprintf("🚧  %s\n", err.Error())
		} else {
			logs = "🚀 Started " + container.name + "\n"
		}
	} else {
		logs = "🚧  " + container.name + " already running\n"
	}
	return logs
}

func removeAndWriteLog(container Container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	id := container.id

	var logs string
	go removeContainer(client, id)
	logs = "🗑️  Remove " + container.name + "\n"
	return logs
}

func restartAndWriteLog(container Container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	state := container.state
	id := container.id

	var logs string
	if state == "running" {
		go restartContainer(client, id)
		logs = "🔃 Restarted " + container.name + "\n"
	} else {
		logs = "🚧  " + container.name + " not running\n"
	}
	return logs
}

func killAndWriteLog(container Container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	state := container.state
	id := container.id

	var logs string
	if state == "running" {
		killContainer(client, id)
		logs = "🔪 Killed " + container.name + "\n"
	} else {
		logs = "🚧 " + container.name + " already stopped\n"
	}
	return logs
}

func stopAndWriteLog(container Container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	state := container.state
	id := container.id

	var logs string
	if state == "running" || state == "restarting" {
		go stopContainer(client, id)
		logs = "🛑 Stop " + container.name + "\n"
	} else {
		logs = "🚧  " + " unable to stop " + container.name + "\n"
	}
	return logs
}

func getContainers() []Container {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	containers := []Container{}
	for _, c := range listContainers(client, true) {
		name := c.Names[0][1:]
		status := c.State
		c := Container{name: name, state: status, id: c.ID, ancestor: c.Image}
		containers = append(containers, c)
	}
	return containers
}

func getImages() []Image {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	images := []Image{}
	for _, c := range listImages(client, true) {
		tags := c.RepoTags
		var name string
		var size int64
		if len(tags) > 0 {
			name = tags[0]
			size = c.Size
			// format size (GB, MB, KB)
			if size > 1000000000 {
				size = size / 1000000000
				name = fmt.Sprintf("%s (%dGB)", name, size)
			} else if size > 1000000 {
				size = size / 1000000
				name = fmt.Sprintf("%s (%dMB)", name, size)
			} else if size > 1000 {
				size = size / 1000
				name = fmt.Sprintf("%s (%dKB)", name, size)
			} else {
				name = fmt.Sprintf("%s (%dB)", name, size)
			}
			c := Image{name: name, id: c.ID}
			images = append(images, c)
		}
	}
	return images
}
