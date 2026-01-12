package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/dev-pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type statusOptions struct {
	output string
	target string
	task   string
	branch string
	abs    bool
}

type statusRow struct {
	Task       string `json:"task"`
	Branch     string `json:"branch"`
	Path       string `json:"path"`
	Base       string `json:"base"`
	Target     string `json:"target"`
	LastCommit string `json:"last_commit"`
	Dirty      bool   `json:"dirty"`
	Ahead      int    `json:"ahead"`
	Behind     int    `json:"behind"`
}

func newStatusCommand(state *runState) *cobra.Command {
	opts := &statusOptions{output: "table"}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show detailed worktree status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			runner := defaultRunner()
			repoRoot, err := repoRoot(ctx, runner)
			if err != nil {
				return err
			}
			repo := repoName(repoRoot)

			target := opts.target
			if target == "" {
				current, err := git.CurrentBranch(ctx, runner)
				if err != nil {
					return err
				}
				target = current
			}

			worktrees, err := worktree.List(ctx, runner, repoRoot)
			if err != nil {
				return err
			}

			rows := make([]statusRow, 0, len(worktrees))
			for _, wt := range worktrees {
				branch := strings.TrimPrefix(wt.Branch, "refs/heads/")
				task, _ := worktree.TaskFromPath(repo, wt.Path)
				if task == "" {
					task = "-"
				}
				if opts.task != "" && task != opts.task {
					continue
				}
				if opts.branch != "" && branch != opts.branch {
					continue
				}

				statusInfo, err := worktree.Status(ctx, runner, wt.Path, target)
				if err != nil {
					return err
				}

				rows = append(rows, statusRow{
					Task:       task,
					Branch:     branch,
					Path:       displayPath(repoRoot, wt.Path, opts.abs),
					Base:       statusInfo.Base,
					Target:     target,
					LastCommit: statusInfo.LastCommit,
					Dirty:      statusInfo.Dirty,
					Ahead:      statusInfo.Ahead,
					Behind:     statusInfo.Behind,
				})
			}

			return renderStatus(cmd, opts.output, rows)
		},
	}

	cmd.Flags().StringVar(&opts.output, "output", opts.output, "output format: table or json")
	cmd.Flags().StringVar(&opts.target, "target", "", "target branch for ahead/behind comparison")
	cmd.Flags().StringVar(&opts.task, "task", "", "filter by task name")
	cmd.Flags().StringVar(&opts.branch, "branch", "", "filter by branch name")
	cmd.Flags().BoolVar(&opts.abs, "absolute-path", false, "show absolute paths instead of relative")
	cmd.Flags().BoolVar(&opts.abs, "abs", false, "alias for --absolute-path")

	return cmd
}

func renderStatus(cmd *cobra.Command, format string, rows []statusRow) error {
	switch format {
	case "table":
		columns := []tableColumn{
			{Header: "TASK"},
			{Header: "BRANCH", Style: func(value string) lipgloss.Style { return ui.AccentStyle }},
			{Header: "PATH"},
			{Header: "BASE", Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
			{Header: "TARGET", Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
			{Header: "LAST_COMMIT", Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
			{Header: "DIRTY", Style: func(value string) lipgloss.Style {
				if value == "true" {
					return ui.WarningStyle
				}
				return ui.SuccessStyle
			}},
			{Header: "AHEAD", Style: func(value string) lipgloss.Style {
				if value != "0" {
					return ui.WarningStyle
				}
				return ui.MutedStyle
			}},
			{Header: "BEHIND", Style: func(value string) lipgloss.Style {
				if value != "0" {
					return ui.ErrorStyle
				}
				return ui.MutedStyle
			}},
		}
		tableRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			tableRows = append(tableRows, []string{
				row.Task,
				row.Branch,
				row.Path,
				row.Base,
				row.Target,
				row.LastCommit,
				strconv.FormatBool(row.Dirty),
				strconv.Itoa(row.Ahead),
				strconv.Itoa(row.Behind),
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
