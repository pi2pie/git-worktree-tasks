---
title: "Clarify cleanup modes and branch visibility"
date: 2026-01-12
status: completed
agent: codex
---

## Goal
Clarify expected behavior around worktree creation, cleanup modes, and branch visibility in list/status.

## Scope
- Confirm desired semantics for cleanup (default dual removal, worktree-only, branch-only).
- Define behavior when a worktree is missing but a task branch exists.
- Decide whether to surface branch info in list/status output or document the current behavior.
- Update README/CLI help to reduce confusion around worktree paths and cleanup.

## Proposed Approach
1. Inspect current CLI behavior for create/list/status/cleanup and identify gaps.
2. Specify updated cleanup flow, including confirmations and no-worktree handling.
3. Implement CLI changes (flags, prompts, outputs) and adjust tests if present.
4. Update README examples and notes to document worktree/branch visibility.

## Open Questions
- Should list/status include branch columns, or should README clarify that they only show worktrees?
- If no worktree exists but the task branch does, should cleanup with default mode remove the branch after confirmation?
- When creating a worktree, should the CLI suggest `cd <path>` or add an optional flag to auto-print a `cd` command?

## Resolution
- Added explicit branch output guidance in README (list/status already include branch columns).
- Cleanup defaults to removing both worktree and branch with a second confirmation for the branch, including missing-worktree messaging.
- Added a create flag to copy a ready-to-run `cd` command and adjusted create output to show relative paths.
- Added create raw output mode for piping the path into `cd`.
- Added create behavior to reuse existing worktree paths and a `--path` override option.
