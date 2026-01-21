---
title: "Release workflow improvements"
date: 2026-01-21
status: completed
agent: codex
---

## Summary
- Upgraded GitHub Actions versions in the release workflow to current major pins.
- Added logic to compute the previous stable tag for stable releases so changelogs include all pre-release history.
- Grouped changelog entries in GoReleaser output for clearer release notes.
- Added tag-branch enforcement and optional shallow-since fetching for light checkouts.

## Rationale
- Stable releases should aggregate changes since the last stable tag, even when multiple pre-releases occurred.
- Grouped release notes make it easier to scan and align with expected formatting.

## Files Touched
- `.github/workflows/release.yml`
- `.goreleaser.yml`
