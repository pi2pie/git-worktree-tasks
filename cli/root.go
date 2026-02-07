package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/config"
	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

var Version = "0.1.2-alpha.1"

var (
	errCanceled     = errors.New("git worktree task process canceled")
	errThemesListed = errors.New("themes listed")
)

func Execute() int {
	cmd, state := gitWorkTreeCommand()
	if err := cmd.Execute(); err != nil {
		if errors.Is(err, errThemesListed) {
			return 0
		}
		if errors.Is(err, errCanceled) {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), ui.WarningStyle.Render("git worktree task process canceled"))
			return 3
		}
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), ui.ErrorStyle.Render(err.Error()))
		return 1
	}
	if state.hasWarnings && state.exitOnWarning {
		return 2
	}
	return 0
}

func RootCommand() *cobra.Command {
	cmd, _ := gitWorkTreeCommand()
	return cmd
}

type runState struct {
	hasWarnings   bool
	exitOnWarning bool
	noColor       bool
	theme         string
	mode          string
	listThemes    bool
}

func gitWorkTreeCommand() (*cobra.Command, *runState) {
	state := &runState{}
	cmd := &cobra.Command{
		Use:           "git-worktree-tasks",
		Aliases:       []string{"gwtt"},
		Short:         "Task-based git worktree helper",
		Long:          "Create, manage, and clean up git worktrees based on task names.",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if state.listThemes {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(ui.ThemeNames(), "\n")); err != nil {
					return err
				}
				return errThemesListed
			}
			return cmd.Help()
		},
	}
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.PersistentFlags().BoolVar(&state.noColor, "nocolor", false, "disable color output")
	cmd.PersistentFlags().StringVar(&state.theme, "theme", ui.DefaultThemeName(), "color theme: "+strings.Join(ui.ThemeNames(), ", "))
	cmd.PersistentFlags().StringVar(&state.mode, "mode", "classic", "execution mode: classic or codex")
	cmd.PersistentFlags().BoolVar(&state.listThemes, "themes", false, "print available themes and exit")
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if state.listThemes {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(ui.ThemeNames(), "\n")); err != nil {
				return err
			}
			return errThemesListed
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		mode := cfg.Mode
		if cmd.Flags().Changed("mode") {
			mode = state.mode
		}
		mode, err = normalizeMode(mode)
		if err != nil {
			return err
		}
		cfg.Mode = mode
		cmd.SetContext(withConfig(cmd.Context(), &cfg))
		themeName := state.theme
		if !cmd.Flags().Changed("theme") {
			if strings.TrimSpace(cfg.Theme.Name) != "" {
				themeName = cfg.Theme.Name
			}
		}
		if err := ui.SetTheme(themeName); err != nil {
			return err
		}
		colorEnabled := cfg.UI.ColorEnabled
		if cmd.Flags().Changed("nocolor") {
			colorEnabled = !state.noColor
		}
		ui.SetColorEnabled(colorEnabled)
		return nil
	}

	cmd.AddCommand(
		newCreateCommand(),
		newFinishCommand(),
		newCleanupCommand(),
		newListCommand(),
		newStatusCommand(),
		newApplyCommand(),
		newTUICommand(),
	)

	return cmd, state
}
