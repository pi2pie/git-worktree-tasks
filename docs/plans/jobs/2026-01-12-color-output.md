---
title: "Implement CLI color output and --nocolor flag"
created-date: 2026-01-12
status: completed
agent: codex
---

## Scope

- Add a global `--nocolor` flag to disable ANSI styling.
- Use lipgloss for CLI prompts, success messages, and table output.

## Notes

- Preserve JSON/raw outputs without ANSI codes.
