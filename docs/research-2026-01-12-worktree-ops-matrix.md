---
title: "Worktree Operations Matrix"
date: 2026-01-12
status: in-progress
agent: codex
---

## Goal
Define detailed, consistent operation behaviors for creating, merging, and cleaning up git worktrees tied to task-name branches.

## Key Findings
- A clean separation of concerns keeps the CLI predictable: (1) create worktree+branch, (2) merge into target, (3) cleanup (worktree, branch) with explicit flags.
- The shell reference implies a naming convention: `../<repo>_<task>` for worktree paths and a sanitized `<task>` for branch names. We should codify the same logic with explicit validation.
- Avoid implicit destructive actions: deletion of worktrees and branches should be explicit, or controlled with flags on `finish`.
- A `list`/`status` command is needed to show active worktrees and task mapping, with simple and detailed modes.

## Operations Matrix

### Create (worktree + branch)
- Input: `task` (required), `base` (optional, default `main`), `path` (optional).
- Behavior:
  - Sanitize branch name from task (replace non `[A-Za-z0-9_/-]` with `-`).
  - Derive worktree path: default `../<repo>_<task>` (same as shell example).
  - Validate: base branch exists locally or remotely; worktree path does not exist; branch does not already exist (unless `--reuse-branch`).
  - Command: `git worktree add -b <branch> <path> <base>`.

### Finish (merge + optional cleanup)
- Input: `task` (required), `target` (default `main`).
- Behavior:
  - Ensure clean index in target branch before merge.
  - Ensure the worktree for task is not currently checked out.
  - Checkout target branch; merge task branch.
  - Cleanup controlled by flags:
    - `--remove-worktree` (default: true) => `git worktree remove <path>`.
    - `--remove-branch` (default: true) => `git branch -d <branch>`.
  - Always `git worktree prune` when removing worktree.

### Cleanup (no merge)
- Input: `task` (required).
- Flags:
  - `--remove-worktree` (default true)
  - `--remove-branch` (default false)
- Behavior:
  - If `--remove-worktree`, run `git worktree remove <path>` then `git worktree prune`.
  - If `--remove-branch`, delete branch (`-d` by default, `-D` with `--force`).

### List / Status
- Input: none; optional filters by task or branch.
- Modes:
  - Simple: show task, branch, path, and whether the worktree is present.
  - Detailed: include base/target, last commit, and cleanliness (dirty/clean).
- Output format:
  - Default `table`.
  - `--output json` for machine-readable output.
- Behavior: read from `git worktree list --porcelain` and map paths to tasks using the fixed naming convention.

## Safety Checks
- Ensure task worktree is not the current working directory before removal.
- Provide `--dry-run` for every destructive command.
- Require confirmation (or `--yes`) when deleting branches or worktrees.

## Edge Cases
- Branch exists but worktree missing: allow `cleanup --remove-branch`.
- Worktree exists but branch deleted: allow `cleanup --remove-worktree` and print warning.
- Merge conflicts: abort `finish` and keep worktree/branch intact.
- Base/target branch doesn't exist locally: optionally fetch or error.

## Implications or Recommendations
- Provide a `validate` subcommand (or internal validator) to preflight operations.
- Surface a `status`/`list` command to show active task worktrees.

## Open Questions
- (resolved) Worktree path is fixed to `../<repo>_<task>`; task name should be slugified consistently.
- (resolved) `finish` requires a second confirmation before deleting worktree/branch; provide bypass flag for full cleanup. Support combinations like remove-worktree-only. Confirm which worktree to remove is scoped to the task, not a global prune.
- (resolved) Merge modes should support `--squash`, `--no-ff`, and `--rebase`; default is standard `git merge` behavior.

## Related Plans
- docs/plans/plan-2026-01-12-init-phase.md
