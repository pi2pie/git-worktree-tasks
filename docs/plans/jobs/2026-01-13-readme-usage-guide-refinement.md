---
title: "README Usage Guide Refinement"
created-date: 2026-01-13
status: completed
agent: GitHub Copilot
---

## Summary

Restructured and refined the README to provide a clearer usage guide with improved documentation for output formats and piping workflows.

## Changes

### Structure

- Added **Quick Start** section at the top for immediate onboarding
- Added **Table of Contents** for easy navigation
- Reorganized content into logical sections with clear headings

### Usage Documentation

- Created **Commands Overview** table showing all commands with aliases
- Added per-command documentation with:
  - Code examples for common use cases
  - Flag tables with short forms and descriptions
  - Clear separation between `create`, `list`, `status`, `finish`, `cleanup`

### Output Formats & Piping (New Section)

Documented the `--output` (`-o`) and `--field` (`-f`) flag combinations:

| Format  | Description                    | Available In     |
| ------- | ------------------------------ | ---------------- |
| `table` | Human-readable table (default) | `list`, `status` |
| `json`  | JSON array                     | `list`, `status` |
| `csv`   | CSV with headers               | `list`, `status` |
| `raw`   | Single value, no decoration    | `create`, `list` |
| `text`  | Styled text output (default)   | `create`         |

Added practical piping examples:

- Navigate to worktree: `cd "$(gwtt list my-task -o raw)"`
- Copy to clipboard: `gwtt list my-task -o raw -f task | pbcopy`
- Open in editor: `code "$(gwtt list my-task -o raw)"`
- JSON filtering with `jq`
- Shell function examples for Bash/Zsh/Fish

### Minor Fixes

- Corrected `examples/` description from "Usage examples" to "Example configs and shell functions"
- Condensed Installation, Shell Configuration, Development, and Troubleshooting sections
- Added Notes section with key behaviors at a glance

## Files Changed

- `README.md` — full restructure and content refinement
