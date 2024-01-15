package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/mattn/go-runewidth"
)

func formatImageVolumes(volumeMap map[string]struct{}) string {
	var s string
	for vol := range volumeMap {
		s += fmt.Sprintf("%s\n", vol)
	}
	s = strings.TrimSuffix(s, "\n")
	return s

}

func formatMounts(mounts []docker.Mount) string {
	var s string
	for _, mount := range mounts {
		s += fmt.Sprintf("(source) %s\n", mount.Source)
		s += fmt.Sprintf("        (des)    %s\n", mount.Destination)
	}
	s = strings.TrimSuffix(s, "\n")
	return s
}

func formatPortsMapping(portsMap map[docker.Port][]docker.PortBinding) string {
	var s string
	for containerPort, hostMachinePorts := range portsMap {
		s += fmt.Sprintf("%s (container)\n", containerPort)
		for _, port := range hostMachinePorts {
			s += fmt.Sprintf("        -> %s:%s (host)\n", port.HostIP, port.HostPort)
		}
	}
	s = strings.TrimSuffix(s, "\n")
	return s
}

func buildContainerDescShort(id string) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	container, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: id,
	})
	if err != nil {
		log.Fatal(err)
	}
	desc := fmt.Sprintf("ID    : %v\n", runewidth.Truncate(container.ID, fixedBodyRWidth-8, "..."))
	desc += fmt.Sprintf("Image : %s\n", container.Config.Image)
	desc += fmt.Sprintf("Cmd   : %s\n", strings.Join(container.Config.Cmd, " "))
	desc += fmt.Sprintf("State : %s\n", container.State.String())
	desc += fmt.Sprintf("Ports : %v\n", formatPortsMapping(container.NetworkSettings.Ports))
	desc += fmt.Sprintf("IP    : %s\n", container.NetworkSettings.IPAddress)
	return desc
}
func buildContainerDescFull(id string) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	container, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{
		ID: id,
	})
	if err != nil {
		log.Fatal(err)
	}
	desc := fmt.Sprintf("ID              : %v\n", runewidth.Truncate(container.ID, fixedBodyRWidth-8, "..."))
	desc += fmt.Sprintf("Image           : %s\n", container.Config.Image)
	desc += fmt.Sprintf("Cmd             : %s\n", strings.Join(container.Config.Cmd, " "))
	desc += fmt.Sprintf("Created         : %s\n", container.Created.Format("2006-01-02 15:04:05"))
	desc += fmt.Sprintf("State           : %s\n", container.State.String())
	desc += fmt.Sprintf("Ports           : %v\n", formatPortsMapping(container.NetworkSettings.Ports))
	desc += fmt.Sprintf("Mounts          : %v\n", formatMounts(container.Mounts))
	desc += fmt.Sprintf("Labels          : %v\n", container.Config.Labels)
	desc += fmt.Sprintf("Env             : %v\n", container.Config.Env)
	desc += fmt.Sprintf("HostConfig      : %v\n", container.HostConfig)
	desc += fmt.Sprintf("NetworkSettings : %v\n", container.NetworkSettings)
	desc += fmt.Sprintf("LogPath         : %s\n", container.LogPath)
	desc += fmt.Sprintf("RestartCount    : %d\n", container.RestartCount)
	desc += fmt.Sprintf("Driver          : %s\n", container.Driver)
	desc += fmt.Sprintf("Platform        : %s\n", container.Platform)
	desc += fmt.Sprintf("ProcessLabel    : %s\n", container.ProcessLabel)
	desc += fmt.Sprintf("IP              : %s\n", container.NetworkSettings.IPAddress)
	return desc
}

func buildImageDescShort(id string) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	image, err := client.InspectImage(id)
	if err != nil {
		log.Fatal(err)
	}
	desc := fmt.Sprintf("ID      : %v\n", runewidth.Truncate(image.ID, fixedBodyRWidth-6, "..."))
	desc += fmt.Sprintf("Created : %s\n", image.Created.Format("2006-01-02 15:04:05"))
	desc += fmt.Sprintf("Size    : %s\n", convertSizeToHumanRedable(image.Size))
	desc += fmt.Sprintf("Cmd     : %s\n", strings.Join(image.Config.Cmd, " "))
	desc += fmt.Sprintf("Volumes : %v\n", formatImageVolumes(image.Config.Volumes))
	return desc
}

func buildImageView(m model) (string, string) {
	var bodyL, bodyR string
	for i, choice := range m.images {
		cursor := " " // default cursor
		check := " "
		if m.cursor == i {
			cursor = "❯"
			bodyR = buildImageDescShort(choice.id)
		}
		name := choice.name
		bodyL += fmt.Sprintf("%s %s %s", cursor, check, name) + "\n"
	}
	padBodyHeight(&bodyL, len(m.images)+2)
	return bodyLStyle.Render(bodyL), bodyRStyle.Render(bodyR)
}

func buildLogView(m model) string {
	var s string
	s += m.logs
	logStyle.MarginLeft((fixedWidth - lipgloss.Width(s)) / 2)
	logStyle.AlignHorizontal(lipgloss.Center)
	return logStyle.Render(s)
}

func buildContainerView(m model) (string, string) {
	var bodyL, bodyR string
	for i, choice := range m.containers {
		cursor := " " // default cursor
		check := " "
		if m.cursor == i {
			cursor = "❯"
			bodyR = buildContainerDescShort(choice.id)
		}

		isProcessing := checkProcess(choice.id, m.processes)
		stateStyle := stateStyleMap[choice.state]
		if isProcessing && m.blinkSwitch == on {
			stateStyle = stateStyle.Copy().Foreground(pitchBlack)
		}
		state := stateStyle.Render("●")
		name := choice.name
		if _, ok := m.selected[i]; ok {
			check = checkStyle.Render("✔")
		}
		bodyL += fmt.Sprintf("%s %s %s %s", cursor, check, state, name) + "\n"
	}

	// pad body height
	padBodyHeight(&bodyL, len(m.containers)+2)
	return bodyLStyle.Render(bodyL), bodyRStyle.Render(bodyR)
}

func (m model) View() string {

	var final string
	var bodyL, bodyR, body, bottom string

	// body L
	switch m.page {
	case pageContainer:
		bodyL, bodyR = buildContainerView(m)
	case pageImage:
		bodyL, bodyR = buildImageView(m)
	}

	//  title
	title := "🐳 Docker" + "  " + m.spinner.View()
	titleStyle.MarginLeft((m.width / 2) - (lipgloss.Width(title) / 2))
	title = titleStyle.Render(title)

	// join left + right component
	body = lipgloss.JoinHorizontal(lipgloss.Left, bodyL, bodyR)
	body = bodyStyle.Render(body)

	// bottom
	bottom = buildLogView(m)

	// joing title + body + log + help
	final += lipgloss.JoinVertical(lipgloss.Top, body, bottom)
	appStyle.MarginLeft((m.width - fixedWidth) / 2)

	return title + "\n" + appStyle.Render(final) + "\n"
}

func padBodyHeight(s *string, itemCount int) {
	if itemCount < minHeightPerView {
		*s += strings.Repeat("\n", minHeightPerView-itemCount)
	}
}
