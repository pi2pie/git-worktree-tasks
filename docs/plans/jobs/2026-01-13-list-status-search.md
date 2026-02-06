---
title: "Add list/status search and repo base name resolution"
date: 2026-01-13
status: completed
agent: codex
---

## Summary

- Added repo base name resolution via git common dir so task derivation works from any worktree.
- Added positional `[task]` filtering for `list`/`status` with contains matching and `--strict` for exact match.
- Updated docs with new list/status filtering behavior.

## Rationale

- Ensure task names resolve consistently when running from a worktree.
- Provide zoxide-like, low-friction task search without extra dependencies.

## Notes

- `--task` remains exact and now normalizes via slugification.
- `list --output raw` still prints the first match when filters are used.
