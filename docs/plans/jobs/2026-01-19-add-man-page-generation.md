---
title: "Add man(1) generation and install wiring"
created-date: 2026-01-19
status: completed
agent: Codex
---

## Overview

Added Cobra-based man(1) generation and wired it into the Go install workflow so man pages are built and installed alongside binaries.

## Changes

- Added a man page generator (`scripts/generate-man.go`) using Cobra doc helpers (run twice for `git-worktree-tasks` and `gwtt`).
- Generated and committed man pages under `man/man1` for installation packaging.
- Exposed a `cli.RootCommand()` helper for documentation tooling.
- Added `make man` to generate man pages into `man/man1`.
- Updated `scripts/go-install.sh` to build and install man pages to `<install_prefix>/share/man/man1` (override via `MAN_DIR`).
- Updated `scripts/go-uninstall.sh` to remove installed man pages.

## Files Touched

- `cli/root.go`
- `scripts/generate-man.go`
- `Makefile`
- `scripts/go-install.sh`
- `scripts/go-uninstall.sh`
- `man/man1/*.1`
