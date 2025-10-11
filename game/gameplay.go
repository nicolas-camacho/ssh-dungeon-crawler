package game

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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

func (m model) renderGameView() string {
	emptyCell := lipgloss.NewStyle().Width(3).SetString(" ")

	var mapRows []string
	for y, row := range m.worldMap {
		var mapRow strings.Builder
		for x, room := range row {
			if x == m.playerMapX && y == m.playerMapY {
				mapRow.WriteString(m.styles.Player.String())
			} else if room != nil {
				symbol := room.getRoomSymbol()
				style := m.styles.Room
				if room.Type == Tresure || room.Type == Shop || room.Type == StairsUp {
					style = style.Inherit(m.styles.RoomSpecial)
				}
				mapRow.WriteString(style.Render(fmt.Sprintf("[%s]", symbol)))
			} else {
				mapRow.WriteString(emptyCell.String())
			}
		}
		mapRows = append(mapRows, mapRow.String())
	}

	mapContent := lipgloss.JoinVertical(lipgloss.Center, mapRows...)
	mapView := m.styles.MapBorder.Width(45).Align(lipgloss.Center).Render(mapContent)

	//mapHeight := lipgloss.Height(mapView)
	cameraWidth := m.width - lipgloss.Width(mapView) - 4

	statsArt := m.styles.StatsArt.Render(playerArt)
	statsText := fmt.Sprintf(
		"HP: %d\nMana: %d\nSpeed: %d\nMagic: %d\nStrength: %d",
		m.stats.hp,
		m.stats.mana,
		m.stats.speed,
		m.stats.magic,
		m.stats.strength,
	)
	statsContent := lipgloss.JoinHorizontal(lipgloss.Top, statsArt, statsText)
	statsView := m.styles.Panel.Width(cameraWidth).Render(statsContent)

	cameraHeight := 2

	currentRoom := m.worldMap[m.playerMapY][m.playerMapX]
	cameraContent := fmt.Sprintf("%s", currentRoom.getRoomDescription())
	cameraView := m.styles.Panel.Width(cameraWidth).Height(cameraHeight).Render(cameraContent)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, cameraView, statsView)

	help := m.styles.Faint.Padding(0, 1).Render("Arrows/wasd: move | 'q': quit")
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, mapView)
	finalView := lipgloss.JoinVertical(lipgloss.Left, mainView, help)

	return finalView
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
