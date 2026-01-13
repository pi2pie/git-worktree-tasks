package cli

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

func TestRepoBaseName(t *testing.T) {
	runner := fakeRunner{
		responses: map[string]fakeResponse{
			"rev-parse --show-toplevel":  {stdout: "/tmp/example"},
			"rev-parse --git-common-dir": {stdout: ".git"},
		},
	}
	got, err := repoBaseName(context.Background(), runner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "example" {
		t.Fatalf("repoBaseName() = %q, want %q", got, "example")
	}
}
