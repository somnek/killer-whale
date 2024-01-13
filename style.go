package main

import (
	"github.com/charmbracelet/lipgloss"
)

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
		"created":    lipgloss.NewStyle().Foreground(midPurple),
		"running":    lipgloss.NewStyle().Foreground(green),
		"paused":     lipgloss.NewStyle().Foreground(yellow),
		"restarting": lipgloss.NewStyle().Foreground(orange),
		"exited":     lipgloss.NewStyle().Foreground(midPink),
		"dead":       lipgloss.NewStyle().Foreground(black),
	}

	bodyLStyle = lipgloss.NewStyle().
			Padding(1, 2, 0, 4).
			Border(lipgloss.RoundedBorder(), true, false, true, true).
			BorderForeground(black)
	bodyRStyle = lipgloss.NewStyle().
			Padding(1, 4, 0, 2).
			PaddingLeft(4).
			Border(lipgloss.RoundedBorder(), true, true, true, false).
			Foreground(black).
			BorderForeground(black)

	bodyStyle = lipgloss.NewStyle().
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

	selectedNameStyle = lipgloss.NewStyle().
				Foreground(black).
				Background(grey)
	styleCheck = lipgloss.NewStyle().
			Foreground(hotGreen)
)