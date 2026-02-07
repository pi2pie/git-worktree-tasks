---
title: "Colorized CLI output and --nocolor flag"
created-date: 2026-01-12
status: completed
agent: codex
---

## Goal

Add colorful default CLI output with a `--nocolor` escape hatch while using lipgloss styling consistently.

## Plan

- Define shared styles with a runtime color toggle.
- Apply styles to CLI messages and table output without breaking JSON/raw outputs.
- Update docs to reflect the implementation work.
