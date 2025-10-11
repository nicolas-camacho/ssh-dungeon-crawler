package game

import "github.com/charmbracelet/lipgloss"

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

func newStyles() styles {
	return styles{
		Title:       lipgloss.NewStyle().Foreground(orange).Bold(true),
		Selected:    lipgloss.NewStyle().Foreground(indigo).Bold(true),
		Faint:       lipgloss.NewStyle().Faint(true),
		Help:        lipgloss.NewStyle().Foreground(darkPurple),
		Panel:       lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(indigo),
		MapBorder:   lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(indigo),
		Player:      lipgloss.NewStyle().Width(3).Align(lipgloss.Center).Foreground(orange).SetString("[@]"),
		Room:        lipgloss.NewStyle().Width(3).Align(lipgloss.Center),
		RoomSpecial: lipgloss.NewStyle().Foreground(indigo),
		StatsArt:    lipgloss.NewStyle().Foreground(orange).Bold(true).Margin(1, 2),
	}
}
