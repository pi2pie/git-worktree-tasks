package worktree

import (
	"context"
	"fmt"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
)

const shortHashDefaultLen = 7

// ShortHashLength returns the short-hash length Git would use for this repository.
// It queries Git directly via `rev-parse --short HEAD` to get the actual abbreviated
// length Git computes based on core.abbrev (auto by default). This respects both
// user configuration and Git's collision-avoidance algorithm.
func ShortHashLength(ctx context.Context, runner git.Runner, repoPath string) (int, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoPath, "rev-parse", "--short", "HEAD")
	if err != nil {
		if stderr != "" {
			return shortHashDefaultLen, fmt.Errorf("short hash length: %w: %s", err, stderr)
		}
		return shortHashDefaultLen, fmt.Errorf("short hash length: %w", err)
	}
	shortHash := strings.TrimSpace(stdout)
	if len(shortHash) == 0 {
		return shortHashDefaultLen, nil
	}
	return len(shortHash), nil
}

func ShortHash(hash string, length int) string {
	if length <= 0 {
		length = shortHashDefaultLen
	}
	if len(hash) <= length {
		return hash
	}
	return hash[:length]
}
