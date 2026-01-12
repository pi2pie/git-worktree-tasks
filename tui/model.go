package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/dev-pi2pie/git-worktree-tasks/ui"
)

type model struct {
	ready bool
}

func NewModel() tea.Model {
	return &model{}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string {
	title := ui.TitleStyle.Render("git-worktree-tasks")
	body := ui.MutedStyle.Render("TUI scaffold: create, finish, cleanup, list, status")
	return ui.BorderStyle.Render(title + "\n" + body + "\n\nPress q to quit.")
}
