package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// projectConfigRoot returns the root directory for project config resolution.
// It walks up from cwd looking for a .git marker (directory or file for worktrees/submodules).
// Falls back to cwd if no .git is found.
func projectConfigRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	root, ok, err := findRepoRoot(cwd)
	if err != nil {
		return "", err
	}
	if ok {
		return root, nil
	}
	return cwd, nil
}

// findRepoRoot walks up from start looking for a .git marker.
// Returns (root, true, nil) if found, ("", false, nil) if not found.
func findRepoRoot(start string) (string, bool, error) {
	dir := filepath.Clean(start)
	for {
		if ok, err := hasGitDir(dir); err != nil {
			return "", false, err
		} else if ok {
			return dir, true, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}

// hasGitDir checks if dir contains a .git entry (directory or file).
// Git worktrees and submodules use a .git file containing "gitdir: /path/to/...".
func hasGitDir(dir string) (bool, error) {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
