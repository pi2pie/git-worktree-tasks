package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Runner abstracts git command execution for testability.
type Runner interface {
	Run(ctx context.Context, args ...string) (stdout string, stderr string, err error)
}

// ExecRunner executes git commands using os/exec.
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, args ...string) (string, string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), nil
}
