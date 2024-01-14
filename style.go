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
)

const (
	lastPage         = 4
	minHeightPerView = 8  // 6 item
	maxHeightPerView = 12 // 10 item
	fixedWidth       = 88 // 92
	fixedBodyLWidth  = 36 // exclude padding
	fixedBodyRWidth  = 44 // exclude padding
)

var (
	stateStyleMap = map[string]lipgloss.Style{
		"created":    lipgloss.NewStyle().Foreground(midPurple),
		"running":    lipgloss.NewStyle().Foreground(green),
		"paused":     lipgloss.NewStyle().Foreground(yellow),
		"restarting": lipgloss.NewStyle().Foreground(orange),
		"exited":     lipgloss.NewStyle().Foreground(midPink),
		"dead":       lipgloss.NewStyle().Foreground(black),
	}

	bodyLStyle = lipgloss.NewStyle().
			Padding(1, 0, 0, 4).
			BorderForeground(black)

	bodyRStyle = lipgloss.NewStyle().
			Padding(1, 4, 0, 0).
			PaddingLeft(4).
			Foreground(black).
			BorderForeground(black)

	appStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Foreground(black)
	bodyStyle = lipgloss.NewStyle().
			Align(lipgloss.Left)

	titleStyle = lipgloss.NewStyle().
			Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(grey).
			Align(lipgloss.Left)

	logStyle = lipgloss.NewStyle().
			Foreground(black)

	selectedNameStyle = lipgloss.NewStyle().
				Foreground(black).
				Background(grey)

	checkStyle = lipgloss.NewStyle().
			Foreground(hotGreen)

	itemCountStyle = lipgloss.NewStyle().
			Foreground(frenchBlue).
			Bold(true)
)
