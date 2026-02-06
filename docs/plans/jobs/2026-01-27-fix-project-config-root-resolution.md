---
title: "Fix project config root resolution"
date: 2026-01-27
status: completed
agent: codex
modified-date: 2026-01-27
---

# Fix project config root resolution

## Summary

- Resolve project config from the repo root (identified by a `.git` marker) instead of only the current working directory.
- Share repo-root resolution between the main config loader and theme resolution to keep behavior consistent.
- Add a regression test to ensure project config is applied when invoked from a subdirectory.

## Why

The config loader was anchored to `os.Getwd()`, so running `gwtt` from a repo subdirectory ignored `gwtt.config.toml` at the repo root, contradicting documented behavior and default expectations.

## Files

- `internal/config/project_root.go`
- `internal/config/config.go`
- `internal/config/theme.go`
- `internal/config/config_test.go`
- `docs/plans/plan-2026-01-27-fix-project-config-root-resolution.md`
