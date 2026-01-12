---
title: "Create existing worktree handling and path override"
date: 2026-01-12
status: completed
agent: codex
---

## Summary
- Create now returns existing task worktree path instead of failing.
- Added `--path/-p` to override worktree location (relative to repo root or absolute).
- Documented the new behavior and option in README.
