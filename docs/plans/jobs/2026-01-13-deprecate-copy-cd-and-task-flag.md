---
title: "Deprecate copy-cd and list --task"
date: 2026-01-13
status: completed
agent: codex
---

# Goal
Deprecate redundant CLI flags: `create --copy-cd` and `list --task`.

# Rationale
- Clipboard copy is better served via shell piping of raw output.
- `list <arg>` already covers task filtering; `--task` is redundant.

# Notes
- No user-facing deprecation notice requested at this time.
- Commit message will be added separately by the user.
 - Deprecations were fully removed in implementation.

# Related Plans
- docs/plans/plan-2026-01-13-list-copy-flag.md

# Related Jobs
- docs/plans/jobs/2026-01-13-remove-copy-cd-and-list-task.md
