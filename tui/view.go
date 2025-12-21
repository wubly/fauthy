package tui

import (
	"fmt"
	"strings"
	"time"

	"fauthy/totp"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch m.mode {
	case passphrasemode:
		title := "enter passphrase"
		if m.needssetup {
			title = "create passphrase"
		}
		msg := "enter to confirm"
		if m.errormsg != "" {
			msg = m.errormsg
		}
		return inputview(title, maskinput(m.input), msg)
	case confirmresetmode:
		return resetview(m.input)
	case addlabelmode:
		return inputview("add label", m.input, "enter to confirm")
	case addsecretmode:
		return inputview("add secret", m.input, "enter to confirm")
	default:
		return listview(m)
	}
}

func resetview(input string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titlestyle.Render("can't remember passphrase?"),
		"",
		"no way to recover without it.",
		"start over?",
		"",
		inputstyle.Render(input),
		mutedstyle.Render("y = yes, n = quit"),
	)

	return lipgloss.Place(50, 10, lipgloss.Center, lipgloss.Center, content)
}

func maskinput(input string) string {
	masked := ""
	for range input {
		masked += "•"
	}
	return masked
}

func inputview(title, input, hint string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titlestyle.Render(title),
		inputstyle.Render(input),
		mutedstyle.Render(hint),
	)

	return lipgloss.Place(50, 10, lipgloss.Center, lipgloss.Center, content)
}

func listview(m model) string {
	var rows []string
	now := time.Now()

	for _, e := range m.entries {
		code, rem := totp.Generate(e.secret, now)
		row := fmt.Sprintf(
			"%-12s %s %2ds",
			e.label,
			codestyle.Render(code),
			rem,
		)
		rows = append(rows, rowstyle.Render(row))
	}

	if len(rows) == 0 {
		rows = append(rows, mutedstyle.Render("press 'a' to add a secret"))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titlestyle.Render("fauthy"),
		strings.Join(rows, "\n"),
		"",
		mutedstyle.Render("a add • q quit"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
