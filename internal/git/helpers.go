package git

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotRepo  = errors.New("not a git repository (run inside a git repository)")
	ErrNoCommits = errors.New("no commits yet (empty history)")
)

func RepoRoot(ctx context.Context, runner Runner) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "rev-parse", "--show-toplevel")
	if err != nil {
		if classified := classifyGitStderr(stderr); classified != nil {
			return "", fmt.Errorf("repo root: %w", classified)
		}
		return "", fmt.Errorf("repo root: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

func CommonDir(ctx context.Context, runner Runner) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "rev-parse", "--git-common-dir")
	if err != nil {
		if classified := classifyGitStderr(stderr); classified != nil {
			return "", fmt.Errorf("git common dir: %w", classified)
		}
		return "", fmt.Errorf("git common dir: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

func CurrentBranch(ctx context.Context, runner Runner) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		if classified := classifyGitStderr(stderr); classified != nil {
			return "", fmt.Errorf("current branch: %w", classified)
		}
		return "", fmt.Errorf("current branch: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

func CurrentBranchAt(ctx context.Context, runner Runner, repoRoot string) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		if classified := classifyGitStderr(stderr); classified != nil {
			return "", fmt.Errorf("current branch: %w", classified)
		}
		return "", fmt.Errorf("current branch: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

func classifyGitStderr(stderr string) error {
	lower := strings.ToLower(stderr)
	if strings.Contains(lower, "not a git repository") || strings.Contains(lower, "bad git dir") {
		return ErrNotRepo
	}
	if strings.Contains(lower, "unknown revision") || strings.Contains(lower, "needed a single revision") {
		return ErrNoCommits
	}
	return nil
}
