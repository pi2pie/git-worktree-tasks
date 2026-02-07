package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestDetectApplyConflicts(t *testing.T) {
	runner := fakeRunner{
		responses: map[string]fakeResponse{
			"-C /repo status --porcelain":                    {stdout: " M file.txt\n"},
			"-C /repo diff --name-only HEAD":                 {stdout: "file.txt\n"},
			"-C /repo ls-files --others --exclude-standard":  {stdout: ""},
			"-C /codex diff --name-only HEAD":                {stdout: "file.txt\n"},
			"-C /codex ls-files --others --exclude-standard": {stdout: ""},
		},
	}

	reasons, err := detectApplyConflicts(context.Background(), runner, "/repo", "local checkout", "/codex")
	if err != nil {
		t.Fatalf("detectApplyConflicts error: %v", err)
	}
	if len(reasons) != 2 {
		t.Fatalf("expected 2 conflict reasons, got %d", len(reasons))
	}

	var hasDirty, hasOverlap bool
	for _, reason := range reasons {
		if strings.Contains(reason, "uncommitted changes") {
			hasDirty = true
		}
		if strings.Contains(reason, "both sides modified") {
			hasOverlap = true
		}
	}
	if !hasDirty {
		t.Fatalf("expected uncommitted changes reason, got %v", reasons)
	}
	if !hasOverlap {
		t.Fatalf("expected overlap reason, got %v", reasons)
	}
}

func TestDetectApplyConflictsNone(t *testing.T) {
	runner := fakeRunner{
		responses: map[string]fakeResponse{
			"-C /repo status --porcelain":                    {stdout: ""},
			"-C /repo diff --name-only HEAD":                 {stdout: "local.txt\n"},
			"-C /repo ls-files --others --exclude-standard":  {stdout: ""},
			"-C /codex diff --name-only HEAD":                {stdout: "other.txt\n"},
			"-C /codex ls-files --others --exclude-standard": {stdout: ""},
		},
	}

	reasons, err := detectApplyConflicts(context.Background(), runner, "/repo", "local checkout", "/codex")
	if err != nil {
		t.Fatalf("detectApplyConflicts error: %v", err)
	}
	if len(reasons) != 0 {
		t.Fatalf("expected no conflict reasons, got %v", reasons)
	}
}

func TestResolveTransferPlan(t *testing.T) {
	repoRoot := "/repo"
	worktreePath := "/codex"

	tests := []struct {
		name      string
		to        string
		wantSrc   string
		wantDst   string
		wantSrcN  string
		wantDstN  string
		expectErr bool
	}{
		{
			name:     "to local",
			to:       transferToLocal,
			wantSrc:  worktreePath,
			wantDst:  repoRoot,
			wantSrcN: "Codex worktree",
			wantDstN: "local checkout",
		},
		{
			name:     "to worktree",
			to:       transferToWorktree,
			wantSrc:  repoRoot,
			wantDst:  worktreePath,
			wantSrcN: "local checkout",
			wantDstN: "Codex worktree",
		},
		{
			name:      "invalid destination",
			to:        "somewhere",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveTransferPlan(repoRoot, worktreePath, tt.to)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveTransferPlan() error: %v", err)
			}
			if got.sourceRoot != tt.wantSrc || got.destinationRoot != tt.wantDst {
				t.Fatalf("resolveTransferPlan() roots = (%q -> %q), want (%q -> %q)", got.sourceRoot, got.destinationRoot, tt.wantSrc, tt.wantDst)
			}
			if got.sourceName != tt.wantSrcN || got.destinationName != tt.wantDstN {
				t.Fatalf("resolveTransferPlan() names = (%q -> %q), want (%q -> %q)", got.sourceName, got.destinationName, tt.wantSrcN, tt.wantDstN)
			}
		})
	}
}

func TestConflictReasonsForApply(t *testing.T) {
	reasons := conflictReasonsForApply(transferPreflight{
		destinationDirty: true,
		overlappingFiles: 3,
	}, "local checkout")

	if len(reasons) != 2 {
		t.Fatalf("expected 2 reasons, got %d", len(reasons))
	}
	if !strings.Contains(reasons[0], "local checkout has uncommitted changes") {
		t.Fatalf("unexpected dirty reason: %v", reasons)
	}
	if !strings.Contains(reasons[1], "both sides modified 3 overlapping file(s)") {
		t.Fatalf("unexpected overlap reason: %v", reasons)
	}
}

func TestDryRunActions(t *testing.T) {
	plan := transferPlan{
		destinationRoot: "/repo",
		sourceRoot:      "/codex",
	}
	preflight := transferPreflight{
		trackedPatch:   true,
		untrackedFiles: []string{"a.txt"},
	}

	actions := dryRunActions(handoffOverwrite, plan, preflight, false)
	if len(actions) != 4 {
		t.Fatalf("expected 4 actions, got %d (%v)", len(actions), actions)
	}
	if !strings.Contains(actions[0], "[destructive] git -C /repo reset --hard") {
		t.Fatalf("expected destructive reset action, got %q", actions[0])
	}
	if !strings.Contains(actions[1], "[destructive] git -C /repo clean -fd") {
		t.Fatalf("expected destructive clean action, got %q", actions[1])
	}
	if !strings.Contains(actions[2], "apply <temp-patch>") {
		t.Fatalf("expected apply action, got %q", actions[2])
	}
	if strings.Contains(actions[2], "--check") {
		t.Fatalf("overwrite dry-run should not include apply --check, got %q", actions[2])
	}
	if !strings.Contains(actions[3], "copy /codex/a.txt -> /repo/a.txt") {
		t.Fatalf("expected copy action, got %q", actions[3])
	}
}

func TestDryRunActionsApplyIncludesPatchCheck(t *testing.T) {
	plan := transferPlan{
		destinationRoot: "/repo",
		sourceRoot:      "/codex",
	}
	preflight := transferPreflight{
		trackedPatch: true,
	}

	actions := dryRunActions(handoffApply, plan, preflight, false)
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d (%v)", len(actions), actions)
	}
	if !strings.Contains(actions[0], "apply --check <temp-patch>") {
		t.Fatalf("expected apply --check action, got %q", actions[0])
	}
	if !strings.Contains(actions[1], "apply <temp-patch>") {
		t.Fatalf("expected apply action, got %q", actions[1])
	}
}

func TestPrintDryRunPlanMasksPaths(t *testing.T) {
	t.Setenv("HOME", "/Users/alice")
	plan := transferPlan{
		to:              transferToLocal,
		sourceRoot:      "/Users/alice/codex/repo",
		destinationRoot: "/Users/alice/repo",
	}
	preflight := transferPreflight{
		trackedPatch:   true,
		untrackedFiles: []string{"a.txt"},
	}

	var out bytes.Buffer
	if err := printDryRunPlan(&out, handoffApply, plan, preflight, true); err != nil {
		t.Fatalf("printDryRunPlan() error = %v", err)
	}
	text := out.String()
	if !strings.Contains(text, "source: $HOME/codex/repo") {
		t.Fatalf("expected masked source path, got:\n%s", text)
	}
	if !strings.Contains(text, "destination: $HOME/repo") {
		t.Fatalf("expected masked destination path, got:\n%s", text)
	}
	if strings.Contains(text, "/Users/alice/codex/repo") || strings.Contains(text, "/Users/alice/repo") {
		t.Fatalf("did not expect raw home paths when masking is enabled, got:\n%s", text)
	}
}

func TestPrintDryRunPlanMaskDisabled(t *testing.T) {
	t.Setenv("HOME", "/Users/alice")
	plan := transferPlan{
		to:              transferToLocal,
		sourceRoot:      "/Users/alice/codex/repo",
		destinationRoot: "/Users/alice/repo",
	}
	preflight := transferPreflight{
		trackedPatch: true,
	}

	var out bytes.Buffer
	if err := printDryRunPlan(&out, handoffApply, plan, preflight, false); err != nil {
		t.Fatalf("printDryRunPlan() error = %v", err)
	}
	text := out.String()
	if !strings.Contains(text, "source: /Users/alice/codex/repo") {
		t.Fatalf("expected raw source path, got:\n%s", text)
	}
	if !strings.Contains(text, "destination: /Users/alice/repo") {
		t.Fatalf("expected raw destination path, got:\n%s", text)
	}
}

type overwriteRunner struct {
	source    string
	dest      string
	seenCheck bool
	seenApply bool
}

func (r *overwriteRunner) Run(_ context.Context, args ...string) (string, string, error) {
	switch {
	case len(args) == 4 && args[0] == "-C" && args[1] == r.dest && args[2] == "reset" && args[3] == "--hard":
		return "", "", nil
	case len(args) == 4 && args[0] == "-C" && args[1] == r.dest && args[2] == "clean" && args[3] == "-fd":
		return "", "", nil
	case len(args) == 5 && args[0] == "-C" && args[1] == r.source && args[2] == "diff" && args[3] == "--binary" && args[4] == "HEAD":
		return "diff --git a/a.txt b/a.txt\n", "", nil
	case len(args) == 5 && args[0] == "-C" && args[1] == r.source && args[2] == "ls-files" && args[3] == "--others" && args[4] == "--exclude-standard":
		return "", "", nil
	case len(args) == 5 && args[0] == "-C" && args[1] == r.dest && args[2] == "apply" && args[3] == "--check":
		r.seenCheck = true
		return "", "", fmt.Errorf("unexpected apply --check in overwrite mode")
	case len(args) == 4 && args[0] == "-C" && args[1] == r.dest && args[2] == "apply":
		r.seenApply = true
		return "", "", nil
	default:
		return "", "", fmt.Errorf("unexpected args: %s", strings.Join(args, " "))
	}
}

func TestTransferChangesOverwriteSkipsPatchCheck(t *testing.T) {
	runner := &overwriteRunner{
		source: "/codex",
		dest:   "/repo",
	}
	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := transferChanges(context.Background(), cmd, runner, "/codex", "/repo", false, true)
	if err != nil {
		t.Fatalf("transferChanges() error = %v", err)
	}
	if runner.seenCheck {
		t.Fatalf("overwrite flow should skip apply --check")
	}
	if !runner.seenApply {
		t.Fatalf("expected overwrite flow to apply patch")
	}
}

type overwriteFallbackRunner struct {
	source string
	dest   string
}

func (r *overwriteFallbackRunner) Run(_ context.Context, args ...string) (string, string, error) {
	switch {
	case len(args) == 4 && args[0] == "-C" && args[1] == r.dest && args[2] == "reset" && args[3] == "--hard":
		return "", "", nil
	case len(args) == 4 && args[0] == "-C" && args[1] == r.dest && args[2] == "clean" && args[3] == "-fd":
		return "", "", nil
	case len(args) == 5 && args[0] == "-C" && args[1] == r.source && args[2] == "diff" && args[3] == "--binary" && args[4] == "HEAD":
		return "diff --git a/tracked.txt b/tracked.txt\n", "", nil
	case len(args) == 5 && args[0] == "-C" && args[1] == r.dest && args[2] == "apply":
		return "", "", fmt.Errorf("simulated apply failure")
	case len(args) == 5 && args[0] == "-C" && args[1] == r.source && args[2] == "diff" && args[3] == "--name-status" && args[4] == "HEAD":
		return "M\ttracked.txt\n", "", nil
	case len(args) == 5 && args[0] == "-C" && args[1] == r.source && args[2] == "ls-files" && args[3] == "--others" && args[4] == "--exclude-standard":
		return "", "", nil
	default:
		return "", "", fmt.Errorf("unexpected args: %s", strings.Join(args, " "))
	}
}

func TestTransferChangesOverwriteFallsBackAfterApplyFailure(t *testing.T) {
	sourceRoot := t.TempDir()
	destinationRoot := t.TempDir()
	writeFile := func(root, rel, content string) {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("MkdirAll(%s): %v", path, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("WriteFile(%s): %v", path, err)
		}
	}
	writeFile(sourceRoot, "tracked.txt", "new content\n")
	writeFile(destinationRoot, "tracked.txt", "old content\n")

	runner := &overwriteFallbackRunner{
		source: sourceRoot,
		dest:   destinationRoot,
	}
	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := transferChanges(context.Background(), cmd, runner, sourceRoot, destinationRoot, false, true)
	if err != nil {
		t.Fatalf("transferChanges() error = %v", err)
	}
	got, err := os.ReadFile(filepath.Join(destinationRoot, "tracked.txt"))
	if err != nil {
		t.Fatalf("ReadFile(tracked.txt): %v", err)
	}
	if string(got) != "new content\n" {
		t.Fatalf("destination tracked.txt = %q, want %q", string(got), "new content\n")
	}
}
