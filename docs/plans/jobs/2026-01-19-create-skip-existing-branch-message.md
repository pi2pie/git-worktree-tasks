---
title: "Create skip-existing branch messaging"
date: 2026-01-19
status: completed
agent: Codex
---

## Summary
- Updated `create --skip-existing` output to show the actual branch associated with the existing worktree.
- Added a worktree lookup helper to resolve the branch by path, falling back to "detached" or the task name.

## Rationale
- The previous message always echoed the task name, which could mislead when the worktree was on a different branch.
