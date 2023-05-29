package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/fsouza/go-dockerclient"
)

type model struct {
	containers []container
	images     []image
	cursor     int
	selected   map[int]struct{}
	logs       string
	page       int
}

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

/* STYLING */
const (
	white        = lipgloss.Color("#F5F5F5")
	green        = lipgloss.Color("#D0F1BF")
	hotGreen     = lipgloss.Color("#73F59F")
	lightBlue    = lipgloss.Color("#C1E0F7")
	midBlue      = lipgloss.Color("#A4DEF9")
	frenchBlue   = lipgloss.Color("#0072BB")
	celesBlue    = lipgloss.Color("#1E91D6")
	electricBlue = lipgloss.Color("#2DE1FC")
	lightPurple  = lipgloss.Color("#CFBAE1")
	midPurple    = lipgloss.Color("#C59FC9")
	yellow       = lipgloss.Color("#F4E3B2")
	orange       = lipgloss.Color("#EFC88B")
	red          = lipgloss.Color("#FF5A5F")
	grey         = lipgloss.Color("#A0A0A0")
	black        = lipgloss.Color("#3C3C3C")
	lightPink    = lipgloss.Color("#F9CFF2")
	midPink      = lipgloss.Color("#F786AA")
)

var (
	stateStyle = map[string]lipgloss.Style{
		"created":    lipgloss.NewStyle().Background(midPurple),
		"running":    lipgloss.NewStyle().Background(green),
		"paused":     lipgloss.NewStyle().Background(yellow),
		"restarting": lipgloss.NewStyle().Background(orange),
		"exited":     lipgloss.NewStyle().Background(midPink),
		"dead":       lipgloss.NewStyle().Background(black),
	}

	wrapStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 5, 1).
			Align(lipgloss.Left).
			MarginLeft(5)

	titleStyle = lipgloss.NewStyle().
			Background(orange).
			Foreground(black).Bold(true).
			Align(lipgloss.Center).
			Blink(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(grey).
			Align(lipgloss.Center)

	logStyle = lipgloss.NewStyle().
			Foreground(black).
			Align(lipgloss.Left)
)

func getContainers() []container {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	containers := []container{}
	for _, c := range listContainers(client, true) {
		name := c.Names[0][1:]
		status := c.State
		c := container{name: name, state: status, id: c.ID, ancestor: c.Image}
		containers = append(containers, c)
	}
	return containers
}

func getImages() []image {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	images := []image{}
	for _, c := range listImages(client, true) {
		tags := c.RepoTags
		var name string
		var size int64
		if len(tags) > 0 {
			name = tags[0]
			size = c.Size
			// format size (GB, MB, KB)
			if size > 1000000000 {
				size = size / 1000000000
				name = fmt.Sprintf("%s (%dGB)", name, size)
			} else if size > 1000000 {
				size = size / 1000000
				name = fmt.Sprintf("%s (%dMB)", name, size)
			} else if size > 1000 {
				size = size / 1000
				name = fmt.Sprintf("%s (%dKB)", name, size)
			} else {
				name = fmt.Sprintf("%s (%dB)", name, size)
			}
			c := image{name: name, id: c.ID}
			images = append(images, c)
		}
	}
	return images
}

func initialModel() model {
	containers := getContainers()
	images := getImages()
	return model{
		containers: containers,
		images:     images,
		selected:   make(map[int]struct{}),
	}
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
			}

			// force for now  aka include running (TODO: opts)
			for k := range m.selected {
				container := m.containers[k]
				id := container.id
				go removeContainer(client, id)
				m.logs += "🗑️  Remove " + container.name + "\n"
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
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "running" {
					go restartContainer(client, id)
					m.logs += "🔃 Restarted " + container.name + "\n"
				} else {
					m.logs += "🚧  " + container.name + " not running\n"
				}
			}
			return m, nil

		case "K": // kill
			m.logs = ""
			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "running" {
					killContainer(client, id)
					m.logs += "🔪 Killed " + container.name + "\n"
				} else {
					m.logs += "🚧 " + container.name + " already stopped\n"
				}
			}
			return m, nil

		case "s": // stop
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "running" {
					go stopContainer(client, id)
					m.logs += "🛑 Stop " + container.name + "\n"
				} else {
					m.logs += "🚧  " + container.name + " already stopped\n"
				}
			}
			return m, nil

		case "u": // up
			m.logs = ""
			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "exited" || state == "created" {
					go startContainer(client, id)
					if err != nil {
						m.logs += fmt.Sprintf("🚧  %s\n", err.Error())
					} else {
						m.logs += "🚀 Started " + container.name + "\n"
					}
				} else {
					m.logs += "🚧  " + container.name + " already running\n"
				}
			}
			return m, nil

		case "p": // pause
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
			}

			for i, choice := range m.containers {
				id := choice.id
				state := choice.state
				if _, ok := m.selected[i]; ok {
					if state == "running" {
						pauseContainer(client, id)
						m.logs += "⏳ Paused " + choice.name + "\n"
					} else {
						m.logs += "🚧  " + choice.name + " is not running\n"
					}
				}
			}
			return m, nil

		case "P": // unpause
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
			}

			for i, choice := range m.containers {
				id := choice.id
				state := choice.state
				if _, ok := m.selected[i]; ok {
					if state == "paused" {
						unPauseContainer(client, id)
						m.logs += "✅ Unpaused " + choice.name + "\n"
					} else {
						m.logs += "🚧  " + choice.name + " is not running\n"
					}
				}
			}
			return m, nil

		case "ctrl+a", "A": // select all
			for i := range m.containers {
				m.selected[i] = struct{}{}
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
			if m.page == 0 {
				m.page = 2
			} else {
				m.page = 0
			}
			return m, nil
		case "tab":
			// should not include controls page
			if m.page == 0 {
				m.page = 1
			} else {
				m.page = 0
			}
		}

	}
	return m, nil
}

func (m model) View() string {

	var s string
	title := "    🐳 Docker Containers    "
	s += titleStyle.Render(title)
	s += "\n\n"

	if m.page == 0 {
		for i, choice := range m.containers {
			cursor := "  " // default cursor
			if m.cursor == i {
				cursor = "👉"
			}
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			state := stateStyle[choice.state].Render(" ")
			name := choice.name
			s += fmt.Sprintf("%s [%s] %s %s\n", cursor, checked, state, name)
		}
	} else if m.page == 1 {
		for i, choice := range m.images {
			cursor := "  " // default cursor
			if m.cursor == i {
				cursor = "👉"
			}
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			name := choice.name
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, name)
		}
	} else if m.page == 2 {
		controls := `
x   - remove
r   - restart
K   - kill
s   - stop
u   - start
p   - pause
P   - unpause
esc - clear
? - hide controls`
		s += controls + "\n"

	}
	hint := "\n'q' quit | '?' toggle controls\n"

	s += hintStyle.Render(hint)
	s += "\n"
	s += strings.Repeat("─", lipgloss.Width(title))
	s += "\n"
	s += logStyle.Render(m.logs)

	wrapAll := wrapStyle.Render(s)
	return wrapAll
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Alas, it's all over. Error: ", err.Error())
		os.Exit(1)
	}
}
