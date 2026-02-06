---
title: "Remove create copy-cd and list --task"
date: 2026-01-13
status: completed
agent: codex
---

# Goal

Remove deprecated CLI flags `create --copy-cd` and `list --task`.

# Summary

- Removed `create --copy-cd` flag and clipboard helper.
- Removed `list --task` flag (use `list <arg>` instead).
- Updated README usage examples to prefer piping raw output.

# Related Plans

- docs/plans/plan-2026-01-13-list-copy-flag.md
