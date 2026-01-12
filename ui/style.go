package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	MutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	BorderStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
)
