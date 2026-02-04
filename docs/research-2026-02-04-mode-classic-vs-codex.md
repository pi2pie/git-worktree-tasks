---
title: "Mode Flag: classic vs codex"
date: 2026-02-04
modified-date: 2026-02-04
status: in-progress
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

### Decisions for `--mode=codex` (CLI alignment)
To keep `classic` stable and keep `codex` aligned with Codex App:
- **Identity & mapping via registry (not per-worktree metadata files):**
  - Store the minimal mapping needed to resolve `<task> -> worktree path` in a registry under `$CODEX_HOME` (rather than scattering metadata files inside worktrees).
  - “Metadata” should remain derivable from the worktree itself (e.g., via `gwtt status`-equivalent logic).
- **Create is detached-only:** in `codex` mode, `create` should not offer a `--branch` escape hatch; the default stays detached to avoid future complexity.
- **Finish is classic-only:** in `codex` mode, `finish` is not a good fit; use a dedicated `sync` command instead.
- **Cleanup follows current mode:** `gwtt cleanup` should operate on the worktrees owned by the active mode (`classic` naming vs `$CODEX_HOME/worktrees` + registry), rather than mixing behaviors.

### Practical restrictions implied by `--mode=codex`
- **No “task branch” assumption:** detached worktrees mean we can’t infer branch names from task names.
- **No arbitrary `--path` override:** codex-mode worktrees live under `$CODEX_HOME/worktrees`; allowing arbitrary paths would complicate cleanup, display, and registry invariants.
- **Different command surface:** branch-merge workflows (`finish`) are replaced by sync workflows (`sync apply` / `sync overwrite`).

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
- Implement `mode` as a first-class config value (not flag-only):
  - Flag: `--mode` (highest precedence).
  - Env var: `GWTT_MODE`.
  - Config: `gwtt.config.toml`/`gwtt.toml` and `$HOME/.config/gwtt/config.toml`.
  - Default: `classic`.
- Keep ignored-file behavior aligned with Codex App in `codex` mode (do not add “include ignored” options initially).

## Open Questions
- **Registry schema & location:** Where exactly under `$CODEX_HOME` should the registry live, and what format should it use?
  - Example options: `$CODEX_HOME/gwtt/registry.json`, `$CODEX_HOME/gwtt/registry.toml`, or `$CODEX_HOME/gwtt/worktrees/registry.json`.
  - Minimum recommended keys per entry: `task`, `repoRoot`, `worktreePath`, `createdAt`, and the “source ref” used to create it (branch/ref/commit).
- **Sync UX:** What should the CLI surface look like?
  - `gwtt sync <task> --apply|--overwrite` vs `gwtt sync apply <task>` / `gwtt sync overwrite <task>`.
  - Confirmations: treat “apply/overwrite” as a second, explicit confirmation (skippable with `--yes`), similar to the existing destructive confirmations pattern.
- **Restoration support:** Do we want a `restore` operation in the CLI (to mirror Codex App), or keep scope to create/sync/cleanup only?

## References
- Codex App worktrees documentation: https://developers.openai.com/codex/app/worktrees/
- Git worktree manual: https://git-scm.com/docs/git-worktree

## Related Plans
- (none)
