package cli

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
)

type fakeRunner struct {
	responses map[string]fakeResponse
}

type fakeResponse struct {
	stdout string
	stderr string
	err    error
}

func (f fakeRunner) Run(_ context.Context, args ...string) (string, string, error) {
	key := strings.Join(args, " ")
	if resp, ok := f.responses[key]; ok {
		return resp.stdout, resp.stderr, resp.err
	}
	return "", "", fmt.Errorf("unexpected args: %s", key)
}

func TestNormalizeTaskQuery(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "trims and slugifies", input: "  My Task  ", want: "my-task"},
		{name: "already slugified", input: "my-task", want: "my-task"},
		{name: "empty", input: "   ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeTaskQuery(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("normalizeTaskQuery(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMatchesTask(t *testing.T) {
	task := "my-feature-task"
	if !matchesTask(task, "feature", false) {
		t.Fatalf("expected contains match")
	}
	if matchesTask(task, "feature", true) {
		t.Fatalf("expected strict mismatch")
	}
	if !matchesTask(task, "my-feature-task", true) {
		t.Fatalf("expected strict match")
	}
	if !matchesTask(task, "FEATURE", false) {
		t.Fatalf("expected case-insensitive match")
	}
}

func TestMainWorktreePathFromCommonDir(t *testing.T) {
	tests := []struct {
		name      string
		repoRoot  string
		commonDir string
		want      string
	}{
		{
			name:      "relative git dir",
			repoRoot:  "/tmp/repo",
			commonDir: ".git",
			want:      "/tmp/repo",
		},
		{
			name:      "absolute git dir",
			repoRoot:  "/tmp/linked",
			commonDir: "/tmp/main/.git",
			want:      "/tmp/main",
		},
		{
			name:      "bare repo common dir",
			repoRoot:  "/tmp/repo",
			commonDir: "/tmp/repo",
			want:      "/tmp/repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mainWorktreePathFromCommonDir(tt.repoRoot, tt.commonDir)
			if got != tt.want {
				t.Fatalf("mainWorktreePathFromCommonDir(%q, %q) = %q, want %q", tt.repoRoot, tt.commonDir, got, tt.want)
			}
		})
	}
}

func TestFallbackPathForBranch(t *testing.T) {
	runner := fakeRunner{
		responses: map[string]fakeResponse{
			"-C /tmp/repo branch --list feature": {stdout: "  feature\n"},
			"rev-parse --git-common-dir":         {stdout: ".git"},
		},
	}
	path, ok, err := fallbackPathForBranch(context.Background(), runner, "/tmp/repo", "feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected fallback path to be available")
	}
	if path != "/tmp/repo" {
		t.Fatalf("fallbackPathForBranch() = %q, want %q", path, "/tmp/repo")
	}
}

func TestDeriveClassicTask(t *testing.T) {
	defaultRepoRoot := "/tmp/repo"
	defaultMainWorktree := "/tmp/repo"
	repo := "repo"
	tests := []struct {
		name         string
		repoRoot     string
		mainWorktree string
		wt           worktree.Worktree
		want         string
	}{
		{
			name: "path naming convention wins",
			wt: worktree.Worktree{
				Path:   "/tmp/repo_task-from-path",
				Branch: "refs/heads/other-branch",
			},
			want: "task-from-path",
		},
		{
			name: "fallback to branch task for custom path",
			wt: worktree.Worktree{
				Path:   "/tmp/repo/.claude/worktrees/new-task",
				Branch: "refs/heads/new-task",
			},
			want: "new-task",
		},
		{
			name: "fallback branch task is slugified",
			wt: worktree.Worktree{
				Path:   "/tmp/repo/.claude/worktrees/release",
				Branch: "refs/heads/release/1.0",
			},
			want: "release/1-0",
		},
		{
			name: "main worktree path stays empty task",
			wt: worktree.Worktree{
				Path:   "/tmp/repo",
				Branch: "refs/heads/main",
			},
			want: "",
		},
		{
			name:         "main worktree stays empty when invoked from linked worktree",
			repoRoot:     "/tmp/repo/.claude/worktrees/new-task",
			mainWorktree: "/tmp/repo",
			wt: worktree.Worktree{
				Path:   "/tmp/repo",
				Branch: "refs/heads/main",
			},
			want: "",
		},
		{
			name:         "linked worktree still infers task when invoked from linked worktree",
			repoRoot:     "/tmp/repo/.claude/worktrees/new-task",
			mainWorktree: "/tmp/repo",
			wt: worktree.Worktree{
				Path:   "/tmp/repo/.claude/worktrees/new-task",
				Branch: "refs/heads/new-task",
			},
			want: "new-task",
		},
		{
			name: "detached stays empty task",
			wt: worktree.Worktree{
				Path: "/tmp/repo/.claude/worktrees/detached",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoRoot := tt.repoRoot
			if repoRoot == "" {
				repoRoot = defaultRepoRoot
			}
			mainWorktree := tt.mainWorktree
			if mainWorktree == "" {
				mainWorktree = defaultMainWorktree
			}
			got, err := deriveClassicTask(repoRoot, mainWorktree, repo, tt.wt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("deriveClassicTask() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatGitCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "simple args",
			args: []string{"status", "--porcelain"},
			want: "git status --porcelain",
		},
		{
			name: "space in arg",
			args: []string{"-C", "/tmp/my repo", "status"},
			want: "git -C '/tmp/my repo' status",
		},
		{
			name: "single quote in arg",
			args: []string{"commit", "-m", "it's done"},
			want: "git commit -m 'it'\\'\"\\'\"\\'s done'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatGitCommand(tt.args)
			if got != tt.want {
				t.Fatalf("formatGitCommand(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
