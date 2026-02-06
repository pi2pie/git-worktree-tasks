---
title: "Apply table grid fallback per config source"
created-date: 2026-01-27
status: completed
agent: codex
---

## Summary

- Updated config loading so table grid defaults apply per config source, allowing higher-precedence table settings to override lower-precedence list/status grids.

## Rationale

- Reviewer feedback noted that list/status grid flags could block later table defaults; the change ensures precedence is respected.
