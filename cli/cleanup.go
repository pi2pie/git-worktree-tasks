package cli

import (
	"context"
	"fmt"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/spf13/cobra"
)

type cleanupOptions struct {
	removeWorktree bool
	removeBranch   bool
	forceBranch    bool
	yes            bool
	dryRun         bool
}

func newCleanupCommand(state *runState) *cobra.Command {
	opts := &cleanupOptions{removeWorktree: true}
	cmd := &cobra.Command{
		Use:     "cleanup <task>",
		Short:   "Remove a task worktree and/or branch",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"rm"},
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
			branch := task

			if !opts.removeWorktree && !opts.removeBranch {
				return fmt.Errorf("nothing to clean: enable --remove-worktree and/or --remove-branch")
			}

			if !opts.yes {
				ok, err := confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), "Remove worktree/branch?")
				if err != nil {
					return err
				}
				if !ok {
					return errCanceled
				}
			}

			if opts.removeWorktree {
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "worktree", "remove", path); err != nil {
					return err
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "worktree", "prune"); err != nil {
					return err
				}
			}

			if opts.removeBranch {
				deleteFlag := "-d"
				if opts.forceBranch {
					deleteFlag = "-D"
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "branch", deleteFlag, branch); err != nil {
					return err
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), "cleanup complete")
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.removeWorktree, "remove-worktree", opts.removeWorktree, "remove the task worktree")
	cmd.Flags().BoolVar(&opts.removeBranch, "remove-branch", false, "remove the task branch")
	cmd.Flags().BoolVar(&opts.forceBranch, "force-branch", false, "force delete branch when removing")
	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")

	return cmd
}
