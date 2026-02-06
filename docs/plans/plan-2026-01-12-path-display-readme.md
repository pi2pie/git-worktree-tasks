---
title: "Relative Paths Flag and README"
date: 2026-01-12
status: completed
agent: codex
---

## Goal

Default list/status output to relative paths with a flag to show absolute paths, add a root README describing the project, and introduce unit tests for key logic.

## Scope

- Add a CLI flag (name TBD) to toggle absolute vs relative path display.
- Default to relative paths in list/status outputs.
- Add a root-level README.md with a concise project overview, install/run basics, and command summary.
- Add unit tests for core helpers (path derivation, slugification, and list/status mapping as appropriate).

## Out of Scope

- Changing worktree creation path behavior.
- Full command reference or extensive docs overhaul.
- End-to-end CLI tests or integration tests.

## Related Research

- docs/research-2026-01-12-worktree-ops-matrix.md

## Plan

1. Review current list/status output formatting and path rendering.
2. Decide flag name, placement, and default behavior; update list/status output accordingly.
3. Draft README.md content and add at repo root.
4. Identify core helpers to test; add unit tests with focused fixtures.
5. Verify CLI output for both relative and absolute paths and update docs if needed.

## Notes

- Keep relative paths anchored to repo root for readability.
