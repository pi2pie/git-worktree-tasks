package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func codexHomeDir() (string, error) {
	if raw, ok := os.LookupEnv("CODEX_HOME"); ok {
		value := strings.TrimSpace(raw)
		if value != "" {
			if strings.HasPrefix(value, "~") {
				home, err := os.UserHomeDir()
				if err != nil || strings.TrimSpace(home) == "" {
					return "", fmt.Errorf("resolve CODEX_HOME: %w", err)
				}
				if value == "~" {
					value = home
				} else if strings.HasPrefix(value, "~/") || strings.HasPrefix(value, `~\`) {
					value = filepath.Join(home, value[2:])
				}
			}
			abs, err := filepath.Abs(value)
			if err != nil {
				return "", fmt.Errorf("resolve CODEX_HOME: %w", err)
			}
			return filepath.Clean(abs), nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "", fmt.Errorf("resolve CODEX_HOME: %w", err)
	}
	return filepath.Join(home, ".codex"), nil
}

func codexWorktreesRoot(codexHome string) string {
	return filepath.Join(codexHome, "worktrees")
}

func codexWorktreeInfo(codexWorktreesRoot, worktreePath string) (opaqueID, relative string, ok bool) {
	if strings.TrimSpace(codexWorktreesRoot) == "" || strings.TrimSpace(worktreePath) == "" {
		return "", "", false
	}
	if !isUnderDir(codexWorktreesRoot, worktreePath) {
		return "", "", false
	}
	rel, err := filepath.Rel(codexWorktreesRoot, worktreePath)
	if err != nil {
		return "", "", false
	}
	rel = filepath.Clean(rel)
	if rel == "." || rel == string(filepath.Separator) {
		return "", "", false
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" || parts[0] == "." {
		return "", "", false
	}
	return parts[0], rel, true
}
