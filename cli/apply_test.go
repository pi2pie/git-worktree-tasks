package cli

import (
	"context"
	"strings"
	"testing"
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

	actions := dryRunActions(handoffOverwrite, plan, preflight)
	if len(actions) != 5 {
		t.Fatalf("expected 5 actions, got %d (%v)", len(actions), actions)
	}
	if !strings.Contains(actions[0], "[destructive] git -C /repo reset --hard") {
		t.Fatalf("expected destructive reset action, got %q", actions[0])
	}
	if !strings.Contains(actions[1], "[destructive] git -C /repo clean -fd") {
		t.Fatalf("expected destructive clean action, got %q", actions[1])
	}
	if !strings.Contains(actions[2], "apply --check <temp-patch>") {
		t.Fatalf("expected apply check action, got %q", actions[2])
	}
	if !strings.Contains(actions[3], "apply <temp-patch>") {
		t.Fatalf("expected apply action, got %q", actions[3])
	}
	if !strings.Contains(actions[4], "copy /codex/a.txt -> /repo/a.txt") {
		t.Fatalf("expected copy action, got %q", actions[4])
	}
}
