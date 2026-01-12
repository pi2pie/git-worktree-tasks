package worktree

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/dev-pi2pie/git-worktree-tasks/internal/git"
)

type Worktree struct {
	Path     string
	Head     string
	Branch   string
	Bare     bool
	Locked   bool
	Prunable bool
}

func List(ctx context.Context, runner git.Runner, repoRoot string) ([]Worktree, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "worktree", "list", "--porcelain")
	if err != nil {
		if stderr != "" {
			return nil, fmt.Errorf("list worktrees: %w: %s", err, stderr)
		}
		return nil, fmt.Errorf("list worktrees: %w", err)
	}

	var worktrees []Worktree
	var current *Worktree
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "worktree ") {
			if current != nil {
				worktrees = append(worktrees, *current)
			}
			current = &Worktree{Path: strings.TrimSpace(strings.TrimPrefix(line, "worktree "))}
			continue
		}
		if current == nil {
			continue
		}

		switch {
		case strings.HasPrefix(line, "HEAD "):
			current.Head = strings.TrimSpace(strings.TrimPrefix(line, "HEAD "))
		case strings.HasPrefix(line, "branch "):
			current.Branch = strings.TrimSpace(strings.TrimPrefix(line, "branch "))
		case line == "bare":
			current.Bare = true
		case line == "locked":
			current.Locked = true
		case strings.HasPrefix(line, "locked "):
			current.Locked = true
		case line == "prunable":
			current.Prunable = true
		case strings.HasPrefix(line, "prunable "):
			current.Prunable = true
		}
	}
	if current != nil {
		worktrees = append(worktrees, *current)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parse worktree list: %w", err)
	}
	return worktrees, nil
}
