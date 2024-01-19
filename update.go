package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/fsouza/go-dockerclient"
)

func (m model) togglePageKey() keyMap {
	switch m.page {
	case pageImage:
		m.keys.Restart.Unbind()
		m.keys.Kill.Unbind()
		m.keys.Stop.Unbind()
		m.keys.Start.Unbind()
		m.keys.Pause.Unbind()
		m.keys.Unpause.Unbind()
	case pageContainer:
		m.keys = keys // keys is default (container)
	}
	return m.keys
}

func getCurrentViewItemCount(m model) int {
	var itemCount int
	switch m.page {
	case pageContainer:
		itemCount = len(m.containers)
	case pageImage:
		itemCount = len(m.images)
	}
	return itemCount
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// case spinner.TickMsg:
	// 	m.spinner, cmd = m.spinner.Update(msg)
	// 	return m, cmd

	case TickMsg:
		// containers
		containers := getContainers()
		m.containers = containers

		// images
		images := getImages()
		m.images = images

		// cursor
		if m.cursor == -1 {
			m.cursor = 0
		}

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
		switch m.page {
		case pageContainer:
			if getCurrentViewItemCount(m) > 0 {
				return handleContainerKeys(m, msg)
			}
			return handleCommonKeys(&m, msg)
		case pageImage:
			if getCurrentViewItemCount(m) > 0 {
				return handleImageKeys(m, msg)
			}
			return handleCommonKeys(&m, msg)
		}

		handleCommonKeys(&m, msg)
	}

	return m, cmd
}

// getContainers return a list of Container that are created using
// this image
func (img Image) findAssociatedContainersInUse(m model) []Container {
	containers := []Container{}
	for _, c := range m.containers {
		if c.ancestor == img.name && (c.state == "running" || c.state == "paused") {
			containers = append(containers, c)
		}
	}
	return containers
}

func handleImageKeys(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// handle 0 images
	if getCurrentViewItemCount(m) == 0 {
		return handleCommonKeys(&m, msg)
	}

	switch {
	case key.Matches(msg, m.keys.Remove): // remove
		client, err := docker.NewClientFromEnv()
		if err != nil {
			m.logs += "Failed to create Docker client"
		}

		targets := []Image{}
		if len(m.selected) == 0 {
			targets = append(targets, m.images[m.cursor])
		} else {
			for k := range m.selected {
				targets = append(targets, m.images[k])
			}
		}

		res := actionResultImages{}
		// for now show 1 dependent erorr at a time
		for _, img := range targets {
			containersInUse := img.findAssociatedContainersInUse(m)
			if len(containersInUse) > 0 {
				res.failed = append(res.failed, img)
				res.associatedContainers = containersInUse
			} else {
				go removeImage(client, img.id)
				desiredState := "x"
				addProcess(&m, img.id, desiredState)
				res.success = append(res.success, img)
			}
		}

		var logs string
		successCount, failedCount := len(res.success), len(res.failed)

		if successCount > 0 {
			logs += fmt.Sprintf(
				"🗑️ Remove %v image(s)\n",
				itemCountStyle.Render(fmt.Sprintf("%d", successCount)))
		}

		if failedCount > 0 {
			logs += fmt.Sprintf(
				"🚧 Skip removing %v image(s), can only remove image that are not in use...\n",
				itemCountStyle.Render(fmt.Sprintf("%d", failedCount)))
		}

		m.logs = logs
		m.selected = make(map[int]struct{})
		m.cursor = -1
		return m, cmd

	default:
		return handleCommonKeys(&m, msg)
	}
}

func handleContainerKeys(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	// handle 0 images
	if getCurrentViewItemCount(m) == 0 {
		return handleCommonKeys(&m, msg)
	}

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
	default:
		return handleCommonKeys(&m, msg)
	}

}

func handleCommonKeys(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.SelectAll): // select all
		// select/clear based on current page
		var items []any // container|image
		switch m.page {
		case pageContainer:
			items = make([]any, len(m.containers))
			for i, container := range m.containers {
				items[i] = container
			}
		case pageImage:
			items = make([]any, len(m.images))
			for i, image := range m.images {
				items[i] = image
			}
		}

		if len(items) == len(m.selected) {
			m.selected = make(map[int]struct{})
		} else {
			for i := range items {
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
		itemCount := getCurrentViewItemCount(*m)
		// increment cursor unless we're at the beginning of the list
		if m.cursor > 0 {
			m.cursor--
		} else {
			m.cursor = itemCount - 1
		}

	case key.Matches(msg, m.keys.Down): // move cursor down
		itemCount := getCurrentViewItemCount(*m)

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
		m.keys = m.togglePageKey()
		m.selected = make(map[int]struct{})
		m.cursor = 0
		m.logs = ""
	}
	return *m, nil
}
