---
title: "Fix cleanup worktree resolution"
created-date: 2026-01-19
status: completed
agent: codex
---

## Summary

- Resolve cleanup worktree path via branch-backed worktree list before deleting.
- Fall back to path matching when no branch match is found.

## Rationale

- Cleanup should delete the correct worktree even when the path was overridden or differs from the default naming convention.

## Notes

- No user-facing flags changed; behavior is stricter about finding existing worktrees by branch.
