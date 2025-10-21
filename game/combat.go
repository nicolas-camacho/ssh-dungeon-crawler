package game

import (
	"fmt"
	"math/rand"
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

	playerStatsContentHeight := lipgloss.Height(playerStatsContent)

	topSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.styles.Panel.Width(turnOrderWidth).Height(playerStatsContentHeight).Render(turnOrderContent),
		m.styles.Panel.Width(enemiesWitdh).Height(playerStatsContentHeight).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Render(enemiesContent),
		m.styles.Panel.Width(playerStatsWidth).Render(playerStatsContent),
	)

	var middleSection string
	if m.combat.isEnemyTurnInProgress {
		enemyName := m.combat.turnOrder[m.combat.turnIndex].GetName()
		actionText := fmt.Sprintf("%s is attacking!", enemyName)
		progressBar := m.combat.enemyActionProgress.View()
		middleSection = lipgloss.JoinVertical(lipgloss.Center, actionText, progressBar)
	} else {
		middleSection = m.renderActionMenu()
	}

	actionMenu := m.styles.Panel.
		Width(m.width - m.styles.Panel.GetHorizontalFrameSize()).
		AlignHorizontal(lipgloss.Center).
		Render(middleSection)

	var helpText string
	switch m.combat.actionState {
	case ActionSelect:
		helpText = "Use ← → to select an action. Press Enter to confirm."
	case AttackSelect, MagicSelect, ItemSelect:
		helpText = "Use ← → to select an option. Press Enter to confirm. Esc to go back."
	case TargetSelect:
		helpText = "Use ← → to select a target. Press Enter to confirm. Esc to go back."
	}

	if m.combat.isEnemyTurnInProgress {
		helpText = "Enemy is taking action..."
	}

	helpView := m.styles.Help.Padding(0, 1).Render(helpText)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topSection,
		actionMenu,
		helpView,
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
	switch msg.String() {
	case "left", "a":
		if m.combat.subActionCursor > 0 {
			m.combat.subActionCursor--
		}
	case "right", "d":
		if m.combat.actionState == AttackSelect && m.combat.subActionCursor < len(m.combat.player.Attacks)-1 {
			m.combat.subActionCursor++
		} else if m.combat.actionState == MagicSelect && m.combat.subActionCursor < len(m.combat.player.Magics)-1 {
			m.combat.subActionCursor++
		}
	case "esc":
		m.combat.actionState = ActionSelect
		m.combat.subActionCursor = 0
		return m, nil
	case "enter":
		if m.combat.actionState == AttackSelect {
			selectedAttack := m.combat.player.Attacks[m.combat.subActionCursor]

			if selectedAttack.Sides == 0 {
				m.applyEffects(m.combat.player, nil, selectedAttack.Effects)
				m = m.advanceTurn()
				if len(m.combat.turnOrder) > 0 && !m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
					return m.startEnemyTurn()
				}
				return m, nil
			}
		}

		if m.combat.actionState == MagicSelect {
			selectedMagic := m.combat.player.Magics[m.combat.subActionCursor]

			if selectedMagic.Sides == 0 {
				if m.combat.player.data.stats.mana < selectedMagic.Cost {
					return m, nil
				}
				m.combat.player.data.stats.mana -= selectedMagic.Cost
				m.applyEffects(m.combat.player, nil, selectedMagic.Effects)
				m = m.advanceTurn()
				if len(m.combat.turnOrder) > 0 && !m.combat.turnOrder[m.combat.turnIndex].IsPlayer() {
					return m.startEnemyTurn()
				}
				return m, nil
			}
		}

		if m.combat.actionState == ItemSelect {
			var itemIDs []string
			for id := range m.combat.player.data.inventory {
				itemIDs = append(itemIDs, id)
			}

			if len(itemIDs) == 0 {
				return m, nil
			}

			selectedItemID := itemIDs[m.combat.subActionCursor]
			item := ItemTemplates[selectedItemID]

			if item.Effect == "heal" {
				m.combat.player.data.stats.hp += item.Value
				if m.combat.player.data.stats.hp > m.combat.player.GetMaxHP() {
					m.combat.player.data.stats.hp = m.combat.player.GetMaxHP()
				}
			}

			m.combat.player.data.inventory[selectedItemID]--
			if m.combat.player.data.inventory[selectedItemID] <= 0 {
				delete(m.combat.player.data.inventory, selectedItemID)
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
	case "esc":
		switch m.combat.actionCursor {
		case 0:
			m.combat.actionState = AttackSelect
		case 1:
			m.combat.actionState = MagicSelect
		case 3:
			m.combat.actionState = ItemSelect
		}
		m.combat.targetCursor = 0
		return m, nil
	case "enter":
		target := aliveEnemies[m.combat.targetCursor]

		switch m.combat.actionCursor {
		case 0:
			selectedAttack := m.combat.player.Attacks[m.combat.subActionCursor]
			roll := rand.Intn(selectedAttack.Sides) + 1

			damage := roll + m.combat.player.data.stats.strength - (target.Defense / 2)
			if damage < 1 {
				damage = 1
			}
			target.TakeDamage(damage)

			m.applyEffects(m.combat.player, target, selectedAttack.Effects)
		case 1:
			selectedMagic := m.combat.player.Magics[m.combat.subActionCursor]

			if m.combat.player.data.stats.mana < selectedMagic.Cost {
				return m, nil
			}

			m.combat.player.data.stats.mana -= selectedMagic.Cost

			roll := rand.Intn(selectedMagic.Sides) + 1
			damage := roll + m.combat.player.data.stats.magic - (target.Defense / 2)
			if damage < 1 {
				damage = 1
			}
			target.TakeDamage(damage)

			m.applyEffects(m.combat.player, target, selectedMagic.Effects)
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

	statsText := fmt.Sprintf(
		"\n\nHP: %d/%d\nMP: %d\n\nSTR: %d\nMAG: %d\nSPD: %d\nDEF: %d",
		m.combat.player.GetHP(),
		m.combat.player.GetMaxHP(),
		m.combat.player.data.stats.mana,
		m.combat.player.data.stats.strength,
		m.combat.player.data.stats.magic,
		m.combat.player.data.stats.speed,
		m.combat.player.data.stats.defense,
	)

	s += statsText

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

	currentTurnEntity := m.combat.turnOrder[m.combat.turnIndex]

	for i, enemy := range aliveEnemies {
		hp := fmt.Sprintf("HP: %d/%d", enemy.GetHP(), enemy.GetMaxHP())
		name := enemy.GetName()

		view := lipgloss.JoinVertical(lipgloss.Center, hp, name)
		enemyStyle := lipgloss.NewStyle().Margin(0, 2)

		if enemy == currentTurnEntity {
			enemyStyle = enemyStyle.Foreground(lipgloss.Color("#F25912")).Bold(true)
		}
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
		var attackOptions []string
		for i, attack := range m.combat.player.Attacks {
			name := attack.Name
			if i == m.combat.subActionCursor {
				name = m.styles.Selected.Render(name)
			}
			attackOptions = append(attackOptions, name)
		}
		content = "Attacks: " + strings.Join(attackOptions, " | ")
	case MagicSelect:
		var magicOptions []string
		for i, magic := range m.combat.player.Magics {
			optionText := fmt.Sprintf("%s (%d)", magic.Name, magic.Cost)

			style := lipgloss.NewStyle()
			if i == m.combat.subActionCursor {
				style = m.styles.Selected
			}

			if magic.Cost > m.combat.player.data.stats.mana {
				style = m.styles.Faint
			}
			magicOptions = append(magicOptions, style.Render(optionText))
		}
		content = "Magics: " + strings.Join(magicOptions, " | ")
	case ItemSelect:
		var itemOptions []string
		var itemIDs []string
		for id := range m.combat.player.data.inventory {
			itemIDs = append(itemIDs, id)
		}
		if len(itemIDs) == 0 {
			content = "Items: (Empty)"
		} else {
			for i, id := range itemIDs {
				item := ItemTemplates[id]
				count := m.combat.player.data.inventory[id]
				optionText := fmt.Sprintf("%s (x%d)", item.Name, count)
				if i == m.combat.subActionCursor {
					optionText = m.styles.Selected.Render(optionText)
				}
				itemOptions = append(itemOptions, optionText)
			}
			content = "Items: " + strings.Join(itemOptions, " | ")
		}
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

func (m *model) applyEffects(source *Player, target *Foe, effect []Effect) {
	for _, effect := range effect {
		var value int
		if effect.Sides > 0 {
			value = rand.Intn(effect.Sides) + 1
		} else if effect.Sides < 0 {
			value = -(rand.Intn(-effect.Sides) + 1)
		}

		var affectedEntity CombatEntity
		if effect.Target == "Self" {
			affectedEntity = source
		} else {
			affectedEntity = target
		}

		affectedEntity.ModifyStat(effect.Stat, value)
	}
}
