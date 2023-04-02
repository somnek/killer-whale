package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/fsouza/go-dockerclient"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	choices := []string{}
	for _, c := range listContainers(client) {
		name := c.Names[0][1:]
		choices = append(choices, name)
	}

	return model{
		// choices:  []string{"🍎 Apple", "🍐 Pear", "🍊 Orange", "🍌 Banana", "🍉 Watermelon", "🍇 Grape", "🍓 Strawberry", "🍈 Melon", "🍒 Cherry", "🍑 Peach"},
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
	s := "What should we order for lunch?\n\n"
	for i, choice := range m.choices {
		cursor := "  " // default cursor
		if m.cursor == i {
			cursor = "👉"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
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
