package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type createOptions struct {
	base         string
	path         string
	output       string
	dryRun       bool
	skipExisting bool
}

func newCreateCommand() *cobra.Command {
	opts := &createOptions{output: "text"}
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
			currentBranch, err := git.CurrentBranchAt(ctx, runner, repoRoot)
			if err != nil {
				return err
			}
			base, err := resolveCreateBase(currentBranch, opts.base)
			if err != nil {
				return err
			}
			repo, err := git.RepoBaseName(ctx, runner)
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
					branch, err := existingWorktreeBranch(ctx, runner, repoRoot, path, task)
					if err != nil {
						return err
					}
					return handleExistingWorktree(cmd, repoRoot, path, branch, opts)
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
			gitArgs := buildCreateWorktreeArgs(repoRoot, path, branch, base, branchExists)
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
			case "raw":
				fmt.Fprintln(cmd.OutOrStdout(), display)
			default:
				return fmt.Errorf("unsupported output format: %s", opts.output)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.base, "base", opts.base, "base branch to create from (default: current branch)")
	cmd.Flags().StringVarP(&opts.path, "path", "p", "", "override worktree path (relative to repo root or absolute)")
	cmd.Flags().StringVarP(&opts.output, "output", "o", opts.output, "output format: text or raw")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")
	cmd.Flags().BoolVar(&opts.skipExisting, "skip-existing", false, "reuse an existing worktree path if present")
	cmd.Flags().BoolVar(&opts.skipExisting, "skip", false, "alias for --skip-existing")

	return cmd
}

func stringSlice(args []string) string {
	return fmt.Sprintf("%s", args)
}

func resolveCreateBase(currentBranch, override string) (string, error) {
	if override != "" {
		return override, nil
	}
	if currentBranch == "HEAD" {
		return "", fmt.Errorf("detached HEAD: specify --base to create from a branch")
	}
	return currentBranch, nil
}

func worktreePathOverride(repoRoot, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(repoRoot, path)
}

func existingWorktreeBranch(ctx context.Context, runner git.Runner, repoRoot, path, fallback string) (string, error) {
	wt, ok, err := worktree.LookupByPath(ctx, runner, repoRoot, path)
	if err != nil {
		return "", err
	}
	if !ok {
		return fallback, nil
	}
	if wt.Branch != "" {
		return strings.TrimPrefix(wt.Branch, "refs/heads/"), nil
	}
	return "detached", nil
}

func handleExistingWorktree(cmd *cobra.Command, repoRoot, path, branch string, opts *createOptions) error {
	display := displayPath(repoRoot, path, false)
	switch opts.output {
	case "text":
		fmt.Fprintf(cmd.OutOrStdout(), "%s: %s (branch: %s)\n",
			ui.WarningStyle.Render("worktree exists"),
			ui.AccentStyle.Render(display),
			ui.AccentStyle.Render(branch),
		)
	case "raw":
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
