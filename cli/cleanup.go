package cli

import (
	"context"
	"fmt"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/pi2pie/git-worktree-tasks/ui"
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

func newCleanupCommand() *cobra.Command {
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
			if _, err := git.CurrentBranch(ctx, runner); err != nil {
				return err
			}
			repo, err := git.RepoBaseName(ctx, runner)
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

			worktrees, err := worktree.List(ctx, runner, repoRoot)
			if err != nil {
				return err
			}

			resolvedPath := path
			worktreeExists := false
			branchRef := "refs/heads/" + branch
			repoRootPath, err := worktree.NormalizePath(repoRoot, repoRoot)
			if err != nil {
				return err
			}
			for _, wt := range worktrees {
				if wt.Branch != branchRef {
					continue
				}
				wtPath, err := worktree.NormalizePath(repoRoot, wt.Path)
				if err != nil {
					return err
				}
				if wtPath == repoRootPath {
					continue
				}
				resolvedPath = wt.Path
				worktreeExists = true
				break
			}

			if !worktreeExists {
				targetPath, err := worktree.NormalizePath(repoRoot, path)
				if err != nil {
					return err
				}
				for _, wt := range worktrees {
					wtPath, err := worktree.NormalizePath(repoRoot, wt.Path)
					if err != nil {
						return err
					}
					if wtPath == targetPath {
						worktreeExists = true
						break
					}
				}
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
						if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
							ui.WarningStyle.Render(fmt.Sprintf("no worktree found for task %q; branch %q exists", task, branch)),
						); err != nil {
							return err
						}
					} else {
						if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
							ui.WarningStyle.Render(fmt.Sprintf("no worktree found for task %q; no branch %q", task, branch)),
						); err != nil {
							return err
						}
					}
				} else {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
						ui.WarningStyle.Render(fmt.Sprintf("no worktree found for task %q", task)),
					); err != nil {
						return err
					}
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
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "worktree", "remove", resolvedPath); err != nil {
					return err
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "worktree", "prune"); err != nil {
					return err
				}
			}

			if opts.removeBranch {
				if !branchExists {
					if !opts.removeWorktree || worktreeExists {
						if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
							ui.WarningStyle.Render(fmt.Sprintf("no branch %q to remove", branch)),
						); err != nil {
							return err
						}
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

			if _, err := fmt.Fprintln(cmd.OutOrStdout(), ui.SuccessStyle.Render("cleanup complete")); err != nil {
				return err
			}
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
