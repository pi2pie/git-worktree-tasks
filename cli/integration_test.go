package cli_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pi2pie/git-worktree-tasks/cli"
)

type listRow struct {
	Task    string `json:"task"`
	Branch  string `json:"branch"`
	Path    string `json:"path"`
	Present bool   `json:"present"`
	Head    string `json:"head"`
}

type statusRow struct {
	Task       string `json:"task"`
	Branch     string `json:"branch"`
	Path       string `json:"path"`
	Base       string `json:"base"`
	Target     string `json:"target"`
	LastCommit string `json:"last_commit"`
	Dirty      bool   `json:"dirty"`
	Ahead      int    `json:"ahead"`
	Behind     int    `json:"behind"`
}

func TestIntegrationCreateListStatusFinish(t *testing.T) {
	repoDir := initRepo(t, true)
	worktreePath := strings.TrimSpace(runCLI(t, repoDir, "", "--nocolor", "create", "my-task", "--output", "raw"))
	if worktreePath == "" {
		t.Fatalf("expected worktree path output")
	}
	absWorktreePath := filepath.Clean(filepath.Join(repoDir, worktreePath))
	if _, err := os.Stat(absWorktreePath); err != nil {
		t.Fatalf("expected worktree dir to exist: %v", err)
	}

	listOutput := runCLI(t, repoDir, "", "--nocolor", "list", "my-task", "--output", "json")
	var listRows []listRow
	if err := json.Unmarshal([]byte(listOutput), &listRows); err != nil {
		t.Fatalf("parse list json: %v", err)
	}
	if len(listRows) != 1 {
		t.Fatalf("expected 1 list row, got %d", len(listRows))
	}
	if listRows[0].Task != "my-task" || listRows[0].Branch != "my-task" {
		t.Fatalf("unexpected list row: %+v", listRows[0])
	}

	statusOutput := runCLI(t, repoDir, "", "--nocolor", "status", "my-task", "--output", "json")
	var statusRows []statusRow
	if err := json.Unmarshal([]byte(statusOutput), &statusRows); err != nil {
		t.Fatalf("parse status json: %v", err)
	}
	if len(statusRows) != 1 {
		t.Fatalf("expected 1 status row, got %d", len(statusRows))
	}
	if statusRows[0].LastCommit == "" {
		t.Fatalf("expected last_commit to be populated, got empty")
	}

	writeFile(t, absWorktreePath, "task.txt", "task change\n")
	runGit(t, absWorktreePath, "add", "task.txt")
	runGit(t, absWorktreePath, "commit", "-m", "task change")

	runCLI(t, repoDir, "", "--nocolor", "finish", "my-task", "--cleanup", "--yes")
	if _, err := os.Stat(absWorktreePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected worktree removed, stat error: %v", err)
	}
	if branchExists(t, repoDir, "my-task") {
		t.Fatalf("expected branch to be removed")
	}
}

func TestIntegrationStatusNoCommits(t *testing.T) {
	repoDir := initRepo(t, false)
	statusOutput := runCLI(t, repoDir, "", "--nocolor", "status", "--output", "json")
	var rows []statusRow
	if err := json.Unmarshal([]byte(statusOutput), &rows); err != nil {
		t.Fatalf("parse status json: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected at least one status row")
	}
	if rows[0].LastCommit != "empty history" {
		t.Fatalf("expected empty history, got %q", rows[0].LastCommit)
	}
	if rows[0].Base != "empty history" {
		t.Fatalf("expected base empty history, got %q", rows[0].Base)
	}
}

func TestIntegrationDetachedHeadAndPrunableCleanup(t *testing.T) {
	repoDir := initRepo(t, true)
	head := strings.TrimSpace(runGit(t, repoDir, "rev-parse", "HEAD"))
	runGit(t, repoDir, "checkout", "--detach", head)

	_, err := runCLIError(t, repoDir, "", "--nocolor", "create", "my-task")
	if err == nil || !strings.Contains(err.Error(), "detached HEAD") {
		t.Fatalf("expected detached HEAD error, got %v", err)
	}

	worktreePath := strings.TrimSpace(runCLI(t, repoDir, "", "--nocolor", "create", "my-task", "--base", "main", "--output", "raw"))
	if worktreePath == "" {
		t.Fatalf("expected worktree path output")
	}
	absWorktreePath := filepath.Clean(filepath.Join(repoDir, worktreePath))
	if err := os.RemoveAll(absWorktreePath); err != nil {
		t.Fatalf("remove worktree to simulate prunable: %v", err)
	}

	listOutput := runCLI(t, repoDir, "", "--nocolor", "list", "my-task", "--output", "json")
	var listRows []listRow
	if err := json.Unmarshal([]byte(listOutput), &listRows); err != nil {
		t.Fatalf("parse list json: %v", err)
	}
	if len(listRows) != 1 {
		t.Fatalf("expected 1 list row, got %d", len(listRows))
	}

	statusOutput := runCLI(t, repoDir, "", "--nocolor", "status", "my-task", "--output", "json")
	var statusRows []statusRow
	if err := json.Unmarshal([]byte(statusOutput), &statusRows); err != nil {
		t.Fatalf("parse status json: %v", err)
	}
	if len(statusRows) != 1 {
		t.Fatalf("expected 1 status row, got %d", len(statusRows))
	}

	runCLI(t, repoDir, "", "--nocolor", "cleanup", "my-task", "--yes")
	if branchExists(t, repoDir, "my-task") {
		t.Fatalf("expected branch to be removed after cleanup")
	}
}

func initRepo(t *testing.T, withCommit bool) string {
	t.Helper()
	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	if err := runGitCmd(repoDir, "init", "-b", "main"); err != nil {
		runGit(t, repoDir, "init")
		runGit(t, repoDir, "branch", "-m", "main")
	}
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	if withCommit {
		writeFile(t, repoDir, "README.md", "hello\n")
		runGit(t, repoDir, "add", "README.md")
		runGit(t, repoDir, "commit", "-m", "initial commit")
	}
	return repoDir
}

func runCLI(t *testing.T, cwd string, input string, args ...string) string {
	t.Helper()
	stdout, stderr, err := runCLIWithErr(t, cwd, input, args...)
	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}
	return stdout
}

func runCLIError(t *testing.T, cwd string, input string, args ...string) (string, error) {
	t.Helper()
	stdout, _, err := runCLIWithErr(t, cwd, input, args...)
	return stdout, err
}

func runCLIWithErr(t *testing.T, cwd string, input string, args ...string) (string, string, error) {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(original) }()
	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Setenv("GWTT_THEME", "default")

	cmd := cli.RootCommand()
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs(args)
	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s: %v (stderr: %s)", strings.Join(args, " "), err, strings.TrimSpace(errBuf.String()))
	}
	return strings.TrimSpace(out.String())
}

func runGitCmd(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}

func branchExists(t *testing.T, dir, branch string) bool {
	t.Helper()
	output := runGit(t, dir, "branch", "--list", branch)
	return strings.TrimSpace(output) != ""
}
