package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pi2pie/git-worktree-tasks/internal/config"
	"github.com/spf13/cobra"
)

type neverRunRunner struct{}

func (neverRunRunner) Run(_ context.Context, _ ...string) (string, string, error) {
	return "", "", nil
}

func TestRunGitDryRunMasksPathsByDefault(t *testing.T) {
	t.Setenv("HOME", "/Users/alice")

	cmd := &cobra.Command{}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := runGit(context.Background(), cmd, true, neverRunRunner{}, "-C", "/Users/alice/repo", "status"); err != nil {
		t.Fatalf("runGit() error = %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != `git -C "$HOME/repo" status` {
		t.Fatalf("runGit() output = %q, want %q", got, `git -C "$HOME/repo" status`)
	}
}

func TestRunGitDryRunRespectsConfigMaskDisabled(t *testing.T) {
	t.Setenv("HOME", "/Users/alice")

	cmd := &cobra.Command{}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	cfg := &config.Config{
		DryRun: config.DryRunConfig{
			MaskSensitivePaths: false,
		},
	}
	ctx := withConfig(context.Background(), cfg)

	if err := runGit(ctx, cmd, true, neverRunRunner{}, "-C", "/Users/alice/repo", "status"); err != nil {
		t.Fatalf("runGit() error = %v", err)
	}
	got := strings.TrimSpace(out.String())
	if got != "git -C /Users/alice/repo status" {
		t.Fatalf("runGit() output = %q, want %q", got, "git -C /Users/alice/repo status")
	}
}
