package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

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

	// pad body height
	padBodyHeight(&s, len(m.containers)+2)
	return bodyLStyle.Render(s)
}

func (m model) View() string {

	var final string
	var bodyL, bodyR, body string

	// body L
	switch m.page {
	case pageContainer:
		bodyL = buildContainerView(m)
	case pageImage:
		bodyL = buildImageView(m)
	}

	// body R
	bodyR = bodyRStyle.Render(buildLogView(m))

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
