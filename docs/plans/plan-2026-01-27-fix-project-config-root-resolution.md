---
title: "Fix project config root resolution"
date: 2026-01-27
status: completed
agent: codex
modified-date: 2026-01-27
---

## Goal

Ensure project config (`gwtt.config.toml`) is resolved from the repo root even when `gwtt` is run from a subdirectory, preserving documented precedence and defaults.

## Scope

- Update config loader to resolve project config relative to the Git repo root (or upward search) instead of only `os.Getwd()`.
- Keep precedence order consistent with existing docs and behavior.
- Add/adjust tests to cover subdirectory invocation.

## Plan

1. Identify current config resolution flow and where `os.Getwd()` anchors project config paths.
2. Reuse existing repo-root discovery logic (or add a shared helper) to resolve project config from the repo root or an upward search.
3. Update config loader to use the new root resolution while preserving precedence and error handling.
4. Add tests for running from subdirectories and ensure project config is applied.
5. Update documentation if needed to clarify repo-root resolution rules.

## Success Criteria

- Running `gwtt` from any subdirectory in a repo applies `gwtt.config.toml` located at the repo root.
- Config precedence matches documented behavior.
- Tests cover subdirectory execution and pass.
