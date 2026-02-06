---
title: "Implement --mode and codex-mode commands (Phase 2)"
date: 2026-02-04
status: completed
agent: codex
---

## Summary

Implemented Phase 2 of `classic` vs `codex` mode support:

- Added global `--mode` flag and mode normalization/validation.
- Added config/env support for `mode` (`GWTT_MODE`, `mode = "..."` in TOML).
- Implemented codex-mode behavior for `list`, `status` (repo-scoped, `$CODEX_HOME` path display), and added `modified_time` to status output.
- Added `sync` command for codex mode (`apply` by default; prompts to overwrite on conflicts).
- Updated `cleanup` to support codex-mode worktree removal under `$CODEX_HOME/worktrees/<opaque-id>` with additional warnings/confirmation.

## Notes

- Tests were executed with `GOCACHE` pointed at a writable location due to sandbox constraints.
