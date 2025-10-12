package game

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateCombat(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
		enemy := m.combat.turnOrder[m.combat.turnIndex].(*Foe)
		m.combat.player.TakeDamage(enemy.Attack)

		if m.combat.player.GetHP() <= 0 {
			m.state = stateMenu
			m.combat = nil
			return m, nil
		}

		m = m.advanceTurn()
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg { return nil })
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.combat.actionState {
		case ActionSelect:
			return m.handleActionSelect(msg)
		case AttackSelect, MagicSelect, ItemSelect:
			return m.handleSubActionSelect(msg)
		case TargetSelect:
			return m.handleTargetSelect(msg)
		}
	}
	return m, nil
}

func (m model) renderCombatView() string {
	turnOrderView := m.renderTurnOrder()

	playerStatsView := m.renderPlayerStatsCombat()
	enemiesView := m.renderEnemies()
	contentTop := lipgloss.JoinHorizontal(lipgloss.Top, playerStatsView, enemiesView)

	actionMenuView := m.renderActionMenu()

	mainContent := lipgloss.JoinVertical(lipgloss.Left, contentTop, actionMenuView)
	combatView := lipgloss.JoinHorizontal(lipgloss.Top, turnOrderView, mainContent)

	return m.styles.Panel.Render(combatView)
}

func (m model) handleActionSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "a":
		if m.combat.actionCursor > 0 {
			m.combat.actionCursor--
		}
	case "right", "d":
		if m.combat.actionCursor < 3 {
			m.combat.actionCursor++
		}
	case "enter":
		switch m.combat.actionCursor {
		case 0:
			m.combat.actionState = AttackSelect
			m.combat.subActionCursor = 0
		case 1:
			m.combat.actionState = MagicSelect
			m.combat.subActionCursor = 0
		case 2:
			m.combat.player.isDefending = true
			return m.advanceTurn(), nil
		case 3:
			m.combat.actionState = ItemSelect
			m.combat.subActionCursor = 0
		}
	}
	return m, nil
}

func (m model) handleSubActionSelect(msg tea.KeyMsg) (model, tea.Cmd) {
	if msg.String() == "enter" {
		m.combat.actionState = TargetSelect
		m.combat.targetCursor = 0
	}
	return m, nil
}

func (m model) handleTargetSelect(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "up", "w":
		if m.combat.targetCursor > 0 {
			m.combat.targetCursor--
		}
	case "down", "s":
		if m.combat.targetCursor < len(m.combat.enemies)-1 {
			m.combat.targetCursor++
		}
	case "enter":
		target := m.combat.enemies[m.combat.targetCursor]

		if m.combat.actionCursor == 0 {
			target.TakeDamage(10)
		} else if m.combat.actionCursor == 1 {
			target.TakeDamage(15)
		}

		if target.GetHP() <= 0 {
			var newEnemies []*Foe
			for i, e := range m.combat.enemies {
				if i != m.combat.targetCursor {
					newEnemies = append(newEnemies, e)
				}
			}
			m.combat.enemies = newEnemies

			if len(m.combat.enemies) == 0 {
				m.state = stateMenu
				m.combat = nil
				return m, nil
			}
		}
		return m.advanceTurn(), nil
	}
	return m, nil
}

func (m model) renderTurnOrder() string {
	var s strings.Builder
	for i, entity := range m.combat.turnOrder {
		name := entity.GetName()
		if i == m.combat.turnIndex {
			name = lipgloss.NewStyle().Background(m.styles.Selected.GetBackground()).Render(name)
		}
		s.WriteString(name + "\n")
	}
	return lipgloss.NewStyle().Padding(1, 2).Render(s.String())
}

func (m model) renderPlayerStatsCombat() string {
	stats := fmt.Sprintf("HP: %d/%d\nMana: %d", m.combat.player.GetHP(), m.combat.player.GetMaxHP(), m.stats.mana)
	return m.styles.Panel.Render(stats)
}

func (m model) renderEnemies() string {
	var enemyViews []string
	for i, enemy := range m.combat.enemies {
		hp := fmt.Sprintf("%d/%d", enemy.GetHP(), enemy.GetMaxHP())
		name := enemy.GetName()
		if m.combat.actionState == TargetSelect && i == m.combat.targetCursor {
			name = m.styles.Selected.Render("> " + name)
		}
		enemyViews = append(enemyViews, lipgloss.JoinVertical(lipgloss.Center, hp, name))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, enemyViews...)
}

func (m model) renderActionMenu() string {
	var content string
	switch m.combat.actionState {
	case ActionSelect:
		options := []string{"Attack", "Magic", "Defend", "Item"}
		var styledOptions []string
		for i, opt := range options {
			if i == m.combat.actionCursor {
				styledOptions = append(styledOptions, m.styles.Selected.Render(opt))
			} else {
				styledOptions = append(styledOptions, opt)
			}
		}
		content = lipgloss.JoinHorizontal(lipgloss.Top, styledOptions...)
	case AttackSelect:
		content = m.styles.Selected.Render("> Basic Attack")
	case MagicSelect:
		content = m.styles.Selected.Render("> Fireball")
	case ItemSelect:
		content = m.styles.Selected.Render("> Potion (+5 HP)")
	case TargetSelect:
		content = "Choose a target..."
	}
	return m.styles.Panel.Width(80).Align(lipgloss.Center).Render(content)
}

func (m model) advanceTurn() model {
	m.combat.turnIndex = (m.combat.turnIndex + 1) % len(m.combat.turnOrder)

	if m.combat.turnIndex == 0 {
		m.combat.turnOrder = calculateTurnOrder(m.combat.player, m.combat.enemies)
	}
	m.combat.actionState = ActionSelect
	m.combat.player.isDefending = false
	return m
}
