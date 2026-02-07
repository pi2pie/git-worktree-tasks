package cli

import (
	"context"
	"fmt"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/spf13/cobra"
)

func runGit(ctx context.Context, cmd *cobra.Command, dryRun bool, runner git.Runner, args ...string) error {
	if dryRun {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), formatGitCommand(args)); err != nil {
			return err
		}
		return nil
	}
	_, stderr, err := runner.Run(ctx, args...)
	if err != nil {
		if stderr != "" {
			return fmt.Errorf("%s: %w: %s", formatGitCommand(args), err, stderr)
		}
		return fmt.Errorf("%s: %w", formatGitCommand(args), err)
	}
	return nil
}
