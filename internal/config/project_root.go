package config

import (
	"fmt"
	"os"
	"path/filepath"
)

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
