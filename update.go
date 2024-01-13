package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/fsouza/go-dockerclient"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TickMsg:
		// containers
		containers := getContainers()
		m.containers = containers
		// images
		images := getImages()
		m.images = images
		return m, doTick()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "x": // remove
			m.logs = ""
			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			// force for now
			for k := range m.selected {
				container := m.containers[k]
				id := container.id
				go removeContainer(client, id)
				m.logs = "🗑️  Remove " + container.name + "\n"
			}
			m.selected = make(map[int]struct{})
			m.cursor = 0
			return m, nil

		case "r": // restart
			m.logs = ""
			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "running" {
					go restartContainer(client, id)
					m.logs = "🔃 Restarted " + container.name + "\n"
				} else {
					m.logs = "🚧  " + container.name + " not running\n"
				}
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "K": // kill
			m.logs = ""
			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "running" {
					killContainer(client, id)
					m.logs = "🔪 Killed " + container.name + "\n"
				} else {
					m.logs = "🚧 " + container.name + " already stopped\n"
				}
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "s": // stop
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "running" || state == "restarting" {
					go stopContainer(client, id)
					m.logs = "🛑 Stop " + container.name + "\n"
				} else {
					m.logs = "🚧  " + " unable to stop " + container.name + "\n"
				}
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "u": // up
			m.logs = ""
			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "exited" || state == "created" {
					go startContainer(client, id)
					if err != nil {
						m.logs = fmt.Sprintf("🚧  %s\n", err.Error())
					} else {
						m.logs = "🚀 Started " + container.name + "\n"
					}
				} else {
					m.logs = "🚧  " + container.name + " already running\n"
				}
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "p": // pause
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			for i, choice := range m.containers {
				id := choice.id
				state := choice.state
				if _, ok := m.selected[i]; ok {
					if state == "running" {
						pauseContainer(client, id)
						m.logs = "⏳ Paused " + choice.name + "\n"
					} else {
						m.logs = "🚧  " + choice.name + " is not running\n"
					}
				}
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "P": // unpause
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			for i, choice := range m.containers {
				id := choice.id
				state := choice.state
				if _, ok := m.selected[i]; ok {
					if state == "paused" {
						unPauseContainer(client, id)
						m.logs = "✅ Unpaused " + choice.name + "\n"
					} else {
						m.logs = "🚧  " + choice.name + " is not running\n"
					}
				}
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "ctrl+a", "A": // select all
			if len(m.containers) == len(m.selected) {
				m.selected = make(map[int]struct{})
			} else {
				for i := range m.containers {
					m.selected[i] = struct{}{}
				}
			}
			return m, nil

		case "esc": // clear selection
			m.logs = ""
			m.selected = make(map[int]struct{})
			return m, nil

		case "ctrl+c", "q": // quit
			return m, tea.Quit

		case "up", "k": // move cursor up
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.containers) - 1
			}

		case "down", "j": // move cursor down
			if m.cursor < len(m.containers)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case "enter", " ": // toggle selection
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			m.logs = ""

		case "?": // controls page
			if m.page != 3 {
				m.page = 3
			} else {
				m.page = 0
			}
			return m, nil

		case "tab":
			if m.page == pageContainer {
				m.page = pageImage
			} else {
				m.page = pageContainer
			}
		}

	}
	return m, nil
}
