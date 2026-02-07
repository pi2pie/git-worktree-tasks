package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

const (
	transferToLocal    = "local"
	transferToWorktree = "worktree"
)

type handoffMode string

const (
	handoffApply     handoffMode = "apply"
	handoffOverwrite handoffMode = "overwrite"
)

type applyOptions struct {
	yes    bool
	dryRun bool
	to     string
	force  bool
}

type handoffOptions struct {
	yes    bool
	dryRun bool
	to     string
}

type transferPlan struct {
	to              string
	sourceRoot      string
	sourceName      string
	destinationRoot string
	destinationName string
}

type transferPreflight struct {
	destinationDirty bool
	overlappingFiles int
	trackedPatch     bool
	untrackedFiles   []string
}

type applyConflictError struct {
	reason string
	err    error
}

func (e *applyConflictError) Error() string {
	if e.err == nil {
		return e.reason
	}
	return fmt.Sprintf("%s: %v", e.reason, e.err)
}

func (e *applyConflictError) Unwrap() error { return e.err }

func newApplyCommand() *cobra.Command {
	opts := &applyOptions{}
	cmd := &cobra.Command{
		Use:   "apply <task>",
		Short: "Apply non-destructive changes between a Codex worktree and local checkout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mode := handoffApply
			if opts.force {
				mode = handoffOverwrite
			}
			return runCodexHandoff(cmd, strings.TrimSpace(args[0]), handoffOptions{
				yes:    opts.yes,
				dryRun: opts.dryRun,
				to:     opts.to,
			}, mode)
		},
	}

	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")
	cmd.Flags().StringVar(&opts.to, "to", transferToLocal, "transfer destination: local or worktree")
	cmd.Flags().BoolVar(&opts.force, "force", false, "compatibility alias for overwrite behavior")
	return cmd
}

func newOverwriteCommand() *cobra.Command {
	opts := &handoffOptions{}
	cmd := &cobra.Command{
		Use:   "overwrite <task>",
		Short: "Overwrite destination with source changes in codex mode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCodexHandoff(cmd, strings.TrimSpace(args[0]), *opts, handoffOverwrite)
		},
	}

	cmd.Flags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show git commands without executing")
	cmd.Flags().StringVar(&opts.to, "to", transferToLocal, "transfer destination: local or worktree")
	return cmd
}

func runCodexHandoff(cmd *cobra.Command, opaqueID string, opts handoffOptions, mode handoffMode) error {
	ctx := cmd.Context()
	modeCtx, err := resolveModeContext(cmd, false)
	if err != nil {
		return err
	}
	cfg, ok := configFromContext(ctx)
	if !ok || modeCtx.mode != modeCodex {
		return fmt.Errorf("%s is only supported in --mode=codex", mode)
	}

	if !cmd.Flags().Changed("yes") {
		opts.yes = !cfg.Cleanup.Confirm
	}

	runner := defaultRunner()
	plan, err := resolveCodexHandoffPlan(ctx, runner, modeCtx, opaqueID, opts.to)
	if err != nil {
		return err
	}

	preflight := transferPreflight{}
	if opts.dryRun || mode == handoffApply {
		preflight, err = collectTransferPreflight(ctx, runner, plan.sourceRoot, plan.destinationRoot, opts.dryRun)
		if err != nil {
			return err
		}
	}

	if opts.dryRun {
		if err := printDryRunPlan(cmd.OutOrStdout(), mode, plan, preflight); err != nil {
			return err
		}
	}

	if mode == handoffApply {
		reasons := conflictReasonsForApply(preflight, plan.destinationName)
		if len(reasons) > 0 {
			if err := printConflictReasons(cmd.OutOrStdout(), reasons); err != nil {
				return err
			}
			if err := printOverwriteHint(cmd.OutOrStdout(), plan.to, opaqueID); err != nil {
				return err
			}
			return errApplyBlocked
		}
	}

	if mode == handoffOverwrite && !opts.dryRun {
		if err := confirmOverwrite(cmd, opts.yes, plan); err != nil {
			return err
		}
	}

	err = transferChanges(ctx, cmd, runner, plan.sourceRoot, plan.destinationRoot, opts.dryRun, mode == handoffOverwrite)
	if err != nil {
		if mode == handoffApply {
			var conflictErr *applyConflictError
			if errors.As(err, &conflictErr) {
				if err := printConflictReasons(cmd.OutOrStdout(), []string{conflictErr.reason}); err != nil {
					return err
				}
				if err := printOverwriteHint(cmd.OutOrStdout(), plan.to, opaqueID); err != nil {
					return err
				}
				return errApplyBlocked
			}
		}
		return err
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), ui.SuccessStyle.Render(fmt.Sprintf("%s complete", mode))); err != nil {
		return err
	}
	return nil
}

func printConflictReasons(out io.Writer, reasons []string) error {
	if _, err := fmt.Fprintf(out, "%s\n", ui.WarningStyle.Render("apply blocked (non-destructive mode):")); err != nil {
		return err
	}
	for _, reason := range reasons {
		if _, err := fmt.Fprintf(out, "- %s\n", ui.WarningStyle.Render(reason)); err != nil {
			return err
		}
	}
	return nil
}

func printOverwriteHint(out io.Writer, to, opaqueID string) error {
	if _, err := fmt.Fprintf(out, "%s\n", ui.WarningStyle.Render(fmt.Sprintf("next step: gwtt overwrite --to %s %s", to, opaqueID))); err != nil {
		return err
	}
	_, err := fmt.Fprintf(out, "%s\n", ui.WarningStyle.Render("add --yes to skip overwrite confirmation prompts"))
	return err
}

func confirmOverwrite(cmd *cobra.Command, yes bool, plan transferPlan) error {
	if yes {
		return nil
	}
	ok, err := confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), fmt.Sprintf("Overwrite the %s from the %s?", plan.destinationName, plan.sourceName))
	if err != nil {
		return err
	}
	if !ok {
		return errCanceled
	}
	ok, err = confirmPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), fmt.Sprintf("This will discard %s changes. Continue?", plan.destinationName))
	if err != nil {
		return err
	}
	if !ok {
		return errCanceled
	}
	return nil
}
