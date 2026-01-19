package worktree

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
)

func Exists(ctx context.Context, runner git.Runner, repoRoot, path string) (bool, error) {
	worktrees, err := List(ctx, runner, repoRoot)
	if err != nil {
		return false, err
	}
	targetPath, err := normalizePath(repoRoot, path)
	if err != nil {
		return false, err
	}
	for _, wt := range worktrees {
		wtPath, err := normalizePath(repoRoot, wt.Path)
		if err != nil {
			return false, err
		}
		if wtPath == targetPath {
			return true, nil
		}
	}
	return false, nil
}

func LookupByPath(ctx context.Context, runner git.Runner, repoRoot, path string) (*Worktree, bool, error) {
	worktrees, err := List(ctx, runner, repoRoot)
	if err != nil {
		return nil, false, err
	}
	targetPath, err := normalizePath(repoRoot, path)
	if err != nil {
		return nil, false, err
	}
	for _, wt := range worktrees {
		wtPath, err := normalizePath(repoRoot, wt.Path)
		if err != nil {
			return nil, false, err
		}
		if wtPath == targetPath {
			found := wt
			return &found, true, nil
		}
	}
	return nil, false, nil
}

func normalizePath(repoRoot, path string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(repoRoot, path)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("normalize path: %w", err)
	}
	return filepath.Clean(absPath), nil
}
