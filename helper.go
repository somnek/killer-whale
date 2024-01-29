package main

import (
	"fmt"
	"log"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/fsouza/go-dockerclient"
)

type actionResultContainers struct {
	success []Container
	failed  []Container
}

type actionResultImages struct {
	success              []Image
	failed               []Image
	associatedContainers []Container
}

const (
	on int = iota
	off
)

// checkProcess check if the container is in m.processes
// (a.k.a process is in progress)
func checkProcess(id string, processes map[string]string) bool {
	if _, ok := processes[id]; ok {
		return true
	}
	return false
}

// updatePendingProcesses cross check the actual state of container
// with the desired state in m.processes
// if the state match, remove the container from m.processes
// implies that the action is completed
// return the updated m.processes
func updatePendingProcesses(m model) map[string]string {
	containers := m.containers
	for _, c := range containers {
		id := c.id
		state := c.state
		desiredState := m.processes[id]
		if _, ok := m.processes[id]; ok {
			if state == desiredState {
				delete(m.processes, id)
			}
		}
	}
	return m.processes
}

// addProcess add Process to m.processes
// container Processes are used to control the blinkSwitch
func addProcess(m *model, id, desiredState string) {
	m.processes[id] = desiredState
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

	res := actionResultContainers{}
	for _, c := range targets {
		if c.state == "paused" {
			go unpauseContainer(client, c.id)
			desiredState := "running"
			addProcess(&m, c.id, desiredState)
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

	res := actionResultContainers{}
	for _, c := range targets {
		if c.state == "running" {
			go pauseContainer(client, c.id)
			desiredState := "paused"
			addProcess(&m, c.id, desiredState)
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
			"🚧 Unable to pause %v container(s), skipping...\n",
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

	res := actionResultContainers{}
	for _, c := range targets {
		if c.state == "running" || c.state == "restarting" {
			go stopContainer(client, c.id)
			desiredState := "exited"
			addProcess(&m, c.id, desiredState)
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
			"🚧 Unable to stop %v container(s), skipping...\n",
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

	res := actionResultContainers{}
	for _, c := range targets {
		if c.state == "exited" || c.state == "created" {
			go startContainer(client, c.id)
			desiredState := "running"
			addProcess(&m, c.id, desiredState)
			res.success = append(res.success, c)

		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🚀 Start %v container(s)\n",
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

	res := actionResultContainers{}
	for _, c := range targets {
		removeContainer(client, c.id)
		desiredState := "x"
		addProcess(&m, c.id, desiredState)
		res.success = append(res.success, c)
	}

	var logs string
	successCount := len(res.success)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🔫 Removed %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	m.logs = logs
	m.selected = make(map[int]struct{})
	// prevent pointing to an nil index
	m.cursor = -1
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

	res := actionResultContainers{}
	for _, c := range targets {
		if c.state == "running" {
			go restartContainer(client, c.id)
			desiredState := "running"
			addProcess(&m, c.id, desiredState)
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

	res := actionResultContainers{}
	for _, c := range targets {
		if c.state == "running" {
			go killContainer(client, c.id)
			desiredState := "exited"
			addProcess(&m, c.id, desiredState)
			res.success = append(res.success, c)
		} else {
			res.failed = append(res.failed, c)
		}
	}

	var logs string
	successCount, failedCount := len(res.success), len(res.failed)

	if successCount > 0 {
		logs += fmt.Sprintf(
			"🔪 Killing %v container(s)\n",
			itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
	}

	if failedCount > 0 {
		logs += fmt.Sprintf(
			"🚧 Skip killing %v container(s), can only kill running container...\n",
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
		c := Container{
			name:     name,
			state:    status,
			id:       c.ID,
			ancestor: c.Image,
		}
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
		if len(tags) > 0 {
			name = tags[0]
			c := Image{name: name, id: c.ID}
			images = append(images, c)
		}
	}
	return images
}

// ---------------- Volume ----------------
func filterContainersByVolume(c *docker.Client, volName string) []docker.APIContainers {
	opts := docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"volume": {volName},
		},
	}
	// TODO: handle error
	containers, _ := c.ListContainers(opts)
	return containers
}

func getVolumes() []Volume {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed to create Docker clinet: %v", err)
	}

	volumes := []Volume{}

	for _, v := range listVolumes(client) {
		// find containers using the volume
		containers := filterContainersByVolume(client, v.Name)

		volume := Volume{
			name:       v.Name,
			mountPoint: v.Mountpoint,
			containers: containers,
			createdAt:  v.CreatedAt,
		}

		volumes = append([]Volume{volume}, volumes...)

	}

	// sort: show newest volume at the top
	sort.Slice(volumes, func(i, j int) bool {
		vi, vj := volumes[i], volumes[j]
		// sort by volume name if vi & vj has the same createdAt
		if vi.createdAt.Equal(vj.createdAt) {
			return vi.name < vj.name
		}
		return volumes[i].createdAt.After(volumes[j].createdAt)
	})

	return volumes
}
