---
title: "Dry-run path masking for sensitive local paths"
created-date: 2026-02-07
modified-date: 2026-02-07
status: completed
agent: codex
---

## Goal

Reduce accidental disclosure of user-identifying local paths in `--dry-run` output by masking home-directory prefixes (for example, `/Users/alice/...` -> `$HOME/...` on POSIX, `C:\Users\Alice\...` -> `%USERPROFILE%\...` on Windows) while keeping output readable and executable.

## Scope

- Add a configurable masking behavior for `--dry-run` output.
- Default masking to enabled.
- Apply masking consistently across all current dry-run-enabled commands:
  - `apply`
  - `overwrite`
  - `create`
  - `cleanup`
  - `finish`
- Keep non-dry-run behavior unchanged.

## Proposed Config

- New config section: `[dry_run]`
- New key: `mask_sensitive_paths = true`
- Default value in code: `true`
- CLI behavior:
  - When `true`, paths under user home are rendered with a platform-specific home token in dry-run output:
    - POSIX: `$HOME`
    - Windows: `%USERPROFILE%`
  - When `false`, dry-run output keeps current raw absolute paths.

## Design Notes

- Helper placement: `cli` package (presentation concern only).
- Suggested helper file: `cli/path_mask.go`.
- Suggested API shape:
  - `maskHomePath(path string) string`
  - `formatGitCommandForDryRun(args []string, mask bool) string`
- Matching rules:
  - Exact home path maps to platform home token (`$HOME` or `%USERPROFILE%`).
  - Descendants map to `<home-token>/<relative>` on POSIX and `<home-token>\<relative>` on Windows.
  - Prefix-safe matching only (separator-aware).
  - On Windows, matching is case-insensitive.
  - If home lookup fails, path is left unchanged.

## Implementation Plan

1. Extend config structs/loader/defaults with `DryRun.MaskSensitivePaths`.
2. Add path masking helper(s) in `cli`.
3. Route dry-run command rendering through masking-aware formatter:
   - `cli/git_exec.go`
   - dry-run print in `cli/create.go`
4. Mask path fields in codex transfer dry-run plan/actions:
   - `cli/apply_transfer.go`
5. Add/adjust tests:
   - Unit tests for path masking edge cases.
   - OS-specific masking tests for POSIX and Windows token behavior.
   - Update existing dry-run output tests to assert masking behavior.
   - Add config loading tests for default and explicit false.
6. Update docs:
   - README option docs for new config.
   - Config schema docs at `docs/schemas/config-gwtt.md`.
   - Regenerate/update man page docs via `make man`.

## Non-Goals

- No masking changes for non-dry-run command output.
- No broad secret redaction framework (tokens, hostnames, etc.) in this change.
- No environment variable override for this feature in this phase.

## Acceptance Criteria

- With default config, dry-run output does not expose the local username via home-absolute paths.
- Dry-run output uses `$HOME` on POSIX and `%USERPROFILE%` on Windows for home-path masking.
- Dry-run output remains copy/paste-usable for the platform's common shell conventions.
- Setting `[dry_run].mask_sensitive_paths = false` restores current raw-path dry-run output.
- Existing command behavior outside dry-run remains unchanged.

## Risks / Notes

- Snapshot/integration tests that currently assert raw absolute paths will require updates.
- Home-token replacement should be limited to path arguments and displayed path fields to avoid over-masking unrelated string content.
- Cross-platform behavior should rely on deterministic helper tests to reduce dependence on running CI on every OS for basic masking validation.
