package main

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

type container struct {
	name     string
	state    string
	id       string
	ancestor string
}

type image struct {
	name string
	id   string
}

const (
	pageContainer int = iota
	pageImage
	pageLog
)

type model struct {
	containers []container
	images     []image
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
	containers := getContainers()
	images := getImages()
	return model{
		containers: containers,
		images:     images,
		selected:   make(map[int]struct{}),
		page:       pageContainer,
	}
}
