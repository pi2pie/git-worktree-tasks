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
