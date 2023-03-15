package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Games struct {
	Parent tea.Model
}

func (Games) Init() tea.Cmd {
	return nil
}

func (g Games) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return g.Parent, nil
	}

	return g, nil
}

func (Games) View() string {
	text := "The following games are available:\n"
	text += "----------------------------------\n"
	text += "Enter !play <game> to start a game.\n\n"

	for name, game := range Gamelist {
		text += fmt.Sprintf("- %s (%s)\n", name, game.Desc)
	}

	return text
}
