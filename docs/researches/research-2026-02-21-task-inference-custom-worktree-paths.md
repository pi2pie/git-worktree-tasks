---
title: "Task inference for custom worktree paths"
created-date: 2026-02-21
modified-date: 2026-02-21
status: completed
agent: codex
---

## Goal

Evaluate a safe, backward-compatible way to support task lookup when worktree paths do not follow the default `<repo>_<task>` naming convention (for example, `.worktrees/<task>`), while preserving current `gwtt` shell workflows.

## Milestone Goal

Define a concrete approach that can be implemented in a small patch for the next release without changing the `create` command contract.

## Key Findings

- Classic-mode task discovery is path-derived only: `list`/`status` call `TaskFromPath(repo, wt.Path)`, and custom path layouts therefore produce `task = "-"`. [^cli-list] [^cli-status] [^worktree-naming]
- `list <task> -o raw` only resolves rows by task match in classic mode; when no row matches, it falls back to the main worktree path if the branch exists. This makes task lookups on custom paths return `"."` instead of the target worktree path. [^cli-list] [^cli-common] [^readme-fallback]
- The behavior is reproducible with `create --path ./.worktrees/new-task` and is tracked as issue `#23`. [^issue-23]
- Current docs already acknowledge that path templating can break `TaskFromPath` discovery unless `{task}` is preserved, so this is a known design boundary. [^plan-extensible]

## Option Analysis

### Option A: Keep current behavior and require `--branch`

- Behavior:
  - Keep task inference path-only.
  - Users must use `gwtt list --branch <name> -o raw` for custom layouts.
- Pros:
  - Zero code risk.
  - No output behavior changes.
- Cons:
  - `list <task>` behaves inconsistently across path layouts.
  - Shell flows like `cd "$(gwtt list <task> -o raw)"` are brittle for `.worktrees/*`.

### Option B: Add branch-backed task inference fallback (recommended)

- Behavior:
  1. Keep existing `<repo>_<task>` parse as first priority.
  2. If parse fails and row has `refs/heads/<branch>`, infer task from that row’s branch name.
  3. Preserve `task = "-"` for detached rows.
  4. Preserve `task = "-"` for the main worktree row (repo root) to avoid changing established main-row semantics.
- Pros:
  - Fixes `list <task>` and `status <task>` for custom path layouts.
  - Keeps compatibility for existing default naming users.
  - No new persisted state.
- Risks:
  - More rows become task-searchable, which may change first-match results in edge cases with similar names.
  - Branch names containing unusual characters may interact with current task-query slugification.
- Mitigations:
  - Keep strict mode behavior unchanged for path-derived tasks.
  - Add tests for branch inference and custom-path lookup to lock expected behavior.

### Option C: Persist explicit task metadata at create-time

- Behavior:
  - Store task identity in worktree-local metadata and read it during list/status.
- Pros:
  - Most explicit and layout-independent model.
- Cons:
  - Does not solve externally created worktrees.
  - Adds storage/read complexity and migration behavior.
  - Larger scope than needed for issue `#23`.

## Recommendation

Implement Option B with two compatibility guardrails:

1. Do not infer task for the main worktree row.
2. Keep branch filter (`--branch`) behavior unchanged and authoritative.

This addresses the problem with minimal scope and preserves existing workflows.

## Suggested Implementation Outline

1. Add a helper in `cli` to derive classic-mode task for a row:
   - Try `TaskFromPath`.
   - If no match, use branch-backed fallback with:
     - task normalization via `SlugifyTask(branch)` for matching consistency
     - path-based main-worktree exclusion (`repoRoot`) so the primary checkout still renders as `-`.
2. Use that helper in both `list` and `status` to keep behavior aligned.
3. Add unit/integration tests for:
   - custom path `.worktrees/new-task` + branch `new-task`
   - `list new-task -o raw` returns custom worktree path
   - `status new-task` resolves custom worktree row
   - main worktree remains `task = "-"`.
4. Optional follow-up in the same patch:
   - Fix raw fallback honoring `--field` (`path|task|branch`) when no row matches and fallback branch is used.

## Decision Notes

- Query behavior remains unchanged:
  - default (without `--strict`) stays fuzzy/contains matching
  - `--strict` remains exact matching against normalized task values.
- For inferred branch tasks, do not add extra raw-text matching in strict mode. Keep strict deterministic by matching normalized values only.
- Main-worktree exclusion is path-based only (`repoRoot`), not branch-name based. Branch-name checks can misclassify legitimate non-main worktrees that happen to use names like `main` or `master`.

## References

[^issue-23]: `https://github.com/pi2pie/git-worktree-tasks/issues/23`

[^cli-list]: `cli/list.go`

[^cli-status]: `cli/status.go`

[^cli-common]: `cli/common.go`

[^worktree-naming]: `internal/worktree/naming.go`

[^readme-fallback]: `README.md`

[^plan-extensible]: `docs/plans/plan-2026-01-27-extensible-config.md`

## Related Plans

- `docs/plans/plan-2026-01-13-worktree-raw-fallback.md`
- `docs/plans/plan-2026-01-27-extensible-config.md`
