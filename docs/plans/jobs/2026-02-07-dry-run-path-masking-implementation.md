---
title: "Implement dry-run home path masking with cross-platform tokens"
created-date: 2026-02-07
modified-date: 2026-02-07
status: completed
agent: codex
---

## Summary

- Implemented configurable dry-run masking via `[dry_run].mask_sensitive_paths` with default `true`.
- Added platform-aware home-token masking:
  - POSIX: `$HOME/...`
  - Windows: `%USERPROFILE%\\...`
- Wired masking into all dry-run command paths in scope:
  - shared `runGit` dry-run printing (covers `finish` and `cleanup`)
  - `create` dry-run command output
  - codex `apply`/`overwrite` dry-run plan fields, action commands, and copy/symlink action lines
- Kept non-dry-run behavior unchanged.
- Added direct overrides for masking behavior:
  - CLI: `--mask-sensitive-paths[=true|false]`
  - CLI: `--no-mask-sensitive-paths`
  - ENV: `GWTT_DRY_RUN_MASK_SENSITIVE_PATHS`
- Follow-up lint cleanup: removed unused `maskHomePath` wrapper after golangci-lint `unused` finding.
- Updated README and config schema docs; regenerated man pages (`make man`).

## Files Updated

- `internal/config/config.go`
- `internal/config/config_test.go`
- `cli/path_mask.go`
- `cli/path_mask_test.go`
- `cli/git_exec.go`
- `cli/git_exec_test.go`
- `cli/root.go`
- `cli/dry_run_mask_test.go`
- `cli/create.go`
- `cli/apply_command.go`
- `cli/apply_transfer.go`
- `cli/apply_files.go`
- `cli/apply_test.go`
- `README.md`
- `docs/schemas/config-gwtt.md`

## Verification

- `go test ./...`
- `make man`
- `golangci-lint run`

## Related Plans

- `docs/plans/plan-2026-02-07-dry-run-path-masking.md`
