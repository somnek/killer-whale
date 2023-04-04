package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/fsouza/go-dockerclient"
)

type model struct {
	choices  []container
	cursor   int
	selected map[int]struct{}
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
	titleStyle = lipgloss.NewStyle().
			Background(electricBlue).
			Foreground(black).
			Bold(true).
			Align(lipgloss.Center).
			Faint(true).
			MarginLeft(5).
			Border(border).Blink(true)
)

func initialModel() model {
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

	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	title := titleStyle.Render("     Docker Containers     ")
	s += fmt.Sprintf("%s\n\n", title)
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
	s += "\nPress q to quit\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Alas, it's all over. Error: ", err.Error())
		os.Exit(1)
	}
}
