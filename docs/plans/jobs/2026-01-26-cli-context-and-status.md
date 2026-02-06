---
title: "CLI context propagation and status parsing"
created-date: 2026-01-26
status: completed
agent: codex
---

## Summary

Implemented context propagation for CLI git operations, improved dry-run/error command formatting, and made status ahead/behind parsing explicit with tests.

## Changes

- Threaded `cmd.Context()` through CLI commands and `runGit`.
- Added `formatGitCommand`/`shellQuote` helpers for readable command output.
- Added explicit parse errors for `rev-list --count` output.
- Added tests for git command formatting and status parsing errors.

## Files Touched

- cli/cleanup.go
- cli/common.go
- cli/common_test.go
- cli/create.go
- cli/finish.go
- cli/list.go
- cli/status.go
- internal/worktree/status.go
- internal/worktree/status_test.go

## Related Plan

- docs/plans/plan-2026-01-26-go-cli-context-and-status.md
