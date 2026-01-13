package git

import (
	"context"
	"fmt"
	"strings"
)

func RepoRoot(ctx context.Context, runner Runner) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("repo root: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

func CommonDir(ctx context.Context, runner Runner) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("git common dir: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

func CurrentBranch(ctx context.Context, runner Runner) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("current branch: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

func CurrentBranchAt(ctx context.Context, runner Runner, repoRoot string) (string, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("current branch: %w: %s", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}
