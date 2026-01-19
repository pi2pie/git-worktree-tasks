package worktree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

	ok, err := isWorktreePath(path)
	if err != nil {
		return info, err
	}
	if !ok {
		return info, nil
	}

	stdout, stderr, err := runner.Run(ctx, "-C", path, "status", "--porcelain")
	if err != nil {
		if isNotRepoError(stderr) {
			return info, nil
		}
		return info, fmt.Errorf("status dirty check: %w: %s", err, stderr)
	}
	info.Dirty = strings.TrimSpace(stdout) != ""

	hasHead, err := headExists(ctx, runner, path)
	if err != nil {
		return info, err
	}

	if hasHead {
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
	} else {
		info.LastCommit = "empty history"
		if target != "" {
			info.Base = "empty history"
		}
	}

	return info, nil
}

func isWorktreePath(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat worktree path: %w", err)
	}
	if !info.IsDir() {
		return false, nil
	}
	_, err = os.Stat(filepath.Join(path, ".git"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat worktree gitdir: %w", err)
	}
	return true, nil
}

func headExists(ctx context.Context, runner git.Runner, path string) (bool, error) {
	_, stderr, err := runner.Run(ctx, "-C", path, "rev-parse", "--verify", "HEAD")
	if err == nil {
		return true, nil
	}
	if isNotRepoError(stderr) {
		return false, nil
	}
	trimmed := strings.TrimSpace(stderr)
	if trimmed == "" {
		return false, nil
	}
	lower := strings.ToLower(trimmed)
	if strings.Contains(lower, "unknown revision") || strings.Contains(lower, "needed a single revision") {
		return false, nil
	}
	return false, fmt.Errorf("status head check: %w: %s", err, stderr)
}

func isNotRepoError(stderr string) bool {
	lower := strings.ToLower(stderr)
	return strings.Contains(lower, "not a git repository") || strings.Contains(lower, "bad git dir")
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
