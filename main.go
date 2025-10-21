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
	host = "localhost"
	port = 2222
)

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
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath("ssh_host_key"),
		wish.WithMiddleware(
			bubbletea.Middleware(programHandler),
			logging.Middleware(),
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
	log.Println("Stopping SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
