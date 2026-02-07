package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
)

func detectApplyConflicts(ctx context.Context, runner git.Runner, destinationRoot, destinationName, sourceRoot string) ([]string, error) {
	preflight, err := collectTransferPreflight(ctx, runner, sourceRoot, destinationRoot, false)
	if err != nil {
		return nil, err
	}
	return conflictReasonsForApply(preflight, destinationName), nil
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

func intersectCount(left, right map[string]struct{}) int {
	if len(left) == 0 || len(right) == 0 {
		return 0
	}
	if len(left) > len(right) {
		left, right = right, left
	}
	count := 0
	for key := range left {
		if _, ok := right[key]; ok {
			count++
		}
	}
	return count
}

func collectTransferPreflight(ctx context.Context, runner git.Runner, sourceRoot, destinationRoot string, includeTransferState bool) (transferPreflight, error) {
	destinationDirty, err := isDirty(ctx, runner, destinationRoot)
	if err != nil {
		return transferPreflight{}, err
	}

	sourceModified, err := modifiedFiles(ctx, runner, sourceRoot)
	if err != nil {
		return transferPreflight{}, err
	}
	destinationModified, err := modifiedFiles(ctx, runner, destinationRoot)
	if err != nil {
		return transferPreflight{}, err
	}

	preflight := transferPreflight{
		destinationDirty: destinationDirty,
		overlappingFiles: intersectCount(sourceModified, destinationModified),
	}

	if includeTransferState {
		patch, err := gitDiff(ctx, runner, sourceRoot)
		if err != nil {
			return transferPreflight{}, err
		}
		preflight.trackedPatch = patch != ""

		untracked, err := listUntracked(ctx, runner, sourceRoot)
		if err != nil {
			return transferPreflight{}, err
		}
		preflight.untrackedFiles = untracked
	}

	return preflight, nil
}

func conflictReasonsForApply(preflight transferPreflight, destinationName string) []string {
	var reasons []string
	if preflight.destinationDirty {
		reasons = append(reasons, fmt.Sprintf("%s has uncommitted changes", destinationName))
	}
	if preflight.overlappingFiles > 0 {
		reasons = append(reasons, fmt.Sprintf("both sides modified %d overlapping file(s)", preflight.overlappingFiles))
	}
	return reasons
}
