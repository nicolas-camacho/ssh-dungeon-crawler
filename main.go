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

	"ssh-dungeon-crawler/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding without it")
	}
}

func programHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Asegurar que el PTY tenga las capacidades correctas
	ptyReq, winCh, isPty := s.Pty()
	if !isPty {
		log.Println("No PTY requested, forcing PTY mode")
		wish.Println(s, "Error: PTY required for this application")
		return nil, nil
	}

	// Log de información del terminal
	log.Printf("PTY Info - Term: %s, Window: %dx%d", ptyReq.Term, ptyReq.Window.Width, ptyReq.Window.Height)

	cmd := s.Command()
	var initialState game.GameState

	if len(cmd) > 0 && cmd[0] == "test-combat" {
		log.Println("Starting test combat session...")
		initialState = game.StateCombat
	} else {
		log.Println("Starting normal game session...")
		initialState = game.StateLoading
	}

	model, options := game.CreateTeaProgram(s, initialState)

	// Agregar opciones adicionales para el terminal SSH
	options = append(options,
		tea.WithAltScreen(),       // Usar pantalla alternativa
		tea.WithMouseCellMotion(), // Habilitar mouse si es posible
	)

	// Manejar cambios de tamaño de ventana
	go func() {
		for win := range winCh {
			log.Printf("Window resized to: %dx%d", win.Width, win.Height)
		}
	}()

	return model, options
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
			// Configuraciones adicionales para mejorar compatibilidad
			wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
				// Aceptar todas las conexiones (para desarrollo)
				// En producción, implementa autenticación real
				return true
			}),
			wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
				// Aceptar todas las contraseñas (para desarrollo)
				return true
			}),
		)
		if err != nil {
			log.Fatalf("failed to create server: %s", err)
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Printf("Starting SSH server on %s:%s", host, port)
		log.Printf("Connect with: ssh -p %s %s", port, host)
		log.Printf("Or with better terminal support: ssh -p %s -o 'SetEnv TERM=xterm-256color' %s", port, host)

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
