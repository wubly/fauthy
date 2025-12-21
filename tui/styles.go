package tui

import "github.com/charmbracelet/lipgloss"

var (
	titlestyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	rowstyle = lipgloss.NewStyle().
			Padding(0, 1)

	codestyle = lipgloss.NewStyle().
			Bold(true)

	mutedstyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	inputstyle = lipgloss.NewStyle().
			Underline(true)
)
