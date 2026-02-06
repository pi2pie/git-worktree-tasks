---
title: "Empty-history handling and friendly errors"
date: 2026-01-19
status: completed
agent: Codex
---

## What changed

- Added Git stderr classification for friendlier "not a repo" and "no commits yet" errors in core helpers.
- Made short-hash lookup tolerant of empty-history repos (returns default length without error).
- Labeled status output as "empty history" for last commit/base when HEAD is missing.
- Standardized empty-history behavior across subcommands by adding current-branch checks where needed.

## Why

- Avoid exposing raw Git stderr while keeping users informed about missing history.
- Keep CLI behavior consistent when repositories have no commits.

## Current behavior

- When run outside a Git repo, commands return: "not a git repository (run inside a git repository)".
- In empty-history repos, `status`, `list`, `finish`, `create`, and `cleanup` all error with: "no commits yet (empty history)".
- `status` labels `LastCommit`/`Base` as "empty history" only if it reaches worktree evaluation.
