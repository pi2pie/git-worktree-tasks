package git

import (
	"context"
	"fmt"
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

func TestRepoBaseName(t *testing.T) {
	runner := fakeRunner{
		responses: map[string]fakeResponse{
			"rev-parse --show-toplevel":  {stdout: "/tmp/example"},
			"rev-parse --git-common-dir": {stdout: ".git"},
		},
	}
	got, err := RepoBaseName(context.Background(), runner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "example" {
		t.Fatalf("RepoBaseName() = %q, want %q", got, "example")
	}
}
