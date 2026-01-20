package tui

import "github.com/charmbracelet/lipgloss"

var (
	titlestyle lipgloss.Style
	rowstyle   lipgloss.Style
	codestyle  lipgloss.Style
	mutedstyle lipgloss.Style
	inputstyle lipgloss.Style
)

type theme struct {
	title string
	code  string
	muted string
	input string
}

var themes = []theme{
	{title: "#7aa2f7", code: "#9ece6a", muted: "#808080", input: "#bb9af7"},
	{title: "#268bd2", code: "#2aa198", muted: "#93a1a1", input: "#b58900"},
	{title: "#bd93f9", code: "#50fa7b", muted: "#6272a4", input: "#ff79c6"},
	{title: "#fabd2f", code: "#b8bb26", muted: "#928374", input: "#83a598"},
	{title: "#ff6ac1", code: "#29d398", muted: "#9aa5b1", input: "#7e5bef"},
	{title: "#cba6f7", code: "#a6e3a1", muted: "#7f849c", input: "#b4befe"},
}

var themeIndex int

func applyTheme(t theme) {
	titlestyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(t.title)).
		MarginBottom(1)

	rowstyle = lipgloss.NewStyle().
		Padding(0, 1)

	codestyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(t.code))

	mutedstyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.muted))

	inputstyle = lipgloss.NewStyle().
		Underline(true).
		Foreground(lipgloss.Color(t.input))
}

func NextTheme() {
	themeIndex = (themeIndex + 1) % len(themes)
	applyTheme(themes[themeIndex])
}

func init() {
	themeIndex = 0
	applyTheme(themes[themeIndex])
}
