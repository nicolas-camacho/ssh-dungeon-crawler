package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"ssh-dungeon-crawler/game"
)

const (
	host = "0.0.0.0"
	port = "2222" // El puerto interno donde escucha la app
)

// --- Modelo de Prueba Súper Simple ---
type helloModel struct{}

func (m helloModel) Init() tea.Cmd { return nil }

func (m helloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m helloModel) View() string {
	return "\n\n  ¡Hola Mundo desde Fly.io! Si ves esto, la conexión funciona.\n\n  Presiona 'q' para salir.\n\n"
}

func simpleHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Esta función no depende de nada más, solo crea el modelo simple.
	return helloModel{}, []tea.ProgramOption{tea.WithAltScreen()}
}

func programHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	cmd := s.Command()

	if len(cmd) > 0 && cmd[0] == "test-combat" {
		log.Println("Starting test combat session...")
		return game.CreateTeaProgram(s, game.StateCombat)
	}

	log.Println("Starting normal game session...")
	return game.CreateTeaProgram(s, game.StateLoading)
}

func main() {

	if err := game.LoadGameData(); err != nil {
		log.Fatalf("Failed to load game data: %v", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%s", host, port)),
		wish.WithHostKeyPath("ssh_host_key"),
		wish.WithMiddleware(
			bubbletea.Middleware(simpleHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%s", host, port)

	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
