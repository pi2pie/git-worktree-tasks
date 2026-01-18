---
title: "Open-source readiness improvements"
date: 2026-01-18
status: draft
agent: Codex
---

## Goal
Prepare the project for an open-source release by addressing correctness gaps, UX alignment, and documentation/CI needs ahead of the next version.

## Proposed Work
- Add LICENSE (MIT) file (completed).
- Harden `status` against stale/prunable worktrees and missing worktree paths.
- Handle empty-history repositories (no commits yet) without failing `status`.
- Align TUI task detection with CLI repo name resolution (use git common dir logic).
- Improve `create --skip-existing` messaging to reflect the actual branch in the existing worktree.
- Add integration tests covering create/list/status/finish/cleanup, with edge cases (no commits, detached HEAD, prunable entries).
- Add CI workflow for `go test ./...` and a linter.

## Rationale
- Open-source release requires clear licensing.
- Current `status` behavior can fail on missing/prunable worktrees or empty-history repos.
- TUI should be consistent with CLI task detection to avoid confusion.
- Tests and CI reduce regressions before the next version.

## Current Artifacts
- LICENSE (MIT) added.
