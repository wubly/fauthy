package tui

import (
	"fauthy/storage"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func tostorageentries(entries []entry) []storage.Entry {
	result := make([]storage.Entry, len(entries))
	for i, e := range entries {
		result[i] = storage.Entry{Label: e.label, Secret: e.secret}
	}
	return result
}

func fromstorageentries(entries []storage.Entry) []entry {
	result := make([]entry, len(entries))
	for i, e := range entries {
		result[i] = entry{label: e.Label, secret: e.Secret}
	}
	return result
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, tick()

	case tea.KeyMsg:
		if m.mode == viewmode {
			m.lastactivity = time.Now()
		}
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.mode == viewmode {
				return m, tea.Quit
			}
			if len(msg.Runes) > 0 {
				for _, r := range msg.Runes {
					m.input += string(r)
				}
			}
			return m, tick()

		case "a":
			if m.mode == viewmode {
				m.mode = addlabelmode
				m.input = ""
				return m, tick()
			}
			if len(msg.Runes) > 0 {
				for _, r := range msg.Runes {
					m.input += string(r)
				}
			}
			return m, tick()

		case "t":
			if m.mode == viewmode {
				NextTheme()
				return m, tick()
			}
			m.input += "t"
			return m, tick()

		case "enter":
			switch m.mode {
			case passphrasemode:
				if len(m.input) == 0 {
					return m, tick()
				}
				m.passphrase = m.input
				m.input = ""

				if m.needssetup {
					if err := m.store.Save([]storage.Entry{}, m.passphrase); err != nil {
						m.errormsg = "failed to create storage"
						m.mode = passphrasemode
						m.input = ""
						return m, tick()
					}
					m.mode = viewmode
					m.lastactivity = time.Now()
				} else {
					entries, err := m.store.Load(m.passphrase)
					if err != nil {
						m.failedattempts++
						if m.failedattempts >= 5 {
							m.mode = confirmresetmode
							m.input = ""
						} else {
							m.errormsg = "wrong passphrase"
							m.mode = passphrasemode
							m.input = ""
						}
						return m, tick()
					}
					m.entries = fromstorageentries(entries)
					m.mode = viewmode
					m.errormsg = ""
					m.failedattempts = 0
					m.lastactivity = time.Now()
				}
			case confirmresetmode:
				input := m.input
				m.input = ""
				if input == "y" {
					m.store.Delete()
					m.needssetup = true
					m.failedattempts = 0
					m.errormsg = ""
					m.mode = passphrasemode
				} else if input == "n" {
					return m, tea.Quit
				}
				return m, tick()
			case addlabelmode:
				m.buffer.label = m.input
				m.input = ""
				m.mode = addsecretmode
			case addsecretmode:
				m.buffer.secret = m.input
				m.entries = append(m.entries, m.buffer)

				if err := m.store.Save(tostorageentries(m.entries), m.passphrase); err != nil {
					m.errormsg = "failed to save"
				}

				m.buffer = entry{}
				m.input = ""
				m.mode = viewmode
				m.lastactivity = time.Now()
			}
			return m, tick()

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, tick()

		default:
			if m.mode != viewmode && len(msg.Runes) > 0 {
				for _, r := range msg.Runes {
					m.input += string(r)
				}
			}
			return m, tick()
		}

	case tickmsg:
		if m.mode == viewmode {
			if !m.lastactivity.IsZero() && time.Since(m.lastactivity) >= 2*time.Minute {
				m.passphrase = ""
				m.entries = []entry{}
				m.mode = passphrasemode
				m.input = ""
				m.lastactivity = time.Time{}
			}
		}
		return m, tick()
	}

	return m, tick()
}
