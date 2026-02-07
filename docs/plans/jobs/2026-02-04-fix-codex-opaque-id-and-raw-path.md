---
title: "Fix codex opaque-id mapping and raw output paths"
created-date: 2026-02-04
status: completed
agent: codex
---

## Summary

Adjusted codex-mode worktree handling to match Codex App directory layout:

- Derive `<opaque-id>` from the first path segment under `$CODEX_HOME/worktrees`.
- Hide codex-owned worktrees in classic mode.
- Make `--output raw` in codex mode return a composable path relative to `$CODEX_HOME`.
- Updated sync/cleanup resolution to use the corrected opaque-id mapping.

## Notes

- Tests run with `GOCACHE=/tmp/gocache` due to sandbox cache restrictions.
