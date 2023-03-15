package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
)

func main() {
	var (
		host string
		port int
	)
	flag.StringVar(&host, "host", "", "host to listen on")
	flag.IntVar(&port, "port", 2222, "port to listen on")
	flag.Parse()

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}
	m := model{
		term:   pty.Term,
		width:  pty.Window.Width,
		height: pty.Window.Height,
		user:   s.User(),
		input:  textinput.New(),
	}
	m.input.Focus()
	// m.input.EchoMode = textinput.EchoPassword

	pkey := "???"
	t := "???"
	if s.PublicKey() != nil {
		pkey = base64.RawStdEncoding.EncodeToString(s.PublicKey().Marshal())
		t = s.PublicKey().Type()
	}
	log.Printf("%q: (%s) %s", s.User(), t, pkey)

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term   string
	width  int
	height int
	user   string
	input  textinput.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			v := strings.TrimSpace(m.input.Value())
			switch {
			case strings.HasPrefix(v, "!name ") && len(v) > 9:
				m.user = strings.TrimSpace(strings.TrimPrefix(v, "!name "))

			case v == "!quit" || v == "!exit":
				return m, tea.Quit

			case v == "!help":
				m.input.SetValue("")
				return Help{Parent: m}, nil

			case v == "!games":
				m.input.SetValue("")
				return Games{Parent: m}, nil

			case strings.HasPrefix(v, "!play "):
				game := strings.TrimPrefix(v, "!play ")
				m.input.SetValue("")
				return Run(game, m.user, m), nil
			}
			m.input.SetValue("")
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := fmt.Sprintf("Hello %s!", m.user) + "\n"
	s += "\n"
	s += fmt.Sprintf("Your term is %s", m.term) + "\n"
	s += fmt.Sprintf("Your window size is x: %d y: %d", m.width, m.height) + "\n"
	s += "\n"
	s += "Press 'Ctrl+C' to quit" + "\n"
	s += "Enter !help for help" + "\n"
	s += "\n"
	s += m.input.View()

	return s
}
