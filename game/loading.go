package game

import (
	"strings"
	"time"

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
		if m.loadingProgress >= 1.0 {
			m.state = stateMenu
			return m, nil
		}
		m.loadingProgress += 0.01
		return m, tickCmd()
	}
	return m, nil
}

func (m model) renderLoadingView() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("228")).Render("SSH Dungeon Crawler")

	barWidth := m.width / 2
	progress := int(m.loadingProgress * float64(barWidth))
	progressBar := strings.Repeat("█", progress) + strings.Repeat("░", barWidth-progress)

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", progressBar)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
