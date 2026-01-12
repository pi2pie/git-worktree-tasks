---
title: "Cleanup behavior and branch output updates"
date: 2026-01-12
status: completed
agent: codex
---

## Summary
- Updated cleanup defaults to remove both worktree and branch with separate confirmations.
- Added worktree-only cleanup mode and improved missing-worktree messaging.
- Added create flag to print a ready-to-run `cd` command.
- Clarified README usage and notes for branch output and cleanup modes.

## Why
- Align cleanup behavior with expected confirmation flow and missing-worktree scenarios.
- Make it easier to jump into new worktrees and understand branch visibility.

## Notes
- Cleanup now checks branch existence and worktree presence before prompting.
