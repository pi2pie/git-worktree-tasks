---
title: "Skip main worktree when resolving cleanup by branch"
date: 2026-01-19
status: completed
agent: codex
---

## Summary
- skip the main checkout when matching worktrees by branch during cleanup so branch-only cleanup still works when a task branch is checked out in the root worktree.

## Rationale
- branch-first resolution can resolve to the repo root, which Git refuses to remove, blocking cleanup of the branch.
