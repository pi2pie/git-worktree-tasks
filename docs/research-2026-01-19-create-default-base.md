---
title: "Create default base selection"
date: 2026-01-19
status: completed
agent: Codex
---

## Goal
Define how `create` selects a base branch by default when invoked from different contexts, including older Git default branches (`master`) and detached HEAD cases.

## Key Findings
- Current behavior: `create` defaults to the current branch when on a named branch.
- Detached HEAD requires an explicit `--base`.
- Worktrees always share the same Git common directory; creating a worktree from another worktree does not create a new repo, it creates another worktree in the same repo.

## Implications or Recommendations
- Prefer base = current branch when on a named branch.
- If detached HEAD, allow `create` only when `--base` is explicitly provided (flexible, explicit intent).
- `--base` always overrides defaults.

## Open Questions
- None. (Detached HEAD: allow only with explicit `--base`; default base: follow current local branch.)

## Rationale Notes
- Detached HEAD behavior options:
  - Allow explicit `--base`: `create --base dev my-task` succeeds even when detached; explicit base removes ambiguity.
  - Require checkout: `create` fails on detached HEAD even with `--base`; stricter safety, less flexible.
- Current behavior matches the explicit-`--base` rule, so detached HEAD without `--base` should produce a clear error.

## References
- None.

## Related Plans
- docs/plans/plan-2026-01-18-open-source-readiness.md
