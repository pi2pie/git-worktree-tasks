package cli

import (
	"context"

	"github.com/pi2pie/git-worktree-tasks/internal/config"
	"github.com/spf13/cobra"
)

type configKey struct{}

func withConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, configKey{}, cfg)
}

func configFromContext(ctx context.Context) (*config.Config, bool) {
	cfg, ok := ctx.Value(configKey{}).(*config.Config)
	return cfg, ok
}

func flagChangedAny(cmd *cobra.Command, names ...string) bool {
	for _, name := range names {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}
