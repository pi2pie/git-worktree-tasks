package worktree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestStatusAheadBehindParseError(t *testing.T) {
	t.Helper()
	root := t.TempDir()
	worktreePath := filepath.Join(root, "wt")
	if err := os.MkdirAll(filepath.Join(worktreePath, ".git"), 0o755); err != nil {
		t.Fatalf("setup worktree: %v", err)
	}

	runner := fakeRunner{
		responses: map[string]fakeResponse{
			"-C " + worktreePath + " status --porcelain":                        {stdout: ""},
			"-C " + worktreePath + " rev-parse --verify HEAD":                   {},
			"-C " + worktreePath + " rev-parse --short HEAD":                    {stdout: "abc1234"},
			"-C " + worktreePath + " log -1 --pretty=format:%H %s":              {stdout: "abcdef1234567890 message"},
			"-C " + worktreePath + " merge-base HEAD main":                      {stdout: "abcdef1234567890"},
			"-C " + worktreePath + " rev-list --left-right --count main...HEAD": {stdout: "x 1"},
		},
	}

	_, err := Status(context.Background(), runner, worktreePath, "main")
	if err == nil {
		t.Fatalf("expected parse error")
	}
	if !strings.Contains(err.Error(), "status ahead/behind parse") {
		t.Fatalf("unexpected error: %v", err)
	}
}
