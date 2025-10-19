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

		damage := enemy.Attack
		if m.combat.player.isDefending {
			damage /= 2
		}
		m.combat.player.TakeDamage(damage)

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
	turnOrderContent := m.renderTurnOrder()
	playerStatsContent := m.renderPlayerStatsCombat()
	enemiesContent := m.renderEnemies()
	actionMenuContent := m.renderActionMenu()

	turnOrderWidth := 20
	playerStatsWidth := 25

	enemiesWitdh := m.width - turnOrderWidth - playerStatsWidth - 6

	topSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.Panel.Width(turnOrderWidth).Render(turnOrderContent),
		m.styles.Panel.Width(enemiesWitdh).Align(lipgloss.Center).Render(enemiesContent),
		m.styles.Panel.Width(playerStatsWidth).Render(playerStatsContent),
	)

	actionMenu := m.styles.Panel.Width(m.width - m.styles.Panel.GetHorizontalFrameSize()).Render(actionMenuContent)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topSection,
		actionMenu,
	)
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
	s.WriteString(m.styles.Title.Render("Turnos") + "\n\n")
	for i, entity := range m.combat.turnOrder {
		if entity.GetHP() <= 0 {
			continue
		}
		name := entity.GetName()
		if i == m.combat.turnIndex {
			name = m.styles.Selected.Render("> " + name)
		} else {
			name = "  " + name
		}
		s.WriteString(name + "\n")
	}
	return s.String()
}

func (m model) renderPlayerStatsCombat() string {
	s := m.styles.Title.Render(m.combat.player.GetName())
	s += fmt.Sprintf("\n\nHP: %d/%d", m.combat.player.GetHP(), m.combat.player.GetMaxHP())
	s += fmt.Sprintf("\nMana: %d", m.stats.mana)
	if m.combat.player.isDefending {
		s += "\n\n" + m.styles.Selected.Render("Defending!")
	}
	return s
}

func (m model) renderEnemies() string {
	var enemyViews []string

	aliveEnemies := []*Foe{}
	for _, e := range m.combat.enemies {
		if e.GetHP() > 0 {
			aliveEnemies = append(aliveEnemies, e)
		}
	}

	for i, enemy := range aliveEnemies {
		hp := fmt.Sprintf("HP: %d/%d", enemy.GetHP(), enemy.GetMaxHP())
		name := enemy.GetName()

		view := lipgloss.JoinVertical(lipgloss.Center, hp, name)

		if m.combat.actionState == TargetSelect && i == m.combat.targetCursor {
			view = m.styles.Selected.Border(lipgloss.RoundedBorder(), true).Padding(0, 1).Render(view)
		}

		enemyViews = append(enemyViews, view)
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, enemyViews...)
}

func (m model) renderActionMenu() string {
	var content string
	switch m.combat.actionState {
	case ActionSelect:
		options := []string{"Attack", "Magic", "Defend", "Item"}
		var styledOptions []string
		for i, opt := range options {
			if i == m.combat.actionCursor {
				styledOptions = append(styledOptions, m.styles.Selected.Render("[ "+opt+" ]"))
			} else {
				styledOptions = append(styledOptions, "[ "+opt+" ]")
			}
		}
		content = lipgloss.JoinHorizontal(lipgloss.Top, styledOptions...)
	case AttackSelect:
		content = "Attacks: " + m.styles.Selected.Render("> Basic Attack")
	case MagicSelect:
		content = "Spells: " + m.styles.Selected.Render("> Fireball")
	case ItemSelect:
		content = "Items: " + m.styles.Selected.Render("> Health Potion(+5 HP)")
	case TargetSelect:
		content = "Select Target..."
	}
	return lipgloss.NewStyle().Align(lipgloss.Center).Render(content)
}

func (m model) advanceTurn() model {
	var aliveInTurnOrder []CombatEntity
	for _, entity := range m.combat.turnOrder {
		if entity.GetHP() > 0 {
			aliveInTurnOrder = append(aliveInTurnOrder, entity)
		}
	}
	m.combat.turnOrder = aliveInTurnOrder

	if len(m.combat.turnOrder) > 0 {
		m.combat.turnIndex = (m.combat.turnIndex + 1) % len(m.combat.turnOrder)
	}

	m.combat.actionState = ActionSelect
	m.combat.player.isDefending = false
	return m
}
