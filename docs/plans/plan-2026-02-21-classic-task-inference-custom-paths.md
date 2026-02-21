---
title: "Classic task inference for custom worktree paths"
created-date: 2026-02-21
status: active
agent: codex
---

## Context

Issue `#23` shows that `list <task> -o raw` can return the main worktree path when the target worktree uses a custom path layout (for example, `.claude/worktrees/<task>`) instead of the default `<repo>_<task>` naming convention.

Current classic-mode behavior derives task names from path only. When derivation fails, task becomes `-`, task filtering misses the row, and raw mode falls back to the main worktree path.

## Goal

Make `list` and `status` task lookup reliable for custom worktree paths in classic mode, without breaking default naming behavior or existing `--branch` workflows.

## Non-Goals

- Redesigning `create` path behavior in this phase.
- Introducing persistent task metadata storage.
- Changing codex-mode task lookup behavior.

## Decisions

- Keep current query semantics:
  - default remains fuzzy (`contains`) when `--strict` is not set
  - `--strict` remains exact matching on normalized task values
- For inferred branch tasks, do not add separate raw-text strict matching.
- Exclude main-worktree inference by path only (`repoRoot`), not by branch name (`main`, `master`).

## Proposed Behavior

1. In classic mode, task derivation for each worktree row should be:
   - first: existing `TaskFromPath(repo, wt.Path)` logic
   - fallback: derive from the row branch name when available
2. Main worktree row (repo root path) should stay `task = "-"`.
3. Detached rows should stay `task = "-"`.
4. `list <task>` and `status <task>` should use the same derivation behavior.

## Implementation Plan

### Phase 1: Shared Task Derivation Helper

- [ ] Add a helper in `cli` for classic-mode row task derivation.
- [ ] Keep current path-based extraction as first priority.
- [ ] Add branch-backed fallback for non-main, non-detached rows.
- [ ] Normalize inferred branch task values consistently for matching.

### Phase 2: Command Integration

- [ ] Replace direct `TaskFromPath(...)` usage in `cli/list.go` with the helper.
- [ ] Replace direct `TaskFromPath(...)` usage in `cli/status.go` with the same helper.
- [ ] Preserve current codex-mode behavior unchanged.

### Phase 3: Verification Tests

- [ ] Add/extend tests to cover custom-path worktree + branch inference.
- [ ] Verify `list new-task -o raw` resolves to custom worktree path.
- [ ] Verify `status new-task` resolves custom worktree row.
- [ ] Verify main worktree row remains `task = "-"`.
- [ ] Verify fuzzy/default and `--strict` semantics remain unchanged.

### Phase 4: Optional Follow-up Fix

- [ ] Evaluate and, if approved, patch raw fallback to honor `--field` when no rows match and fallback branch is used.

### Phase 5: Documentation

- [ ] Update `README.md` to reflect classic-mode task inference fallback for custom worktree paths.
- [ ] Update `README.md` examples for `list <task>` / `status <task>` to align with post-enhancement behavior.
- [ ] Clarify in `README.md` that `--branch` remains explicit/authoritative filtering.
- [ ] Review and revise `docs/schemas/config-gwtt.md` wording around `create.path.format` and task discovery constraints.
- [ ] Add/refresh a job record under `docs/plans/jobs/` when implementation starts.

## Acceptance Criteria

- `list <task>` works for classic-mode worktrees whose paths do not match `<repo>_<task>` when the worktree has a branch.
- `list <task> -o raw` returns the target custom worktree path (not main fallback) when that worktree exists.
- `status <task>` resolves the same task consistently.
- Main worktree still displays `task = "-"`.
- Existing default-path behavior remains unchanged.

## Risks / Notes

- Branch-backed inference increases the number of rows that are task-searchable, so first-match outcomes should be validated with tests.
- Normalization and matching rules must stay aligned with current `normalizeTaskQuery` and `matchesTask` behavior to avoid subtle regressions.

## Related Research

- `docs/researches/research-2026-02-21-task-inference-custom-worktree-paths.md`
