---
title: "Implement codex apply Phase 4 and Phase 5"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary

- Implemented structured dry-run output for codex handoff commands with:
  - operation header (`plan`)
  - transfer preflight summary (`preflight`)
  - ordered action list (`actions`) including destructive markers for overwrite
- Updated apply conflict messages to be explicit, non-destructive, and action-oriented:
  - shows why apply was blocked
  - prints concrete next-step command for `overwrite --to ...`
  - includes `--yes` hint for confirmation bypass
- Added tests covering dry-run schema and next-step guidance.
- Updated README command semantics for `apply`/`overwrite`.
- Regenerated man pages, including `gwtt_overwrite(1)` and `git-worktree-tasks_overwrite(1)`.

## Files Updated

- `cli/apply.go`
- `cli/apply_test.go`
- `cli/integration_test.go`
- `README.md`
- `man/man1/gwtt_apply.1`
- `man/man1/git-worktree-tasks_apply.1`
- `man/man1/gwtt_overwrite.1`
- `man/man1/git-worktree-tasks_overwrite.1`
- `man/man1/gwtt.1`
- `man/man1/git-worktree-tasks.1`
- `docs/plans/plan-2026-02-07-codex-apply-two-way-directions.md`
- `docs/research-2026-02-07-codex-apply-direction-and-source-checkout.md`

## Verification

- `GOCACHE=/tmp/go-build go test ./...`
- `GOCACHE=/tmp/go-build make man`

## Related Plans

- `docs/plans/plan-2026-02-07-codex-apply-two-way-directions.md`

## Related Research

- `docs/research-2026-02-07-codex-apply-direction-and-source-checkout.md`
