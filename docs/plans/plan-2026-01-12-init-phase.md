---
title: "Init Phase Plan"
created-date: 2026-01-12
status: active
agent: codex
---

## Goal

Establish the initial CLI structure, UX patterns, and repository scaffolding for git worktree task management.

## Scope

- CLI skeleton using cobra with core command groups and flags.
- Shared styling utilities using lipgloss.
- Initial Bubbletea TUI entry point (basic navigation and layout only).
- Documentation and examples to onboard contributors.

## Non-Goals

- Full feature-complete worktree workflows.
- Persistent config or integration with external services.
- Advanced TUI flows beyond a simple shell.

## Milestones

1. Define command tree and CLI UX conventions, including fixed worktree path rules, slugification, and list/status output.
2. Create base packages/modules for git/worktree operations and UI, aligned with the operations matrix.
3. Add a minimal TUI with placeholder screens for create/merge/cleanup and confirmation flows.
4. Add docs for the init phase and initial usage examples.

## Deliverables

- Cobra command structure with placeholders for `create`, `finish`, `cleanup`, and `list`, plus explicit remove-worktree/remove-branch flags.
- Lipgloss style package with consistent typography and colors.
- Bubbletea app with navigation between placeholder views and double-confirmation for destructive actions.
- Basic README updates and examples folder entries.
- `list` output defaults to `table` with `--output json` option.

## Decisions

- Worktree path is fixed to `../<repo>_<task>` with consistent task slugification.
- `finish` requires a second confirmation before deleting worktree/branch, with a bypass flag for full cleanup.
- Merge modes support `--squash`, `--no-ff`, and `--rebase`; default follows standard `git merge` behavior.

## Checklist

- [x] Confirm final command names and subcommand hierarchy.
- [x] Confirm slugification rules and worktree path derivation.
- [x] Confirm `finish` confirmation flow and bypass flag naming.
- [x] Confirm cleanup combinations (remove worktree only, remove branch only, both).
- [x] Confirm list/status output fields for simple vs detailed modes.
- [x] Confirm list output formats (`table` default, `--output json`).
- [x] Confirm merge strategy flag mapping and defaults.

## Risks / Open Questions

- Final command naming and subcommand hierarchy.
- Merge strategy defaults and how they map to CLI flags.
- Git worktree edge cases (existing branches, detached HEAD, conflicts).

## Related Research

- docs/researches/research-2026-01-12-worktree-ops-matrix.md
