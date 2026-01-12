package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/spf13/cobra"
)

type createOptions struct {
	base   string
	dryRun bool
}

func newCreateCommand(state *runState) *cobra.Command {
	opts := &createOptions{base: "main"}
	cmd := &cobra.Command{
		Use:   "create <task>",
		Short: "Create a worktree and branch for a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			runner := defaultRunner()

			repoRoot, err := repoRoot(ctx, runner)
			if err != nil {
				return err
			}
			repo := repoName(repoRoot)
			task := worktree.SlugifyTask(args[0])
			path := worktree.WorktreePath(repoRoot, repo, task)

			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("worktree path already exists: %s", path)
			}

			branch := task
			gitArgs := []string{"-C", repoRoot, "worktree", "add", "-b", branch, path, opts.base}
			if opts.dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), "git", stringSlice(gitArgs))
				return nil
			}

			_, stderr, err := runner.Run(ctx, gitArgs...)
			if err != nil {
				if stderr != "" {
					return fmt.Errorf("create worktree: %w: %s", err, stderr)
				}
				return fmt.Errorf("create worktree: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "worktree ready: %s (branch: %s)\n", path, branch)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.base, "base", opts.base, "base branch to create from")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")

	return cmd
}

func stringSlice(args []string) string {
	return fmt.Sprintf("%s", args)
}
