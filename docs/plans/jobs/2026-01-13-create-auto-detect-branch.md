---
title: "Auto-detect existing branch on create"
created-date: 2026-01-13
status: completed
agent: codex
---

## Summary

- `create` now detects an existing local branch for the task and reuses it when adding the worktree.
- Updated README to document the behavior.

## Rationale

- Simplify workflow when a task branch already exists but no worktree is present.
