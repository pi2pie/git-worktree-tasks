package cli

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type listOptions struct {
	output string
	branch string
	field  string
	abs    bool
	grid   bool
	strict bool
}

type listRow struct {
	Task    string `json:"task"`
	Branch  string `json:"branch"`
	Path    string `json:"path"`
	Present bool   `json:"present"`
	Head    string `json:"head"`
}

func newListCommand() *cobra.Command {
	opts := &listOptions{output: "table"}
	cmd := &cobra.Command{
		Use:     "list [task]",
		Short:   "List task worktrees",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			runner := defaultRunner()
			repoRoot, err := repoRoot(ctx, runner)
			if err != nil {
				return err
			}
			repo, err := git.RepoBaseName(ctx, runner)
			if err != nil {
				return err
			}
			if _, err := git.CurrentBranch(ctx, runner); err != nil {
				return err
			}
			var query string
			if len(args) == 1 {
				query, err = normalizeTaskQuery(args[0])
				if err != nil {
					return err
				}
			}
			if opts.output == "raw" && query == "" && opts.branch == "" {
				return fmt.Errorf("raw output requires a task or branch filter")
			}
			field, err := normalizeListField(opts.field)
			if err != nil {
				return err
			}

			worktrees, err := worktree.List(ctx, runner, repoRoot)
			if err != nil {
				return err
			}
			shortHashLen, err := worktree.ShortHashLength(ctx, runner, repoRoot)
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
					Head:    worktree.ShortHash(wt.Head, shortHashLen),
				}
				if query != "" && !matchesTask(row.Task, query, opts.strict) {
					continue
				}
				if opts.branch != "" && row.Branch != opts.branch {
					continue
				}
				rows = append(rows, row)
			}

			if opts.output == "raw" && len(rows) == 0 {
				fallbackBranch := opts.branch
				if fallbackBranch == "" {
					fallbackBranch = query
				}
				path, ok, err := fallbackPathForBranch(ctx, runner, repoRoot, fallbackBranch)
				if err != nil {
					return err
				}
				if ok {
					if _, err := fmt.Fprintln(cmd.OutOrStdout(), displayPath(repoRoot, path, opts.abs)); err != nil {
						return err
					}
					return nil
				}
			}

			return renderList(cmd, opts.output, field, rows, opts.grid)
		},
	}

	cmd.Flags().StringVarP(&opts.output, "output", "o", opts.output, "output format: table, json, csv, or raw")
	cmd.Flags().StringVar(&opts.branch, "branch", "", "filter by branch name")
	cmd.Flags().StringVarP(&opts.field, "field", "f", "", "raw output field: path, task, or branch (default path)")
	cmd.Flags().BoolVar(&opts.abs, "absolute-path", false, "show absolute paths instead of relative")
	cmd.Flags().BoolVar(&opts.abs, "abs", false, "alias for --absolute-path")
	cmd.Flags().BoolVar(&opts.grid, "grid", false, "render table with grid borders")
	cmd.Flags().BoolVar(&opts.strict, "strict", false, "require exact task match (after trimming and slugifying)")

	return cmd
}

func renderList(cmd *cobra.Command, format, field string, rows []listRow, grid bool) error {
	switch format {
	case "table":
		columns := []tableColumn{
			{Header: "TASK", MinWidth: 6},
			{Header: "BRANCH", MinWidth: 10, Flexible: true, Truncate: true, Style: func(value string) lipgloss.Style { return ui.AccentStyle }},
			{Header: "PATH", MinWidth: 16, Flexible: true, Truncate: true},
			{Header: "PRESENT", MinWidth: 7, Style: func(value string) lipgloss.Style {
				if value == "true" {
					return ui.SuccessStyle
				}
				return ui.ErrorStyle
			}},
			{Header: "HEAD", MinWidth: 7, Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
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
		renderTable(cmd, columns, tableRows, grid)
		return nil
	case "json":
		payload, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(payload)); err != nil {
			return err
		}
		return nil
	case "csv":
		writer := csv.NewWriter(cmd.OutOrStdout())
		if err := writer.Write([]string{"task", "branch", "path", "present", "head"}); err != nil {
			return err
		}
		for _, row := range rows {
			if err := writer.Write([]string{
				row.Task,
				row.Branch,
				row.Path,
				strconv.FormatBool(row.Present),
				row.Head,
			}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	case "raw":
		if len(rows) == 0 {
			return fmt.Errorf("no matching worktrees found")
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), listFieldValue(rows[0], field)); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func normalizeListField(value string) (string, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return "path", nil
	}
	switch value {
	case "path", "task", "branch":
		return value, nil
	default:
		return "", fmt.Errorf("unsupported field: %s (use path, task, or branch)", value)
	}
}

func listFieldValue(row listRow, field string) string {
	switch field {
	case "task":
		return row.Task
	case "branch":
		return row.Branch
	default:
		return row.Path
	}
}
