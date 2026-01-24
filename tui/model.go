package tui

import (
	"time"

	"fauthy/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type entry struct {
	label  string
	secret string
}

type mode int

const (
	passphrasemode mode = iota
	viewmode
	addlabelmode
	addsecretmode
	confirmresetmode
)

type tickmsg time.Time

type model struct {
	entries        []entry
	store          *storage.Store
	passphrase     string
	needssetup     bool
	errormsg       string
	failedattempts int

	mode         mode
	input        string
	buffer       entry
	lastactivity time.Time

	width  int
	height int
}

func Newmodel(store *storage.Store) model {
	needssetup := !store.Exists()
	mode := passphrasemode
	if needssetup {
		mode = passphrasemode
	}

	return model{
		store:      store,
		mode:       mode,
		needssetup: needssetup,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickmsg(t)
	})
}
