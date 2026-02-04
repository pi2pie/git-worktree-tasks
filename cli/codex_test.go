package cli

import (
	"path/filepath"
	"testing"
)

func TestCodexWorktreeInfo(t *testing.T) {
	root := filepath.Join(t.TempDir(), "codex", "worktrees")
	worktreePath := filepath.Join(root, "bf15", "repo")

	opaqueID, rel, ok := codexWorktreeInfo(root, worktreePath)
	if !ok {
		t.Fatalf("expected codex worktree info to be detected")
	}
	if opaqueID != "bf15" {
		t.Fatalf("opaque id = %q, want %q", opaqueID, "bf15")
	}
	if rel != filepath.Join("bf15", "repo") {
		t.Fatalf("relative path = %q, want %q", rel, filepath.Join("bf15", "repo"))
	}

	if _, _, ok := codexWorktreeInfo(root, root); ok {
		t.Fatalf("expected root path to be rejected")
	}
	if _, _, ok := codexWorktreeInfo(root, filepath.Join(t.TempDir(), "other")); ok {
		t.Fatalf("expected outside path to be rejected")
	}
}

func TestDisplayPathForModeCodex(t *testing.T) {
	repoRoot := filepath.Join(t.TempDir(), "repo")
	codexHome := filepath.Join(t.TempDir(), "codex")
	worktreePath := filepath.Join(codexHome, "worktrees", "bf15", "repo")

	got := displayPathForMode(repoRoot, worktreePath, false, modeCodex, codexHome)
	want := filepath.Join("$CODEX_HOME", "worktrees", "bf15", "repo")
	if got != want {
		t.Fatalf("display path = %q, want %q", got, want)
	}
}
