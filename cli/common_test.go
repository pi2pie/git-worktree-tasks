package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDisplayPathRelative(t *testing.T) {
	root := t.TempDir()
	repoRoot := filepath.Join(root, "repo")
	if err := os.MkdirAll(repoRoot, 0o755); err != nil {
		t.Fatalf("mkdir repoRoot: %v", err)
	}
	absPath := filepath.Join(root, "repo_task")
	if err := os.MkdirAll(absPath, 0o755); err != nil {
		t.Fatalf("mkdir absPath: %v", err)
	}

	got := displayPath(repoRoot, absPath, false)
	want, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		t.Fatalf("rel: %v", err)
	}
	if got != want {
		t.Fatalf("relative path = %q, want %q", got, want)
	}
}

func TestDisplayPathAbsolute(t *testing.T) {
	root := t.TempDir()
	repoRoot := filepath.Join(root, "repo")
	if err := os.MkdirAll(repoRoot, 0o755); err != nil {
		t.Fatalf("mkdir repoRoot: %v", err)
	}
	absPath := filepath.Join(root, "repo_task")
	if err := os.MkdirAll(absPath, 0o755); err != nil {
		t.Fatalf("mkdir absPath: %v", err)
	}

	relInput := filepath.Join("..", "repo_task")
	got := displayPath(repoRoot, relInput, true)
	if got != absPath {
		t.Fatalf("absolute path = %q, want %q", got, absPath)
	}
}
