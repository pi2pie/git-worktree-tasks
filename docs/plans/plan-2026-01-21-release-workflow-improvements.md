---
title: "Release workflow improvements"
date: 2026-01-21
status: completed
agent: codex
---

## Goal
Improve the GitHub Actions release workflow for GoReleaser by upgrading action versions, supporting optional shallow checkouts, enforcing allowed tag branches, and ensuring final releases include cumulative pre-release changes.

## Scope
- Update action versions in `.github/workflows/release.yml` (checkout, setup-go, goreleaser action).
- Add optional `shallow_since` fetching for lighter checkouts while preserving tag history.
- Enforce tags only from `main` and `dev*/beta*/alpha*/canary*` branches.
- Adjust changelog/release behavior so stable releases include all changes since the last stable release, including pre-release history.

## Non-Goals
- Replacing GoReleaser with another release tool.
- Changing the semantic versioning scheme or tag naming conventions.
- Altering build matrix or artifact formats unless required by release notes behavior.

## Plan
1. Audit the current release workflow and GoReleaser config to identify where changelog data is sourced and how tags are selected.
2. Verify current supported versions for `actions/checkout`, `actions/setup-go`, and `goreleaser/goreleaser-action`, and note any required config changes for upgrades (e.g., permissions, default inputs).
3. Decide on shallow history strategy:
   - Use `fetch-depth: 1` and explicit tag fetches.
   - Add optional `shallow_since` for manual runs to cap history depth.
4. Design changelog behavior for formal releases:
   - Determine how to include changes from all pre-releases since the last stable tag (e.g., set `previous_tag` explicitly or adjust tag selection rules).
   - Update GoReleaser `changelog`/`release` config or add a pre-step to compute the correct comparison range.
5. Implement workflow and config changes, then verify with a dry-run or a local GoReleaser release in CI (if available).
6. Document rationale in the workflow comments or a short note in the plan job record.

## Open Questions
- Are there existing release notes expectations (format or tool output) that we must preserve beyond grouped sections?

## Success Criteria
- Workflow uses updated actions with confirmed compatibility.
- Release notes for a stable tag include all changes since the last stable tag, even if intermediate pre-releases existed.
- Changelog generation is deterministic and documented.
