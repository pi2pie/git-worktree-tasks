---
title: "Themes and responsive tables"
created-date: 2026-01-12
status: completed
agent: codex
---

## Goal

Deliver multiple selectable themes and improve CLI/TUI table responsiveness for narrow terminals.

## Scope

- Add named themes with role-based colors and a `--theme` selector.
- Add `--themes` to print available themes.
- Make CLI tables width-aware and optionally render grid borders.
- Provide a TUI list view using `bubbles/table` with resize-aware column widths.

## Out of Scope

- Config file, env var, or auto-detect theme selection.
- User-defined/custom theme definitions.

## Related Research

- ../research-2026-01-12-themes-and-table.md
