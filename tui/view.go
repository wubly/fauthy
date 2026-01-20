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
		return inputview(m, title, maskinput(m.input), msg)
	case confirmresetmode:
		return resetview(m)
	case addlabelmode:
		return inputview(m, "add label", m.input, "enter to confirm")
	case addsecretmode:
		return inputview(m, "add secret", m.input, "enter to confirm")
	default:
		return listview(m)
	}
}

func resetview(m model) string {
	w, h := m.width, m.height
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titlestyle.Render("can't remember passphrase?"),
		"",
		"no way to recover without it.",
		"start over?",
		"",
		inputstyle.Render(m.input),
		mutedstyle.Render("y = yes, n = quit"),
	)

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

func maskinput(input string) string {
	masked := ""
	for range input {
		masked += "•"
	}
	return masked
}

func inputview(m model, title, input, hint string) string {
	w, h := m.width, m.height
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titlestyle.Render(title),
		inputstyle.Render(input),
		mutedstyle.Render(hint),
	)

	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}

func listview(m model) string {
	var rows []string
	now := time.Now()
	maxLabel := 8
	for _, e := range m.entries {
		if l := len(e.label); l > maxLabel {
			maxLabel = l
		}
	}

	for _, e := range m.entries {
		code, rem := totp.Generate(e.secret, now)
		row := fmt.Sprintf(
			"%-*s %s %2ds",
			maxLabel, e.label,
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
		mutedstyle.Render("a add • t theme • q quit"),
	)

	w, h := m.width, m.height
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}
	return lipgloss.Place(
		w,
		h,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
