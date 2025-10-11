package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

const playerArt = `👤YOU`

type playerStats struct {
	hp, mana, speed, magic, strength int
}

type roomType int

const (
	Empty roomType = iota
	Enemy
	Tresure
	Shop
	StairsUp
	StairsDown
)

type room struct {
	Type roomType
}

type model struct {
	width      int
	height     int
	playerMapX int
	playerMapY int
	worldMap   [][]*room
	stats      playerStats
}

func generateMap(width, height, maxRooms int) ([][]*room, int, int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	worldMap := make([][]*room, height)
	for i := range worldMap {
		worldMap[i] = make([]*room, width)
	}

	startX, startY := width/2, height/2
	currentX, currentY := startX, startY
	roomsCreated := 0

	commonRoomTypes := []roomType{Empty, Enemy, Tresure, Shop}

	for roomsCreated < maxRooms {
		if worldMap[currentY][currentX] == nil {
			randomIndex := rand.Intn(len(commonRoomTypes))
			randomType := commonRoomTypes[randomIndex]
			worldMap[currentY][currentX] = &room{Type: randomType}
			roomsCreated++
		}

		dx, dy := 0, 0
		switch rand.Intn(4) {
		case 0:
			dy = -1
		case 1:
			dy = 1
		case 2:
			dx = -1
		case 3:
			dx = 1
		}

		if currentX+dx >= 0 && currentX+dx < width && currentY+dy >= 0 && currentY+dy < height {
			currentX += dx
			currentY += dy
		}
	}

	worldMap[startY][startX].Type = Empty

	for {
		randX, randY := rand.Intn(width), rand.Intn(height)
		if worldMap[randY][randX] != nil && (randX != startX || randY != startY) {
			worldMap[randY][randX].Type = StairsUp
			break
		}
	}

	return worldMap, startX, startY
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {

	worldMap, startX, startY := generateMap(9, 9, 15)

	initialModel := model{
		playerMapX: startX,
		playerMapY: startY,
		worldMap:   worldMap,
		stats: playerStats{
			hp:       100,
			mana:     25,
			speed:    10,
			magic:    12,
			strength: 12,
		},
	}

	return initialModel, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "w":
			if m.playerMapY > 0 && m.worldMap[m.playerMapY-1][m.playerMapX] != nil {
				m.playerMapY--
			}

		case "down", "s":
			if m.playerMapY < len(m.worldMap)-1 && m.worldMap[m.playerMapY+1][m.playerMapX] != nil {
				m.playerMapY++
			}

		case "left", "a":
			if m.playerMapX > 0 && m.worldMap[m.playerMapY][m.playerMapX-1] != nil {
				m.playerMapX--
			}

		case "right", "d":
			if m.playerMapX < len(m.worldMap[0])-1 && m.worldMap[m.playerMapY][m.playerMapX+1] != nil {
				m.playerMapX++
			}
		}
	}

	return m, nil
}

func (r *room) getRoomSymbol() string {
	switch r.Type {
	case Empty:
		return " "
	case Enemy:
		return "E"
	case Tresure:
		return "T"
	case Shop:
		return "$"
	case StairsUp:
		return "▲"
	case StairsDown:
		return "▼"
	default:
		return "?"
	}
}

func (r *room) getRoomDescription() string {
	switch r.Type {
	case Empty:
		return "An empty room.\nHere you can rest."
	case Enemy:
		return "You feel the danger!\nGet ready for the battle."
	case Tresure:
		return "You see a big chest in the middle of the room!"
	case Shop:
		return "A suspicious merchan greets you.\n'I have some thing to offer you, take a look.'"
	case StairsUp:
		return "Some stone stairs, they take you to the darkness\n You wanna go up?"
	default:
		return "Unknown room type."
	}
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	cameraStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	mapStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("57"))

	playerCell := lipgloss.NewStyle().Width(3).Align(lipgloss.Center).Foreground(lipgloss.Color("214")).SetString("[@]")
	emptyCell := lipgloss.NewStyle().Width(3).SetString(" ")

	var mapRows []string
	for y, row := range m.worldMap {
		var mapRow strings.Builder
		for x, room := range row {
			if x == m.playerMapX && y == m.playerMapY {
				mapRow.WriteString(playerCell.String())
			} else if room != nil {
				symbol := room.getRoomSymbol()
				style := lipgloss.NewStyle().Width(3).Align(lipgloss.Center).SetString(fmt.Sprintf("[%s]", symbol))
				mapRow.WriteString(style.String())
			} else {
				mapRow.WriteString(emptyCell.String())
			}
		}
		mapRows = append(mapRows, mapRow.String())
	}

	mapContent := lipgloss.JoinVertical(lipgloss.Center, mapRows...)
	mapView := mapStyle.Width(45).Align(lipgloss.Center).Render(mapContent)

	//mapHeight := lipgloss.Height(mapView)
	cameraWidth := m.width - lipgloss.Width(mapView) - 4

	statsArt := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Margin(1, 2).Render(playerArt)
	statsText := fmt.Sprintf(
		"HP: %d\nMana: %d\nSpeed: %d\nMagic: %d\nStrength: %d",
		m.stats.hp,
		m.stats.mana,
		m.stats.speed,
		m.stats.magic,
		m.stats.strength,
	)
	statsContent := lipgloss.JoinHorizontal(lipgloss.Top, statsArt, statsText)
	statsView := panelStyle.Width(cameraWidth).Render(statsContent)

	cameraHeight := 2

	currentRoom := m.worldMap[m.playerMapY][m.playerMapX]
	cameraContent := fmt.Sprintf("%s", currentRoom.getRoomDescription())
	cameraView := cameraStyle.Width(cameraWidth).Height(cameraHeight).Render(cameraContent)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, cameraView, statsView)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1).
		Render("Arrows/wasd: move | 'q': quit")

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, mapView)
	finalView := lipgloss.JoinVertical(lipgloss.Left, mainView, help)

	return finalView
}
