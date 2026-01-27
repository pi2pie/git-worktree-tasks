---
title: "gwtt configuration schema"
date: 2026-01-27
status: draft
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

## Schema
### [theme]
- `name` (string, default: `"default"`)

### [ui]
- `color_enabled` (bool, default: `true`)

### [table]
- `grid` (bool, default: `false`)

### [create]
- `output` (string enum: `text`, `raw`; default: `text`)
- `skip_existing` (bool, default: `false`)

#### [create.path]
- `root` (string, default: `"../"`)
- `format` (string, default: `"{repo}_{task}"`)
  - Must include `{task}`.
  - `{repo}` is optional.

### [list]
- `output` (string enum: `table`, `json`, `csv`, `raw`; default: `table`)
- `field` (string enum: `path`, `task`, `branch`; default: `path`)
- `absolute_path` (bool, default: `false`)
- `grid` (bool, default: `false`)
- `strict` (bool, default: `false`)

### [status]
- `output` (string enum: `table`, `json`, `csv`; default: `table`)
- `absolute_path` (bool, default: `false`)
- `grid` (bool, default: `false`)
- `strict` (bool, default: `false`)

### [finish]
- `cleanup` (bool, default: `false`)
- `remove_worktree` (bool, default: `false`)
- `remove_branch` (bool, default: `false`)
- `force_branch` (bool, default: `false`)
- `merge_mode` (string enum: `ff`, `no-ff`, `squash`, `rebase`; default: `ff`)
- `confirm` (bool, default: `true`)

### [cleanup]
- `remove_worktree` (bool, default: `true`)
- `remove_branch` (bool, default: `true`)
- `worktree_only` (bool, default: `false`)
- `force_branch` (bool, default: `false`)
- `confirm` (bool, default: `true`)

## Decisions
- `create.path.format` must include `{task}` to preserve task discovery.
- `merge_mode` is exclusive; only one strategy may be active at a time.
- No config defaults for `create.base` or `status/finish.target`.

## Examples
```toml
[theme]
name = "nord"

[ui]
color_enabled = true

[table]
grid = false

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
