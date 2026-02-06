---
title: "Go CLI context propagation and status robustness"
created-date: 2026-01-26
status: completed
modified-date: 2026-01-27
agent: codex
---

## Goal

Improve CLI correctness and UX by propagating contexts, making dry-run output readable, and handling status parsing errors explicitly.

## Scope

- Thread `cmd.Context()` through CLI commands and `runGit`.
- Replace slice-style git arg printing with readable, shell-like output.
- Make `worktree.Status` parsing explicit (error or warning path) when git output is unexpected.

## Issues Memo

- Context is hardcoded to `context.Background()` in CLI command handlers and `runGit`, preventing cancellation/timeouts from `cmd.Context()`.
- Dry-run and git error messages print args as Go slices, which are not copy/paste-friendly for users.
- `worktree.Status` ignores `strconv.Atoi` errors, silently mapping unexpected output to zero ahead/behind values.

## Non-Goals

- Changing command semantics or defaults beyond the above.
- Adding new CLI flags unrelated to the issues.

## Proposed Steps

1. Introduce context plumbing in CLI entry points and helper functions.
2. Implement a small helper to format git commands for dry-run and error messages.
3. Update `worktree.Status` to handle `Atoi` errors deterministically.
4. Add or update tests to cover the new behavior (as needed).

## Task Checklist

- [x] Propagate `cmd.Context()` through CLI commands and `runGit`.
- [x] Add git command formatting helper for dry-run/error output.
- [x] Handle `rev-list` parsing errors explicitly in `worktree.Status`.
- [x] Update/add tests for context usage and status parsing behavior.

## Related Jobs

- docs/plans/jobs/2026-01-26-cli-context-and-status.md

## Risks

- Minor behavior changes if callers relied on silent parsing failures.
- Any quoting changes may affect copy/paste expectations; keep it conservative.
- Context propagation should not affect theme config or theme listing, since those are handled in `PersistentPreRunE` before subcommand execution.

## Acceptance Criteria

- Git commands in dry-run and errors are readable and copy/paste-friendly.
- CLI respects cancellation via `cmd.Context()` where applicable.
- Status parsing errors are surfaced in a consistent and documented way.

## Related Research

- None.
