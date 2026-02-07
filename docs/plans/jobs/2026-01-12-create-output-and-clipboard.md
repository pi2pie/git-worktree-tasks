---
title: "Create output and clipboard support"
created-date: 2026-01-12
status: completed
agent: codex
---

## Summary

- Changed create output to show relative worktree paths.
- Replaced print-cd with copy-cd flag to copy a `cd` command to clipboard.
- Added cross-platform clipboard helper with common OS commands.

## Notes

- Clipboard support uses pbcopy (macOS), clip (Windows), wl-copy or xclip (Linux).
