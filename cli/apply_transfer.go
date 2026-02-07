package cli

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/spf13/cobra"
)

func printDryRunPlan(out io.Writer, mode handoffMode, plan transferPlan, preflight transferPreflight) error {
	if _, err := fmt.Fprintf(out, "%s plan\n", mode); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  to: %s\n", plan.to); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  source: %s\n", plan.sourceRoot); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  destination: %s\n", plan.destinationRoot); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  overwrite: %t\n", mode == handoffOverwrite); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(out, ""); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, "preflight"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  destination_dirty: %t\n", preflight.destinationDirty); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  overlapping_files: %d\n", preflight.overlappingFiles); err != nil {
		return err
	}
	trackedPatch := "none"
	if preflight.trackedPatch {
		trackedPatch = "present"
	}
	if _, err := fmt.Fprintf(out, "  tracked_patch: %s\n", trackedPatch); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  untracked_files: %d\n", len(preflight.untrackedFiles)); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(out, ""); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, "actions"); err != nil {
		return err
	}
	actions := dryRunActions(mode, plan, preflight)
	for idx, action := range actions {
		if _, err := fmt.Fprintf(out, "  %d. %s\n", idx+1, action); err != nil {
			return err
		}
	}
	return nil
}

func dryRunActions(mode handoffMode, plan transferPlan, preflight transferPreflight) []string {
	actions := make([]string, 0, len(preflight.untrackedFiles)+4)

	if mode == handoffOverwrite {
		actions = append(actions, "[destructive] "+formatGitCommand([]string{"-C", plan.destinationRoot, "reset", "--hard"}))
		actions = append(actions, "[destructive] "+formatGitCommand([]string{"-C", plan.destinationRoot, "clean", "-fd"}))
	}

	if preflight.trackedPatch {
		actions = append(actions, formatGitCommand([]string{"-C", plan.destinationRoot, "apply", "--check", "<temp-patch>"}))
		actions = append(actions, formatGitCommand([]string{"-C", plan.destinationRoot, "apply", "<temp-patch>"}))
	}

	for _, rel := range preflight.untrackedFiles {
		actions = append(actions, fmt.Sprintf("copy %s -> %s", filepath.Join(plan.sourceRoot, rel), filepath.Join(plan.destinationRoot, rel)))
	}

	if len(actions) == 0 {
		actions = append(actions, "no tracked or untracked changes detected")
	}

	return actions
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
