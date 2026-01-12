package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/dev-pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

var Version = "0.0.4"

var errCanceled = errors.New("git worktree task process canceled")

func Execute() int {
	cmd, state := gitWorkTreeCommand()
	if err := cmd.Execute(); err != nil {
		if errors.Is(err, errCanceled) {
			fmt.Fprintln(cmd.ErrOrStderr(), ui.WarningStyle.Render("git worktree task process canceled"))
			return 3
		}
		fmt.Fprintln(cmd.ErrOrStderr(), ui.ErrorStyle.Render(err.Error()))
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
	noColor       bool
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
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.PersistentFlags().BoolVar(&state.noColor, "nocolor", false, "disable color output")
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		ui.SetColorEnabled(!state.noColor)
		return nil
	}

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
