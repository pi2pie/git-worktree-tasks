---
title: "Mode Flag: classic vs codex"
date: 2026-02-04
status: draft
agent: codex
---

## Goal
Define what a new global `--mode` flag should mean for this CLI, so we can support Codex App-style worktrees without breaking the current (“classic”) behavior.

## Key Findings

### Current CLI behavior (“classic”)
- **Create** couples “task” to both branch and path:
  - Branch: `<task>` (slugified).
  - Path: `../<repo>_<task>` (relative to repo root’s parent).
- **List/status/cleanup/finish** assume the above naming convention to map paths back to tasks.
- **Paths are displayed** relative to the repo root by default; `--abs`/`--absolute-path` shows absolute paths.

### Codex App worktree behavior (“codex”)
Based on Codex App documentation, the worktree model is intentionally different from this CLI’s task/branch model:
- **Worktree location is not user-chosen**: worktrees are created under `$CODEX_HOME/worktrees` so the app can manage them consistently.
- **Worktrees start in detached HEAD** by default (to avoid Git’s restriction that a branch cannot be checked out in two worktrees at once).
- **Local changes may be applied** when the worktree is created from an existing local branch with uncommitted changes.
- **Sync is a first-class operation** for getting changes between the local checkout and the worktree:
  - “Apply” worktree changes into local checkout.
  - “Overwrite” worktree from local checkout.
  - Sync does not transfer ignored files (and the resulting state may not match a full re-clone).
- **Worktree restoration** is a distinct concept (recreate a worktree from a Codex snapshot, rather than from the current local checkout).

### Practical restrictions we likely need in `--mode=codex`
To avoid accidental behavior drift and to reflect the Codex App constraints, `codex` mode likely implies:
- **No implicit “task branch”**: default worktrees should be detached and may not have a stable `<task>` branch name.
- **Different identity model**: a worktree might be identified by an ID (or metadata) rather than the `<repo>_<task>` path convention.
- **Limited/changed support for merge flows**:
  - `finish` (merge branch into target) only makes sense if a branch exists; detached worktrees need either (1) an explicit branch creation step or (2) a new “sync/apply” flow.
- **No arbitrary `--path` override** (or an explicit escape hatch), since Codex App doesn’t allow it and allowing it would complicate cleanup and display rules.

### Path display differences (UX)
Codex App uses a “variable-aware” presentation of paths (and the user specifically called out `$CODEX_HOME`):
- In `codex` mode, prefer showing worktree paths under `$CODEX_HOME` as `$CODEX_HOME/...` rather than a long absolute path.
- In both modes, consider shortening the home directory to `~` (or `$HOME`) when rendering absolute paths.

## Implications or Recommendations
- Add a global flag `--mode` with values:
  - `classic` (default): current behavior and naming convention.
  - `codex`: Codex App-aligned behavior (detached worktrees, `$CODEX_HOME` root, ID-based mapping, sync-oriented workflow).
- Treat `codex` mode as additive:
  - Keep existing commands and semantics intact in `classic`.
  - Introduce new behavior behind `--mode=codex` (and/or new codex-only subcommands like `sync`) rather than changing defaults.
- Define a “codex worktree root”:
  - Required env var: `$CODEX_HOME` (error clearly if missing), or a documented default fallback (e.g., `~/.codex`) if we want to be permissive.
- Introduce a path rendering helper that can:
  - Render relative-to-repo paths (classic default).
  - Render `$CODEX_HOME`-relative paths (codex default).
  - Still respect existing `--abs` behavior.

## Open Questions
- **Identity & mapping:** How should `gwtt` map `<task>` to a Codex-style worktree (ID, metadata file, a `.gwtt/` registry under `$CODEX_HOME`, etc.)?
- **Command support matrix:** Which commands should be supported in `codex` mode?
  - `create`: detached by default? allow `--branch` to opt into a branch?
  - `finish`: allowed only when branch exists? replaced by `sync apply`?
  - `cleanup`: should it clean only `$CODEX_HOME/worktrees` entries or also classic paths?
- **Sync semantics in a CLI context:** Do we implement Codex-style “apply/overwrite” as:
  - new `sync` command, or
  - an option on existing commands (riskier for UX), or
  - both (with `sync` as the primary entry point)?
- **Ignored files:** Do we explicitly document that ignored files are not synced in `codex` mode, and do we provide an opt-in mechanism (e.g. tar/rsync) or keep behavior aligned with Codex App?
- **Config & precedence:** Should mode be configurable via `GWTT_MODE` and `config.toml` (similar to theme), or remain flag-only to reduce ambiguity?

## References
- Codex App worktrees documentation: https://developers.openai.com/codex/app/worktrees/
- Git worktree manual: https://git-scm.com/docs/git-worktree

## Related Plans
- (none)

