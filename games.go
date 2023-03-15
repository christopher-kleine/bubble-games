package main

import (
	tea "github.com/charmbracelet/bubbletea"
	tictactoe "github.com/christopher-kleine/bubble-tictactoe"
)

type Game struct {
	New  func(string, tea.Model) tea.Model
	Desc string
}

var (
	Gamelist = map[string]Game{
		"tictactoe": {
			New:  tictactoe.New,
			Desc: "A simple TicTacToe",
		},
	}
)

func Run(game string, user string, parent tea.Model) tea.Model {
	if f, ok := Gamelist[game]; ok {
		return f.New(user, parent)
	}

	return parent
}
