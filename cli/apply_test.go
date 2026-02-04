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

	reasons, err := detectApplyConflicts(context.Background(), runner, "/repo", "/codex")
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

	reasons, err := detectApplyConflicts(context.Background(), runner, "/repo", "/codex")
	if err != nil {
		t.Fatalf("detectApplyConflicts error: %v", err)
	}
	if len(reasons) != 0 {
		t.Fatalf("expected no conflict reasons, got %v", reasons)
	}
}
