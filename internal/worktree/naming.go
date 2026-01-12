package worktree

import (
	"path/filepath"
	"regexp"
	"strings"
)

var invalidBranchChar = regexp.MustCompile(`[^A-Za-z0-9_/-]+`)

// SlugifyTask converts a task name into a safe branch name.
func SlugifyTask(task string) string {
	slug := invalidBranchChar.ReplaceAllString(task, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "task"
	}
	return slug
}

// WorktreePath returns the default worktree path: ../<repo>_<task>.
func WorktreePath(repoRoot, repoName, task string) string {
	parent := filepath.Dir(repoRoot)
	return filepath.Join(parent, repoName+"_"+task)
}

// TaskFromPath derives the task name from a worktree path using the fixed naming convention.
func TaskFromPath(repoName, path string) (string, bool) {
	base := filepath.Base(path)
	prefix := repoName + "_"
	if !strings.HasPrefix(base, prefix) {
		return "", false
	}
	return strings.TrimPrefix(base, prefix), true
}
