package main

import (
	"fmt"
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

func getContainers() []container {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	containers := []container{}
	for _, c := range listContainers(client, true) {
		name := c.Names[0][1:]
		status := c.State
		c := container{name: name, state: status, id: c.ID, ancestor: c.Image}
		containers = append(containers, c)
	}
	return containers
}

func getImages() []image {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	images := []image{}
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
			c := image{name: name, id: c.ID}
			images = append(images, c)
		}
	}
	return images
}
