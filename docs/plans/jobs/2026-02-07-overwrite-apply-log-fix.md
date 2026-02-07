---
title: "Fix overwrite/apply conflict logging and fallback behavior"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary

- Fixed misleading failure-style output for expected `apply` conflict guardrails by returning a dedicated sentinel (`errApplyBlocked`) and handling it as warning-level flow in root command execution.
- Adjusted overwrite behavior to skip `git apply --check` (check remains for non-destructive `apply`).
- Added overwrite fallback sync for tracked changes when direct `git apply` fails after reset/clean, so overwrite can still complete when patch application is brittle but file-level sync is possible.
- Updated dry-run action rendering and tests to reflect the overwrite flow changes.

## User-Visible Outcome

- `gwtt --mode codex apply ...` conflict hints no longer produce an additional scary red error line for expected non-destructive blocking.
- `gwtt --mode codex apply ... --force` / `gwtt --mode codex overwrite ...` no longer fails early on `apply patch check failed`.
- In overwrite mode, if `git apply` fails but fallback sync succeeds, command completes and emits a warning instead of an error exit.

## Files Updated

- `cli/root.go`
- `cli/apply_command.go`
- `cli/apply_transfer.go`
- `cli/apply_test.go`

## Verification

- `go test ./cli`
- `go test ./...`

