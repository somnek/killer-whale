package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/wrap"
)

// ----------------------------- render utils -----------------------------

func padBodyHeight(s *string, itemCount int) {
	if itemCount < minHeightPerView {
		*s += strings.Repeat("\n", minHeightPerView-itemCount)
	}
}

func padHelpWidth(s *string, windowWidth, maxAppWidth int) {
	var outerPad, innerPad, longest int

	// get width of longer help string (fullHelp)
	split := strings.Split(*s, "\n")
	for _, line := range split {
		if lipgloss.Width(line) > longest {
			longest = lipgloss.Width(line)
		}
	}
	sWidth := longest

	if windowWidth > 0 && longest < maxAppWidth-4 {
		outerPad = (windowWidth - maxAppWidth) / 2
		innerPad = ((maxAppWidth - 4) - sWidth) / 2
	}

	var newS string
	for _, line := range split {
		newS += strings.Repeat(" ", outerPad+innerPad) + line + "\n"
	}
	*s = newS
}

func padItemWidth(s *string, maxWidth int) {
	sWidth := lipgloss.Width(*s)
	if sWidth < maxWidth-10 {
		*s = *s + strings.Repeat(" ", maxWidth-sWidth)
	}
	*s += "\n"
}

func formatCmd(cmd []string) string {
	s := wrap.String(strings.Join(cmd, " "), fixedBodyRWidth-10)
	split := strings.Split(s, "\n")

	if len(split) > 1 {
		s = split[0] + "\n"
		for _, line := range split[1:] {
			s += strings.Repeat(" ", 10) + line + "\n"
		}
	}
	s = strings.TrimSuffix(s, "\n")
	return s
}

func formatImageVolumes(volumeMap map[string]struct{}) string {
	var s string
	for vol := range volumeMap {
		s += fmt.Sprintf("%s\n", vol)
	}
	s = strings.TrimSuffix(s, "\n")
	return s
}

func getSortedPort(portsMap map[docker.Port][]docker.PortBinding) []docker.Port {
	l := []docker.Port{}
	for k := range portsMap {
		l = append(l, k)
	}

	sort.Slice(l, func(i, j int) bool {
		// convert to int before compare
		iInt, _ := strconv.Atoi(l[i].Port())
		jInt, _ := strconv.Atoi(l[j].Port())
		return iInt > jInt
	})
	return l
}

func formatPortsMapping(portsMap map[docker.Port][]docker.PortBinding) string {
	s := "\n"
	sortedPorts := getSortedPort(portsMap)

	for _, containerPort := range sortedPorts {
		// find matching host machine port
		var portBindings []docker.PortBinding // host machine
		for k, v := range portsMap {
			if k.Port() == containerPort.Port() {
				portBindings = v
			}
		}

		var joinedPortBindingStr string
		if len(portBindings) == 0 {
			joinedPortBindingStr = "null"
		}

		for i, bindings := range portBindings {
			IP, Port := bindings.HostIP, bindings.HostPort
			joinedPortBindingStr += fmt.Sprintf("%s:%s", IP, Port)
			// has more than 1 port binding per container port
			if i > 0 {
				joinedPortBindingStr += "\n"
			}
		}
		containerPortStr := fmt.Sprintf("%5s/%s", containerPort.Port(), containerPort.Proto())
		s += fmt.Sprintf("        %s -> %s\n", containerPortStr, joinedPortBindingStr)
	}

	// add column (container | hostmachine)
	if len(sortedPorts) > 0 {
		s = fmt.Sprintf("%s -> %s", PortMapColStyle.Render("container"), PortMapColStyle.Render("host machine")) + s
	}

	s = strings.TrimSuffix(s, "\n")
	return s
}

// ----------------------------- log view -----------------------------

func buildLogView(m model) string {
	var s string
	s += m.logs
	logStyle.MarginLeft(((fixedWidth - 4) - lipgloss.Width(s)) / 2)
	logStyle.AlignHorizontal(lipgloss.Center)
	return logStyle.Render(s)
}

// ----------------------------- image view -----------------------------

func buildImageDescShort(id string) string {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	image, err := client.InspectImage(id)
	if err != nil {
		log.Fatal(err)
	}
	desc := fmt.Sprintf("ID      : %v\n", runewidth.Truncate(image.ID, fixedBodyRWidth-8, "..."))
	desc += fmt.Sprintf("Created : %s\n", image.Created.Format("2006-01-02 15:04:05"))
	desc += fmt.Sprintf("Size    : %s\n", convertSizeToHumanRedable(image.Size))
	desc += fmt.Sprintf("Cmd     : %v\n", formatCmd(image.Config.Cmd))
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
		if _, ok := m.selected[i]; ok {
			check = checkStyle.Render("✔")
		}
		row := fmt.Sprintf("%s %s %s", cursor, check, name)
		padItemWidth(&row, fixedBodyLWidth-8)
		bodyL += row
	}
	padBodyHeight(&bodyL, len(m.images)+2)
	return bodyLStyle.Render(bodyL), bodyRStyle.Render(bodyR)
}

// ----------------------------- container view -----------------------------

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
	desc += fmt.Sprintf("IP    : %s\n", container.NetworkSettings.IPAddress)
	desc += fmt.Sprintf("Ports : %v\n", formatPortsMapping(container.NetworkSettings.Ports))
	return desc
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
		name = runewidth.Truncate(name, maxItemNameWidth, "...")
		if _, ok := m.selected[i]; ok {
			check = checkStyle.Render("✔")
		}
		row := fmt.Sprintf("%s %s %s %s", cursor, check, state, name)
		padItemWidth(&row, fixedBodyLWidth-10)
		bodyL += row
	}

	// pad body height
	padBodyHeight(&bodyL, len(m.containers)+2)
	return bodyLStyle.Render(bodyL), bodyRStyle.Render(bodyR)
}

// ----------------------------- main view -----------------------------

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
	title := "🐳 Killer Whale" + "  "
	titleStyle.MarginLeft((m.width / 2) - (lipgloss.Width(title) / 2))
	title = titleStyle.Render(title)

	// join left + right component
	body = lipgloss.JoinHorizontal(lipgloss.Left, bodyL, bodyR)
	body = bodyStyle.Render(body)

	// bottom
	bottom = buildLogView(m)

	// help
	help := m.help.View(m.keys)
	padHelpWidth(&help, m.width, fixedWidth)

	// join title + body + log + help
	final += lipgloss.JoinVertical(lipgloss.Top, body, bottom)
	appStyle.MarginLeft((m.width - fixedWidth) / 2)

	// 0 containers/ image
	if len(m.containers) == 0 && m.page == pageContainer {
		body = bodyStyle.Render("No containers found")
		return title + "\n" + titleStyle.Render(body) + "\n"
	} else if len(m.images) == 0 && m.page == pageImage {
		body = bodyStyle.Render("No images found")
		return title + "\n" + titleStyle.Render(body) + "\n"
	}

	return title + "\n" + appStyle.Render(final) + "\n" + help
}
