package cli

import (
	"context"
	"fmt"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type finishOptions struct {
	target         string
	removeWorktree bool
	removeBranch   bool
	cleanup        bool
	forceBranch    bool
	noFF           bool
	squash         bool
	rebase         bool
	yes            bool
	dryRun         bool
}

func newFinishCommand() *cobra.Command {
	opts := &finishOptions{}
	cmd := &cobra.Command{
		Use:   "finish <task>",
		Short: "Merge a task branch into a target branch",
		Args:  cobra.ExactArgs(1),
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

			if opts.squash && opts.rebase {
				return fmt.Errorf("--squash and --rebase cannot be used together")
			}
			if opts.noFF && opts.rebase {
				return fmt.Errorf("--no-ff and --rebase cannot be used together")
			}

			task := worktree.SlugifyTask(args[0])
			branch := task
			path := worktree.WorktreePath(repoRoot, repo, task)

			target := opts.target
			if target == "" {
				current, err := git.CurrentBranch(ctx, runner)
				if err != nil {
					return err
				}
				target = current
			}

			if opts.cleanup {
				opts.removeBranch = true
				opts.removeWorktree = true
			}

			if opts.rebase {
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "checkout", branch); err != nil {
					return err
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "rebase", target); err != nil {
					return err
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "checkout", target); err != nil {
					return err
				}
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "merge", "--ff-only", branch); err != nil {
					return err
				}
			} else {
				if err := runGit(cmd, opts.dryRun, runner, "-C", repoRoot, "checkout", target); err != nil {
					return err
				}

				mergeArgs := []string{"-C", repoRoot, "merge"}
				if opts.noFF {
					mergeArgs = append(mergeArgs, "--no-ff")
				}
				if opts.squash {
					mergeArgs = append(mergeArgs, "--squash")
				}
				mergeArgs = append(mergeArgs, branch)
				if err := runGit(cmd, opts.dryRun, runner, mergeArgs...); err != nil {
					return err
				}
			}

			if opts.removeWorktree || opts.removeBranch {
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
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s\n",
				ui.SuccessStyle.Render("merged"),
				ui.AccentStyle.Render(branch),
				ui.MutedStyle.Render(fmt.Sprintf("into %s", target)),
			)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.target, "target", "", "target branch (default: current)")
	cmd.Flags().BoolVar(&opts.cleanup, "cleanup", false, "remove worktree and branch after merge")
	cmd.Flags().BoolVar(&opts.removeWorktree, "remove-worktree", false, "remove the task worktree after merge")
	cmd.Flags().BoolVar(&opts.removeBranch, "remove-branch", false, "remove the task branch after merge")
	cmd.Flags().BoolVar(&opts.forceBranch, "force-branch", false, "force delete branch when removing")
	cmd.Flags().BoolVar(&opts.noFF, "no-ff", false, "use --no-ff merge")
	cmd.Flags().BoolVar(&opts.squash, "squash", false, "use --squash merge")
	cmd.Flags().BoolVar(&opts.rebase, "rebase", false, "rebase task branch onto target before merging")
	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")

	return cmd
}

func runGit(cmd *cobra.Command, dryRun bool, runner git.Runner, args ...string) error {
	if dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "git", stringSlice(args))
		return nil
	}
	_, stderr, err := runner.Run(context.Background(), args...)
	if err != nil {
		if stderr != "" {
			return fmt.Errorf("git %s: %w: %s", stringSlice(args), err, stderr)
		}
		return fmt.Errorf("git %s: %w", stringSlice(args), err)
	}
	return nil
}
