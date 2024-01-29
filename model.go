package main

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/fsouza/go-dockerclient"
)

type Container struct {
	name     string
	state    string
	id       string
	ancestor string
	desc     string
}

type Volume struct {
	name       string
	mountPoint string
	containers []docker.APIContainers
	createdAt  time.Time
}

type Image struct {
	id   string
	name string
}

const (
	pageContainer int = iota
	pageImage
	pageVolume
	pageLog
)

type model struct {
	containers  []Container
	images      []Image
	volumes     []Volume
	cursor      int
	selected    map[int]struct{}
	blinkSwitch int
	// TODO: merge process into Container struct
	processes map[string]string // map[containerID]desiredState
	keys      keyMap
	help      help.Model
	logs      string
	page      int
	width     int
	height    int
}

// fast tick rate doesn't seems to affect performance (average 20 container)
// TODO: change tickrate to 1s if theres no running docker action
const tickRate = 300 * time.Millisecond

// const tickRate = time.Second

type TickMsg struct {
	Time time.Time
}

func doTick() tea.Cmd {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return TickMsg{Time: t}
	})
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func initialModel() model {
	cursor := 0

	// containers
	containers := getContainers()
	images := getImages()
	volumes := getVolumes()

	// descriptions of container at cursor
	if len(containers) > 0 {
		containers[cursor].desc = buildContainerDescShort(containers[cursor].id)
	}

	// help
	h := help.New()
	h.Width = fullWidth

	// processes
	processes := make(map[string]string)
	return model{
		cursor:     0,
		containers: containers,
		images:     images,
		volumes:    volumes,
		selected:   make(map[int]struct{}),
		processes:  processes,
		page:       pageContainer,
		keys:       keys,
		help:       h,
	}
}
