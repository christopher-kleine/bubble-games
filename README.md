# bubble-games

Play games over an SSH connection. It uses the [Bubbletea](https://github.com/charmbracelet/bubbletea) for it's display. For SSH, the [Wish](https://github.com/charmbracelet/wish) library is used. Both from [charm.sh](https://charm.sh).

## Install

```sh
go install github.com/christopher-kleine/bubble-games@latest
```

## Options / Flags

Currently the following options / flags are available:

- port (default: 2222): The port the server listens on
- host (default: ""): The host the server binds itself to

## Usage

After connecting to the serer, you have the following options:

- `!help`        = Show the help
- `!quit`, `!exit` = Disconnect from server
- `!name <NAME>` = Change your name to `<NAME>`
- `!games`       = Show available games
- `!play <GAME>` = Start a game

## Games

Currently there is only one game available: [TicTacToe](https://github.com/christop-kleine/bubble-tictactoe).

## Roadmap

- [ ] Add more games
- [ ] Add Multiplayer
- [ ] Add Chat

## Adding a custom game

To add a custom game, you currently need to edit the file `games.go`. The important part looks like this:

```go
	Gamelist = map[string]Game{
		"tictactoe": {
			New:  tictactoe.New,
			Desc: "A simple TicTacToe",
		},
	}
```

If you were to add a new game - let's say chess - you would add it like this:

```go
	Gamelist = map[string]Game{
		"tictactoe": {
			New:  tictactoe.New,
			Desc: "A simple TicTacToe",
		},
        "chess": {
            New: chess.New,
            Desc: "Play a game of chess",
        }
	}
```

An entry in the `Gamelist` map looks like this:

```go
type Game struct {
	New  func(string, tea.Model) tea.Model
	Desc string
}
```