---
title: "Mode Flag: classic vs codex"
created-date: 2026-02-04
modified-date: 2026-02-04
status: completed
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

Based on Codex App documentation (and current app UI), the worktree model is intentionally different from this CLI’s task/branch model:

- **Worktree location is not per-worktree user-chosen**: worktrees are created under `$CODEX_HOME/worktrees` so the app can manage them consistently.
- **Worktrees start in detached HEAD** by default (to avoid Git’s restriction that a branch cannot be checked out in two worktrees at once).
- **Local changes may be applied** when the worktree is created from an existing local branch with uncommitted changes.
- **“Hand off changes” is a first-class operation** for getting changes between the local checkout and the worktree:
  - UI labels are now “Hand off changes” with directions “To local” / “From local”.
  - The official docs still use “Sync with local” for the same action and describe “Apply” and “Overwrite” modes.
  - Sync does not transfer ignored files (and the resulting state may not match a full re-clone).
  - Terminology overlap: “Apply” also exists in the Codex CLI as `codex apply <task_id>` (Codex Cloud task diff). This can be confusing when discussing “apply” in the app UI vs the CLI.
- **Worktree restoration** is a distinct concept (recreate a worktree from a Codex snapshot, rather than from the current local checkout).
- **Cleanup is app-governed and tied to threads**: the Codex app cleans up worktrees when you archive threads (or on startup for worktrees with no associated threads), and it preserves a snapshot for later restore.

### Codex App FAQ takeaways (constraints we should mirror)

- Worktrees are created under `$CODEX_HOME/worktrees` so Codex can manage them consistently.
- Sessions cannot be moved between worktrees: to change environments, you start a new thread in the target environment and restate the prompt.
- Threads can remain even if the worktree directory is cleaned up; Codex snapshots work before cleanup and can offer restore when reopening the thread.

### Decisions for `--mode=codex` (CLI alignment)

To keep `classic` stable and keep `codex` aligned with Codex App:

- **Identity & mapping via inspection (no extra registry files):**
  - Do not introduce new state files like `registry.json`/TOML for codex mode.
  - For `list`/`status`, inspect `$CODEX_HOME/worktrees/**` (and/or `git worktree list --porcelain` scoped to the local checkout) to discover worktrees and compute status.
  - In codex mode, `<task>` is the **opaque ID directory** directly under `$CODEX_HOME/worktrees`.
    - Example path: `~/.codex/worktrees/bf15/git-worktree-tasks`
    - `<task>` is `bf15` (the opaque ID), **not** `git-worktree-tasks`.
- **Create is detached-only:** in `codex` mode, `create` should not offer a `--branch` escape hatch; the default stays detached to avoid future complexity.
- **Finish is classic-only:** in `codex` mode, `finish` is not a good fit; use a dedicated `apply` command instead.
- **Cleanup follows current mode:** `gwtt cleanup` should operate on the worktrees owned by the active mode (`classic` naming vs `$CODEX_HOME/worktrees`), rather than mixing behaviors.
- **Apply UX:** `gwtt apply <opaque-id>` defaults to “apply”; if a conflict is detected, prompt to “overwrite” (second confirmation), skippable with `--yes`.
- **Cleanup in codex mode:** free disk by deleting the on-disk directory under `$CODEX_HOME/worktrees/<opaque-id>`, but only when it is safe:
  - Skip worktrees that fall under Codex App’s “never clean up if …” restrictions.
  - If we cannot verify a restriction (e.g., pinned/sidebar linkage), still allow deletion but show a prominent warning and require a second confirmation (skippable with `--yes`).
  - Treat “Codex can restore later” as best-effort: Codex App snapshots before _its own_ cleanup; `gwtt` cannot guarantee a snapshot exists before manual deletion.

### Practical restrictions implied by `--mode=codex`

- **No “task branch” assumption:** detached worktrees mean we can’t infer branch names from task names.
- **No arbitrary `--path` override:** codex-mode worktrees live under `$CODEX_HOME/worktrees`; allowing arbitrary paths would complicate cleanup, display, and registry invariants.
- **Different command surface:** branch-merge workflows (`finish`) are replaced by apply workflows (`apply` / `overwrite`).
- **Cleanup restrictions from Codex App:** Codex App’s auto-cleanup is disabled in some cases (e.g., pinned conversation, added to sidebar, age > 4 days, worktree count > 10).
  - Note: the “age > 4 days” / “count > 10” conditions are counterintuitive, but this is the wording in the official docs as of 2026-02-04.
- **Detached HEAD as codex marker:** codex-mode worktrees are detached; classic-mode worktrees are expected to be on a branch. Use this to keep classic commands from “seeing” codex worktrees.

### Path display differences (UX)

Codex App uses a “variable-aware” presentation of paths (and the user specifically called out `$CODEX_HOME`):

- In `codex` mode, prefer showing worktree paths under `$CODEX_HOME` as `$CODEX_HOME/...` rather than a long absolute path.
- In both modes, consider shortening the home directory to `~` (or `$HOME`) when rendering absolute paths.

## Implications or Recommendations

- Add a global flag `--mode` with values:
  - `classic` (default): current behavior and naming convention.
  - `codex`: Codex App-aligned behavior (detached worktrees, `$CODEX_HOME` root, ID-based mapping, apply-oriented workflow).
- Treat `codex` mode as additive:
  - Keep existing commands and semantics intact in `classic`.
  - Introduce new behavior behind `--mode=codex` (and/or new codex-only subcommands like `apply`) rather than changing defaults.
- Define a “codex worktree root”:
  - In codex mode, treat `$CODEX_HOME/worktrees` as the only allowable root for “managed” worktrees.
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
- Codex-mode cleanup should be **disk-focused and conservative**:
  - Only target paths under `$CODEX_HOME/worktrees/<opaque-id>` (no arbitrary deletion).
  - Attempt to mirror Codex App’s “never clean up if …” rules; if we cannot verify a rule (e.g., pinned/sidebar linkage), show a prominent warning and require a second confirmation (skippable with `--yes`).
- Repo scoping in codex mode should come from Git, not naming:
  - To list/status only the Codex worktrees for the _current repo_, run `git -C <repoRoot> worktree list --porcelain` and include only entries whose `worktree` path is under `$CODEX_HOME/worktrees/`.
  - This avoids relying on the opaque directory name to encode repo identity.
- Consider adding `modified_time` to `status` (and optionally `list`) rows:
  - Use the filesystem `mtime` of the worktree directory as a pragmatic “last touched” signal.
  - Output format recommendation: RFC3339 in UTC for JSON/CSV; table output can display the same value (no additional config initially).

## CLI Spec (Draft)

### Mode resolution

- Flag: `--mode` (`classic` or `codex`).
- Env: `GWTT_MODE`.
- Config: top-level `mode = "classic"|"codex"`.
- Default: `classic`.

### Codex home + worktrees root

- Resolve Codex home from the `CODEX_HOME` env var; if unset, default to `~/.codex` (Codex App default).
- Managed codex-mode worktrees are always under `$CODEX_HOME/worktrees/`.
- Display paths under this root as `$CODEX_HOME/...` by default (unless `--abs` forces absolute).

### Selection model (`<task>` in codex mode)

- `<task>` is the **exact opaque ID** directory name under `$CODEX_HOME/worktrees` (the first path segment).
- `gwtt list/status --mode=codex` should render `TASK=<opaque-id>` to make the identifier discoverable and copy/paste friendly.

### Repo scoping (`list/status` in codex mode)

- Use Git as the source of truth for “worktrees belonging to this repo”:
  - Run `git -C <repoRoot> worktree list --porcelain`.
  - Filter entries whose worktree path is under `$CODEX_HOME/worktrees/`.
- Do not attempt to infer repo identity from `<opaque-id>` naming.

### `apply` (codex mode)

- CLI: `gwtt apply <opaque-id>` (default operation: apply worktree changes into the local checkout).
- Conflict detection signals (predictable, conservative):
  - Local checkout is dirty, or
  - the apply/merge step fails, or
  - both sides modified the same file (where detectable).
- On conflict: prompt whether to “overwrite” (local -> worktree) and require a second confirmation; `--yes` bypasses the overwrite confirmation.
- Keep Codex App parity: ignored files are not transferred.

### `cleanup` (codex mode)

- Scope: only delete the on-disk directory at `$CODEX_HOME/worktrees/<opaque-id>` (free disk; do not touch classic paths).
- Since pinned/sidebar/thread linkage is not reliably detectable without reading Codex App state:
  - Always show a prominent warning (“may break pinned/sidebar restore expectations”) and require an extra confirmation (skippable with `--yes`).
  - Treat restore as best-effort: Codex saves a snapshot **before its own cleanup**; `gwtt`-initiated deletion cannot guarantee a snapshot exists unless Codex already took one.

### `status` metadata addition: `modified_time`

- Add `modified_time` derived from filesystem `mtime` of the worktree directory.
- Format: RFC3339 UTC for JSON/CSV; table output prints the same string.
- No date-format config in the first iteration.

### Raw output (codex mode)

- `--output raw` should return a **composable path** (relative to `$CODEX_HOME`), e.g. `worktrees/bf15/git-worktree-tasks`.
- Display output (table/text) should continue to render `$CODEX_HOME/...` for readability.

## Open Questions (Remaining)

- Can we (safely) detect any Codex cleanup-restriction signals from disk without coupling `gwtt` to Codex’s internal storage formats?
  - Likely no; default to warnings + a second confirmation (skippable with `--yes`) for codex cleanup.
- What is the most user-friendly confirmation wording for “overwrite” (apply) and “yolo delete” (cleanup) that still prevents accidents?
  - Pending: depends on user confidence that Codex can restore a worktree after manual deletion (docs only guarantee snapshots before **Codex-managed** cleanup).

## Notes

- Restoration remains out of scope: keep to create/apply/list/status only for now; a future `restore` likely needs to integrate with Codex App snapshot state.

## Open Issues (Worktrees + Shell/Run Scripts)

- App-side issues only (not CLI behavior). Track open worktree-related issues: https://github.com/openai/codex/issues?q=is%3Aissue%20state%3Aopen%20worktree
- Recent example: “Codex app: Worktrees keep forgetting the "Run" script” (open, Feb 3, 2026). https://github.com/openai/codex/issues/10476

## References

- Codex App worktrees documentation (includes cleanup + FAQ): https://developers.openai.com/codex/app/worktrees/
- Git worktree manual: https://git-scm.com/docs/git-worktree

## Related Plans

- `docs/plans/plan-2026-02-04-mode-classic-and-codex.md`
