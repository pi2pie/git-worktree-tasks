---
title: "Dynamic short hash length"
date: 2026-01-19
status: completed
agent: Codex
---

## Summary
- Updated short hash rendering to use a dynamic 7/8/10 length based on repository size.
- Applied the dynamic length to list, TUI, and status output (merge-base + last commit).

## Rationale
- Large repositories need longer abbreviations to stay unambiguous while still keeping output compact.
- Commit count is a practical proxy for repo size and matches the intended 7/8/10 policy.

## References
- gitrevisions(7)
- git-rev-parse(1)
