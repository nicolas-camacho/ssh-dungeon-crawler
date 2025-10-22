package main

import (
	"context"
	"flag"
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
	"github.com/joho/godotenv"

	"ssh-dungeon-crawler/game"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding without it")
	}
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
	sshMode := flag.Bool("ssh", false, "Run in SSH mode")
	startMode := flag.String("mode", "normal", "Starting mode: normal or test-combat")
	flag.Parse()

	if err := game.LoadGameData(); err != nil {
		log.Fatalf("Failed to load game data: %v", err)
	}

	if *sshMode {
		log.Println("Running in SSH mode...")

		host := os.Getenv("SSH_HOST")
		if host == "" {
			host = "0.0.0.0"
		}
		port := os.Getenv("SSH_PORT")
		if port == "" {
			port = "2222"
		}

		s, err := wish.NewServer(
			wish.WithAddress(fmt.Sprintf("%s:%s", host, port)),
			wish.WithHostKeyPath("ssh_host_key"),
			wish.WithMiddleware(
				bubbletea.Middleware(programHandler),
				logging.Middleware(),
			),
		)
		if err != nil {
			log.Fatalf("failed to create server: %s", err)
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
	} else {
		log.Println("Running in local terminal mode...")

		var startState game.GameState
		switch *startMode {
		case "test-combat":
			startState = game.StateCombat
		default:
			startState = game.StateLoading
		}

		initialModel, options := game.CreateTeaProgram(nil, startState)
		p := tea.NewProgram(initialModel, options...)
		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running program: %v", err)
		}
	}
}
