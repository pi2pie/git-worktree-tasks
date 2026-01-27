---
title: "Extensible config options beyond theme"
date: 2026-01-27
status: draft
agent: codex
---

## Goal
Identify optional configuration areas beyond theme that could be safely exposed to users, and outline a TOML shape that stays extensible and idiomatic for Go.

## Key Findings
- **Current config scope is theme-only** with fixed file locations and precedence (flag > env > project config > user config > default). This flow can be generalized for other settings without adding new dependencies. [^theme-config] [^readme-config]
- **Most CLI flags map cleanly to defaults** that users might want to set once: output format (`list`, `status`, `create`), table grid, absolute paths, strict task matching, default targets, and cleanup behaviors. [^cli-create] [^cli-list] [^cli-status] [^cli-finish] [^cli-cleanup]
- **Naming and path conventions are centralized** in `internal/worktree/naming.go`, which makes them a natural candidate for a configurable strategy (prefix, separator, or template) but requires compatibility safeguards because `TaskFromPath` assumes a fixed prefix format. [^worktree-naming]
- **UI toggles are global** (`--nocolor`, theme selection) and could live under a `[ui]` or `[theme]` block; doing so keeps visual settings separated from command behavior. [^cli-root] [^readme-config]
- **Create base should remain current-branch-only** per product constraint; avoid exposing a config default for `--base`. [^cli-create]
- **Target defaults are context-dependent** (`status`/`finish` use the current branch when `--target` is unset), so making a config default for `target` can be surprising if users expect “current branch” semantics. [^cli-status] [^cli-finish]
- **Avoid config defaults for `target`**; keep target selection bound to the user’s current branch unless explicitly provided at runtime. [^cli-status] [^cli-finish]
- **Path templates should remain discoverable** by requiring a `{task}` placeholder; unrestricted templates make reverse lookup unreliable. [^worktree-naming]
- **Merge mode should be exclusive**; only one of `ff` (default), `no-ff`, `squash`, or `rebase` should be active to avoid ambiguous or conflicting semantics. [^cli-finish]

## Candidate Optional Config Sections
The following are low-risk candidates because they are already CLI flags or deterministic defaults:

1. **[ui]**
   - `theme.name` (existing)
   - `color.enabled` (mirrors `--nocolor`)
   - `table.grid` (default for list/status grid tables)

2. **[create]**
   - `output` (`text` or `raw`)
   - `skip_existing` (default for `--skip-existing`)
   - `path.template` or `path.root` (override default `../<repo>_<task>`; require `{task}` placeholder)

3. **[list]**
   - `output` (`table`, `json`, `csv`, `raw`)
   - `field` (default for raw output)
   - `absolute_path`
   - `grid`
   - `strict`

4. **[status]**
   - `output` (`table`, `json`, `csv`)
   - `absolute_path`
   - `grid`
   - `strict`

5. **[finish]**
   - `cleanup` (default merge cleanup toggle)
   - `remove_worktree`, `remove_branch`, `force_branch`
   - `merge.mode` (`ff`, `no-ff`, `squash`, `rebase`)
   - `confirm` (default for prompt behavior)

6. **[cleanup]**
   - `remove_worktree`, `remove_branch`, `worktree_only`
   - `force_branch`
   - `confirm`

## Extensible TOML Shape (Draft)
A minimal, additive structure that avoids breaking existing `theme` parsing:

```toml
[theme]
name = "nord"

[ui]
color_enabled = true

table_grid = false

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

Notes:
- Keep `theme` as-is to preserve compatibility.
- `create.path.format` would require updating `TaskFromPath` logic or providing a parallel lookup strategy to avoid breaking list/status discovery for non-default naming.

## Implications or Recommendations
- **Start with config for existing CLI defaults** (output format, grid, absolute paths, confirm toggles) because they are low-risk and align with current option handling.
- **Treat naming/path customization as a second phase** due to the coupling in `TaskFromPath` and potential mismatch with existing worktrees; consider supporting a templated suffix/prefix that still allows parsing. [^worktree-naming]
- **Keep precedence consistent with theme** (flags > env > project > user > default) to reduce user confusion and keep logic predictable. [^theme-config]
- **Avoid config defaults that override context-derived values** like current-branch `target`, unless the product explicitly wants to change that behavior. [^cli-status] [^cli-finish]
- **Enforce a single merge strategy** if `merge_mode` is introduced (exactly one of `ff`, `no-ff`, `squash`, `rebase`) to avoid conflicting or ambiguous behavior; update CLI validation to match. [^cli-finish]
- **Prefer explicit, typed config structs** for each section to keep Go code simple and avoid "magic" config behavior. (Inference based on current code style.)

## Open Questions
- None.

## References
[^theme-config]: `internal/config/theme.go`
[^cli-root]: `cli/root.go`
[^cli-create]: `cli/create.go`
[^cli-list]: `cli/list.go`
[^cli-status]: `cli/status.go`
[^cli-finish]: `cli/finish.go`
[^cli-cleanup]: `cli/cleanup.go`
[^worktree-naming]: `internal/worktree/naming.go`
[^readme-config]: `README.md`
