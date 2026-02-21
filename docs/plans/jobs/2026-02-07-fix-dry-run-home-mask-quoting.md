---
title: "Fix dry-run masked home path quoting"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary

Updated dry-run command formatting so masked POSIX home paths remain executable when copied from output.

## What Changed

- Added dry-run-specific argument formatting in `cli/path_mask.go`.
- Kept masked POSIX home paths in expandable form (`"$HOME/..."`) instead of single-quoted literals.
- Escaped unsafe characters in the suffix after `$HOME` so only the home token expands.
- Updated/added tests in `cli/git_exec_test.go` and `cli/path_mask_test.go`.

## Why

Single-quoting the home token in dry-run output prevented shell expansion, so copy-pasted commands failed. The new formatting preserves masking while keeping commands runnable.

## Validation

- `GOCACHE=/tmp/gocache-gwtt go test ./cli -run 'TestRunGitDryRunMasksPathsByDefault|TestFormatGitCommandForDryRunWithContext' -v`
- `GOCACHE=/tmp/gocache-gwtt go test ./...`
