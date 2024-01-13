package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func buildContainerView(m model) string {
	var s string
	for i, choice := range m.containers {
		cursor := " " // default cursor
		check := " "
		if m.cursor == i {
			cursor = "❯"
		}
		state := stateStyle[choice.state].Render("●")
		name := choice.name
		name = runewidth.Truncate(name, 25, "...")
		if _, ok := m.selected[i]; ok {
			check = styleCheck.Render("✔")
		}
		s += fmt.Sprintf("%s %s %s %s", cursor, check, state, name) + "\n"
	}
	return s
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
	return s
}

func (m model) View() string {

	var final string
	var bodyL, bodyR, body string

	// body L
	switch m.page {
	case pageContainer:
		bodyL = bodyLStyle.Render(buildContainerView(m))
	case pageImage:
		bodyL = bodyLStyle.Render(buildImageView(m))
	}

	// body R
	bodyR = bodyRStyle.Render(m.logs)

	//  title
	// title := titleStyle.Render("        🐳 Docker Containers        ")
	title := strings.Repeat(" ", 36) + "🐳 Docker"

	// join left + right component
	body = lipgloss.JoinHorizontal(lipgloss.Left, bodyL, bodyR)
	bodyStyle = bodyStyle.MarginLeft(m.width/2 - 36)
	body = bodyStyle.Render(body)

	// joing title + body + help
	final += lipgloss.JoinVertical(lipgloss.Top, title, body)
	return final + "\n"
}

// var title string
// *title = "        🐳 Docker Containers        " // 30 characters
// s += titleStyle.Render(*title)
// s += "\n\n"
// *title = "          🐳 Docker Images          " // 30 characters
// s += titleStyle.Render(*title)
// s += "\n\n"

// s += logStyle.Render(m.logs)
// wrapStyle = wrapStyle.MarginLeft(m.width/2 - 36)
