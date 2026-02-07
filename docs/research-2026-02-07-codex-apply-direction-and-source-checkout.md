---
title: "Codex apply direction and source checkout behavior"
created-date: 2026-02-07
modified-date: 2026-02-07
status: in-progress
agent: codex
---

## Goal

Define a clearer, less surprising `gwtt apply` model in `--mode=codex` by making direction explicit, defining overwrite semantics for both directions, and redesigning `--dry-run` output so users can see what will happen before mutation.

## Key Findings

### Current `apply` behavior is directionally ambiguous

- `gwtt apply <opaque-id>` currently starts as worktree -> local, but on conflict it can switch to local -> worktree after confirmation.
- This fallback is implemented in `cli/apply.go` via:
  - `applyWorktreeChanges(...)` for worktree -> local.
  - `overwriteWorktreeChanges(...)` for local -> worktree (hard reset + clean on worktree first).
- Result: users can enter `apply` expecting "to local" but end in an overwrite flow that discards worktree changes.

### Direction flag fits the Codex App mental model

- Codex App exposes two explicit directions ("To local" / "From local") for the same handoff concept.
- CLI parity is straightforward with a destination-oriented flag:
  - `gwtt apply <opaque-id>` as safe default (`--to local`).
  - `gwtt apply --to worktree <opaque-id>` for reverse flow.

### Overwrite should be destination-scoped, not direction-switching

- In a two-way model, overwrite must mean "replace destination with source" regardless of direction.
- This is safer and more predictable than "on conflict, flip direction."
- The current hard reset + clean logic can be reused, but should only run when overwrite is explicitly requested.
- Since Codex App presents "Apply" and "Overwrite" as separate top-level actions, CLI parity is best achieved with separate commands rather than hiding overwrite behind a flag.

### The "source local checkout removed" case should be scoped narrowly for now

- Current CLI design is intentionally registry-free for codex mode and identifies worktrees by Git + `$CODEX_HOME/worktrees`.
- App-style "wire this worktree to a new/existing branch" implies extra state and UX that `gwtt` does not currently persist.
- Adding branch wiring now would expand scope beyond `apply` direction clarification.

### Current `--dry-run` output lacks operation context

- Today, dry-run mostly prints raw git commands and per-file copy lines.
- It does not clearly tell users:
  - selected direction,
  - source/destination paths,
  - whether overwrite/reset will occur,
  - summary of what is expected to change.

### `cli/apply.go` is doing too much in one file

- The file contains command wiring, worktree resolution, conflict detection, transfer logic, patch I/O, and file copying.
- There is duplicated transfer logic between forward/reverse paths, increasing maintenance cost and making direction changes riskier.

## Implications or Recommendations

- Add `--to` with enum values `local|worktree` (default: `local`).
- Keep backward compatibility:
  - `gwtt apply <opaque-id>` behaves as today’s safe default (worktree -> local).
- Introduce explicit overwrite behavior as a sibling command:
  - `gwtt apply ...` for non-destructive transfer attempt.
  - `gwtt overwrite ...` for destructive destination replacement.
- Optional compatibility alias:
  - `gwtt apply --force ...` can be supported as shorthand that dispatches to overwrite behavior.
  - Document it as compatibility/convenience, not the primary UX.
- Stop implicit direction switching:
  - On conflict in `apply`, return a clear conflict error.
  - Suggest rerun with `overwrite` for the same `--to` direction.
- Make overwrite direction-agnostic:
  - `overwrite --to local` discards local destination changes and applies worktree -> local.
  - `overwrite --to worktree` discards worktree destination changes and applies local -> worktree.
- For "source local checkout gone", do not add wiring state in this phase:
  - If selected `<opaque-id>` is not resolvable from the current repository context, fail with an actionable message.
  - Message should instruct user to run from a checkout sharing the same Git common dir or recreate/relink manually.
- Redesign `--dry-run` output to be plan-oriented:
  - Print operation header (direction, source, destination, overwrite mode).
  - Print preflight summary (destination dirty, overlap, tracked patch state, untracked count).
  - Print ordered action plan (check/apply/reset/clean/copy steps).
  - Keep command echo for transparency, but make it secondary to the plan summary.
- Split `apply` implementation into focused files before/with feature work:
  - `cli/apply_command.go`: Cobra command, flags, and top-level flow.
  - `cli/apply_resolve.go`: codex worktree resolution and validation.
  - `cli/apply_conflicts.go`: conflict detection helpers.
  - `cli/apply_transfer.go`: direction-agnostic transfer operations.
  - `cli/apply_files.go`: temp patch + copy helpers.
- Refactor transfer logic to one core function:
  - `transferChanges(sourceRoot, destRoot, opts)` with optional `resetDest` behavior.
  - This keeps direction implementation DRY and lowers regression risk.

## Proposed CLI Spec (Draft)

```bash
# safe default
gwtt apply <opaque-id>                        # same as --to local

# non-destructive apply, explicit direction
gwtt apply --to local <opaque-id>             # worktree -> local
gwtt apply --to worktree <opaque-id>          # local -> worktree

# destructive overwrite as peer command
gwtt overwrite --to local <opaque-id>         # destination local reset/clean first
gwtt overwrite --to worktree <opaque-id>      # destination worktree reset/clean first

# optional compatibility alias
gwtt apply --force --to worktree <opaque-id>  # shorthand for overwrite behavior
```

Apply vs overwrite matrix:

- `apply`:
  - never reset/clean destination;
  - on destination dirtiness or patch conflict, exit non-zero and print next-step hint.
- `overwrite`:
  - require confirmation unless `--yes`;
  - reset and clean destination before transfer;
  - then apply tracked diff + copy untracked files from source.

## Phase 1 Decisions (Locked)

- Ship `overwrite` as a peer command and keep `apply --force` as compatibility alias.
- `apply` remains non-destructive and direction-stable; conflicts now fail with overwrite guidance.
- Keep dry-run redesign text-first for now; JSON planning output is deferred.
- Keep relink/wiring out of scope in this phase.

## Implementation Status (Current)

- Implemented:
  - `apply --to local|worktree`
  - `overwrite --to local|worktree`
  - `apply --force` compatibility alias
  - no implicit direction switching on `apply` conflicts
- Pending:
  - dry-run plan-style output redesign
  - split `cli/apply.go` into focused files (`apply_command`, `apply_resolve`, `apply_conflicts`, `apply_transfer`, `apply_files`)

## Dry-Run Output Redesign (Draft)

Proposed dry-run output shape:

```text
apply plan
  to: local
  source: $CODEX_HOME/worktrees/bf15/repo
  destination: /path/to/repo
  overwrite: false

preflight
  destination_dirty: true
  overlapping_files: 2
  tracked_patch: present
  untracked_files: 3

actions
  1. git -C /path/to/repo apply --check <temp-patch>
  2. git -C /path/to/repo apply <temp-patch>
  3. copy <src>/a.txt -> <dst>/a.txt
  4. copy <src>/b.txt -> <dst>/b.txt
```

Overwrite mode should show destructive steps explicitly:

```text
actions
  1. git -C <destination> reset --hard
  2. git -C <destination> clean -fd
  ...
```

Output requirements:

- Must clearly indicate whether destination will be reset/cleaned.
- Must include source and destination path labels.
- Must preserve command transparency for debugging.
- Must avoid ambiguous "apply complete" style messaging in dry-run mode.

## References

- Current implementation: [^apply-code]
- Prior codex mode research: [^mode-research]
- Related mode plan: [^mode-plan]

[^apply-code]: `cli/apply.go`
[^mode-research]: `docs/research-2026-02-04-mode-classic-vs-codex.md`
[^mode-plan]: `docs/plans/plan-2026-02-04-mode-classic-and-codex.md`

## Related Plans

- `docs/plans/plan-2026-02-04-mode-classic-and-codex.md`
