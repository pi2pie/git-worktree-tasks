package cli

import (
	"fmt"
	"strings"
)

const (
	modeClassic = "classic"
	modeCodex   = "codex"
)

func normalizeMode(raw string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return modeClassic, nil
	}
	switch value {
	case modeClassic, modeCodex:
		return value, nil
	default:
		return "", fmt.Errorf("unsupported mode %q (use classic or codex)", raw)
	}
}
