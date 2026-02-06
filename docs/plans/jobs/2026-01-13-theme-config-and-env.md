---
title: "Implement theme config and env"
date: 2026-01-13
status: completed
agent: codex
---

## Summary

- Added TOML-backed theme config resolution with env and file precedence.
- Wired config resolution into CLI theme selection logic.
- Documented configuration usage in README and recorded TOML library research.

## Changes

- `internal/config/theme.go` — resolve theme name from `GWTT_THEME`, project config, and user config.
- `cli/root.go` — respect config/env when `--theme` is not explicitly set.
- `README.md` — configuration section with precedence and examples.
- `docs/research-2026-01-13-toml-config-options.md` — TOML library options and pros/cons.
- `internal/config/theme_test.go` — tests for precedence and parsing errors.

## Rationale

Users need a default theme without always passing `--theme`. Config/env-based selection keeps behavior explicit and consistent with CLI overrides.
