package game

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

var (
	orange         = lipgloss.Color("#F25912")
	indigo         = lipgloss.Color("#5C3E94")
	darkPurple     = lipgloss.Color("#412B6B")
	veryDarkPurple = lipgloss.Color("#211832")
)

type styles struct {
	Title       lipgloss.Style
	Selected    lipgloss.Style
	Faint       lipgloss.Style
	Help        lipgloss.Style
	Panel       lipgloss.Style
	MapBorder   lipgloss.Style
	Player      lipgloss.Style
	Room        lipgloss.Style
	RoomSpecial lipgloss.Style
	StatsArt    lipgloss.Style
}

func newStyles(s ssh.Session) styles {
	renderer := bubbletea.MakeRenderer(s)
	return styles{
		Title:       renderer.NewStyle().Foreground(orange).Bold(true),
		Selected:    renderer.NewStyle().Foreground(indigo).Bold(true),
		Faint:       renderer.NewStyle().Faint(true),
		Help:        renderer.NewStyle().Foreground(orange),
		Panel:       renderer.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(indigo),
		MapBorder:   renderer.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(indigo),
		Player:      renderer.NewStyle().Width(3).Align(lipgloss.Center).Foreground(orange).SetString("[@]"),
		Room:        renderer.NewStyle().Width(3).Align(lipgloss.Center),
		RoomSpecial: renderer.NewStyle().Foreground(indigo),
		StatsArt:    renderer.NewStyle().Foreground(orange).Bold(true).Margin(1, 2),
	}
}
