package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
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
	viewport   viewport.Model
	altscreen  bool
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

	lastPage = 4
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
			Align(lipgloss.Left)

	titleStyle = lipgloss.NewStyle().
			Background(orange).
			Foreground(black).Bold(true).
			Align(lipgloss.Center).
			Blink(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(grey).
			Align(lipgloss.Left)

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

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		// m.logs += fmt.Sprintf("resize: %dx%d\n", msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "f": // toggle fullscreen
			var cmd tea.Cmd
			if m.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altscreen = !m.altscreen
			return m, cmd

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
				return m, nil
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
				return m, nil
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
				return m, nil
			}

			for k := range m.selected {
				container := m.containers[k]
				state := container.state
				id := container.id
				if state == "running" || state == "restarting" {
					go stopContainer(client, id)
					m.logs += "🛑 Stop " + container.name + "\n"
				} else {
					m.logs += "🚧  " + " unable to stop " + container.name + "\n"
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
				return m, nil
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
				return m, nil
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
				return m, nil
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
			// should not include controls page
			if m.page < lastPage-1 {
				m.page++
			} else {
				m.page = 0
			}
		}

	}
	return m, nil
}

func (m model) View() string {

	var s string
	var title string

	if m.page == 0 {
		// container page
		title = "        🐳 Docker Containers        " // 30 characters
		s += titleStyle.Render(title)
		s += "\n\n"

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
			// limit to 25 characters for now
			// TODO: make this dynamic
			if len(name) > 25 {
				name = name[:25] + "..."
			}
			s += fmt.Sprintf("%s [%s] %s %s\n", cursor, checked, state, name)
		}
	} else if m.page == 1 {
		// image page
		title = "          🐳 Docker Images          " // 30 characters
		s += titleStyle.Render(title)
		s += "\n\n"

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

	} else if m.page == 2 {
		// utilities like docker volume prune, docker system prune, etc
		title = "         🐳 Docker Utilities        "
		s += titleStyle.Render(title)
		s += "\n\n"
		s += " 1. docker system prune\n"
		s += " 2. docker volume prune\n"
	} else if m.page == 3 {
		// controls page
		title = "            🔧 Controls             "
		s += titleStyle.Render(title)
		s += "\n"

		controls := `
 x  - remove
 r  - restart
 K  - kill
 s  - stop
 u  - start
 p  - pause
 P  - unpause
esc - clear
C-a - select all
 f  - toggle fullscreen
 ?  - hide controls`
		s += controls + "\n"

	}

	hint := "\n'q' quit | '?' controls | ' ' select\n'h/j/k/l' move | 'tab' switch page"

	s += hintStyle.Render(hint)
	s += "\n"
	s += strings.Repeat("─", lipgloss.Width(title))
	s += "\n"
	s += logStyle.Render(m.logs)

	// wrapStyle = wrapStyle.Foreground(lipgloss.Color(celesBlue))
	wrapStyle = wrapStyle.MarginLeft(m.viewport.Width/2 - lipgloss.Width(title))
	wrapAll := wrapStyle.Render(s)
	return wrapAll
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		// tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Alas, it's all over. Error: ", err.Error())
		os.Exit(1)
	}
}
