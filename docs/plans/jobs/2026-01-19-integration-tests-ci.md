---
title: "Integration tests and CI workflow"
date: 2026-01-19
status: completed
agent: Codex
---

## Summary
- Added integration tests for `create`, `list`, `status`, `finish`, and `cleanup`, covering no-commit repos, detached HEAD with explicit `--base`, and prunable worktree entries.
- Added a CI workflow to run `go test ./...` and `golangci-lint`.

## Details
- Tests live in `cli/integration_test.go` and use temporary repos with real git worktree operations.
- CI workflow added at `.github/workflows/ci.yml`.

## Related Plans
- ../plan-2026-01-18-open-source-readiness.md
