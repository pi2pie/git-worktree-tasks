---
title: "gwtt configuration schema"
created-date: 2026-01-27
modified-date: 2026-02-07
status: in-progress
agent: codex
---

## Goal

Define the authoritative configuration schema for `gwtt`, including keys, types, defaults, and precedence.

## Precedence

1. CLI flags
2. Environment variables
3. Project config (`gwtt.config.toml` or `gwtt.toml`)
4. User config (`$HOME/.config/gwtt/config.toml`)
5. Built-in defaults

## Environment variables

- `GWTT_THEME` overrides `[theme].name`.
- `GWTT_COLOR` overrides `[ui].color_enabled`.
- `GWTT_MODE` overrides `mode`.
- `GWTT_DRY_RUN_MASK_SENSITIVE_PATHS` overrides `[dry_run].mask_sensitive_paths`.
- `CODEX_HOME` is consumed in `mode="codex"` to locate `$CODEX_HOME/worktrees` (this is a Codex App/Codex CLI convention, not a `gwtt` config key).
  - Fallback: if `CODEX_HOME` is unset, `gwtt` should assume `~/.codex` (home dir + `/.codex`) to align with Codex defaults.
  - Note: confirm the default path against current Codex App/Codex CLI docs when implementing (and prefer matching their behavior over introducing a new `gwtt`-specific default).

## Schema

### Root

- `mode` (string enum: `classic`, `codex`; default: `classic`)

### `[theme]`

- `name` (string, default: `"default"`)

### `[ui]`

- `color_enabled` (bool, default: `true`)

### `[table]`

- `grid` (bool, default: `false`)

### `[dry_run]`

- `mask_sensitive_paths` (bool, default: `true`)
  - When `true`, home-prefixed paths in `--dry-run` output are masked:
    - POSIX: `$HOME/...`
    - Windows: `%USERPROFILE%\\...`
  - When `false`, `--dry-run` output keeps raw absolute paths.
  - CLI overrides:
    - `--mask-sensitive-paths[=true|false]`
    - `--no-mask-sensitive-paths`

### `[create]`

- `output` (string enum: `text`, `raw`; default: `text`)
- `skip_existing` (bool, default: `false`)

#### `[create.path]`

- `root` (string, default: `"../"`)
- `format` (string, default: `"{repo}_{task}"`)
  - Must include `{task}`.
  - `{repo}` is optional.

### `[list]`

- `output` (string enum: `table`, `json`, `csv`, `raw`; default: `table`)
- `field` (string enum: `path`, `task`, `branch`; default: `path`)
- `absolute_path` (bool, default: `false`)
- `grid` (bool, default: `false`)
- `strict` (bool, default: `false`)

### `[status]`

- `output` (string enum: `table`, `json`, `csv`; default: `table`)
- `absolute_path` (bool, default: `false`)
- `grid` (bool, default: `false`)
- `strict` (bool, default: `false`)

### `[finish]`

- `cleanup` (bool, default: `false`)
- `remove_worktree` (bool, default: `false`)
- `remove_branch` (bool, default: `false`)
- `force_branch` (bool, default: `false`)
- `merge_mode` (string enum: `ff`, `no-ff`, `squash`, `rebase`; default: `ff`)
- `confirm` (bool, default: `true`)
  - `false` bypasses prompts (equivalent to `--yes`).

#### Merge strategy mapping

- `ff` (default): no merge flags.
- `no-ff`: adds `--no-ff`.
- `squash`: adds `--squash`.
- `rebase`: uses the rebase flow (`rebase` then `merge --ff-only`).
- Exactly one merge strategy may be active at a time; config and flags must agree.

### `[cleanup]`

- `remove_worktree` (bool, default: `true`)
- `remove_branch` (bool, default: `true`)
- `worktree_only` (bool, default: `false`)
- `force_branch` (bool, default: `false`)
- `confirm` (bool, default: `true`)
  - `false` bypasses prompts (equivalent to `--yes`).

## Decisions

- `create.path.format` must include `{task}` to preserve task discovery.
- `merge_mode` is exclusive; only one strategy may be active at a time.
- Codex-mode uses an `apply` command for hand-off changes; there are no config keys for it yet.
- No config defaults for `create.base` or `status/finish.target`.

## Examples

```toml
mode = "classic"

[theme]
name = "nord"

[ui]
color_enabled = true

[table]
grid = false

[dry_run]
mask_sensitive_paths = true

[create]
output = "text"
skip_existing = false

[create.path]
root = "../"
format = "{repo}_{task}"

[list]
output = "table"
field = "path"
absolute_path = false
grid = false
strict = false

[status]
output = "table"
absolute_path = false
grid = false
strict = false

[finish]
cleanup = false
remove_worktree = false
remove_branch = false
force_branch = false
merge_mode = "ff"
confirm = true

[cleanup]
remove_worktree = true
remove_branch = true
worktree_only = false
force_branch = false
confirm = true
```
