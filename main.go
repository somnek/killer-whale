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
	choices  []container
	cursor   int
	selected map[int]struct{}
	logs     string
}

type container struct {
	name     string
	state    string
	id       string
	ancestor string
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
	border = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	stateStyle = map[string]lipgloss.Style{
		"created":    lipgloss.NewStyle().Background(midPurple),
		"running":    lipgloss.NewStyle().Background(hotGreen),
		"paused":     lipgloss.NewStyle().Background(yellow),
		"restarting": lipgloss.NewStyle().Background(orange),
		"exited":     lipgloss.NewStyle().Background(midPink),
		"dead":       lipgloss.NewStyle().Background(black),
	}

	wrapStyle = lipgloss.NewStyle().
			Border(border).
			Padding(1, 5, 1).
			Align(lipgloss.Left).
			MarginLeft(5)

	titleStyle = lipgloss.NewStyle().
			Background(celesBlue).
			Foreground(black).
			Bold(true).
			Align(lipgloss.Center).
			Blink(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(grey).
			Align(lipgloss.Center)

	logStyle = lipgloss.NewStyle().
			Foreground(black).
			Align(lipgloss.Left)
)

func getChoices() []container {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	choices := []container{}
	for _, c := range listContainers(client, true) {
		name := c.Names[0][1:]
		status := c.State
		c := container{name: name, state: status, id: c.ID, ancestor: c.Image}
		choices = append(choices, c)
	}
	return choices
}

func initialModel() model {
	choices := getChoices()
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
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
		choices := getChoices()
		m.choices = choices
		return m, doTick()

	case tea.KeyMsg:
		switch msg.String() {

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
				container := m.choices[k]
				state := container.state
				id := container.id
				if state == "running" {
					killContainer(client, id)
					m.logs += "🔪 Killed " + container.name + "\n"
				} else {
					m.logs += "❌ " + container.name + " already stopped\n"
				}
			}

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
				container := m.choices[k]
				state := container.state
				id := container.id
				if state == "running" {
					stopContainer(client, id)
					m.logs += "🛑 Stopped " + container.name + "\n"
				} else {
					m.logs += "❌ " + container.name + " already stopped\n"
				}
			}

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
				container := m.choices[k]
				state := container.state
				id := container.id
				if state == "exited" {
					startContainer(client, id)
					m.logs += "🚀 Started " + container.name + "\n"
				} else {
					m.logs += "❌ " + container.name + " already running\n"
				}
			}

		case "p": // pause
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
			}

			for i, choice := range m.choices {
				id := choice.id
				state := choice.state
				if _, ok := m.selected[i]; ok {
					if state == "running" {
						pauseContainer(client, id)
						m.logs += "⏳ Paused " + choice.name + "\n"
					} else {
						m.logs += "❌ " + choice.name + "is not running\n"
					}
				}
			}

		case "P": // unpause
			m.logs = ""

			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			if len(m.selected) == 0 {
				m.logs = "No container selected\n"
			}

			for i, choice := range m.choices {
				id := choice.id
				state := choice.state
				if _, ok := m.selected[i]; ok {
					if state == "paused" {
						unPauseContainer(client, id)
						m.logs += "⏳ Unpaused " + choice.name + "\n"
					} else {
						m.logs += "❌ " + choice.name + "is not running\n"
					}
				}
			}

		case "ctrl+a", "A": // select all
			for i := range m.choices {
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
				m.cursor = len(m.choices) - 1
			}

		case "down", "j": // move cursor down
			if m.cursor < len(m.choices)-1 {
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
		}
	}
	return m, nil
}

func (m model) View() string {

	var s string
	title := "    🐳 Docker Containers    "
	s += titleStyle.Render(title)
	s += "\n\n"

	for i, choice := range m.choices {
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

	hint := "\n'q' quit | '?' controls\n"

	s += hintStyle.MaxWidth(lipgloss.Width(title) * 2).Render(hint)
	s += "\n"
	s += strings.Repeat("─", lipgloss.Width(title))
	s += "\n"
	s += logStyle.MaxWidth(lipgloss.Width(title)).Render(m.logs)

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
