package git

import (
	"context"
	"fmt"
	"strings"
)

func BranchExists(ctx context.Context, runner Runner, repoRoot, branch string) (bool, error) {
	stdout, stderr, err := runner.Run(ctx, "-C", repoRoot, "branch", "--list", branch)
	if err != nil {
		if stderr != "" {
			return false, fmt.Errorf("branch exists: %w: %s", err, stderr)
		}
		return false, fmt.Errorf("branch exists: %w", err)
	}
	return strings.TrimSpace(stdout) != "", nil
}
