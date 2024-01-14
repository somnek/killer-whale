package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/fsouza/go-dockerclient"
)

type actionResult struct {
	success []Container
	failed  []Container
}

func unpauseAndWriteLog(m model) (tea.Model, tea.Cmd) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	targets := []Container{}
	if len(m.selected) == 0 {
		targets = append(targets, m.containers[m.cursor])
	} else {
		for k := range m.selected {
			targets = append(targets, m.containers[k])
		}
	}

	res := actionResult{}
	for _, c := range targets {
		if c.state == "paused" {
			go unpauseContainer(client, c.id)
			res.success = append(res.success, c)
		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"✅ Unpaused %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	if failedCount > 0 {
		logs += fmt.Sprintf(
			"🚧 Skip unpausing %v container(s), can only unpausing paused container...\n",
			itemCountStyle.Render(fmt.Sprintf("%d", failedCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	return m, nil
}

func pauseAndWriteLog(m model) (tea.Model, tea.Cmd) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	targets := []Container{}
	if len(m.selected) == 0 {
		targets = append(targets, m.containers[m.cursor])
	} else {
		for k := range m.selected {
			targets = append(targets, m.containers[k])
		}
	}

	res := actionResult{}
	for _, c := range targets {
		if c.state == "running" {
			go pauseContainer(client, c.id)
			res.success = append(res.success, c)
		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"⏳ Paused  %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	if failedCount > 0 {
		logs += fmt.Sprintf(
			"🚧 %v container(s) is not running, skipping...\n",
			itemCountStyle.Render(fmt.Sprintf("%d", failedCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	return m, nil
}

func stopAndWriteLog(m model) (tea.Model, tea.Cmd) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	targets := []Container{}
	if len(m.selected) == 0 {
		targets = append(targets, m.containers[m.cursor])
	} else {
		for k := range m.selected {
			targets = append(targets, m.containers[k])
		}
	}

	res := actionResult{}
	for _, c := range targets {
		if c.state == "running" || c.state == "restarting" {
			go stopContainer(client, c.id)
			res.success = append(res.success, c)
		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🛑 Stopping %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	if failedCount > 0 {
		logs += fmt.Sprintf(
			"🚧 unable to stop %v container(s), skipping...\n",
			itemCountStyle.Render(fmt.Sprintf("%d", failedCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	return m, nil
}

func startAndWriteLog(m model) (tea.Model, tea.Cmd) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	targets := []Container{}
	if len(m.selected) == 0 {
		targets = append(targets, m.containers[m.cursor])
	} else {
		for k := range m.selected {
			targets = append(targets, m.containers[k])
		}
	}

	res := actionResult{}
	for _, c := range targets {
		if c.state == "exited" || c.state == "created" {
			go startContainer(client, c.id)
			res.success = append(res.success, c)
		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🚀 Starting %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	if failedCount > 0 {
		logs += fmt.Sprintf(
			"🚧 %v container(s) already running, skipping...\n",
			itemCountStyle.Render(fmt.Sprintf("%d", failedCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	return m, nil
}

func removeAndWriteLog(m model) (tea.Model, tea.Cmd) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	targets := []Container{}
	if len(m.selected) == 0 {
		targets = append(targets, m.containers[m.cursor])
	} else {
		for k := range m.selected {
			targets = append(targets, m.containers[k])
		}
	}

	res := actionResult{}
	for _, c := range targets {
		go removeContainer(client, c.id)
		res.success = append(res.success, c)
	}

	var logs string
	successCount := len(res.success)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🗑️  Removing %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	m.cursor = 0
	return m, nil
}

func restartAndWriteLog(m model) (tea.Model, tea.Cmd) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	targets := []Container{}
	if len(m.selected) == 0 {
		targets = append(targets, m.containers[m.cursor])
	} else {
		for k := range m.selected {
			targets = append(targets, m.containers[k])
		}
	}

	res := actionResult{}
	for _, c := range targets {
		if c.state == "running" {
			go restartContainer(client, c.id)
			res.success = append(res.success, c)
		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🌀 Restarted %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	if failedCount > 0 {
		logs += fmt.Sprintf(
			"🚧 Skip restarting %v containera(s), container must be in a running state...\n",
			itemCountStyle.Render(fmt.Sprintf("%d", failedCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	return m, nil
}

func killAndWriteLog(m model) (tea.Model, tea.Cmd) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}

	targets := []Container{}
	if len(m.selected) == 0 {
		targets = append(targets, m.containers[m.cursor])
	} else {
		for k := range m.selected {
			targets = append(targets, m.containers[k])
		}
	}

	res := actionResult{}
	for _, c := range targets {
		if c.state == "running" {
			killContainer(client, c.id)
			res.success = append(res.success, c)
		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🔪 Killed %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	if failedCount > 0 {
		logs += fmt.Sprintf(
			"🚧 skip killing %v container(s), can only kill running container...\n",
			itemCountStyle.Render(fmt.Sprintf("%d", failedCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	return m, nil
}

func getContainers() []Container {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
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
		log.Fatalf("failed to create Docker client: %v", err)
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
