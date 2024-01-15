package main

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
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
	containers  []Container
	images      []Image
	cursor      int
	selected    map[int]struct{}
	spinner     spinner.Model
	blinkSwitch int
	// TODO: merge process into Container struct
	processes map[string]string // map[containerID]desiredState
	keys      keyMap
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
	return tea.Batch(doTick(), m.spinner.Tick)
}

func initialModel() model {
	cursor := 0
	// containers
	containers := getContainers()
	images := getImages()
	containers[cursor].desc = buildContainerDescShort(containers[cursor])

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Jump
	s.Style = spinnerStyle

	// processes
	processes := make(map[string]string)
	return model{
		containers: containers,
		images:     images,
		selected:   make(map[int]struct{}),
		processes:  processes,
		spinner:    s,
		page:       pageContainer,
		keys:       keys,
	}
}
