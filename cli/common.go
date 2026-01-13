package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
	"github.com/dev-pi2pie/git-worktree-tasks/internal/worktree"
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

func repoBaseName(ctx context.Context, runner git.Runner) (string, error) {
	root, err := repoRoot(ctx, runner)
	if err != nil {
		return "", err
	}
	commonDir, err := git.CommonDir(ctx, runner)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(root, commonDir)
	}
	commonDir = filepath.Clean(commonDir)
	return filepath.Base(filepath.Dir(commonDir)), nil
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

func normalizeTaskQuery(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("task query cannot be empty")
	}
	return strings.ToLower(worktree.SlugifyTask(trimmed)), nil
}

func matchesTask(task, query string, strict bool) bool {
	task = strings.ToLower(task)
	query = strings.ToLower(query)
	if strict {
		return task == query
	}
	return strings.Contains(task, query)
}
