package cli

import "github.com/spf13/cobra"

type modeContext struct {
	mode           string
	codexHome      string
	codexWorktrees string
}

// resolveModeContext returns command mode and, when needed, codex path context.
// If includeClassicCodex is true, codex paths are also resolved in classic mode
// on a best-effort basis (errors are ignored to preserve existing behavior).
func resolveModeContext(cmd *cobra.Command, includeClassicCodex bool) (modeContext, error) {
	ctx := modeContext{mode: modeClassic}
	cfg, ok := configFromContext(cmd.Context())
	if !ok {
		return ctx, nil
	}

	ctx.mode = cfg.Mode
	switch ctx.mode {
	case modeCodex:
		home, err := codexHomeDir()
		if err != nil {
			return ctx, err
		}
		ctx.codexHome = home
		ctx.codexWorktrees = codexWorktreesRoot(home)
	default:
		if includeClassicCodex {
			if home, err := codexHomeDir(); err == nil {
				ctx.codexHome = home
				ctx.codexWorktrees = codexWorktreesRoot(home)
			}
		}
	}

	return ctx, nil
}
