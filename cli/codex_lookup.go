package cli

import (
	"context"

	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
)

func resolveCodexWorktreePath(ctx context.Context, runner git.Runner, repoRoot, codexWorktreesRoot, opaqueID string) (string, bool, error) {
	worktrees, err := worktree.List(ctx, runner, repoRoot)
	if err != nil {
		return "", false, err
	}
	return resolveCodexWorktreePathFromList(worktrees, repoRoot, codexWorktreesRoot, opaqueID)
}

func resolveCodexWorktreePathFromList(worktrees []worktree.Worktree, repoRoot, codexWorktreesRoot, opaqueID string) (string, bool, error) {
	for _, wt := range worktrees {
		wtAbs, err := worktree.NormalizePath(repoRoot, wt.Path)
		if err != nil {
			return "", false, err
		}
		id, _, ok := codexWorktreeInfo(codexWorktreesRoot, wtAbs)
		if !ok || id != opaqueID {
			continue
		}
		return wtAbs, true, nil
	}
	return "", false, nil
}
