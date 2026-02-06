---
title: "Fix GoReleaser previous tag detection"
date: 2026-01-19
status: completed
agent: codex
---

## Summary

- Updated release workflow checkout settings to fetch full history and tags so GoReleaser can resolve the correct previous tag for changelog compare links.

## Rationale

- Shallow fetches of a specific tag can omit other tags, which causes GoReleaser to compute an incorrect or empty `PreviousTag` for release compare links.

## Changes

- Set `fetch-depth: 0` and `fetch-tags: true` for release workflow checkouts.
