package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRepoRoot_Found(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	nested := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	got, ok, err := findRepoRoot(nested)
	if err != nil {
		t.Fatalf("findRepoRoot() error = %v", err)
	}
	if !ok {
		t.Fatalf("findRepoRoot() ok = false, want true")
	}
	if got != root {
		t.Fatalf("findRepoRoot() = %q, want %q", got, root)
	}
}

func TestFindRepoRoot_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, ok, err := findRepoRoot(dir)
	if err != nil {
		t.Fatalf("findRepoRoot() error = %v", err)
	}
	if ok {
		t.Fatalf("findRepoRoot() ok = true, want false")
	}
}

func TestFindRepoRoot_GitFile(t *testing.T) {
	// In worktrees/submodules, .git can be a file (gitdir: /path/to/...)
	root := t.TempDir()
	gitFile := filepath.Join(root, ".git")
	if err := os.WriteFile(gitFile, []byte("gitdir: /some/path\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(.git) error = %v", err)
	}
	nested := filepath.Join(root, "sub")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("Mkdir(sub) error = %v", err)
	}

	got, ok, err := findRepoRoot(nested)
	if err != nil {
		t.Fatalf("findRepoRoot() error = %v", err)
	}
	if !ok {
		t.Fatalf("findRepoRoot() ok = false, want true")
	}
	if got != root {
		t.Fatalf("findRepoRoot() = %q, want %q", got, root)
	}
}

func TestHasGitDir_Directory(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}

	ok, err := hasGitDir(dir)
	if err != nil {
		t.Fatalf("hasGitDir() error = %v", err)
	}
	if !ok {
		t.Fatalf("hasGitDir() = false, want true")
	}
}

func TestHasGitDir_File(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /path\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(.git) error = %v", err)
	}

	ok, err := hasGitDir(dir)
	if err != nil {
		t.Fatalf("hasGitDir() error = %v", err)
	}
	if !ok {
		t.Fatalf("hasGitDir() = false, want true")
	}
}

func TestHasGitDir_NotExists(t *testing.T) {
	dir := t.TempDir()

	ok, err := hasGitDir(dir)
	if err != nil {
		t.Fatalf("hasGitDir() error = %v", err)
	}
	if ok {
		t.Fatalf("hasGitDir() = true, want false")
	}
}

func TestProjectConfigRoot_WithGit(t *testing.T) {
	root := t.TempDir()
	// Resolve symlinks for macOS /var -> /private/var
	root, _ = filepath.EvalSymlinks(root)
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	subdir := filepath.Join(root, "nested")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}

	restore := chdir(t, subdir)
	defer restore()

	got, err := projectConfigRoot()
	if err != nil {
		t.Fatalf("projectConfigRoot() error = %v", err)
	}
	if got != root {
		t.Fatalf("projectConfigRoot() = %q, want %q", got, root)
	}
}

func TestProjectConfigRoot_NoGit(t *testing.T) {
	dir := t.TempDir()
	// Resolve symlinks for macOS /var -> /private/var
	dir, _ = filepath.EvalSymlinks(dir)

	restore := chdir(t, dir)
	defer restore()

	got, err := projectConfigRoot()
	if err != nil {
		t.Fatalf("projectConfigRoot() error = %v", err)
	}
	if got != dir {
		t.Fatalf("projectConfigRoot() = %q, want %q (cwd fallback)", got, dir)
	}
}
