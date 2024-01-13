package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Container struct {
	name     string
	state    string
	id       string
	ancestor string
	desc     string
}

type Image struct {
	name string
	id   string
}

const (
	pageContainer int = iota
	pageImage
	pageLog
)

type model struct {
	containers []Container
	images     []Image
	cursor     int
	selected   map[int]struct{}
	logs       string
	page       int
	width      int
	height     int
}

type TickMsg struct {
	Time time.Time
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg{Time: t}
	})
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func initialModel() model {
	cursor := 0
	containers := getContainers()
	images := getImages()
	containers[cursor].desc = buildContainerDescShort(containers[cursor])
	return model{
		containers: containers,
		images:     images,
		selected:   make(map[int]struct{}),
		page:       pageContainer,
	}
}
