---
title: "Implement codex apply Phase 6 and Phase 7"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary

- Completed Phase 6 by splitting monolithic `apply` logic into focused files:
  - command wiring and handoff flow in `cli/apply_command.go`
  - codex worktree resolution and transfer planning in `cli/apply_resolve.go`
  - conflict/preflight helpers in `cli/apply_conflicts.go`
  - transfer and dry-run action logic in `cli/apply_transfer.go`
  - patch/file copy helpers in `cli/apply_files.go`
- Removed `cli/apply.go` after migrating all behavior and preserving existing test-targeted function signatures.
- Completed Phase 7 verification and docs pass:
  - full Go test suite passed
  - `apply`/`overwrite` help text re-verified
  - README/man alignment re-checked
  - plan and research docs updated to reflect completion

## Files Updated

- `cli/apply_command.go`
- `cli/apply_resolve.go`
- `cli/apply_conflicts.go`
- `cli/apply_transfer.go`
- `cli/apply_files.go`
- `cli/apply.go` (removed)
- `docs/plans/plan-2026-02-07-codex-apply-two-way-directions.md`
- `docs/researches/research-2026-02-07-codex-apply-direction-and-source-checkout.md`

## Verification

- `go test ./...`
- `go run . --nocolor apply --help`
- `go run . --nocolor overwrite --help`

## Related Plans

- `docs/plans/plan-2026-02-07-codex-apply-two-way-directions.md`

## Related Research

- `docs/researches/research-2026-02-07-codex-apply-direction-and-source-checkout.md`
