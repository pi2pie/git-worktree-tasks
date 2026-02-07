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
