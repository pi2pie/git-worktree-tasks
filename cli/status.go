package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/spf13/cobra"
)

type statusOptions struct {
	output string
	target string
	task   string
	branch string
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
					Path:       filepath.Clean(wt.Path),
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

	return cmd
}

func renderStatus(cmd *cobra.Command, format string, rows []statusRow) error {
	switch format {
	case "table":
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "TASK\tBRANCH\tPATH\tBASE\tTARGET\tLAST_COMMIT\tDIRTY\tAHEAD\tBEHIND")
		for _, row := range rows {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%t\t%d\t%d\n",
				row.Task, row.Branch, row.Path, row.Base, row.Target, row.LastCommit, row.Dirty, row.Ahead, row.Behind)
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
