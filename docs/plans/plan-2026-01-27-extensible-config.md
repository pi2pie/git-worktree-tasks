---
title: "Extensible config rollout"
date: 2026-01-27
status: active
agent: codex
---

## Goal
Ship an extensible config system beyond theme, prioritizing low-risk defaults and preserving current-branch semantics.

## Scope
- Config for UI and output defaults, confirmation toggles, and merge strategy.
- No config default for `create.base` or `status/finish.target`.
- Path templating only if `{task}` placeholder is enforced.

## Non-Goals
- No change to theme selection precedence.
- No dynamic target defaults tied to repo state beyond current behavior.
- No new dependency for config parsing beyond existing TOML usage.

## Plan
### Phase 1: Design
- [x] Create schema doc in `docs/schemas/config-gwtt.md`.
- [x] Define config schema for `ui`, `create`, `list`, `status`, `finish`, `cleanup` (exclude `base`/`target`).
- [x] Decide whether to include `create.path.format` in this phase; if yes, require `{task}` and document reversibility constraints.
- [x] Define merge strategy enum and mapping from flags + config.
- [x] Document env/config/flag precedence for new settings.

### Phase 2: Implementation
- [ ] Add config loader with precedence: flags > env > project > user > defaults.
- [ ] Introduce `GWTT_COLOR` env (mirrors `--nocolor`, inverted boolean).
- [ ] Implement config-backed defaults for list/status output, grid, absolute path, confirm toggles.
- [ ] Add `merge_mode` enum and enforce exclusivity across CLI flags.

### Phase 3: Verification & Docs
- [ ] Add tests for config resolution (env/project/user) and merge-mode validation.
- [ ] Update examples and README config section.
- [ ] Verify related docs status and update as phases complete:
  - [ ] Plan status aligns with phase completion.
  - [ ] Research status is updated (`draft` -> `in-progress`/`completed`).
  - [ ] Schema doc status reflects readiness.
  - [ ] Job records reflect actual work status.

## Dependencies
- Current TOML parsing via `internal/config` (BurntSushi/toml).
- CLI flag wiring in `cli/*` must remain backward compatible.
- Worktree naming utilities in `internal/worktree/naming.go` if path templating is introduced.

## Acceptance Criteria
- CLI behavior is unchanged when no config is present.
- `GWTT_COLOR` and config honor existing precedence rules.
- Only one merge strategy is allowed at a time (`ff`, `no-ff`, `squash`, `rebase`).
- Config defaults never override current-branch target selection.

## Risks / Notes
- Path templating can break `TaskFromPath` discovery unless `{task}` is enforced.
- Merge-mode validation must stay consistent across flags and config to avoid drift.

## Process Notes
- Update doc statuses immediately after each task completes.
- Use Phase 3 verification as a final guardrail to confirm status alignment.

## Related Research
- `docs/research-2026-01-27-extensible-config.md`
