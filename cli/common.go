package cli

import (
	"context"
	"path/filepath"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
)

func defaultRunner() git.Runner {
	return git.ExecRunner{}
}

func repoRoot(ctx context.Context, runner git.Runner) (string, error) {
	return git.RepoRoot(ctx, runner)
}

func repoName(root string) string {
	return filepath.Base(root)
}
