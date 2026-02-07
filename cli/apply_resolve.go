package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
)

func resolveCodexHandoffPlan(ctx context.Context, runner git.Runner, modeCtx modeContext, opaqueID, to string) (transferPlan, error) {
	if strings.TrimSpace(opaqueID) == "" {
		return transferPlan{}, fmt.Errorf("task query cannot be empty")
	}

	repoRoot, err := repoRoot(ctx, runner)
	if err != nil {
		return transferPlan{}, err
	}
	if _, err := git.CurrentBranch(ctx, runner); err != nil {
		return transferPlan{}, err
	}

	wtPath, found, err := resolveCodexWorktreePath(ctx, runner, repoRoot, modeCtx.codexWorktrees, opaqueID)
	if err != nil {
		return transferPlan{}, err
	}
	if !found {
		return transferPlan{}, fmt.Errorf("no Codex worktree found for %q under %s", opaqueID, filepath.Join("$CODEX_HOME", "worktrees"))
	}

	return resolveTransferPlan(repoRoot, wtPath, to)
}

func resolveTransferPlan(repoRoot, worktreePath, to string) (transferPlan, error) {
	switch strings.TrimSpace(to) {
	case transferToLocal:
		return transferPlan{
			to:              transferToLocal,
			sourceRoot:      worktreePath,
			sourceName:      "Codex worktree",
			destinationRoot: repoRoot,
			destinationName: "local checkout",
		}, nil
	case transferToWorktree:
		return transferPlan{
			to:              transferToWorktree,
			sourceRoot:      repoRoot,
			sourceName:      "local checkout",
			destinationRoot: worktreePath,
			destinationName: "Codex worktree",
		}, nil
	default:
		return transferPlan{}, fmt.Errorf("invalid --to value %q (expected local or worktree)", to)
	}
}
