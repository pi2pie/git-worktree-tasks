---
title: "Align output flags and list filtering"
date: 2026-01-13
status: completed
agent: codex
---

## Summary
- Added `--skip-existing`/`--skip` to `create` and defaulted existing-worktree handling to a non-zero error.
- Added `-o` alias plus `csv` output to `list` and `status`, and `raw` output to `list` only.
- Enabled `list [task]` arg filtering with slugification to streamline single-task lookups.

## Rationale
- Make output behavior consistent and predictable across subcommands.
- Support composition with shell workflows (raw path for `cd`, CSV for tooling).

## Notes
- `list --output raw` now requires a task or branch filter and prints the first match.
