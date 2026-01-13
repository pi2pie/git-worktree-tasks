package cli

import (
	"github.com/dev-pi2pie/git-worktree-tasks/tui"
	"github.com/spf13/cobra"
)

func newTUICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "ui",
		Short:  "Launch the TUI (preview)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run()
		},
	}

	return cmd
}
