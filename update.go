package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func getCurrentViewItemCount(m model) int {
	var itemCount int
	if m.page == pageContainer {
		itemCount = len(m.containers)
	} else {
		itemCount = len(m.images)
	}
	return itemCount
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case TickMsg:
		// containers
		containers := getContainers()
		m.containers = containers
		// images
		images := getImages()
		m.images = images

		// blink switch
		if m.blinkSwitch == on {
			m.blinkSwitch = off
		} else {
			m.blinkSwitch = on
		}

		// processes
		m.processes = updatePendingProcesses(m)

		return m, doTick()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Remove): // remove
			return removeAndWriteLog(m)

		case key.Matches(msg, m.keys.Restart): // restart
			return restartAndWriteLog(m)

		case key.Matches(msg, m.keys.Kill): // kill
			return killAndWriteLog(m)

		case key.Matches(msg, m.keys.Stop): // stop
			return stopAndWriteLog(m)

		case key.Matches(msg, m.keys.Start): // start
			return startAndWriteLog(m)

		case key.Matches(msg, m.keys.Pause): // pause
			return pauseAndWriteLog(m)

		case key.Matches(msg, m.keys.Unpause): // unpause
			return unpauseAndWriteLog(m)

		case key.Matches(msg, m.keys.SelectAll): // slecet all
			if len(m.containers) == len(m.selected) {
				m.selected = make(map[int]struct{})
			} else {
				for i := range m.containers {
					m.selected[i] = struct{}{}
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Clear): // clear selection
			m.logs = ""
			m.selected = make(map[int]struct{})
			return m, nil

		case key.Matches(msg, m.keys.Quit): // quit
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up): // move cursor up
			itemCount := getCurrentViewItemCount(m)
			// increment cursor unless we're at the beginning of the list
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = itemCount - 1
			}

		case key.Matches(msg, m.keys.Down): // move cursor down
			itemCount := getCurrentViewItemCount(m)

			// decrement cursor unless we're at the end of the list
			if m.cursor < itemCount-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case key.Matches(msg, m.keys.Toggle): // toggle selection
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			m.logs = ""

		case key.Matches(msg, m.keys.Help): // toggle help
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keys.Tab): // switch tab
			if m.page == pageContainer {
				m.page = pageImage
			} else {
				m.page = pageContainer
			}
			m.cursor = 0
		}

	}
	return m, nil
}
