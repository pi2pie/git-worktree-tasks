---
title: "CI/CD Release Flow with Pre-Release Gate"
date: 2026-01-13
status: draft
agent: Codex
---

## Goal
Design a GitHub Actions release pipeline that adds a pre-release verification stage before publishing, supports manual invocation, and aligns Go tooling with the repository’s Go version.

## Scope
- Release workflow structure (jobs, triggers, permissions, prerelease detection)
- Pre-release checks (lint/test/vet, snapshot build)
- Manual trigger inputs and behavior
- Local developer commands for lint/test/vet

## Plan
1. Review the current release workflow and identify required changes (permissions, Go version, prerelease logic, job separation).
2. Define the pre-release verification job (lint/test/vet and GoReleaser snapshot build).
3. Define the publish job (GoReleaser release) and prerelease tagging logic.
4. Add manual workflow_dispatch inputs and ref handling for manual runs.
5. Document local lint/test/vet commands and how they mirror CI.

## Decisions Needed
- Whether to add golangci-lint config (`.golangci.yml`) and which checks to enable.
- Which GoReleaser config file to use (default `.goreleaser.yml` or a specific file).
- Whether to require tests for release (fail release if tests fail).

## Related Research
None.
