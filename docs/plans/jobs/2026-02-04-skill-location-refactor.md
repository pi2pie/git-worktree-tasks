---
title: "Codex Skill Path Refactor"
date: 2026-02-04
status: completed
agent: codex
---

## Summary
Documented the refactor to align skill paths with Codex v0.94.0, moving references from `.codex/skills` to `.agents/skills`.

## Changes
- Updated skill location references to the new `.agents/skills` path.
- Noted the change in repository documentation and context.

## Rationale
Codex v0.94.0 now expects skills under `.agents/skills`, so references must match the new location to avoid broken lookups.

## References
- [Codex v0.94.0 release](https://github.com/openai/codex/releases/tag/rust-v0.94.0)
