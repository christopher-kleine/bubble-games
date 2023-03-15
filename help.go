package main

import tea "github.com/charmbracelet/bubbletea"

type Help struct {
	Parent tea.Model
}

func (Help) Init() tea.Cmd {
	return nil
}

func (h Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return h.Parent, nil
	}

	return h, nil
}

func (Help) View() string {
	return `You can use the following commands:

- !help        = Show this help
- !quit, !exit = Disconnect from server
- !name <NAME> = Change your name to <NAME>
- !games       = Show available games
- !play <GAME> = Start a game

Press any key to continue...`
}
