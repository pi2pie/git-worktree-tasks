---
title: "Implement classic task inference for custom paths (phase 1-4)"
created-date: 2026-02-21
status: completed
agent: codex
---

## Summary

- Added shared classic-mode task derivation that keeps path-based parsing first and falls back to branch-based inference for custom worktree paths.
- Kept main-worktree exclusion path-based (`repoRoot`) and left detached rows as `-`.
- Applied the shared derivation to both `list` and `status` for consistent behavior.
- Patched raw fallback output to honor `--field` when no worktree row matches but branch fallback is available.

## Changes

- `cli/common.go`
  - Added `deriveClassicTask(repoRoot, repo, wt)` helper.
- `cli/list.go`
  - Replaced direct `TaskFromPath` use with `deriveClassicTask`.
  - Updated raw fallback to output selected field (`path|task|branch`) via synthetic fallback row.
- `cli/status.go`
  - Replaced direct `TaskFromPath` use with `deriveClassicTask`.
- `cli/common_test.go`
  - Added `TestDeriveClassicTask` coverage for path-first derivation, branch fallback, slugification, main-path exclusion, and detached behavior.
- `cli/integration_test.go`
  - Added custom-path inference integration test covering `list/status` query behavior and strict/fuzzy expectations.
  - Added raw-fallback field integration test (`path`, `branch`, `task`).

## Verification

- `GOCACHE=/tmp/gocache-gwtt go test ./cli -run 'TestDeriveClassicTask|TestIntegrationListStatusCustomPathTaskInference|TestIntegrationListRawFallbackHonorsField' -v`
- `GOCACHE=/tmp/gocache-gwtt go test ./...`

## Related Plan

- `docs/plans/plan-2026-02-21-classic-task-inference-custom-paths.md`
