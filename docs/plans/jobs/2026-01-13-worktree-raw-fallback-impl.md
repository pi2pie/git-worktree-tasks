---
title: "Implement raw output fallback"
date: 2026-01-13
status: completed
agent: codex
---

## Summary

Start implementing the raw output fallback for `list` and `status` so missing worktree paths fall back to the main worktree path.

## Work Notes

- Source plan: `docs/plans/plan-2026-01-13-worktree-raw-fallback.md`
- Goal: Align `list --output raw` and `status` path resolution with a shared fallback rule.
- Added fallback helper to resolve main worktree path when a task branch exists without a worktree.
- Updated `list` raw output and `status` to use the fallback path.
- Added unit tests for main worktree path resolution and fallback selection.
- Documented fallback behavior in README.
- Tests: `go test ./...` (reported ok).
