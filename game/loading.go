package game

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tickMsg:
		if m.progress.Percent() >= 1.0 {
			m.state = StateMenu
			return m, nil
		}
		progressCmd := m.progress.IncrPercent(0.02)
		return m, tea.Batch(tickCmd(), progressCmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m model) renderLoadingView() string {
	title := m.styles.Title.Render("SSH Dungeon Crawler")

	progressBar := m.progress.View()

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", progressBar)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
