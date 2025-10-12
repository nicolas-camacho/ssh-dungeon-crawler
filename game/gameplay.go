package game

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentMap := m.floors[m.currentFloor].worldMap

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		prevX, prevY := m.playerMapX, m.playerMapY

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "w":
			if m.playerMapY > 0 && currentMap[m.playerMapY-1][m.playerMapX] != nil {
				m.playerMapY--
				currentMap[m.playerMapY][m.playerMapX].Visited = true
			}
		case "down", "s":
			if m.playerMapY < len(currentMap)-1 && currentMap[m.playerMapY+1][m.playerMapX] != nil {
				m.playerMapY++
				currentMap[m.playerMapY][m.playerMapX].Visited = true
			}
		case "left", "a":
			if m.playerMapX > 0 && currentMap[m.playerMapY][m.playerMapX-1] != nil {
				m.playerMapX--
				currentMap[m.playerMapY][m.playerMapX].Visited = true
			}
		case "right", "d":
			if m.playerMapX < len(currentMap[0])-1 && currentMap[m.playerMapY][m.playerMapX+1] != nil {
				m.playerMapX++
				currentMap[m.playerMapY][m.playerMapX].Visited = true
			}
		case "enter", "x":
			currentRoom := currentMap[m.playerMapY][m.playerMapX]
			switch currentRoom.Type {

			case StairsUp:
				m.currentFloor++
				if m.currentFloor >= len(m.floors) {
					newFloor, startX, startY := generateMap(9, 9, 15, m.currentFloor)
					m.floors = append(m.floors, *newFloor)
					m.playerMapX, m.playerMapY = startX, startY
					m.floors[m.currentFloor].worldMap[startY][startX].Visited = true
				} else {
					for y, row := range m.floors[m.currentFloor].worldMap {
						for x, room := range row {
							if room != nil && room.Type == StairsDown {
								m.playerMapX, m.playerMapY = x, y
								break
							}
						}
					}
				}
			case StairsDown:
				if m.currentFloor > 0 {
					m.currentFloor--
					for y, row := range m.floors[m.currentFloor].worldMap {
						for x, room := range row {
							if room != nil && room.Type == StairsUp {
								m.playerMapX, m.playerMapY = x, y
								break
							}
						}
					}
				}
			}
		}

		if prevX != m.playerMapX || prevY != m.playerMapY {
			newRoom := currentMap[m.playerMapY][m.playerMapX]
			newRoom.Visited = true

			if newRoom.Type == Enemy {
				m.state = stateCombat

				playerEntity := &Player{stats: &m.stats}
				numEnemies := 1 + rand.Intn(3)
				enemies := make([]*Foe, numEnemies)
				for i := range enemies {
					enemies[i] = newGoblin()
				}

				turnOrder := calculateTurnOrder(playerEntity, enemies)

				m.combat = &CombatState{
					player:      playerEntity,
					enemies:     enemies,
					turnOrder:   turnOrder,
					turnIndex:   0,
					actionState: ActionSelect,
				}

				newRoom.Type = Empty
			}
		}
	}
	return m, nil
}

func (m model) renderGameView() string {
	currentMap := m.floors[m.currentFloor].worldMap
	currentRoom := currentMap[m.playerMapY][m.playerMapX]

	emptyCell := lipgloss.NewStyle().Width(3).SetString(" ")

	var mapRows []string
	for y, row := range currentMap {
		var mapRow strings.Builder
		for x, room := range row {
			if x == m.playerMapX && y == m.playerMapY {
				mapRow.WriteString(m.styles.Player.String())
			} else if room != nil {
				var symbol string
				if room.Visited {
					symbol = room.getRoomSymbol()
				} else {
					symbol = "?"
				}

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
		"HP: %d\nMana: %d\nSpeed: %d\nMagic: %d\nStrength: %d \nDefense: %d",
		m.stats.hp,
		m.stats.mana,
		m.stats.speed,
		m.stats.magic,
		m.stats.strength,
		m.stats.defense,
	)
	statsContent := lipgloss.JoinHorizontal(lipgloss.Top, statsArt, statsText)
	statsView := m.styles.Panel.Width(cameraWidth).Render(statsContent)

	cameraHeight := 2

	cameraContent := fmt.Sprintf("%s", currentRoom.getRoomDescription())
	cameraView := m.styles.Panel.Width(cameraWidth).Height(cameraHeight).Render(cameraContent)

	leftPanel := lipgloss.JoinVertical(lipgloss.Left, cameraView, statsView)

	helpText := "Arrows/wasd: move | 'q': quit"
	if currentRoom.Type == StairsUp || currentRoom.Type == StairsDown {
		helpText += " | 'enter'/'x': Use Stairs"
	}
	help := m.styles.Faint.Padding(0, 1).Render(helpText)

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
		return "Some stone stairs, they take you to the darkness\nYou wanna go up?"
	case StairsDown:
		return "You can go down again"
	default:
		return "Unknown room type."
	}
}
