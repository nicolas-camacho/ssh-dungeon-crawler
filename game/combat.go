package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type enemyTickMsg time.Time

func enemyTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return enemyTickMsg(t)
	})
}

func (m model) updateCombat(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case enemyTickMsg:
		if m.combat.enemyActionProgress.Percent() >= 1.0 {
			enemy := m.combat.turnOrder[m.combat.turnIndex].(*Foe)
			damage := enemy.Attack
			if m.combat.player.isDefending {
				damage /= 2
			}
			m.combat.player.TakeDamage(damage)

			m.combat.isEnemyTurnInProgress = false
			m.combat.enemyActionProgress.SetPercent(0)

			if m.combat.player.GetHP() <= 0 {
				m.state = StateMenu
				m.combat = nil
				return m, nil
			}

			m = m.advanceTurn()
			if !m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
				return m.startEnemyTurn()
			}
			return m, nil
		}

		cmd := m.combat.enemyActionProgress.IncrPercent(0.025)
		return m, tea.Batch(enemyTickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.combat.enemyActionProgress.Update(msg)
		m.combat.enemyActionProgress = progressModel.(progress.Model)
		return m, cmd

	case tea.KeyMsg:
		if m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
			switch m.combat.actionState {
			case ActionSelect:
				return m.handleActionSelect(msg)
			case AttackSelect, MagicSelect, ItemSelect:
				return m.handleSubActionSelect(msg)
			case TargetSelect:
				return m.handleTargetSelect(msg)
			}
		}
	}
	return m, nil
}

func (m model) renderCombatView() string {
	turnOrderContent := m.renderTurnOrder()
	playerStatsContent := m.renderPlayerStatsCombat()
	enemiesContent := m.renderEnemies()

	turnOrderWidth := 20
	playerStatsWidth := 25

	enemiesWitdh := m.width - turnOrderWidth - playerStatsWidth - 6

	turnOrderContentHeight := lipgloss.Height(turnOrderContent)

	topSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.Panel.Width(turnOrderWidth).Render(turnOrderContent),
		m.styles.Panel.Width(enemiesWitdh).Height(turnOrderContentHeight).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Render(enemiesContent),
		m.styles.Panel.Width(playerStatsWidth).Height(turnOrderContentHeight).Render(playerStatsContent),
	)

	var bottomSection string
	if m.combat.isEnemyTurnInProgress {
		enemyName := m.combat.turnOrder[m.combat.turnIndex].GetName()
		actionText := fmt.Sprintf("%s is attacking!", enemyName)
		progressBar := m.combat.enemyActionProgress.View()
		bottomSection = lipgloss.JoinVertical(lipgloss.Center, actionText, progressBar)
	} else {
		bottomSection = m.renderActionMenu()
	}

	actionMenu := m.styles.Panel.
		Width(m.width - m.styles.Panel.GetHorizontalFrameSize()).
		AlignHorizontal(lipgloss.Center).
		Render(bottomSection)

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
			m = m.advanceTurn()

			if !m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
				return m.startEnemyTurn()
			}
			return m, nil
		case 3:
			m.combat.actionState = ItemSelect
			m.combat.subActionCursor = 0
		}
	}
	return m, nil
}

func (m model) handleSubActionSelect(msg tea.KeyMsg) (model, tea.Cmd) {
	if msg.String() == "enter" {
		if m.combat.actionState == ItemSelect {
			m.stats.hp += 5
			if m.stats.hp > 100 {
				m.stats.hp = 100
			}
			m = m.advanceTurn()

			if !m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
				return m.startEnemyTurn()
			}
			return m, nil
		}
		m.combat.actionState = TargetSelect
		m.combat.targetCursor = 0
	}
	return m, nil
}

func (m model) handleTargetSelect(msg tea.KeyMsg) (model, tea.Cmd) {
	aliveEnemies := []*Foe{}

	for _, e := range m.combat.enemies {
		if e.GetHP() > 0 {
			aliveEnemies = append(aliveEnemies, e)
		}
	}

	if len(aliveEnemies) == 0 {
		m.state = StateGame
		m.combat = nil
		return m, nil
	}

	switch msg.String() {
	case "left", "a":
		if m.combat.targetCursor > 0 {
			m.combat.targetCursor--
		}
	case "right", "d":
		if m.combat.targetCursor < len(aliveEnemies)-1 {
			m.combat.targetCursor++
		}
	case "enter":
		target := aliveEnemies[m.combat.targetCursor]

		if m.combat.actionCursor == 0 {
			target.TakeDamage(10)
		}
		if m.combat.actionCursor == 1 {
			target.TakeDamage(15)
		}

		if target.GetHP() <= 0 {
			hasAliveEnemies := false
			for _, e := range m.combat.enemies {
				if e.GetHP() > 0 {
					hasAliveEnemies = true
					break
				}
			}
			if !hasAliveEnemies {
				m.state = StateGame
				m.combat = nil
				return m, nil
			}
		}
		m = m.advanceTurn()
		if len(m.combat.turnOrder) > 0 && !m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
			return m.startEnemyTurn()
		}
		return m, nil
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

		enemyStyle := lipgloss.NewStyle().Margin(0, 2)
		view = enemyStyle.Render(view)

		if m.combat.actionState == TargetSelect && i == m.combat.targetCursor {
			view = m.styles.Selected.Border(lipgloss.RoundedBorder(), true).Padding(0, 1).Align(lipgloss.Center).Render(view)
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

func (m model) startEnemyTurn() (model, tea.Cmd) {
	m.combat.isEnemyTurnInProgress = true
	m.combat.enemyActionProgress.Width = m.width - m.styles.Panel.GetHorizontalFrameSize()
	return m, enemyTickCmd()
}
