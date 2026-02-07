package cli

import (
	"strings"
	"testing"
)

func TestMaskHomePathWithPOSIX(t *testing.T) {
	home := "/Users/alice"
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "exact home",
			path: "/Users/alice",
			want: "$HOME",
		},
		{
			name: "descendant path",
			path: "/Users/alice/project",
			want: "$HOME/project",
		},
		{
			name: "prefix safe mismatch",
			path: "/Users/alice2/project",
			want: "/Users/alice2/project",
		},
		{
			name: "outside home",
			path: "/tmp/project",
			want: "/tmp/project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maskHomePathWith(tt.path, home, false)
			if got != tt.want {
				t.Fatalf("maskHomePathWith(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestMaskHomePathWithWindows(t *testing.T) {
	home := `C:\Users\Alice`
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "exact home",
			path: `C:\Users\Alice`,
			want: `%USERPROFILE%`,
		},
		{
			name: "descendant path",
			path: `C:\Users\Alice\project`,
			want: `%USERPROFILE%\project`,
		},
		{
			name: "mixed case and separator",
			path: `c:/users/alice/project`,
			want: `%USERPROFILE%\project`,
		},
		{
			name: "prefix safe mismatch",
			path: `C:\Users\Alice2\project`,
			want: `C:\Users\Alice2\project`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maskHomePathWith(tt.path, home, true)
			if got != tt.want {
				t.Fatalf("maskHomePathWith(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestFormatGitCommandForDryRunWithContext(t *testing.T) {
	args := []string{"-C", "/Users/alice/repo", "worktree", "add", "/Users/alice/repo_task", "main"}
	maskCtx := pathMaskContext{
		enabled: true,
		home:    "/Users/alice",
		windows: false,
	}

	got := formatGitCommandForDryRunWithContext(args, maskCtx)
	if !strings.Contains(got, "$HOME/repo") {
		t.Fatalf("expected repo root to be masked, got %q", got)
	}
	if !strings.Contains(got, "$HOME/repo_task") {
		t.Fatalf("expected worktree path to be masked, got %q", got)
	}
	if strings.Contains(got, "/Users/alice/repo") {
		t.Fatalf("did not expect raw home path in output, got %q", got)
	}
}

func TestFormatGitCommandForDryRunWithContextDisabled(t *testing.T) {
	args := []string{"-C", "/Users/alice/repo", "status"}
	maskCtx := pathMaskContext{
		enabled: false,
		home:    "/Users/alice",
		windows: false,
	}

	got := formatGitCommandForDryRunWithContext(args, maskCtx)
	if got != "git -C /Users/alice/repo status" {
		t.Fatalf("formatGitCommandForDryRunWithContext() = %q, want %q", got, "git -C /Users/alice/repo status")
	}
}
