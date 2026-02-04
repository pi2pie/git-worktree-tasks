package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type applyOptions struct {
	yes    bool
	dryRun bool
}

type applyConflictError struct {
	reason string
	err    error
}

func (e *applyConflictError) Error() string {
	if e.err == nil {
		return e.reason
	}
	return fmt.Sprintf("%s: %v", e.reason, e.err)
}

func (e *applyConflictError) Unwrap() error { return e.err }

func newApplyCommand() *cobra.Command {
	opts := &applyOptions{}
	cmd := &cobra.Command{
		Use:   "apply <task>",
		Short: "Apply changes between a Codex worktree and the local checkout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, ok := configFromContext(ctx)
			if !ok || cfg.Mode != modeCodex {
				return fmt.Errorf("apply is only supported in --mode=codex")
			}

			runner := defaultRunner()
			repoRoot, err := repoRoot(ctx, runner)
			if err != nil {
				return err
			}
			if _, err := git.CurrentBranch(ctx, runner); err != nil {
				return err
			}

			if !cmd.Flags().Changed("yes") {
				opts.yes = !cfg.Cleanup.Confirm
			}

			opaqueID := strings.TrimSpace(args[0])
			if opaqueID == "" {
				return fmt.Errorf("task query cannot be empty")
			}

			codexHome, err := codexHomeDir()
			if err != nil {
				return err
			}
			codexWorktrees := codexWorktreesRoot(codexHome)

			wtPath, ok, err := resolveCodexWorktreePath(ctx, runner, repoRoot, codexWorktrees, opaqueID)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("no Codex worktree found for %q under %s", opaqueID, filepath.Join("$CODEX_HOME", "worktrees"))
			}

			conflictReasons, err := detectApplyConflicts(ctx, runner, repoRoot, wtPath)
			if err != nil {
				return err
			}
			if len(conflictReasons) > 0 {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", ui.WarningStyle.Render("apply conflict detected:")); err != nil {
					return err
				}
				for _, reason := range conflictReasons {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", ui.WarningStyle.Render(reason)); err != nil {
						return err
					}
				}

				if !opts.yes {
					ok, err := confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), "Overwrite the Codex worktree from the local checkout?")
					if err != nil {
						return err
					}
					if !ok {
						return errCanceled
					}
					ok, err = confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), "This will discard worktree changes. Continue?")
					if err != nil {
						return err
					}
					if !ok {
						return errCanceled
					}
				}

				return overwriteWorktreeChanges(ctx, cmd, runner, repoRoot, wtPath, opts.dryRun)
			}

			if err := applyWorktreeChanges(ctx, cmd, runner, repoRoot, wtPath, opts.dryRun); err != nil {
				var conflictErr *applyConflictError
				if errors.As(err, &conflictErr) {
					if !opts.yes {
						if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", ui.WarningStyle.Render(conflictErr.reason)); err != nil {
							return err
						}
						ok, err := confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), "Overwrite the Codex worktree from the local checkout?")
						if err != nil {
							return err
						}
						if !ok {
							return errCanceled
						}
						ok, err = confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), "This will discard worktree changes. Continue?")
						if err != nil {
							return err
						}
						if !ok {
							return errCanceled
						}
					}
					return overwriteWorktreeChanges(ctx, cmd, runner, repoRoot, wtPath, opts.dryRun)
				}
				return err
			}
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), ui.SuccessStyle.Render("apply complete")); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")
	return cmd
}

func resolveCodexWorktreePath(ctx context.Context, runner git.Runner, repoRoot, codexWorktreesRoot, opaqueID string) (string, bool, error) {
	worktrees, err := worktree.List(ctx, runner, repoRoot)
	if err != nil {
		return "", false, err
	}
	for _, wt := range worktrees {
		wtAbs, err := worktree.NormalizePath(repoRoot, wt.Path)
		if err != nil {
			return "", false, err
		}
		id, _, ok := codexWorktreeInfo(codexWorktreesRoot, wtAbs)
		if !ok || id != opaqueID {
			continue
		}
		return wtAbs, true, nil
	}
	return "", false, nil
}

func detectApplyConflicts(ctx context.Context, runner git.Runner, repoRoot, worktreePath string) ([]string, error) {
	var reasons []string

	dirty, err := isDirty(ctx, runner, repoRoot)
	if err != nil {
		return nil, err
	}
	if dirty {
		reasons = append(reasons, "local checkout has uncommitted changes")
	}

	localModified, err := modifiedFiles(ctx, runner, repoRoot)
	if err != nil {
		return nil, err
	}
	worktreeModified, err := modifiedFiles(ctx, runner, worktreePath)
	if err != nil {
		return nil, err
	}
	if intersects(localModified, worktreeModified) {
		reasons = append(reasons, "both sides modified the same file(s)")
	}

	return reasons, nil
}

func isDirty(ctx context.Context, runner git.Runner, repoRoot string) (bool, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "status", "--porcelain")
	if err != nil {
		if stderr != "" {
			return false, fmt.Errorf("git status: %w: %s", err, stderr)
		}
		return false, fmt.Errorf("git status: %w", err)
	}
	return strings.TrimSpace(stdout) != "", nil
}

func modifiedFiles(ctx context.Context, runner git.Runner, repoRoot string) (map[string]struct{}, error) {
	files := map[string]struct{}{}

	diffNames, stderr, err := runner.Run(ctx, "-C", repoRoot, "diff", "--name-only", "HEAD")
	if err != nil {
		if stderr != "" {
			return nil, fmt.Errorf("git diff --name-only: %w: %s", err, stderr)
		}
		return nil, fmt.Errorf("git diff --name-only: %w", err)
	}
	for _, line := range strings.Split(diffNames, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		files[trimmed] = struct{}{}
	}

	untracked, stderr, err := runner.Run(ctx, "-C", repoRoot, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		if stderr != "" {
			return nil, fmt.Errorf("git ls-files: %w: %s", err, stderr)
		}
		return nil, fmt.Errorf("git ls-files: %w", err)
	}
	for _, line := range strings.Split(untracked, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		files[trimmed] = struct{}{}
	}

	return files, nil
}

func intersects(left, right map[string]struct{}) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}
	if len(left) > len(right) {
		left, right = right, left
	}
	for key := range left {
		if _, ok := right[key]; ok {
			return true
		}
	}
	return false
}

func applyWorktreeChanges(ctx context.Context, cmd *cobra.Command, runner git.Runner, repoRoot, worktreePath string, dryRun bool) error {
	patch, err := gitDiff(ctx, runner, worktreePath)
	if err != nil {
		return err
	}

	patchFile, err := writeTempPatch(patch)
	if err != nil {
		return err
	}
	defer func() {
		if err := removeTempPatch(patchFile); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: failed to remove temp patch %s: %v\n", patchFile, err)
		}
	}()

	if patch != "" {
		if err := runGit(ctx, cmd, dryRun, runner, "-C", repoRoot, "apply", "--check", patchFile); err != nil {
			return &applyConflictError{reason: "apply patch check failed", err: err}
		}
		if err := runGit(ctx, cmd, dryRun, runner, "-C", repoRoot, "apply", patchFile); err != nil {
			return fmt.Errorf("apply patch: %w", err)
		}
	}

	untracked, err := listUntracked(ctx, runner, worktreePath)
	if err != nil {
		return err
	}
	for _, rel := range untracked {
		if err := copyFile(worktreePath, repoRoot, rel, dryRun, cmd.OutOrStdout()); err != nil {
			return err
		}
	}

	return nil
}

func overwriteWorktreeChanges(ctx context.Context, cmd *cobra.Command, runner git.Runner, repoRoot, worktreePath string, dryRun bool) error {
	if err := runGit(ctx, cmd, dryRun, runner, "-C", worktreePath, "reset", "--hard"); err != nil {
		return err
	}
	if err := runGit(ctx, cmd, dryRun, runner, "-C", worktreePath, "clean", "-fd"); err != nil {
		return err
	}

	patch, err := gitDiff(ctx, runner, repoRoot)
	if err != nil {
		return err
	}
	patchFile, err := writeTempPatch(patch)
	if err != nil {
		return err
	}
	defer func() {
		if err := removeTempPatch(patchFile); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: failed to remove temp patch %s: %v\n", patchFile, err)
		}
	}()

	if patch != "" {
		if err := runGit(ctx, cmd, dryRun, runner, "-C", worktreePath, "apply", "--check", patchFile); err != nil {
			return &applyConflictError{reason: "apply patch check failed", err: err}
		}
		if err := runGit(ctx, cmd, dryRun, runner, "-C", worktreePath, "apply", patchFile); err != nil {
			return err
		}
	}

	untracked, err := listUntracked(ctx, runner, repoRoot)
	if err != nil {
		return err
	}
	for _, rel := range untracked {
		if err := copyFile(repoRoot, worktreePath, rel, dryRun, cmd.OutOrStdout()); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), ui.SuccessStyle.Render("overwrite complete")); err != nil {
		return err
	}
	return nil
}

func gitDiff(ctx context.Context, runner git.Runner, repoRoot string) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "diff", "--binary", "HEAD")
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("git diff: %w: %s", err, stderr)
		}
		return "", fmt.Errorf("git diff: %w", err)
	}
	return stdout, nil
}

func listUntracked(ctx context.Context, runner git.Runner, repoRoot string) ([]string, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		if stderr != "" {
			return nil, fmt.Errorf("git ls-files: %w: %s", err, stderr)
		}
		return nil, fmt.Errorf("git ls-files: %w", err)
	}
	var out []string
	for _, line := range strings.Split(stdout, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out, nil
}

func writeTempPatch(contents string) (string, error) {
	tmp, err := os.CreateTemp("", "gwtt-apply-*.patch")
	if err != nil {
		return "", err
	}
	if _, err := io.WriteString(tmp, contents); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			return "", fmt.Errorf("write patch: %w (close error: %v)", err, closeErr)
		}
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func copyFile(srcRoot, dstRoot, rel string, dryRun bool, out io.Writer) (err error) {
	srcPath := filepath.Join(srcRoot, rel)
	dstPath := filepath.Join(dstRoot, rel)

	info, err := os.Lstat(srcPath)
	if err != nil {
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(srcPath)
		if err != nil {
			return err
		}
		if dryRun {
			_, err := fmt.Fprintf(out, "symlink %s -> %s (%s)\n", srcPath, dstPath, target)
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return err
		}
		if err := os.Remove(dstPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return os.Symlink(target, dstPath)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("unsupported file type for copy: %s", srcPath)
	}

	if dryRun {
		_, err := fmt.Fprintf(out, "copy %s -> %s\n", srcPath, dstPath)
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	outFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if _, err := io.Copy(outFile, in); err != nil {
		return err
	}
	return nil
}

func removeTempPatch(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
