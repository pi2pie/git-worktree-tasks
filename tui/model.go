package tui

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/dev-pi2pie/git-worktree-tasks/ui"
)

type listRow struct {
	Task    string
	Branch  string
	Path    string
	Present string
	Head    string
}

type columnSpec struct {
	Title    string
	MinWidth int
	MaxWidth int
	Flexible bool
}

type model struct {
	table  table.Model
	rows   []listRow
	err    error
	ready  bool
	width  int
	height int
}

type rowsMsg []listRow
type errMsg struct{ err error }

func NewModel() tea.Model {
	return &model{}
}

func (m *model) Init() tea.Cmd {
	return loadRowsCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case rowsMsg:
		m.rows = msg
		m.ready = true
		if m.width == 0 {
			m.width = 120
		}
		m.table = newListTable(m.width, m.height, m.rows)
		return m, nil
	case errMsg:
		m.err = msg.err
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.ready {
			m.table = updateListTableLayout(m.table, m.width, m.height, m.rows)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.ready {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	title := ui.TitleStyle.Render("git-worktree-tasks")
	if m.err != nil {
		message := ui.ErrorStyle.Render(m.err.Error())
		return ui.BorderStyle.Render(title + "\n" + message + "\n\nPress q to quit.")
	}
	if !m.ready {
		message := ui.MutedStyle.Render("Loading worktrees...")
		return ui.BorderStyle.Render(title + "\n" + message + "\n\nPress q to quit.")
	}
	help := ui.MutedStyle.Render("j/k or arrows to move • q to quit")
	body := m.table.View()
	return ui.BorderStyle.Render(title + "\n" + body + "\n" + help)
}

func loadRowsCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		runner := git.ExecRunner{}
		repoRoot, err := git.RepoRoot(ctx, runner)
		if err != nil {
			return errMsg{err: err}
		}
		repo := repoName(repoRoot)
		worktrees, err := worktree.List(ctx, runner, repoRoot)
		if err != nil {
			return errMsg{err: err}
		}
		rows := make([]listRow, 0, len(worktrees))
		for _, wt := range worktrees {
			branch := strings.TrimPrefix(wt.Branch, "refs/heads/")
			task, _ := worktree.TaskFromPath(repo, wt.Path)
			if task == "" {
				task = "-"
			}
			rows = append(rows, listRow{
				Task:    task,
				Branch:  branch,
				Path:    displayPath(repoRoot, wt.Path),
				Present: "true",
				Head:    worktree.ShortHash(wt.Head),
			})
		}
		return rowsMsg(rows)
	}
}

func newListTable(width, height int, rows []listRow) table.Model {
	t := table.New(
		table.WithColumns(listColumns(width, rows)),
		table.WithRows(listTableRows(rows)),
		table.WithFocused(true),
	)
	t = updateListTableLayout(t, width, height, rows)
	styles := table.DefaultStyles()
	styles.Header = ui.HeaderStyle
	styles.Selected = ui.AccentStyle.Copy().Reverse(true)
	t.SetStyles(styles)
	return t
}

func updateListTableLayout(t table.Model, width, height int, rows []listRow) table.Model {
	t.SetColumns(listColumns(width, rows))
	t.SetHeight(tableHeight(height))
	return t
}

func listColumns(width int, rows []listRow) []table.Column {
	specs := []columnSpec{
		{Title: "TASK", MinWidth: 6},
		{Title: "BRANCH", MinWidth: 10, Flexible: true},
		{Title: "PATH", MinWidth: 16, Flexible: true},
		{Title: "PRESENT", MinWidth: 7},
		{Title: "HEAD", MinWidth: 7},
	}
	widths := computeColumnWidths(width, specs, rows)
	columns := make([]table.Column, 0, len(specs))
	for i, spec := range specs {
		columns = append(columns, table.Column{
			Title: spec.Title,
			Width: widths[i],
		})
	}
	return columns
}

func listTableRows(rows []listRow) []table.Row {
	tableRows := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		tableRows = append(tableRows, table.Row{
			row.Task,
			row.Branch,
			row.Path,
			row.Present,
			row.Head,
		})
	}
	return tableRows
}

func computeColumnWidths(totalWidth int, specs []columnSpec, rows []listRow) []int {
	if totalWidth <= 0 {
		totalWidth = 120
	}
	usableWidth := totalWidth - 6
	if usableWidth < 40 {
		usableWidth = totalWidth
	}

	widths := make([]int, len(specs))
	for i, spec := range specs {
		widths[i] = runewidth.StringWidth(spec.Title)
		if spec.MinWidth > widths[i] {
			widths[i] = spec.MinWidth
		}
	}

	for _, row := range rows {
		values := []string{row.Task, row.Branch, row.Path, row.Present, row.Head}
		for i, value := range values {
			if i >= len(widths) {
				continue
			}
			if w := runewidth.StringWidth(value); w > widths[i] {
				widths[i] = w
			}
		}
	}

	for i, spec := range specs {
		if spec.MaxWidth > 0 && widths[i] > spec.MaxWidth {
			widths[i] = spec.MaxWidth
		}
	}

	total := tableWidth(widths)
	if total <= usableWidth {
		return widths
	}

	minWidths := make([]int, len(specs))
	for i, spec := range specs {
		minWidth := spec.MinWidth
		if minWidth <= 0 {
			minWidth = 3
		}
		if minWidth > widths[i] {
			minWidth = widths[i]
		}
		minWidths[i] = minWidth
	}

	for total > usableWidth {
		reduced := false
		for i, spec := range specs {
			if !spec.Flexible {
				continue
			}
			if widths[i] > minWidths[i] {
				widths[i]--
				total--
				reduced = true
				if total <= usableWidth {
					break
				}
			}
		}
		if !reduced {
			break
		}
	}

	return widths
}

func tableWidth(widths []int) int {
	total := 0
	for _, width := range widths {
		total += width
	}
	if len(widths) > 1 {
		total += 2 * (len(widths) - 1)
	}
	return total
}

func tableHeight(height int) int {
	if height <= 0 {
		return 12
	}
	height = height - 6
	if height < 6 {
		return 6
	}
	return height
}

func repoName(root string) string {
	return filepath.Base(root)
}

func displayPath(repoRoot, path string) string {
	clean := filepath.Clean(path)
	absPath := clean
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(repoRoot, absPath)
	}
	absPath = filepath.Clean(absPath)
	rel, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		return clean
	}
	return rel
}
