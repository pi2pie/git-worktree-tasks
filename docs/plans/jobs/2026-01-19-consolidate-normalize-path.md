---
title: "Consolidate duplicate path normalization"
created-date: 2026-01-19
status: completed
agent: Zed Agent
---

## Summary

- Exported `NormalizePath` in `internal/worktree/exists.go` (previously unexported `normalizePath`).
- Removed duplicate `normalizeCleanupPath` function from `cli/cleanup.go`.
- Updated `cli/cleanup.go` to use `worktree.NormalizePath`.
- Removed unused `path/filepath` import from `cli/cleanup.go`.

## Rationale

- DRY principle: both functions had identical logic for converting relative paths to absolute, cleaned paths.
- Single source of truth for path normalization reduces maintenance burden.
- Bug fixes or improvements to path handling now apply in one place.

## Related Plans

- ../plan-2026-01-18-open-source-readiness.md
