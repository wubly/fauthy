package tui

import (
	"fauthy/storage"

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
		return m, nil

	case tea.KeyMsg:
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
			return m, nil

		case "a":
			if m.mode == viewmode {
				m.mode = addlabelmode
				m.input = ""
				return m, nil
			}
			if len(msg.Runes) > 0 {
				for _, r := range msg.Runes {
					m.input += string(r)
				}
			}
			return m, nil

		case "enter":
			switch m.mode {
			case passphrasemode:
				if len(m.input) == 0 {
					return m, nil
				}
				m.passphrase = m.input
				m.input = ""

				if m.needssetup {
					if err := m.store.Save([]storage.Entry{}, m.passphrase); err != nil {
						m.errormsg = "failed to create storage"
						m.mode = passphrasemode
						m.input = ""
						return m, nil
					}
					m.mode = viewmode
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
						return m, nil
					}
					m.entries = fromstorageentries(entries)
					m.mode = viewmode
					m.errormsg = ""
					m.failedattempts = 0
				}
			case confirmresetmode:
				input := m.input
				m.input = ""
			if input == "y" || input == "y" {
				m.store.Delete()
				m.needssetup = true
				m.failedattempts = 0
				m.errormsg = ""
				m.mode = passphrasemode
			} else if input == "n" || input == "n" {
					return m, tea.Quit
				}
				return m, nil
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
			}
			return m, nil

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil

		default:
			if m.mode != viewmode && len(msg.Runes) > 0 {
				for _, r := range msg.Runes {
					m.input += string(r)
				}
			}
			return m, nil
		}

	case tickmsg:
		return m, tick()
	}

	return m, nil
}
