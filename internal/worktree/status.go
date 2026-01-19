package worktree

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
)

type StatusInfo struct {
	Dirty      bool
	LastCommit string
	Ahead      int
	Behind     int
	Base       string
}

func Status(ctx context.Context, runner git.Runner, path string, target string) (StatusInfo, error) {
	var info StatusInfo

	stdout, stderr, err := runner.Run(ctx, "-C", path, "status", "--porcelain")
	if err != nil {
		return info, fmt.Errorf("status dirty check: %w: %s", err, stderr)
	}
	info.Dirty = strings.TrimSpace(stdout) != ""

	shortHashLen, err := ShortHashLength(ctx, runner, path)
	if err != nil {
		return info, err
	}

	stdout, stderr, err = runner.Run(ctx, "-C", path, "log", "-1", "--pretty=format:%H %s")
	if err != nil {
		return info, fmt.Errorf("status last commit: %w: %s", err, stderr)
	}
	info.LastCommit = formatCommitLine(stdout, shortHashLen)

	if target != "" {
		stdout, stderr, err = runner.Run(ctx, "-C", path, "merge-base", "HEAD", target)
		if err == nil {
			info.Base = ShortHash(strings.TrimSpace(stdout), shortHashLen)
		}

		stdout, stderr, err = runner.Run(ctx, "-C", path, "rev-list", "--left-right", "--count", target+"...HEAD")
		if err != nil {
			return info, fmt.Errorf("status ahead/behind: %w: %s", err, stderr)
		}
		parts := strings.Fields(stdout)
		if len(parts) >= 2 {
			behind, _ := strconv.Atoi(parts[0])
			ahead, _ := strconv.Atoi(parts[1])
			info.Behind = behind
			info.Ahead = ahead
		}
	}

	return info, nil
}

func formatCommitLine(line string, shortHashLen int) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}
	parts := strings.SplitN(line, " ", 2)
	if len(parts) == 1 {
		return ShortHash(parts[0], shortHashLen)
	}
	return ShortHash(parts[0], shortHashLen) + " " + parts[1]
}
