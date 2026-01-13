---
title: "Worktree raw output fallback"
date: 2026-01-13
status: completed
agent: codex
---

## Context
The `list --output raw` and `status` subcommands are used in shell flows (e.g., `cd $(...)`). When a task has no worktree path, the current output is not usable and causes unexpected behavior. We need a deterministic fallback to the main worktree path so these flows are stable.

## Goals
- Define a clear fallback rule for tasks with no worktree path.
- Align `list --output raw` and `status` so they resolve paths consistently.
- Keep existing behavior unchanged for tasks with explicit worktree paths.

## Non-Goals
- Redesigning other output formats beyond `--output raw`.
- Changing worktree creation or task discovery logic.

## Proposed Behavior
- If a task has a worktree path, use it.
- If a task has no worktree path (even if the branch exists), fall back to the main worktree path.
- The main worktree path should be derived from the repo’s primary worktree/root, not inferred from the task branch.

## Work Plan
1) Inspect current `list` and `status` path resolution logic and identify where raw output is generated.
2) Add a shared fallback helper (or equivalent) to resolve task path with main worktree fallback.
3) Update `list --output raw` and `status` to use the shared fallback.
4) Add or update tests/fixtures covering:
   - task with worktree path
   - task with branch but no worktree path
   - task without branch and without worktree path
5) Document the fallback behavior in CLI docs/README if applicable.

## Acceptance Criteria
- `list --output raw` returns a usable path for all tasks.
- `status` path reporting matches the same fallback rule.
- No behavior change for tasks that already have a worktree path.

## Risks
- Repos with multiple main worktrees or non-standard default branch names.
- Detached HEAD when determining the primary worktree.

## Open Questions
- Should `--output raw` for `list` ever return blank paths when fallback is unavailable?
- How should we detect the "main worktree path" in a repo with multiple worktrees?

## Related Jobs
- `docs/plans/jobs/2026-01-13-align-output-flags.md`
- `docs/plans/jobs/2026-01-13-list-status-search.md`
