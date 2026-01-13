package cli

import (
	"context"
	"fmt"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/dev-pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type cleanupOptions struct {
	removeWorktree bool
	removeBranch   bool
	forceBranch    bool
	worktreeOnly   bool
	yes            bool
	dryRun         bool
}

func newCleanupCommand(state *runState) *cobra.Command {
	opts := &cleanupOptions{removeWorktree: true, removeBranch: true}
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
			repo, err := repoBaseName(ctx, runner)
			if err != nil {
				return err
			}
			task := worktree.SlugifyTask(args[0])
			path := worktree.WorktreePath(repoRoot, repo, task)
			branch := task

			if opts.worktreeOnly {
				opts.removeWorktree = true
				opts.removeBranch = false
			}

			if !opts.removeWorktree && !opts.removeBranch {
				return fmt.Errorf("nothing to clean: enable --remove-worktree and/or --remove-branch")
			}

			worktreeExists, err := worktree.Exists(ctx, runner, repoRoot, path)
			if err != nil {
				return err
			}

			branchExists := false
			if opts.removeBranch {
				branchExists, err = git.BranchExists(ctx, runner, repoRoot, branch)
				if err != nil {
					return err
				}
			}

			if opts.removeWorktree && !worktreeExists {
				if opts.removeBranch {
					if branchExists {
						fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
							ui.WarningStyle.Render(fmt.Sprintf("no worktree found for task %q; branch %q exists", task, branch)),
						)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
							ui.WarningStyle.Render(fmt.Sprintf("no worktree found for task %q; no branch %q", task, branch)),
						)
					}
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
						ui.WarningStyle.Render(fmt.Sprintf("no worktree found for task %q", task)),
					)
				}
			}

			if opts.removeWorktree && worktreeExists {
				if !opts.yes {
					ok, err := confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), "Remove worktree?")
					if err != nil {
						return err
					}
					if !ok {
						return errCanceled
					}
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "worktree", "remove", path); err != nil {
					return err
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "worktree", "prune"); err != nil {
					return err
				}
			}

			if opts.removeBranch {
				if !branchExists {
					if !(opts.removeWorktree && !worktreeExists) {
						fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
							ui.WarningStyle.Render(fmt.Sprintf("no branch %q to remove", branch)),
						)
					}
				} else {
					if !opts.yes {
						message := fmt.Sprintf("Remove branch %q?", branch)
						if !worktreeExists && opts.removeWorktree {
							message = fmt.Sprintf("No worktree found for task %q. Remove branch %q anyway?", task, branch)
						}
						ok, err := confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), message)
						if err != nil {
							return err
						}
						if !ok {
							return errCanceled
						}
					}
					deleteFlag := "-d"
					if opts.forceBranch {
						deleteFlag = "-D"
					}
					if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "branch", deleteFlag, branch); err != nil {
						return err
					}
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), ui.SuccessStyle.Render("cleanup complete"))
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.removeWorktree, "remove-worktree", opts.removeWorktree, "remove the task worktree")
	cmd.Flags().BoolVar(&opts.removeBranch, "remove-branch", opts.removeBranch, "remove the task branch")
	cmd.Flags().BoolVar(&opts.worktreeOnly, "worktree-only", false, "remove only the task worktree (keep branch)")
	cmd.Flags().BoolVar(&opts.forceBranch, "force-branch", false, "force delete branch when removing")
	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")

	return cmd
}
