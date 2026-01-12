package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/dev-pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type listOptions struct {
	output string
	task   string
	branch string
	abs    bool
}

type listRow struct {
	Task    string `json:"task"`
	Branch  string `json:"branch"`
	Path    string `json:"path"`
	Present bool   `json:"present"`
	Head    string `json:"head"`
}

func newListCommand(state *runState) *cobra.Command {
	opts := &listOptions{output: "table"}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List task worktrees",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			runner := defaultRunner()
			repoRoot, err := repoRoot(ctx, runner)
			if err != nil {
				return err
			}
			repo := repoName(repoRoot)

			worktrees, err := worktree.List(ctx, runner, repoRoot)
			if err != nil {
				return err
			}

			rows := make([]listRow, 0, len(worktrees))
			for _, wt := range worktrees {
				branch := strings.TrimPrefix(wt.Branch, "refs/heads/")
				task, _ := worktree.TaskFromPath(repo, wt.Path)
				if task == "" {
					task = "-"
				}
				row := listRow{
					Task:    task,
					Branch:  branch,
					Path:    displayPath(repoRoot, wt.Path, opts.abs),
					Present: true,
					Head:    worktree.ShortHash(wt.Head),
				}
				if opts.task != "" && row.Task != opts.task {
					continue
				}
				if opts.branch != "" && row.Branch != opts.branch {
					continue
				}
				rows = append(rows, row)
			}

			return renderList(cmd, opts.output, rows)
		},
	}

	cmd.Flags().StringVar(&opts.output, "output", opts.output, "output format: table or json")
	cmd.Flags().StringVar(&opts.task, "task", "", "filter by task name")
	cmd.Flags().StringVar(&opts.branch, "branch", "", "filter by branch name")
	cmd.Flags().BoolVar(&opts.abs, "absolute-path", false, "show absolute paths instead of relative")
	cmd.Flags().BoolVar(&opts.abs, "abs", false, "alias for --absolute-path")

	return cmd
}

func renderList(cmd *cobra.Command, format string, rows []listRow) error {
	switch format {
	case "table":
		columns := []tableColumn{
			{Header: "TASK"},
			{Header: "BRANCH", Style: func(value string) lipgloss.Style { return ui.AccentStyle }},
			{Header: "PATH"},
			{Header: "PRESENT", Style: func(value string) lipgloss.Style {
				if value == "true" {
					return ui.SuccessStyle
				}
				return ui.ErrorStyle
			}},
			{Header: "HEAD", Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
		}
		tableRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			tableRows = append(tableRows, []string{
				row.Task,
				row.Branch,
				row.Path,
				strconv.FormatBool(row.Present),
				row.Head,
			})
		}
		renderTable(cmd, columns, tableRows)
		return nil
	case "json":
		payload, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(payload))
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
