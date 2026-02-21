package cli_test

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	Task         string `json:"task"`
	Branch       string `json:"branch"`
	Path         string `json:"path"`
	ModifiedTime string `json:"modified_time"`
	Base         string `json:"base"`
	Target       string `json:"target"`
	LastCommit   string `json:"last_commit"`
	Dirty        bool   `json:"dirty"`
	Ahead        int    `json:"ahead"`
	Behind       int    `json:"behind"`
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
	if statusRows[0].ModifiedTime == "" {
		t.Fatalf("expected modified_time to be populated, got empty")
	}
	parsedModified, err := time.Parse(time.RFC3339, statusRows[0].ModifiedTime)
	if err != nil {
		t.Fatalf("parse modified_time: %v", err)
	}
	if parsedModified.Location() != time.UTC {
		t.Fatalf("expected modified_time in UTC, got %s", parsedModified.Location())
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

func TestIntegrationListStatusCustomPathTaskInference(t *testing.T) {
	repoDir := initRepo(t, true)
	customRelPath := filepath.Join(".claude", "worktrees", "new-task")

	worktreePath := strings.TrimSpace(runCLI(t, repoDir, "", "--nocolor", "create", "new-task", "--path", customRelPath, "--output", "raw"))
	if worktreePath != customRelPath {
		t.Fatalf("expected custom worktree path %q, got %q", customRelPath, worktreePath)
	}

	rawPath := strings.TrimSpace(runCLI(t, repoDir, "", "--nocolor", "list", "new-task", "--output", "raw"))
	if rawPath != customRelPath {
		t.Fatalf("list raw path = %q, want %q", rawPath, customRelPath)
	}

	listOutput := runCLI(t, repoDir, "", "--nocolor", "list", "new-task", "--output", "json")
	var listRows []listRow
	if err := json.Unmarshal([]byte(listOutput), &listRows); err != nil {
		t.Fatalf("parse list json: %v", err)
	}
	if len(listRows) != 1 {
		t.Fatalf("expected 1 list row, got %d", len(listRows))
	}
	if listRows[0].Task != "new-task" {
		t.Fatalf("expected inferred task new-task, got %q", listRows[0].Task)
	}
	if listRows[0].Branch != "new-task" {
		t.Fatalf("expected branch new-task, got %q", listRows[0].Branch)
	}
	if listRows[0].Path != customRelPath {
		t.Fatalf("expected path %q, got %q", customRelPath, listRows[0].Path)
	}

	statusOutput := runCLI(t, repoDir, "", "--nocolor", "status", "new-task", "--output", "json")
	var statusRows []statusRow
	if err := json.Unmarshal([]byte(statusOutput), &statusRows); err != nil {
		t.Fatalf("parse status json: %v", err)
	}
	if len(statusRows) != 1 {
		t.Fatalf("expected 1 status row, got %d", len(statusRows))
	}
	if statusRows[0].Task != "new-task" {
		t.Fatalf("expected inferred task new-task, got %q", statusRows[0].Task)
	}
	if statusRows[0].Branch != "new-task" {
		t.Fatalf("expected branch new-task, got %q", statusRows[0].Branch)
	}
	if statusRows[0].Path != customRelPath {
		t.Fatalf("expected path %q, got %q", customRelPath, statusRows[0].Path)
	}

	strictMismatch := runCLI(t, repoDir, "", "--nocolor", "list", "new", "--strict", "--output", "json")
	var strictRows []listRow
	if err := json.Unmarshal([]byte(strictMismatch), &strictRows); err != nil {
		t.Fatalf("parse strict list json: %v", err)
	}
	if len(strictRows) != 0 {
		t.Fatalf("expected 0 strict rows for partial query, got %d", len(strictRows))
	}

	fuzzyOutput := runCLI(t, repoDir, "", "--nocolor", "list", "new", "--output", "json")
	var fuzzyRows []listRow
	if err := json.Unmarshal([]byte(fuzzyOutput), &fuzzyRows); err != nil {
		t.Fatalf("parse fuzzy list json: %v", err)
	}
	if len(fuzzyRows) != 1 || fuzzyRows[0].Task != "new-task" {
		t.Fatalf("expected fuzzy match for new-task, got %+v", fuzzyRows)
	}

	allListOutput := runCLI(t, repoDir, "", "--nocolor", "list", "--output", "json")
	var allRows []listRow
	if err := json.Unmarshal([]byte(allListOutput), &allRows); err != nil {
		t.Fatalf("parse full list json: %v", err)
	}
	foundMain := false
	for _, row := range allRows {
		if row.Path == "." {
			foundMain = true
			if row.Task != "-" {
				t.Fatalf("expected main worktree task '-', got %q", row.Task)
			}
		}
	}
	if !foundMain {
		t.Fatalf("expected to find main worktree row")
	}
}

func TestIntegrationListRawFallbackHonorsField(t *testing.T) {
	repoDir := initRepo(t, true)
	runGit(t, repoDir, "branch", "feature")

	pathRaw := strings.TrimSpace(runCLI(t, repoDir, "", "--nocolor", "list", "feature", "--output", "raw"))
	if pathRaw != "." {
		t.Fatalf("expected fallback path '.', got %q", pathRaw)
	}

	branchRaw := strings.TrimSpace(runCLI(t, repoDir, "", "--nocolor", "list", "feature", "--output", "raw", "--field", "branch"))
	if branchRaw != "feature" {
		t.Fatalf("expected fallback branch 'feature', got %q", branchRaw)
	}

	taskRaw := strings.TrimSpace(runCLI(t, repoDir, "", "--nocolor", "list", "feature", "--output", "raw", "--field", "task"))
	if taskRaw != "-" {
		t.Fatalf("expected fallback task '-', got %q", taskRaw)
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
	if rows[0].ModifiedTime == "" {
		t.Fatalf("expected modified_time to be populated, got empty")
	}
	parsedModified, err := time.Parse(time.RFC3339, rows[0].ModifiedTime)
	if err != nil {
		t.Fatalf("parse modified_time: %v", err)
	}
	if parsedModified.Location() != time.UTC {
		t.Fatalf("expected modified_time in UTC, got %s", parsedModified.Location())
	}
}

func TestIntegrationStatusCsvIncludesModifiedTime(t *testing.T) {
	repoDir := initRepo(t, true)
	statusOutput := runCLI(t, repoDir, "", "--nocolor", "status", "--output", "csv")
	reader := csv.NewReader(strings.NewReader(statusOutput))
	header, err := reader.Read()
	if err != nil {
		t.Fatalf("read csv header: %v", err)
	}
	modifiedIndex := indexOf(header, "modified_time")
	if modifiedIndex == -1 {
		t.Fatalf("expected modified_time column in header, got %v", header)
	}
	record, err := reader.Read()
	if err != nil {
		t.Fatalf("read csv record: %v", err)
	}
	if len(record) != len(header) {
		t.Fatalf("csv record length %d != header length %d", len(record), len(header))
	}
	if record[modifiedIndex] == "" {
		t.Fatalf("expected modified_time column to be populated, got empty")
	}
	parsedModified, err := time.Parse(time.RFC3339, record[modifiedIndex])
	if err != nil {
		t.Fatalf("parse modified_time: %v", err)
	}
	if parsedModified.Location() != time.UTC {
		t.Fatalf("expected modified_time in UTC, got %s", parsedModified.Location())
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

func TestIntegrationCodexListStatusFiltering(t *testing.T) {
	repoDir := initRepo(t, true)
	codexHome := setCodexHome(t)
	opaqueID := "bf15"
	addCodexWorktree(t, repoDir, codexHome, opaqueID)

	listOutput := runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "list", "--output", "json")
	var listRows []listRow
	if err := json.Unmarshal([]byte(listOutput), &listRows); err != nil {
		t.Fatalf("parse list json: %v", err)
	}
	if len(listRows) != 1 {
		t.Fatalf("expected 1 codex list row, got %d", len(listRows))
	}
	if listRows[0].Task != opaqueID {
		t.Fatalf("expected codex task %q, got %q", opaqueID, listRows[0].Task)
	}
	wantPath := filepath.Join("$CODEX_HOME", "worktrees", opaqueID, filepath.Base(repoDir))
	if listRows[0].Path != wantPath {
		t.Fatalf("expected codex path %q, got %q", wantPath, listRows[0].Path)
	}

	statusOutput := runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "status", "--output", "json")
	var statusRows []statusRow
	if err := json.Unmarshal([]byte(statusOutput), &statusRows); err != nil {
		t.Fatalf("parse status json: %v", err)
	}
	if len(statusRows) != 1 {
		t.Fatalf("expected 1 codex status row, got %d", len(statusRows))
	}
	if statusRows[0].Task != opaqueID {
		t.Fatalf("expected codex task %q, got %q", opaqueID, statusRows[0].Task)
	}
	if statusRows[0].Path != wantPath {
		t.Fatalf("expected codex path %q, got %q", wantPath, statusRows[0].Path)
	}

	classicListOutput := runCLI(t, repoDir, "", "--nocolor", "list", "--output", "json")
	var classicRows []listRow
	if err := json.Unmarshal([]byte(classicListOutput), &classicRows); err != nil {
		t.Fatalf("parse classic list json: %v", err)
	}
	for _, row := range classicRows {
		if row.Task == opaqueID {
			t.Fatalf("expected codex worktree to be filtered in classic mode")
		}
	}
}

func TestIntegrationApplyConflictRequiresExplicitOverwrite(t *testing.T) {
	repoDir := initRepo(t, true)
	codexHome := setCodexHome(t)
	opaqueID := "apply01"
	codexPath := addCodexWorktree(t, repoDir, codexHome, opaqueID)

	writeFile(t, repoDir, "shared.txt", "local change\n")
	writeFile(t, codexPath, "shared.txt", "codex change\n")

	output, err := runCLIError(t, repoDir, "", "--nocolor", "--mode", "codex", "apply", opaqueID)
	if err == nil || !strings.Contains(err.Error(), "apply aborted due to conflicts") {
		t.Fatalf("expected apply conflict abort, got %v", err)
	}
	if !strings.Contains(output, "next step: gwtt overwrite --to local "+opaqueID) {
		t.Fatalf("expected overwrite next-step guidance, got output:\n%s", output)
	}

	content, err := os.ReadFile(filepath.Join(repoDir, "shared.txt"))
	if err != nil {
		t.Fatalf("read local file: %v", err)
	}
	if string(content) != "local change\n" {
		t.Fatalf("expected local content to remain, got %q", string(content))
	}

	content, err = os.ReadFile(filepath.Join(codexPath, "shared.txt"))
	if err != nil {
		t.Fatalf("read codex file: %v", err)
	}
	if string(content) != "codex change\n" {
		t.Fatalf("expected codex content to remain, got %q", string(content))
	}
}

func TestIntegrationApplyDryRunPlanOutput(t *testing.T) {
	repoDir := initRepo(t, true)
	codexHome := setCodexHome(t)
	opaqueID := "applydry1"
	codexPath := addCodexWorktree(t, repoDir, codexHome, opaqueID)

	writeFile(t, codexPath, "dry.txt", "codex dry-run\n")

	output := runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "apply", opaqueID, "--dry-run")
	for _, want := range []string{
		"apply plan",
		"  to: local",
		"preflight",
		"actions",
		"tracked_patch:",
		"untracked_files:",
		"copy ",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, output)
		}
	}
}

func TestIntegrationOverwriteDryRunPlanOutput(t *testing.T) {
	repoDir := initRepo(t, true)
	codexHome := setCodexHome(t)
	opaqueID := "applydry2"
	codexPath := addCodexWorktree(t, repoDir, codexHome, opaqueID)

	writeFile(t, repoDir, "dry-overwrite.txt", "from local\n")
	writeFile(t, codexPath, "dry-overwrite.txt", "from codex\n")

	output := runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "overwrite", opaqueID, "--to", "worktree", "--dry-run")
	for _, want := range []string{
		"overwrite plan",
		"  to: worktree",
		"  overwrite: true",
		"[destructive] git -C ",
		"reset --hard",
		"clean -fd",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, output)
		}
	}
}

func TestIntegrationOverwriteToLocalConfirmation(t *testing.T) {
	repoDir := initRepo(t, true)
	codexHome := setCodexHome(t)
	opaqueID := "apply02"
	codexPath := addCodexWorktree(t, repoDir, codexHome, opaqueID)

	writeFile(t, repoDir, "shared.txt", "local change\n")
	writeFile(t, codexPath, "shared.txt", "codex change\n")

	_, err := runCLIError(t, repoDir, "no\n", "--nocolor", "--mode", "codex", "overwrite", opaqueID)
	if err == nil || !strings.Contains(err.Error(), "canceled") {
		t.Fatalf("expected overwrite to be canceled, got %v", err)
	}

	content, err := os.ReadFile(filepath.Join(repoDir, "shared.txt"))
	if err != nil {
		t.Fatalf("read local file after cancel: %v", err)
	}
	if string(content) != "local change\n" {
		t.Fatalf("expected local content unchanged after cancel, got %q", string(content))
	}

	runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "overwrite", opaqueID, "--to", "local", "--yes")
	content, err = os.ReadFile(filepath.Join(repoDir, "shared.txt"))
	if err != nil {
		t.Fatalf("read local file after overwrite: %v", err)
	}
	if string(content) != "codex change\n" {
		t.Fatalf("expected local content overwritten from codex, got %q", string(content))
	}
}

func TestIntegrationApplyAndOverwriteToWorktree(t *testing.T) {
	repoDir := initRepo(t, true)
	codexHome := setCodexHome(t)
	opaqueID := "apply03"
	codexPath := addCodexWorktree(t, repoDir, codexHome, opaqueID)

	writeFile(t, repoDir, "shared.txt", "from local\n")

	runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "apply", opaqueID, "--to", "worktree")

	content, err := os.ReadFile(filepath.Join(codexPath, "shared.txt"))
	if err != nil {
		t.Fatalf("read codex file after apply --to worktree: %v", err)
	}
	if string(content) != "from local\n" {
		t.Fatalf("expected worktree content from local apply, got %q", string(content))
	}

	writeFile(t, repoDir, "shared.txt", "source overwrite\n")
	writeFile(t, codexPath, "shared.txt", "dest dirty\n")

	_, err = runCLIError(t, repoDir, "", "--nocolor", "--mode", "codex", "apply", opaqueID, "--to", "worktree")
	if err == nil || !strings.Contains(err.Error(), "apply aborted due to conflicts") {
		t.Fatalf("expected apply --to worktree conflict abort, got %v", err)
	}

	runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "apply", opaqueID, "--to", "worktree", "--force", "--yes")

	content, err = os.ReadFile(filepath.Join(codexPath, "shared.txt"))
	if err != nil {
		t.Fatalf("read codex file after apply --force: %v", err)
	}
	if string(content) != "source overwrite\n" {
		t.Fatalf("expected apply --force to overwrite worktree, got %q", string(content))
	}
}

func TestIntegrationCodexCleanupScopeAndConfirm(t *testing.T) {
	repoDir := initRepo(t, true)
	codexHome := setCodexHome(t)
	opaqueID := "clean01"
	codexPath := addCodexWorktree(t, repoDir, codexHome, opaqueID)
	classicPath := addClassicWorktree(t, repoDir, "classic-task")

	_, err := runCLIError(t, repoDir, "", "--nocolor", "--mode", "codex", "cleanup", opaqueID, "--remove-branch")
	if err == nil || !strings.Contains(err.Error(), "branch cleanup is not supported") {
		t.Fatalf("expected codex branch cleanup error, got %v", err)
	}

	_, err = runCLIError(t, repoDir, "yes\nno\n", "--nocolor", "--mode", "codex", "cleanup", opaqueID)
	if err == nil || !strings.Contains(err.Error(), "canceled") {
		t.Fatalf("expected codex cleanup to be canceled, got %v", err)
	}
	if _, err := os.Stat(codexPath); err != nil {
		t.Fatalf("expected codex worktree to remain after cancel: %v", err)
	}

	runCLI(t, repoDir, "", "--nocolor", "--mode", "codex", "cleanup", opaqueID, "--yes")
	if _, err := os.Stat(codexPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected codex worktree removed, stat error: %v", err)
	}
	if _, err := os.Stat(classicPath); err != nil {
		t.Fatalf("expected classic worktree to remain, stat error: %v", err)
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
	t.Setenv("HOME", t.TempDir())

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

func setCodexHome(t *testing.T) string {
	t.Helper()
	codexHome := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(codexHome); err == nil {
		codexHome = resolved
	}
	t.Setenv("CODEX_HOME", codexHome)
	return codexHome
}

func addCodexWorktree(t *testing.T, repoDir, codexHome, opaqueID string) string {
	t.Helper()
	worktreesRoot := filepath.Join(codexHome, "worktrees", opaqueID)
	if err := os.MkdirAll(worktreesRoot, 0o755); err != nil {
		t.Fatalf("mkdir codex worktrees: %v", err)
	}
	worktreePath := filepath.Join(worktreesRoot, filepath.Base(repoDir))
	runGit(t, repoDir, "worktree", "add", "--detach", worktreePath)
	return worktreePath
}

func addClassicWorktree(t *testing.T, repoDir, branch string) string {
	t.Helper()
	worktreePath := filepath.Join(filepath.Dir(repoDir), filepath.Base(repoDir)+"_"+branch)
	runGit(t, repoDir, "worktree", "add", "-b", branch, worktreePath)
	return worktreePath
}

func indexOf(values []string, target string) int {
	for i, value := range values {
		if value == target {
			return i
		}
	}
	return -1
}
