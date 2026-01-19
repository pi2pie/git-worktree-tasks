package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
)

func defaultRunner() git.Runner {
	return git.ExecRunner{}
}

func repoRoot(ctx context.Context, runner git.Runner) (string, error) {
	return git.RepoRoot(ctx, runner)
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

func mainWorktreePathFromCommonDir(repoRoot, commonDir string) string {
	if commonDir == "" {
		return repoRoot
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(repoRoot, commonDir)
	}
	commonDir = filepath.Clean(commonDir)
	if filepath.Base(commonDir) == ".git" {
		return filepath.Dir(commonDir)
	}
	return repoRoot
}

func mainWorktreePath(ctx context.Context, runner git.Runner, repoRoot string) (string, error) {
	commonDir, err := git.CommonDir(ctx, runner)
	if err != nil {
		return "", err
	}
	return mainWorktreePathFromCommonDir(repoRoot, commonDir), nil
}

func fallbackPathForBranch(ctx context.Context, runner git.Runner, repoRoot, branch string) (string, bool, error) {
	if branch == "" {
		return "", false, nil
	}
	exists, err := git.BranchExists(ctx, runner, repoRoot, branch)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	path, err := mainWorktreePath(ctx, runner, repoRoot)
	if err != nil {
		return "", false, err
	}
	return path, true, nil
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
