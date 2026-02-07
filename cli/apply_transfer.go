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
	"github.com/spf13/cobra"
)

func printDryRunPlan(out io.Writer, mode handoffMode, plan transferPlan, preflight transferPreflight, maskPaths bool) error {
	if _, err := fmt.Fprintf(out, "%s plan\n", mode); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  to: %s\n", plan.to); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  source: %s\n", maskPathForDryRun(plan.sourceRoot, maskPaths)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  destination: %s\n", maskPathForDryRun(plan.destinationRoot, maskPaths)); err != nil {
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
	actions := dryRunActions(mode, plan, preflight, maskPaths)
	for idx, action := range actions {
		if _, err := fmt.Fprintf(out, "  %d. %s\n", idx+1, action); err != nil {
			return err
		}
	}
	return nil
}

func dryRunActions(mode handoffMode, plan transferPlan, preflight transferPreflight, maskPaths bool) []string {
	actions := make([]string, 0, len(preflight.untrackedFiles)+4)

	if mode == handoffOverwrite {
		actions = append(actions, "[destructive] "+formatGitCommandForDryRun([]string{"-C", plan.destinationRoot, "reset", "--hard"}, maskPaths))
		actions = append(actions, "[destructive] "+formatGitCommandForDryRun([]string{"-C", plan.destinationRoot, "clean", "-fd"}, maskPaths))
	}

	if preflight.trackedPatch {
		if mode == handoffApply {
			actions = append(actions, formatGitCommandForDryRun([]string{"-C", plan.destinationRoot, "apply", "--check", "<temp-patch>"}, maskPaths))
		}
		actions = append(actions, formatGitCommandForDryRun([]string{"-C", plan.destinationRoot, "apply", "<temp-patch>"}, maskPaths))
	}

	for _, rel := range preflight.untrackedFiles {
		actions = append(actions, fmt.Sprintf(
			"copy %s -> %s",
			maskPathForDryRun(filepath.Join(plan.sourceRoot, rel), maskPaths),
			maskPathForDryRun(filepath.Join(plan.destinationRoot, rel), maskPaths),
		))
	}

	if len(actions) == 0 {
		actions = append(actions, "no tracked or untracked changes detected")
	}

	return actions
}

func transferChanges(ctx context.Context, cmd *cobra.Command, runner git.Runner, sourceRoot, destinationRoot string, dryRun, resetDestination bool) error {
	maskPaths := shouldMaskSensitivePaths(ctx)
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
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: failed to remove temp patch %s: %v\n", maskPathForDryRun(patchFile, maskPaths), err)
		}
	}()

	if patch != "" {
		if !resetDestination {
			if err := runGit(ctx, cmd, dryRun, runner, "-C", destinationRoot, "apply", "--check", patchFile); err != nil {
				return &applyConflictError{reason: "apply patch check failed", err: err}
			}
		}
		if err := runGit(ctx, cmd, dryRun, runner, "-C", destinationRoot, "apply", patchFile); err != nil {
			if resetDestination {
				if fallbackErr := syncTrackedChangesFallback(ctx, runner, sourceRoot, destinationRoot); fallbackErr != nil {
					return fmt.Errorf("overwrite apply patch: %w (fallback sync failed: %v)", err, fallbackErr)
				}
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: overwrite apply patch failed; used tracked-file fallback sync: %v\n", err)
			} else {
				return fmt.Errorf("apply patch: %w", err)
			}
		}
	}

	untracked, err := listUntracked(ctx, runner, sourceRoot)
	if err != nil {
		return err
	}
	for _, rel := range untracked {
		if err := copyFile(sourceRoot, destinationRoot, rel, dryRun, cmd.OutOrStdout(), maskPaths); err != nil {
			return err
		}
	}
	return nil
}

type trackedChange struct {
	status byte
	oldRel string
	newRel string
}

func syncTrackedChangesFallback(ctx context.Context, runner git.Runner, sourceRoot, destinationRoot string) error {
	changes, err := listTrackedChanges(ctx, runner, sourceRoot)
	if err != nil {
		return err
	}
	for _, change := range changes {
		switch change.status {
		case 'D':
			if err := removeTrackedPath(destinationRoot, change.newRel); err != nil {
				return err
			}
		case 'R':
			if err := removeTrackedPath(destinationRoot, change.oldRel); err != nil {
				return err
			}
			if err := copyFile(sourceRoot, destinationRoot, change.newRel, false, io.Discard, false); err != nil {
				return err
			}
		case 'A', 'M', 'T', 'C':
			if err := copyFile(sourceRoot, destinationRoot, change.newRel, false, io.Discard, false); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported tracked change status %q", string(change.status))
		}
	}
	return nil
}

func listTrackedChanges(ctx context.Context, runner git.Runner, repoRoot string) ([]trackedChange, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "diff", "--name-status", "HEAD")
	if err != nil {
		if stderr != "" {
			return nil, fmt.Errorf("git diff --name-status: %w: %s", err, stderr)
		}
		return nil, fmt.Errorf("git diff --name-status: %w", err)
	}
	var out []trackedChange
	for _, line := range strings.Split(stdout, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		change, err := parseTrackedChange(trimmed)
		if err != nil {
			return nil, err
		}
		out = append(out, change)
	}
	return out, nil
}

func parseTrackedChange(line string) (trackedChange, error) {
	parts := strings.Split(line, "\t")
	if len(parts) < 2 {
		return trackedChange{}, fmt.Errorf("invalid tracked change line: %q", line)
	}
	statusText := strings.TrimSpace(parts[0])
	if statusText == "" {
		return trackedChange{}, fmt.Errorf("missing tracked change status: %q", line)
	}
	status := statusText[0]
	switch status {
	case 'R', 'C':
		if len(parts) < 3 {
			return trackedChange{}, fmt.Errorf("invalid rename/copy tracked change line: %q", line)
		}
		return trackedChange{
			status: status,
			oldRel: strings.TrimSpace(parts[1]),
			newRel: strings.TrimSpace(parts[2]),
		}, nil
	default:
		return trackedChange{
			status: status,
			newRel: strings.TrimSpace(parts[1]),
		}, nil
	}
}

func removeTrackedPath(root, rel string) error {
	path := filepath.Join(root, rel)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
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
