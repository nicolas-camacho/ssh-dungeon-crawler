package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "w":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
		case "down", "s":
			if m.menuCursor < 1 {
				m.menuCursor++
			}
		case "enter":
			if m.menuCursor == 0 {
				m.state = StateGame

				firstFloor, startX, startY := generateMap(9, 9, 15, 0)
				m.floors = []floor{*firstFloor}
				m.currentFloor = 0

				m.playerMapX = startX
				m.playerMapY = startY
				m.player = playerData{
					stats: playerStats{
						hp:       100,
						mana:     50,
						speed:    10,
						magic:    12,
						strength: 8,
						defense:  8,
					},
					inventory: make(map[string]int),
				}

				m.floors[m.currentFloor].worldMap[startY][startX].Visited = true
			} else {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) renderMenuView() string {
	title := m.styles.Title.Render("SSH Dungeon Crawler")

	var start, exit string
	if m.menuCursor == 0 {
		start = m.styles.Selected.Render("> Start Game")
		exit = "  Exit"
	} else {
		start = "  Start Game"
		exit = m.styles.Selected.Render("> Exit")
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, start, exit)
	help := m.styles.Faint.Render("Arrows: navigation | 'enter': select")

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", menu, "", help)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
