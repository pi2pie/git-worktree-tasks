---
title: "CLI helper extraction for codex apply prep"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary

Completed no-behavior-change refactor work to prepare for two-way `apply`/`overwrite` implementation.

- Added shared mode/codex context resolver for CLI commands.
- Added shared codex worktree lookup helper (including list-based reuse path).
- Moved shared `runGit` execution helper out of `finish.go` into a neutral helper file.
- Updated `apply`, `cleanup`, `list`, and `status` to use extracted helpers.

## Why

- Reduce duplicated setup logic before introducing directional command behavior.
- Lower implementation risk for upcoming `apply`/`overwrite` feature work.
- Keep command files more focused and easier to evolve.

## Files Updated

- `cli/mode_context.go`
- `cli/codex_lookup.go`
- `cli/git_exec.go`
- `cli/apply.go`
- `cli/cleanup.go`
- `cli/list.go`
- `cli/status.go`
- `cli/finish.go`

## Verification

- Ran `go test ./...` successfully.

## Related Plans

- `docs/plans/plan-2026-02-07-codex-apply-two-way-directions.md`
