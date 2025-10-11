package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
)

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {

	initialModel := model{
		state:      stateLoading,
		menuCursor: 0,
	}

	return initialModel, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	}

	switch m.state {
	case stateLoading:
		return m.updateLoading(msg)
	case stateMenu:
		return m.updateMenu(msg)
	case stateGame:
		return m.updateGame(msg)
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.state {
	case stateLoading:
		return m.renderLoadingView()
	case stateMenu:
		return m.renderMenuView()
	case stateGame:
		return m.renderGameView()
	default:
		return "Unknown state"
	}
}
