package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle   lipgloss.Style
	MutedStyle   lipgloss.Style
	BorderStyle  lipgloss.Style
	SuccessStyle lipgloss.Style
	WarningStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
	HeaderStyle  lipgloss.Style
	PromptStyle  lipgloss.Style
	AccentStyle  lipgloss.Style
)

func init() {
	SetColorEnabled(true)
}

func SetColorEnabled(enabled bool) {
	TitleStyle = lipgloss.NewStyle().Bold(true)
	MutedStyle = lipgloss.NewStyle()
	BorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
	SuccessStyle = lipgloss.NewStyle()
	WarningStyle = lipgloss.NewStyle()
	ErrorStyle = lipgloss.NewStyle().Bold(true)
	HeaderStyle = lipgloss.NewStyle().Bold(true)
	PromptStyle = lipgloss.NewStyle().Bold(true)
	AccentStyle = lipgloss.NewStyle().Bold(true)

	if !enabled {
		return
	}

	TitleStyle = TitleStyle.Foreground(lipgloss.Color("63"))
	MutedStyle = MutedStyle.Foreground(lipgloss.Color("244"))
	SuccessStyle = SuccessStyle.Foreground(lipgloss.Color("35"))
	WarningStyle = WarningStyle.Foreground(lipgloss.Color("214"))
	ErrorStyle = ErrorStyle.Foreground(lipgloss.Color("203"))
	HeaderStyle = HeaderStyle.Foreground(lipgloss.Color("81"))
	PromptStyle = PromptStyle.Foreground(lipgloss.Color("75"))
	AccentStyle = AccentStyle.Foreground(lipgloss.Color("213"))
}
