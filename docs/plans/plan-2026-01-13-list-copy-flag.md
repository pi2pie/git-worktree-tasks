---
title: "List copy flag"
created-date: 2026-01-13
status: completed
agent: codex
---

# Goal

Define a copy-to-clipboard flag for `list`/`ls` that mirrors `create --copy-cd`, covering task worktree paths and/or branch names.

# Context

- `create --copy-cd` copies a `cd <path>` command to the clipboard.
- `list --output raw` prints a single path when filtered (task/branch).

# Proposed UX (draft)

- No clipboard flags; rely on shell piping (e.g., `command | pbcopy`) with `--output raw`.
- Require a specific task or branch filter when using `--output raw` (already enforced).
- Ensure `--output raw` continues to return a single resolved path for the first match.
- Deprecate `create --copy-cd` (redundant with piping) once an alternative is documented.
- Add a field selector for raw output so users can copy task/branch names for cleanup workflows.

# Implementation Plan

1. Confirm `--output raw` requires a specific task or branch name (already enforced).
2. Add `--field`/`-f` (path|task|branch) for `--output raw`, defaulting to `path`.
3. Treat empty `--field=""` as default `path`.
4. Keep first-match behavior for `--output raw` when multiple rows match.
5. Clarify behavior in docs/help: recommend piping `--output raw --field task` to clipboard tools.

# Open Questions

- `--field` applies only to `--output raw` and is ignored for other formats.

# Completion Notes

- Implemented `--field/-f` for `list --output raw` (path/task/branch, default path).
- Removed `create --copy-cd` and `list --task`; updated README to prefer piping raw output.
