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

func buildEmptyBody(text, title string, width int) string {
	emptyBody := text
	padOuterComponent(&emptyBody, width)
	emptyBody = strings.TrimSuffix(emptyBody, "\n")
	return fmt.Sprintf("%s\n%s\n", title, emptyBody)
}

func padBodyHeight(s *string, itemCount int) {
	if itemCount < minHeightPerView {
		*s += strings.Repeat("\n", minHeightPerView-itemCount)
	}
}

func padOuterComponent(s *string, windowWidth int) {
	var outerPad, innerPad, longest int

	// get width of longer help string (fullHelp)
	split := strings.Split(*s, "\n")
	for _, line := range split {
		if lipgloss.Width(line) > longest {
			longest = lipgloss.Width(line)
		}
	}
	sWidth := longest

	if windowWidth > 0 {
		outerPad = (windowWidth - fullWidth) / 2
		innerPad = (fullWidth - sWidth) / 2
	}

	var newS string
	for _, line := range split {
		newS += strings.Repeat(" ", outerPad+innerPad) + line + "\n"
	}
	*s = newS
}

func padItemName(name string, maxLen int) string {
	nameWidth := lipgloss.Width(name)
	if nameWidth < maxLen {
		name = name + strings.Repeat(" ", maxLen-nameWidth)
	}
	name += "\n"
	return name
}

func buildTitleView(m model) string {
	s := "🐳 Killer Whale" + "  "
	padOuterComponent(&s, m.width)
	s = strings.TrimSuffix(s, "\n")
	return s
}

// ----------------------------- log view -----------------------------

func buildLogView(m model) string {
	var s string
	s += m.logs
	logStyle.MarginLeft((fullWidth - lipgloss.Width(s)) / 2)
	logStyle.AlignHorizontal(lipgloss.Center)
	return logStyle.Render(s)
}

// ----------------------------- volume view -----------------------------
func formatVolumeMountPoint(mp string) string {
	s := wrap.String(mp, fixedBodyRWidth-8)
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
func formatVolumeInUse(b bool) string {
	s := fmt.Sprintf("%v", b)
	var style lipgloss.Style

	if b {
		style = inUseStyleTrue
	} else {
		style = inUseStyleFalse
	}
	return style.Render(s)
}

func formatVolumeName(name string) string {
	s := wrap.String(name, fixedBodyRWidth-8)
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

func buildVolumeDescShort(volume Volume) string {
	name := volume.name
	createdAt := volume.createdAt
	containers := volume.containers
	mountPoint := volume.mountPoint

	// get first element from result
	// TODO: display list of containers, might have >1 container per volume
	var inUse bool
	containerName := "null"

	if len(containers) > 0 {
		inUse = true
		containerName = containers[0].Names[0][1:]
	}

	var desc string
	desc += fmt.Sprintf("Created : %s\n", createdAt)
	desc += fmt.Sprintf("Mount   : %s\n", formatVolumeMountPoint(mountPoint))
	desc += fmt.Sprintf("Name    : %s\n", formatVolumeName(name))
	desc += fmt.Sprintf("In Use  : %v\n", formatVolumeInUse(inUse))
	desc += fmt.Sprintf("Use by  : %s\n", containerName)

	return desc
}

func buildVolumeView(m model) (string, string) {
	var bodyL, bodyR string
	for i, choice := range m.volumes {
		cursor := " "
		check := " "
		if m.cursor == i {
			cursor = "❯"
			bodyR = buildVolumeDescShort(choice)
		}

		name := runewidth.Truncate(choice.name, maxImageNameWidth, "")
		if len(choice.containers) > 0 {
			name = volumeItemInUseStyle.Render(name)
		}

		if _, ok := m.selected[i]; ok {
			check = checkStyle.Render("✔")
		}
		row := fmt.Sprintf("%s %s %s", cursor, check, name)
		bodyL += row + "\n"
	}

	return bodyLStyle.Render(bodyL), bodyRStyle.Render(bodyR)
}

// ----------------------------- image view -----------------------------

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
		name = padItemName(name, maxImageNameWidth)
		row := fmt.Sprintf("%s %s %s", cursor, check, name)
		bodyL += row
	}
	padBodyHeight(&bodyL, len(m.images)+2)
	return bodyLStyle.Render(bodyL), bodyRStyle.Render(bodyR)
}

// ----------------------------- container view -----------------------------

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
		s += fmt.Sprintf("          %s -> %s\n", containerPortStr, joinedPortBindingStr)
	}

	// add column (container | hostmachine)
	if len(sortedPorts) > 0 {
		s = fmt.Sprintf(
			"%s -> %s",
			PortMapColStyle.Render("container"),
			PortMapColStyle.Render("host machine"),
		) + s
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
	desc := fmt.Sprintf(
		"ID      : %v\n",
		runewidth.Truncate(container.ID, fixedBodyRWidth-8, "..."),
	)
	desc += fmt.Sprintf("Image   : %s\n", container.Config.Image)
	desc += fmt.Sprintf("Cmd     : %s\n", strings.Join(container.Config.Cmd, " "))
	desc += fmt.Sprintf("State   : %s\n", container.State.String())
	desc += fmt.Sprintf("IP      : %s\n", container.NetworkSettings.IPAddress)
	desc += fmt.Sprintf("Ports   : %v\n", formatPortsMapping(container.NetworkSettings.Ports))
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
		// state := stateStyle.Render("●")
		state := stateStyle.Render("❖")
		name := choice.name
		name = runewidth.Truncate(name, maxContainerNameWidth, "...")
		if _, ok := m.selected[i]; ok {
			check = checkStyle.Render("✔")
		}
		name = padItemName(name, maxContainerNameWidth)
		row := fmt.Sprintf("%s %s %s %s", cursor, check, state, name)
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
	case pageVolume:
		bodyL, bodyR = buildVolumeView(m)
	}

	//  title
	title := buildTitleView(m)
	title = titleStyle.Render(title)

	// join left + right component
	body = lipgloss.JoinHorizontal(lipgloss.Left, bodyL, bodyR)
	body = bodyStyle.Render(body)

	// bottom
	bottom = buildLogView(m)

	// help
	help := m.help.View(m.keys)
	padOuterComponent(&help, m.width)

	// join title + body + log + help
	final += lipgloss.JoinVertical(lipgloss.Top, body, bottom)
	appStyle.MarginLeft((m.width - fullWidth) / 2)

	// 0 containers/ image
	if len(m.containers) == 0 && m.page == pageContainer {
		return buildEmptyBody("\nNo containers found.", title, m.width)
	} else if len(m.images) == 0 && m.page == pageImage {
		return buildEmptyBody("\nNo images found.", title, m.width)
	}

	return title + "\n" + appStyle.Render(final) + "\n" + help
}
