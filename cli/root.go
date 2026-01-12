package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.0.3"

var errCanceled = errors.New("git worktree task process canceled")

func Execute() int {
	cmd, state := gitWorkTreeCommand()
	if err := cmd.Execute(); err != nil {
		if errors.Is(err, errCanceled) {
			fmt.Fprintln(cmd.ErrOrStderr(), "git worktree task process canceled")
			return 3
		}
		fmt.Fprintln(cmd.ErrOrStderr(), err)
		return 1
	}
	if state.hasWarnings && state.exitOnWarning {
		return 2
	}
	return 0
}

type runState struct {
	hasWarnings   bool
	exitOnWarning bool
}

func gitWorkTreeCommand() (*cobra.Command, *runState) {
	state := &runState{}
	cmd := &cobra.Command{
		Use:           "git-worktree-tasks",
		Short:         "Task-based git worktree helper",
		Long:          "Create, manage, and clean up git worktrees based on task names.",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	cmd.AddCommand(
		newCreateCommand(state),
		newFinishCommand(state),
		newCleanupCommand(state),
		newListCommand(state),
		newStatusCommand(state),
		newTUICommand(state),
	)

	return cmd, state
}
