package worktree

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
)

const (
	shortHashDefaultLen = 7
	shortHashMediumLen  = 8
	shortHashLargeLen   = 10

	// Thresholds are commit counts for 7/8/10 character abbreviations.
	shortHashMediumThreshold = 100000
	shortHashLargeThreshold  = 1000000
)

// ShortHashLength returns a dynamic short-hash length based on repository size.
// Rationale: abbreviated object names must remain unambiguous; larger repositories
// benefit from longer abbreviations. We approximate repo size using the commit
// count and use 7/8/10 as coarse steps. References: gitrevisions(7), git-rev-parse(1).
func ShortHashLength(ctx context.Context, runner git.Runner, repoPath string) (int, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoPath, "rev-list", "--count", "--all")
	if err != nil {
		if stderr != "" {
			return shortHashDefaultLen, fmt.Errorf("short hash length: %w: %s", err, stderr)
		}
		return shortHashDefaultLen, fmt.Errorf("short hash length: %w", err)
	}
	count, err := strconv.Atoi(strings.TrimSpace(stdout))
	if err != nil {
		return shortHashDefaultLen, fmt.Errorf("short hash length: %w", err)
	}
	return shortHashLengthForCommitCount(count), nil
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

func shortHashLengthForCommitCount(count int) int {
	switch {
	case count <= shortHashMediumThreshold:
		return shortHashDefaultLen
	case count <= shortHashLargeThreshold:
		return shortHashMediumLen
	default:
		return shortHashLargeLen
	}
}
