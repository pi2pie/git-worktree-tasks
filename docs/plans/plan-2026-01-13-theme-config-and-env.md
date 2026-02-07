---
title: "Theme config and env"
created-date: 2026-01-13
status: completed
agent: codex
---

## Goal

Add config/env-driven theme selection so users can set a default theme without always passing `--theme`.

## Scope

- Define config sources, schema, and precedence for theme selection only.
- Implement config discovery for:
  - User config: `$HOME/.config/gwtt/config.toml`
  - Project config: `gwtt.config.toml` or `gwtt.toml` in the current working directory
- Add env override (e.g., `GWTT_THEME`).
- Update CLI init to merge settings with the existing `--theme` flag.
- Document the config options in `README.md`.

## Out of Scope

- Custom theme palette definitions in config.
- Other configuration options beyond theme selection.

## Related Research

- ../research-2026-01-13-toml-config-options.md

## Design Decisions

- **Precedence (highest → lowest)**: CLI flag `--theme` → `GWTT_THEME` → project config (`gwtt.config.toml`, then `gwtt.toml`) → user config → default theme.
- **Config schema (TOML)**:
  - Prefer a focused `theme` table.
  - Example:
    ```toml
    [theme]
    name = "nord"
    ```
  - If `theme.name` is empty, treat it as unset.
- **Project config resolution**:
  - Look in the current working directory only (no parent traversal) to keep behavior explicit and fast.
  - If both project files exist, `gwtt.config.toml` wins.

## Milestones

1. Add a small config package to parse TOML and resolve precedence.
2. Wire config into CLI pre-run logic.
3. Update docs and examples.
