---
title: "Fix codex raw/absolute path handling"
date: 2026-02-04
status: completed
agent: codex
---

## Summary
Adjusted codex-mode list output so:
- `--output raw` returns a relative path to `$CODEX_HOME` by default.
- `--output raw --abs` returns an absolute path (no `$CODEX_HOME` placeholder).

