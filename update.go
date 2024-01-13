package main

import (
	tea "github.com/charmbracelet/bubbletea"
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
			if len(m.selected) == 0 {
				container := m.containers[m.cursor]
				m.logs = removeAndWriteLog(container)
				return m, nil
			}

			// force for now
			for k := range m.selected {
				container := m.containers[k]
				m.logs = removeAndWriteLog(container)
			}
			m.selected = make(map[int]struct{})
			m.cursor = 0
			return m, nil

		case "r": // restart

			if len(m.selected) == 0 {
				container := m.containers[m.cursor]
				m.logs = restartAndWriteLog(container)
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				m.logs = restartAndWriteLog(container)
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "K": // kill
			if len(m.selected) == 0 {
				container := m.containers[m.cursor]
				m.logs = killAndWriteLog(container)
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				m.logs = killAndWriteLog(container)
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "s": // stop
			if len(m.selected) == 0 {
				container := m.containers[m.cursor]
				m.logs = stopAndWriteLog(container)
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				m.logs = stopAndWriteLog(container)
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "u": // up
			if len(m.selected) == 0 {
				container := m.containers[m.cursor]
				m.logs = startAndWriteLog(container)
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				m.logs = startAndWriteLog(container)
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "p": // pause
			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				m.logs = pauseAndWriteLog(container)
			}
			m.selected = make(map[int]struct{})
			return m, nil

		case "P": // unpause
			if len(m.selected) == 0 {
				container := m.containers[m.cursor]
				m.logs = unpauseAndWriteLog(container)
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				m.logs = unpauseAndWriteLog(container)
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
