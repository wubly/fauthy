package main

import (
	"log"

	"fauthy/storage"
	"fauthy/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	store, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	m := tui.Newmodel(store)
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
