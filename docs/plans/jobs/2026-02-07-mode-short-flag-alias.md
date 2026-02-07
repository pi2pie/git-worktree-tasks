---
title: "Add -m shorthand alias for --mode"
created-date: 2026-02-07
status: completed
agent: codex
---

## Summary

- Added persistent shorthand `-m` for the global `--mode` flag.
- Added mode precedence test coverage for short flag usage (`-m` overriding env/config).
- Updated README global flag references to include `-m`.
- Regenerated man pages so help output documents the short alias.

## Files Updated

- `cli/root.go`
- `cli/mode_test.go`
- `README.md`
- `man/man1/*` (regenerated)

## Verification

- `GOCACHE=/tmp/go-build go test ./...`
- `GOCACHE=/tmp/go-build make man`
