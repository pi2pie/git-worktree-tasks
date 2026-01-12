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

func displayPath(repoRoot, path string, absolute bool) string {
	clean := filepath.Clean(path)
	absPath := clean
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(repoRoot, absPath)
	}
	absPath = filepath.Clean(absPath)
	if absolute {
		return absPath
	}
	rel, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		return clean
	}
	return rel
}
