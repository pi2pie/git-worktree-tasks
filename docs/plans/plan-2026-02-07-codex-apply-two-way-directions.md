---
title: "Codex apply two-way directions and cli extraction prep"
created-date: 2026-02-07
modified-date: 2026-02-07
status: active
agent: codex
---

## Goal

Implement a two-way `apply`/`overwrite` command model in codex mode with explicit direction semantics, while first extracting high-reuse CLI helpers to reduce risk before behavior changes.

## Scope

- Define and implement two-way direction handling for codex handoff flows.
- Align CLI command surface with Codex app-level concepts (`apply` and `overwrite` as peer actions).
- Improve dry-run clarity for direction and destructive operations.
- Perform foundational CLI refactors first (no behavior change).

## Non-Goals

- No registry/state file for branch wiring in this phase.
- No new classic-mode behavior.
- No sub-package split for `cli` command internals yet.

## Phase Checklist

### Phase 1: Spec and UX Lock

- [ ] Finalize command matrix for:
  - `apply --to local|worktree`
  - `overwrite --to local|worktree`
  - optional `apply --force` alias policy
- [ ] Finalize conflict behavior for non-destructive apply (no implicit direction switching).
- [ ] Finalize dry-run output schema (header + preflight + action list).

### Phase 2: Foundational Refactor (No Behavior Change)

- [x] Extract mode/codex context helper used by command handlers.
- [x] Extract/reuse codex worktree lookup helper by opaque ID.
- [x] Relocate shared `runGit` helper from `finish.go` into a neutral CLI helper file.
- [x] Keep output/error semantics unchanged in this phase.

### Phase 3: Apply/Overwrite Implementation

- [ ] Implement `--to` direction support in codex handoff commands.
- [ ] Introduce `overwrite` as peer command (or locked alias strategy per Phase 1).
- [ ] Remove implicit overwrite fallback from `apply` conflict path.
- [ ] Preserve confirmation gating (`--yes`) for destructive paths.

### Phase 4: Dry-Run and Messages

- [ ] Implement structured dry-run plan output with explicit source/destination and destructive markers.
- [ ] Update user-facing conflict and next-step guidance.

### Phase 5: Tests and Docs

- [ ] Add/adjust unit and integration tests for direction + overwrite behaviors.
- [ ] Update README/man pages/help text for new command semantics.
- [ ] Update related research and mark this plan status appropriately.

## Acceptance Criteria

- `apply` and `overwrite` semantics are explicit and direction-stable.
- No hidden mutation of the opposite side on `apply` conflict.
- Dry-run clearly shows operation intent and action sequence.
- Foundational helper extraction lands without behavior regressions.

## Risks / Notes

- Confirmation wording must remain clear when destructive destination reset is involved.
- Refactor-first approach lowers risk but can surface latent coupling in `cli`.

## Related Research

- `docs/research-2026-02-07-codex-apply-direction-and-source-checkout.md`
