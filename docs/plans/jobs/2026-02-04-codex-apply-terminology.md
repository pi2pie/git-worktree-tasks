---
title: "Codex apply terminology updates"
created-date: 2026-02-04
status: completed
agent: codex
---

## Summary

- Updated docs to align Codex App UI terminology with "Hand off changes" and clarified that app-side worktree/shell issues are out of CLI scope.
- Renamed codex-mode command references from `sync` to `apply` across research, plan, and config docs.
- Renamed CLI implementation from `sync` to `apply` (including file rename, symbols, and user-facing messages).

## Why

- Codex App UI now uses "Hand off changes" with directions "To local" / "From local", while official docs still say "Sync with local".
- Avoid confusion with "apply" terminology in the Codex CLI by making the gwtt command name explicit and consistent with the app wording.
- Track app-side worktree issues separately from CLI responsibilities.

## Files Updated

- `docs/researches/research-2026-02-04-mode-classic-vs-codex.md`
- `docs/plans/plan-2026-02-04-mode-classic-and-codex.md`
- `docs/schemas/config-gwtt.md`
- `cli/apply.go`
- `cli/root.go`
- `cli/finish.go`
