---
title: "Implement codex apply/overwrite Phase 1 and Phase 3"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary

- Locked Phase 1 command/UX spec in implementation and documentation:
  - `apply --to local|worktree` (default: `local`)
  - `overwrite --to local|worktree` as a peer command
  - `apply --force` as compatibility alias to overwrite behavior
- Implemented Phase 3 behavior in codex mode:
  - direction-aware transfer for both `local` and `worktree` destinations
  - explicit `overwrite` command with confirmation gating and `--yes` bypass
  - removed implicit overwrite fallback from `apply` conflict flow
  - conflict handling now exits with guidance to rerun `overwrite --to ...`
- Added/updated tests for direction handling, overwrite confirmation, and apply conflict behavior.
- Updated plan/research documents to reflect Phase 1 and Phase 3 completion/decisions.

## Files Updated

- `cli/apply.go`
- `cli/root.go`
- `cli/apply_test.go`
- `cli/integration_test.go`
- `docs/plans/plan-2026-02-07-codex-apply-two-way-directions.md`
- `docs/research-2026-02-07-codex-apply-direction-and-source-checkout.md`

## Verification

- `GOCACHE=/tmp/go-build go test ./...`

## Related Plans

- `docs/plans/plan-2026-02-07-codex-apply-two-way-directions.md`

## Related Research

- `docs/research-2026-02-07-codex-apply-direction-and-source-checkout.md`
