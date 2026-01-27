---
title: "Extensible config rollout"
date: 2026-01-27
status: draft
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
- [ ] Define config schema for `ui`, `create`, `list`, `status`, `finish`, `cleanup` (exclude `base`/`target`).
- [ ] Decide whether to include `create.path.format` in this phase; if yes, require `{task}` and document reversibility constraints.
- [ ] Add config loader with precedence: flags > env > project > user > defaults.
- [ ] Introduce `GWTT_COLOR` env (mirrors `--nocolor`, inverted boolean).
- [ ] Implement config-backed defaults for list/status output, grid, absolute path, confirm toggles.
- [ ] Add `merge_mode` enum and enforce exclusivity across CLI flags.
- [ ] Add tests for config resolution (env/project/user) and merge-mode validation.
- [ ] Update examples and README config section.

## Acceptance Criteria
- CLI behavior is unchanged when no config is present.
- `GWTT_COLOR` and config honor existing precedence rules.
- Only one merge strategy is allowed at a time (`ff`, `no-ff`, `squash`, `rebase`).
- Config defaults never override current-branch target selection.

## Risks / Notes
- Path templating can break `TaskFromPath` discovery unless `{task}` is enforced.
- Merge-mode validation must stay consistent across flags and config to avoid drift.

## Related Research
- `docs/research-2026-01-27-extensible-config.md`
