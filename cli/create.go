package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/dev-pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type createOptions struct {
	base         string
	path         string
	output       string
	copyCd       bool
	dryRun       bool
	skipExisting bool
}

func newCreateCommand(state *runState) *cobra.Command {
	opts := &createOptions{base: "main", output: "text"}
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
			repo, err := repoBaseName(ctx, runner)
			if err != nil {
				return err
			}
			task := worktree.SlugifyTask(args[0])
			path := worktree.WorktreePath(repoRoot, repo, task)
			if opts.path != "" {
				path = worktreePathOverride(repoRoot, opts.path)
			}

			worktreeExists, err := worktree.Exists(ctx, runner, repoRoot, path)
			if err != nil {
				return err
			}
			if worktreeExists {
				if opts.skipExisting {
					return handleExistingWorktree(ctx, cmd, repoRoot, path, task, opts)
				}
				return fmt.Errorf("worktree path already occupied: %s", displayPath(repoRoot, path, false))
			}

			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("worktree path already exists: %s", path)
			}

			branch := task
			branchExists, err := git.BranchExists(ctx, runner, repoRoot, branch)
			if err != nil {
				return err
			}
			gitArgs := buildCreateWorktreeArgs(repoRoot, path, branch, opts.base, branchExists)
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

			display := displayPath(repoRoot, path, false)
			switch opts.output {
			case "text":
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s (branch: %s)\n",
					ui.SuccessStyle.Render("worktree ready"),
					ui.AccentStyle.Render(display),
					ui.AccentStyle.Render(branch),
				)
				if opts.copyCd {
					cdCommand := fmt.Sprintf("cd %s", display)
					if err := copyToClipboard(ctx, cdCommand); err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n",
						ui.SuccessStyle.Render("copied to clipboard"),
						ui.AccentStyle.Render(cdCommand),
					)
				}
			case "raw":
				if opts.copyCd {
					return fmt.Errorf("copy-cd is not supported with --output raw")
				}
				fmt.Fprintln(cmd.OutOrStdout(), display)
			default:
				return fmt.Errorf("unsupported output format: %s", opts.output)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.base, "base", opts.base, "base branch to create from")
	cmd.Flags().StringVarP(&opts.path, "path", "p", "", "override worktree path (relative to repo root or absolute)")
	cmd.Flags().StringVarP(&opts.output, "output", "o", opts.output, "output format: text or raw")
	cmd.Flags().BoolVar(&opts.copyCd, "copy-cd", false, "copy a ready-to-run cd command to the clipboard")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")
	cmd.Flags().BoolVar(&opts.skipExisting, "skip-existing", false, "reuse an existing worktree path if present")
	cmd.Flags().BoolVar(&opts.skipExisting, "skip", false, "alias for --skip-existing")

	return cmd
}

func stringSlice(args []string) string {
	return fmt.Sprintf("%s", args)
}

func worktreePathOverride(repoRoot, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(repoRoot, path)
}

func handleExistingWorktree(ctx context.Context, cmd *cobra.Command, repoRoot, path, task string, opts *createOptions) error {
	display := displayPath(repoRoot, path, false)
	switch opts.output {
	case "text":
		fmt.Fprintf(cmd.OutOrStdout(), "%s: %s (branch: %s)\n",
			ui.WarningStyle.Render("worktree exists"),
			ui.AccentStyle.Render(display),
			ui.AccentStyle.Render(task),
		)
		if opts.copyCd {
			cdCommand := fmt.Sprintf("cd %s", display)
			if err := copyToClipboard(ctx, cdCommand); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n",
				ui.SuccessStyle.Render("copied to clipboard"),
				ui.AccentStyle.Render(cdCommand),
			)
		}
	case "raw":
		if opts.copyCd {
			return fmt.Errorf("copy-cd is not supported with --output raw")
		}
		fmt.Fprintln(cmd.OutOrStdout(), display)
	default:
		return fmt.Errorf("unsupported output format: %s", opts.output)
	}
	return nil
}

func buildCreateWorktreeArgs(repoRoot, path, branch, base string, branchExists bool) []string {
	args := []string{"-C", repoRoot, "worktree", "add"}
	if branchExists {
		return append(args, path, branch)
	}
	return append(args, "-b", branch, path, base)
}
