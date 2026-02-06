---
title: "Update go-install man page fallback"
created-date: 2026-01-19
status: completed
agent: codex
---

## Summary

- Updated the go-install fallback man page generation to emit both `git-worktree-tasks` and `gwtt` man pages.

## Rationale

- Ensure installs from source archives without prebuilt man pages still include the `gwtt` alias documentation.

## Files Touched

- `scripts/go-install.sh`
