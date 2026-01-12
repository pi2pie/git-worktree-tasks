package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/spf13/cobra"
)

type listOptions struct {
	output string
	task   string
	branch string
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
					Path:    filepath.Clean(wt.Path),
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

	return cmd
}

func renderList(cmd *cobra.Command, format string, rows []listRow) error {
	switch format {
	case "table":
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "TASK\tBRANCH\tPATH\tPRESENT\tHEAD")
		for _, row := range rows {
			fmt.Fprintf(w, "%s\t%s\t%s\t%t\t%s\n", row.Task, row.Branch, row.Path, row.Present, row.Head)
		}
		return w.Flush()
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
