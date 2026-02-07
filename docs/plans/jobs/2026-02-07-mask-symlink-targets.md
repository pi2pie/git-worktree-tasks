---
title: "Mask symlink targets in dry-run output"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary
- masked symlink targets in dry-run output when sensitive path masking is enabled

## Rationale
- prevent leaking sensitive symlink targets in logs while `--mask-sensitive-paths` is active

## Result
- dry-run symlink output now uses the same masking helper for the target path
