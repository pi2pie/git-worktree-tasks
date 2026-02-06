---
title: "Lint and Cleanup"
created-date: 2026-01-13
status: completed
agent: GitHub Copilot
---

## Goal

Perform a code quality review to identify and fix linting issues, unused parameters, and ensure consistency in the `cli`, `tui`, `internal`, and `ui` packages.

## Review Findings

### 1. Unused Parameters in CLI

The `state *runState` parameter is passed to all subcommand constructors but is unused in their function bodies.

- **Files affected**:
  - `cli/cleanup.go`: `newCleanupCommand(state *runState)`
  - `cli/create.go`: `newCreateCommand(state *runState)`
  - `cli/finish.go`: `newFinishCommand(state *runState)`
  - `cli/list.go`: `newListCommand(state *runState)`
  - `cli/status.go`: `newStatusCommand(state *runState)`
  - `cli/tui.go`: `newTUICommand(state *runState)`

- **Observation**: The `root.go` establishes global state (flags) and applies them via `ui.SetColorEnabled` in `PersistentPreRun`. The subcommands do not currently need to access `state` directly.

- **Recommendation**: Remove the `state` parameter from these signatures to clean up the API, or rename to `_` if we anticipate future use (though removal is cleaner).

### 2. Error Handling & Output

- **Consistency**: `cli` commands generally use `cmd.OutOrStdout()` and `cmd.ErrOrStderr()` which is correct for Cobra applications.
- **Output**: No direct `fmt.Print` calls found in `cli` package (verified via grep), ensuring output capture is robust.

### 3. TUI & Configuration

- **Color Settings**: `tui.Run()` is called in `cli/tui.go`. It relies on the global `ui` package configuration for styles. Since `root.go` configures `ui` before the command runs, this implicit dependency works but hides the data flow.
- **Context**: `newTUICommand` creates a closure that calls `tui.Run()`. It does not pass `cmd.Context()`. While `tea.NewProgram` handles its own context, passing the command context might be beneficial for cancellation.

## Action Plan

### Core Cleanup

- [x] **Refactor CLI Constructors**: Remove `state *runState` parameter from all `new*Command` functions in `cli/`.
- [x] **Update Root**: Update `cli/root.go` calls to reflect the removed parameter.

### Verification

- [x] Run `go vet ./...` to ensure no regressions.
- [x] Run `go build` to verify compilation.
- [x] Manual check of `ui` command to ensure it still launches.

## Future/Optional

- Consider passing a Context to `tui.Run()` to respect Cobra's cancellation signals if needed.
