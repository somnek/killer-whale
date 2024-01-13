package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/mattn/go-runewidth"
)

func buildContainerDescShort(c container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	container, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: c.id,
	})
	if err != nil {
		log.Fatal(err)
	}
	desc := fmt.Sprintf("ID    : %v\n", runewidth.Truncate(container.ID, fixedBodyRWidth-6, "..."))
	desc += fmt.Sprintf("Image : %s\n", container.Config.Image)
	desc += fmt.Sprintf("Cmd   : %s\n", strings.Join(container.Config.Cmd, " "))
	desc += fmt.Sprintf("State : %s\n", container.State.String())
	desc += fmt.Sprintf("Ports : %v\n", container.NetworkSettings.Ports)
	desc += fmt.Sprintf("IP    : %s\n", container.NetworkSettings.IPAddress)
	return desc
}
func buildContainerDescFull(c container) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	container, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: c.id,
	})
	if err != nil {
		log.Fatal(err)
	}
	desc := fmt.Sprintf("ID    : %v\n", runewidth.Truncate(container.ID, fixedBodyRWidth-8, "..."))
	desc += fmt.Sprintf("Image: %s\n", container.Config.Image)
	desc += fmt.Sprintf("Cmd: %s\n", strings.Join(container.Config.Cmd, " "))
	desc += fmt.Sprintf("Created: %s\n", container.Created)
	desc += fmt.Sprintf("State: %s\n", container.State.String())
	desc += fmt.Sprintf("Ports: %v\n", container.NetworkSettings.Ports)
	desc += fmt.Sprintf("Mounts: %v\n", container.Mounts)
	desc += fmt.Sprintf("Labels: %v\n", container.Config.Labels)
	desc += fmt.Sprintf("Env: %v\n", container.Config.Env)
	desc += fmt.Sprintf("HostConfig: %v\n", container.HostConfig)
	desc += fmt.Sprintf("NetworkSettings: %v\n", container.NetworkSettings)
	desc += fmt.Sprintf("LogPath: %s\n", container.LogPath)
	desc += fmt.Sprintf("RestartCount: %d\n", container.RestartCount)
	desc += fmt.Sprintf("Driver: %s\n", container.Driver)
	desc += fmt.Sprintf("Platform: %s\n", container.Platform)
	desc += fmt.Sprintf("ProcessLabel: %s\n", container.ProcessLabel)
	desc += fmt.Sprintf("IP: %s\n", container.NetworkSettings.IPAddress)
	return desc
}

func buildImageView(m model) string {

	var s string
	// truncate
	shouldTruncate := false
	imageList := m.images
	if len(m.images) > 10 {
		shouldTruncate = true
		imageList = m.images[:10]
	}

	for i, choice := range imageList {
		cursor := "  " // default cursor
		if m.cursor == i {
			cursor = "👉"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		// limit to 25 characters for now
		// TODO: make this dynamic
		name := choice.name
		if len(name) > 25 {
			name = name[:25] + "..."
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, name)
	}
	// add more images message
	if shouldTruncate {
		s += fmt.Sprintf("\n     %d more images... ↓\n", len(m.images)-10)
	}
	return bodyLStyle.Render(s)
}

func buildLogView(m model) string {
	var s string
	s += m.logs
	padBodyHeight(&s, lipgloss.Height(m.logs)+1)
	return logStyle.Render(s)
}

func buildContainerView(m model) (string, string) {
	var bodyL, bodyR string
	for i, choice := range m.containers {
		cursor := " " // default cursor
		check := " "
		if m.cursor == i {
			cursor = "❯"

			bodyR = buildContainerDescShort(choice)
		}
		state := stateStyle[choice.state].Render("●")
		name := choice.name
		name = runewidth.Truncate(name, 25, "...")
		if _, ok := m.selected[i]; ok {
			check = styleCheck.Render("✔")
		}
		bodyL += fmt.Sprintf("%s %s %s %s", cursor, check, state, name) + "\n"
	}

	// pad body height
	padBodyHeight(&bodyL, len(m.containers)+2)
	return bodyLStyle.Render(bodyL), bodyRStyle.Render(bodyR)
}

func (m model) View() string {

	var final string
	var bodyL, bodyR, body string

	// body L
	switch m.page {
	case pageContainer:
		bodyL, bodyR = buildContainerView(m)
	case pageImage:
		bodyL = buildImageView(m)
	}

	//  title
	title := strings.Repeat(" ", 36) + "🐳 Docker"
	titleStyle.MarginLeft((m.width - (fixedWidth + lipgloss.Width(title)/2)) / 2)
	title = titleStyle.Render(title)

	// join left + right component
	body = lipgloss.JoinHorizontal(lipgloss.Left, bodyL, bodyR)
	bodyStyle = bodyStyle.MarginLeft((m.width - fixedWidth) / 2)
	body = bodyStyle.Render(body)

	// joing title + body + help
	final += lipgloss.JoinVertical(lipgloss.Top, title, body)
	return final + "\n"
}

func padBodyHeight(s *string, itemCount int) {
	if itemCount < minHeightPerView {
		*s += strings.Repeat("\n", minHeightPerView-itemCount)
	}
}

// body R
// bodyR = bodyRStyle.Render(buildLogView(m))
