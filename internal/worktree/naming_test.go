package worktree

import (
	"path/filepath"
	"testing"
)

func TestSlugifyTask(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple", input: "feature-123", want: "feature-123"},
		{name: "spaces", input: "feature 123", want: "feature-123"},
		{name: "symbols", input: "feat@123", want: "feat-123"},
		{name: "trim dashes", input: "---", want: "task"},
		{name: "preserve separators", input: "foo/bar_baz", want: "foo/bar_baz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SlugifyTask(tt.input)
			if got != tt.want {
				t.Fatalf("SlugifyTask(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWorktreePath(t *testing.T) {
	repoRoot := filepath.Join("/tmp", "repo")
	repoName := "repo"
	task := "test"

	got := WorktreePath(repoRoot, repoName, task)
	want := filepath.Join(filepath.Dir(repoRoot), repoName+"_"+task)
	if got != want {
		t.Fatalf("WorktreePath() = %q, want %q", got, want)
	}
}

func TestTaskFromPath(t *testing.T) {
	repoName := "repo"
	path := filepath.Join("/tmp", "repo_task")

	task, ok := TaskFromPath(repoName, path)
	if !ok {
		t.Fatal("expected ok true")
	}
	if task != "task" {
		t.Fatalf("task = %q, want %q", task, "task")
	}

	_, ok = TaskFromPath(repoName, filepath.Join("/tmp", "other_task"))
	if ok {
		t.Fatal("expected ok false")
	}
}
