package game

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
)

func CreateTeaProgram(s ssh.Session, startState GameState) (tea.Model, []tea.ProgramOption) {

	prog := progress.New(
		progress.WithGradient(string(orange), string(indigo)),
		progress.WithoutPercentage(),
	)

	initialModel := model{
		state:      startState,
		menuCursor: 0,
		progress:   prog,
		styles:     newStyles(),
	}

	if startState == StateCombat {
		initialModel.combat = newTestCombatState()
		initialModel.stats = *initialModel.combat.player.stats
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
	case StateLoading:
		return m.updateLoading(msg)
	case StateMenu:
		return m.updateMenu(msg)
	case StateGame:
		return m.updateGame(msg)
	case StateCombat:
		return m.updateCombat(msg)
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.state {
	case StateLoading:
		return m.renderLoadingView()
	case StateMenu:
		return m.renderMenuView()
	case StateGame:
		return m.renderGameView()
	case StateCombat:
		return m.renderCombatView()
	default:
		return "Unknown state"
	}
}
