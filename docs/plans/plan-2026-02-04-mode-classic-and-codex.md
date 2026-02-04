---
title: "Mode flag: classic and codex"
date: 2026-02-04
status: active
agent: codex
---

## Goal
Add a new global `--mode` (`classic` default, `codex` optional) to support Codex App-style worktrees while keeping current behavior stable and non-breaking.

## Scope
- `--mode` flag + `GWTT_MODE` env + config default (flag > env > config > default).
- Codex-mode worktree discovery under `$CODEX_HOME/worktrees/**` with repo-scoped filtering via `git worktree list --porcelain`.
- Codex-mode selection model: `<task>` is the exact `<opaque-id>` directory name under `$CODEX_HOME/worktrees`.
- Codex-mode `sync` command (`apply` by default; optional overwrite on conflict with a second confirmation).
- Codex-mode `cleanup` that only targets `$CODEX_HOME/worktrees/<opaque-id>` and is conservative with prominent warnings + confirmations.
- Add `modified_time` to `status` output (RFC3339 UTC).

## Non-Goals
- No registry/state files (no `registry.json` / extra TOML) for codex mode.
- No branch-based workflows in codex mode (no `finish` in codex mode; no `create --branch` in codex mode).
- No `restore` command in this phase.

## Plan
### Phase 1: Docs & Spec Design
- [x] Update `docs/research-2026-02-04-mode-classic-vs-codex.md` to reflect final decisions as implementation progresses.
- [x] Update `docs/schemas/config-gwtt.md` to include `mode` and env var `GWTT_MODE`.
- [x] Write a short CLI spec section (either in the research doc or a new schema doc) covering:
  - [x] Codex-mode worktree selection: `<opaque-id>` resolution rules and error messages.
  - [x] Repo scoping strategy for `list/status` in codex mode (Git-derived, not naming-derived).
  - [x] `sync` conflict detection signals and the overwrite confirmation flow (`--yes` behavior).
  - [x] `cleanup` safety model (scope restriction, warnings, second confirmation, and “restore is best-effort” note).

### Phase 2: Code Implementation
- [x] Add global `--mode` persistent flag on `cli/root.go` and plumb mode into command execution (context/config).
- [x] Add mode to config resolution:
  - [x] Env: `GWTT_MODE`.
  - [x] Config: `mode = "classic"|"codex"`.
  - [x] Validation and error messaging for unsupported values.
- [x] Codex-mode worktree discovery primitives:
  - [x] Determine `$CODEX_HOME` (or `CODEX_HOME`) and `worktreesRoot := $CODEX_HOME/worktrees`.
  - [x] Use `internal/worktree.List(ctx, runner, repoRoot)` and filter entries whose path is under `worktreesRoot`.
  - [x] Derive `<opaque-id>` as `filepath.Base(worktreePath)`.
- [x] Implement codex-mode behavior for read-only commands:
  - [x] `gwtt list` shows codex worktrees for the current repo with `Task=<opaque-id>` and a `$CODEX_HOME/...`-aware display path.
  - [x] `gwtt status` does the same, plus `modified_time`.
- [x] Add `modified_time` to status rows:
  - [x] Use filesystem `mtime` of the worktree directory.
  - [x] Format as RFC3339 UTC for JSON/CSV; table uses the same value.
- [x] Add `gwtt sync <opaque-id>` (codex-mode only):
  - [x] Default to “apply” (worktree -> local checkout).
  - [x] Conflict detection: dirty local checkout, failed apply/merge step, and/or both sides modified the same file (where detectable).
  - [x] On conflict, prompt whether to overwrite; require a second confirmation; `--yes` bypasses overwrite confirmation.
  - [x] Keep behavior aligned with Codex App: ignored files are not synced.
- [x] Re-check `cleanup` behavior for codex mode:
  - [x] Restrict deletions to `$CODEX_HOME/worktrees/<opaque-id>` only.
  - [x] Mirror Codex App “never clean up if …” rules when detectable; otherwise warn prominently and require a second confirmation.
  - [x] Document/communicate that Codex restore is best-effort (not guaranteed by `gwtt`).

### Phase 3: Unit Test Verification
- [ ] Add tests for `mode` precedence and validation (flag/env/config/default).
- [ ] Add tests for codex-mode list/status filtering (repo-scoped via `git worktree list` + `$CODEX_HOME/worktrees` prefix filter).
- [ ] Add tests for `<opaque-id>` derivation and path rendering (`$CODEX_HOME` display).
- [ ] Add tests for `modified_time` formatting (RFC3339 UTC) and JSON/CSV output shape.
- [ ] Add tests for `sync` conflict detection and confirmation gating (including `--yes`).
- [ ] Add tests for codex cleanup scope restriction + confirmation flow.

### Phase 4: README / CLI Docs Update
- [ ] Update `README.md`:
  - [ ] Document `--mode`, `GWTT_MODE`, and config `mode`.
  - [ ] Add codex-mode usage examples for `list/status/sync/cleanup`.
  - [ ] Update `## Notes` “Global flags” list to include `--mode`.
  - [ ] Document `modified_time` in `status` outputs (and the fixed date format).
- [ ] If applicable, update any man/help text sources under `man/` to reflect new commands/flags.

### Phase 5: Verify Doc Statuses
- [ ] Ensure this plan’s `status` matches the actual phase progress (`active` -> `completed` when done).
- [ ] Update the research doc’s `status` to `completed` once decisions are implemented and verified.
- [ ] Ensure any schema/doc updates have consistent status and dates (`modified-date` as needed).

## Acceptance Criteria
- Default behavior (no `--mode`, no `GWTT_MODE`, no config) remains unchanged.
- `--mode=codex` enables codex-specific list/status/sync/cleanup without impacting classic users.
- Codex-mode selection uses `<opaque-id>` reliably and errors clearly when not found/ambiguous.
- `status` includes `modified_time` with RFC3339 UTC formatting for machine outputs.
- Codex-mode cleanup is narrowly scoped and always warns + confirms before deletion.

## Risks / Notes
- Codex App “pinned/sidebar/thread linkage” signals may not be detectable from disk without reading Codex’s internal state; default to warnings + a second confirmation when uncertain.
- `sync` semantics are easy to get subtly wrong; keep the initial implementation conservative and well-tested.

## Related Research
- `docs/research-2026-02-04-mode-classic-vs-codex.md`
