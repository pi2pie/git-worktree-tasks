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
	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

const (
	transferToLocal    = "local"
	transferToWorktree = "worktree"
)

type handoffMode string

const (
	handoffApply     handoffMode = "apply"
	handoffOverwrite handoffMode = "overwrite"
)

type applyOptions struct {
	yes    bool
	dryRun bool
	to     string
	force  bool
}

type handoffOptions struct {
	yes    bool
	dryRun bool
	to     string
}

type transferPlan struct {
	to              string
	sourceRoot      string
	sourceName      string
	destinationRoot string
	destinationName string
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
		Short: "Apply non-destructive changes between a Codex worktree and local checkout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := handoffApply
			if opts.force {
				mode = handoffOverwrite
			}
			return runCodexHandoff(cmd, strings.TrimSpace(args[0]), handoffOptions{
				yes:    opts.yes,
				dryRun: opts.dryRun,
				to:     opts.to,
			}, mode)
		},
	}

	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")
	cmd.Flags().StringVar(&opts.to, "to", transferToLocal, "transfer destination: local or worktree")
	cmd.Flags().BoolVar(&opts.force, "force", false, "compatibility alias for overwrite behavior")
	return cmd
}

func newOverwriteCommand() *cobra.Command {
	opts := &handoffOptions{}
	cmd := &cobra.Command{
		Use:   "overwrite <task>",
		Short: "Overwrite destination with source changes in codex mode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCodexHandoff(cmd, strings.TrimSpace(args[0]), *opts, handoffOverwrite)
		},
	}

	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")
	cmd.Flags().StringVar(&opts.to, "to", transferToLocal, "transfer destination: local or worktree")
	return cmd
}

func runCodexHandoff(cmd *cobra.Command, opaqueID string, opts handoffOptions, mode handoffMode) error {
	ctx := cmd.Context()
	modeCtx, err := resolveModeContext(cmd, false)
	if err != nil {
		return err
	}
	cfg, ok := configFromContext(ctx)
	if !ok || modeCtx.mode != modeCodex {
		return fmt.Errorf("%s is only supported in --mode=codex", mode)
	}

	if !cmd.Flags().Changed("yes") {
		opts.yes = !cfg.Cleanup.Confirm
	}

	if opaqueID == "" {
		return fmt.Errorf("task query cannot be empty")
	}

	runner := defaultRunner()
	repoRoot, err := repoRoot(ctx, runner)
	if err != nil {
		return err
	}
	if _, err := git.CurrentBranch(ctx, runner); err != nil {
		return err
	}

	wtPath, found, err := resolveCodexWorktreePath(ctx, runner, repoRoot, modeCtx.codexWorktrees, opaqueID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("no Codex worktree found for %q under %s", opaqueID, filepath.Join("$CODEX_HOME", "worktrees"))
	}

	plan, err := resolveTransferPlan(repoRoot, wtPath, opts.to)
	if err != nil {
		return err
	}

	if mode == handoffApply {
		reasons, err := detectApplyConflicts(ctx, runner, plan.destinationRoot, plan.destinationName, plan.sourceRoot)
		if err != nil {
			return err
		}
		if len(reasons) > 0 {
			if err := printConflictReasons(cmd.OutOrStdout(), reasons); err != nil {
				return err
			}
			if err := printOverwriteHint(cmd.OutOrStdout(), plan.to, opaqueID); err != nil {
				return err
			}
			return fmt.Errorf("apply aborted due to conflicts")
		}
	}

	if mode == handoffOverwrite {
		if err := confirmOverwrite(cmd, opts.yes, plan); err != nil {
			return err
		}
	}

	err = transferChanges(ctx, cmd, runner, plan.sourceRoot, plan.destinationRoot, opts.dryRun, mode == handoffOverwrite)
	if err != nil {
		if mode == handoffApply {
			var conflictErr *applyConflictError
			if errors.As(err, &conflictErr) {
				if err := printConflictReasons(cmd.OutOrStdout(), []string{conflictErr.reason}); err != nil {
					return err
				}
				if err := printOverwriteHint(cmd.OutOrStdout(), plan.to, opaqueID); err != nil {
					return err
				}
				return fmt.Errorf("apply aborted due to conflicts")
			}
		}
		return err
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), ui.SuccessStyle.Render(fmt.Sprintf("%s complete", mode))); err != nil {
		return err
	}
	return nil
}

func resolveTransferPlan(repoRoot, worktreePath, to string) (transferPlan, error) {
	switch strings.TrimSpace(to) {
	case transferToLocal:
		return transferPlan{
			to:              transferToLocal,
			sourceRoot:      worktreePath,
			sourceName:      "Codex worktree",
			destinationRoot: repoRoot,
			destinationName: "local checkout",
		}, nil
	case transferToWorktree:
		return transferPlan{
			to:              transferToWorktree,
			sourceRoot:      repoRoot,
			sourceName:      "local checkout",
			destinationRoot: worktreePath,
			destinationName: "Codex worktree",
		}, nil
	default:
		return transferPlan{}, fmt.Errorf("invalid --to value %q (expected local or worktree)", to)
	}
}

func printConflictReasons(out io.Writer, reasons []string) error {
	if _, err := fmt.Fprintf(out, "%s\n", ui.WarningStyle.Render("apply conflict detected:")); err != nil {
		return err
	}
	for _, reason := range reasons {
		if _, err := fmt.Fprintf(out, "- %s\n", ui.WarningStyle.Render(reason)); err != nil {
			return err
		}
	}
	return nil
}

func printOverwriteHint(out io.Writer, to, opaqueID string) error {
	_, err := fmt.Fprintf(out, "%s\n", ui.WarningStyle.Render(fmt.Sprintf("rerun with overwrite: gwtt overwrite --to %s %s", to, opaqueID)))
	return err
}

func confirmOverwrite(cmd *cobra.Command, yes bool, plan transferPlan) error {
	if yes {
		return nil
	}
	ok, err := confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), fmt.Sprintf("Overwrite the %s from the %s?", plan.destinationName, plan.sourceName))
	if err != nil {
		return err
	}
	if !ok {
		return errCanceled
	}
	ok, err = confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), fmt.Sprintf("This will discard %s changes. Continue?", plan.destinationName))
	if err != nil {
		return err
	}
	if !ok {
		return errCanceled
	}
	return nil
}

func detectApplyConflicts(ctx context.Context, runner git.Runner, destinationRoot, destinationName, sourceRoot string) ([]string, error) {
	var reasons []string

	dirty, err := isDirty(ctx, runner, destinationRoot)
	if err != nil {
		return nil, err
	}
	if dirty {
		reasons = append(reasons, fmt.Sprintf("%s has uncommitted changes", destinationName))
	}

	sourceModified, err := modifiedFiles(ctx, runner, sourceRoot)
	if err != nil {
		return nil, err
	}
	destinationModified, err := modifiedFiles(ctx, runner, destinationRoot)
	if err != nil {
		return nil, err
	}
	if intersects(sourceModified, destinationModified) {
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

func transferChanges(ctx context.Context, cmd *cobra.Command, runner git.Runner, sourceRoot, destinationRoot string, dryRun, resetDestination bool) error {
	if resetDestination {
		if err := runGit(ctx, cmd, dryRun, runner, "-C", destinationRoot, "reset", "--hard"); err != nil {
			return err
		}
		if err := runGit(ctx, cmd, dryRun, runner, "-C", destinationRoot, "clean", "-fd"); err != nil {
			return err
		}
	}

	patch, err := gitDiff(ctx, runner, sourceRoot)
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
		if err := runGit(ctx, cmd, dryRun, runner, "-C", destinationRoot, "apply", "--check", patchFile); err != nil {
			return &applyConflictError{reason: "apply patch check failed", err: err}
		}
		if err := runGit(ctx, cmd, dryRun, runner, "-C", destinationRoot, "apply", patchFile); err != nil {
			return fmt.Errorf("apply patch: %w", err)
		}
	}

	untracked, err := listUntracked(ctx, runner, sourceRoot)
	if err != nil {
		return err
	}
	for _, rel := range untracked {
		if err := copyFile(sourceRoot, destinationRoot, rel, dryRun, cmd.OutOrStdout()); err != nil {
			return err
		}
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
