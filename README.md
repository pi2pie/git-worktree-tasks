# git-worktree-tasks

A small CLI to manage task-based Git worktrees with predictable naming and cleanup flows.

## Quick Start

```bash
# Install latest release (run from your target folder)
cd $HOME/.local/bin
curl -fsSL https://raw.githubusercontent.com/pi2pie/git-worktree-tasks/main/scripts/install.sh | bash

# Create a task worktree (defaults to current branch)
gwtt create "my-feature"

# List all worktrees
gwtt list

# Show status
gwtt status

# Cleanup when done
gwtt cleanup "my-feature"
```

---

## Table of Contents

- [Installation](#installation)
- [Binary Naming and Shell Configuration](#binary-naming-and-shell-configuration)
- [Configuration](#configuration)
- [Usage Guide](#usage-guide)
  - [Commands Overview](#commands-overview)
  - [Creating Worktrees](#creating-worktrees)
  - [Listing Worktrees](#listing-worktrees)
  - [Checking Status](#checking-status)
  - [Finishing Tasks](#finishing-tasks)
  - [Cleanup](#cleanup)
- [Output Formats & Piping](#output-formats--piping)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [License](#license)

---

## Installation

### Option 1: Install Latest Release (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/pi2pie/git-worktree-tasks/main/scripts/install.sh | bash
# Or with wget:
wget -qO- https://raw.githubusercontent.com/pi2pie/git-worktree-tasks/main/scripts/install.sh | bash
```

Run from the directory where you want `gwtt` to live:

```bash
cd $HOME/.local/bin
curl -fsSL https://raw.githubusercontent.com/pi2pie/git-worktree-tasks/main/scripts/install.sh | bash
```

Custom install directory (default is current folder):

```bash
curl -fsSL https://raw.githubusercontent.com/pi2pie/git-worktree-tasks/main/scripts/install.sh | bash -s -- ~/.local/bin
```

### Option 2: Using Makefile (Release Assets)

```bash
make install
```

Note: `make install` installs into the current directory (the repo root if you run it here).

Directly from the repo:

```bash
./scripts/install.sh
```

### Option 3: Standard Go Install (Build From Source)

```bash
go install github.com/pi2pie/git-worktree-tasks@latest
```

> **Note:** `go install` creates `git-worktree-tasks`. Release assets install `gwtt` and can optionally create a `git-worktree-tasks` symlink.

### Option 4: Build Locally

```bash
make build
# Binaries in dist/
```

### Uninstall

```bash
./scripts/uninstall.sh
# Or with Makefile:
make uninstall
```

If the binary is not in the current directory, pass the install path:

```bash
./scripts/uninstall.sh $HOME/.local/bin
```

Go install removal:

```bash
./scripts/go-uninstall.sh
# Or:
make go-uninstall
```

### Requirements

- **curl or wget** for release install
- **tar/unzip** for release install archives
- **sha256sum or shasum** for release checksum verification
- **Go 1.25.5+** for building from source
- **`$GOPATH/bin` in `$PATH`** for `go-install` targets

> **Windows PATH Note:** If you install `gwtt.exe` into a custom folder (e.g., `C:\Users\<you>\bin`), add that folder to your PATH and open a new terminal to pick it up.

> [!Note]
> **Ownership Change (v0.0.7+)**
>
> The repository ownership changed after v0.0.6. The old `dev-pi2pie` path no longer exists, so use the new module path for all installs and imports:
>
> - **v0.0.7 and later:** `github.com/pi2pie/git-worktree-tasks`

## Binary Naming and Shell Configuration

- Release assets ship the `gwtt` binary.
- `go install` produces `git-worktree-tasks`.
- If you prefer `git-worktree-tasks`, create an alias or symlink manually.

Set up the `gwtt` alias for convenience:

| Shell | Config File                  | Alias Syntax                      |
| ----- | ---------------------------- | --------------------------------- |
| Bash  | `~/.bashrc`                  | `alias gwtt="git-worktree-tasks"` |
| Zsh   | `~/.zshrc`                   | `alias gwtt="git-worktree-tasks"` |
| Fish  | `~/.config/fish/config.fish` | `alias gwtt git-worktree-tasks`   |

After adding, reload your shell (`source ~/.bashrc`, `source ~/.zshrc`, or `exec fish`).

**Alternative:** Create a symlink:

```bash
ln -s $(which git-worktree-tasks) $(dirname $(which git-worktree-tasks))/gwtt
```

---

## Configuration

### Precedence

Settings resolve in this order (highest precedence first):

1. CLI flags (for example `--theme`, `--mode` / `-m`, `--mask-sensitive-paths`, `--no-mask-sensitive-paths`)
2. Environment variables
3. Project config (`gwtt.config.toml` or `gwtt.toml` in repo root)
4. User config (`$HOME/.config/gwtt/config.toml`)
5. Built-in defaults

```bash
# Environment variable
export GWTT_THEME=nord

# Disable color output
export GWTT_COLOR=0

# Mode selection
export GWTT_MODE=codex

# Dry-run path masking (1/true/on/yes to enable, 0/false/off/no to disable)
export GWTT_DRY_RUN_MASK_SENSITIVE_PATHS=1

# List available themes
gwtt --themes
```

### Theme Selection

```toml
[theme]
name = "nord"
```

### Mode Selection

```toml
mode = "classic" # or "codex"
```

### Other Defaults

Common defaults you can set once:

```toml
[ui]
color_enabled = true

[table]
grid = false

[dry_run]
mask_sensitive_paths = true # mask $HOME/%USERPROFILE% prefixes in --dry-run output

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
merge_mode = "ff" # ff, no-ff, squash, rebase
confirm = true # set false to bypass prompts (same as --yes)
cleanup = false
remove_worktree = false
remove_branch = false
force_branch = false

[cleanup]
confirm = true # set false to bypass prompts (same as --yes)
remove_worktree = true
remove_branch = true
worktree_only = false
force_branch = false
```

`[dry_run].mask_sensitive_paths` defaults to `true`. Set it to `false` if you need raw absolute paths in `--dry-run` output.  
When enabled, home-prefixed paths are rendered as `$HOME/...` on POSIX and `%USERPROFILE%\\...` on Windows.
You can override this per-invocation with `--mask-sensitive-paths=true|false`, `--no-mask-sensitive-paths`, or via `GWTT_DRY_RUN_MASK_SENSITIVE_PATHS`.
For bool flags, prefer `--mask-sensitive-paths=false` (with `=`) rather than `--mask-sensitive-paths false`.

### Config File Location

Project: `gwtt.config.toml` or `gwtt.toml` in the repo root  
User: `$HOME/.config/gwtt/config.toml`

**Minimal config:**

```toml
[theme]
name = "nord"
```

---

## Usage Guide

### Commands Overview

| Command   | Alias | Description                                                          |
| --------- | ----- | -------------------------------------------------------------------- |
| `apply`   |       | Apply non-destructive changes between Codex worktree and local checkout (codex mode only) |
| `overwrite` |     | Destructively replace destination with source changes in codex mode |
| `create`  |       | Create a worktree and branch for a task                              |
| `list`    | `ls`  | List task worktrees                                                  |
| `status`  |       | Show detailed worktree status                                        |
| `finish`  |       | Merge a task branch into target                                      |
| `cleanup` | `rm`  | Remove a task worktree and/or branch                                 |

### Creating Worktrees

```bash
# Basic usage (defaults to current branch)
gwtt create "my-task"

# Explicit base override
gwtt create "my-task" --base main

# Reuse existing worktree (no error if exists)
gwtt create "my-task" --skip-existing

# Custom path
gwtt create "my-task" --path ../custom-location

# Preview without executing
gwtt create "my-task" --dry-run
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--base` | | Base branch to create from (default: current branch) |
| `--path` | `-p` | Override worktree path |
| `--output` | `-o` | Output format: `text`, `raw` |
| `--skip-existing` | `--skip` | Reuse existing worktree |
| `--dry-run` | | Show git commands without executing |

**Notes:**

- The default base is the current local branch (for example `main`, `master`, or `dev`).
- If you are in a detached HEAD state, you must pass `--base` explicitly.

### Listing Worktrees

```bash
# List all worktrees
gwtt list

# Filter by task name (contains match)
gwtt list "my-task"

# Exact match
gwtt list "my-task" --strict

# Filter by branch
gwtt list --branch feature-branch

# Show absolute paths
gwtt list --abs

# Grid borders in table
gwtt list --grid

# Codex mode: list Codex-managed worktrees for this repo
gwtt --mode codex list
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Format: `table`, `json`, `csv`, `raw` |
| `--field` | `-f` | Raw output field: `path`, `task`, `branch` |
| `--branch` | | Filter by branch name |
| `--absolute-path` | `--abs` | Show absolute paths |
| `--strict` | | Require exact task match |
| `--grid` | | Render table with grid borders |

### Checking Status

```bash
# Show status of all worktrees
gwtt status

# Filter by task
gwtt status "my-task"

# Compare against specific target branch
gwtt status --target main

# Filter by exact task name
gwtt status --task "my-task"

# Codex mode: show Codex-managed worktree status
gwtt --mode codex status
```

**Status columns:** Task, Branch, Path, Modified Time (RFC3339 UTC), Base, Target, Last Commit, Dirty, Ahead, Behind

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Format: `table`, `json`, `csv` |
| `--target` | | Target branch for ahead/behind comparison |
| `--task` | | Filter by task name (enables strict match) |
| `--branch` | | Filter by branch name |
| `--absolute-path` | `--abs` | Show absolute paths |
| `--strict` | | Require exact task match |
| `--grid` | | Render table with grid borders |

### Finishing Tasks

```bash
# Merge task branch into target
gwtt finish "my-task" --target main

# Merge with cleanup (remove worktree + branch)
gwtt finish "my-task" --target main --cleanup

# Merge strategies
gwtt finish "my-task" --no-ff        # No fast-forward
gwtt finish "my-task" --squash       # Squash commits
gwtt finish "my-task" --rebase       # Rebase before merge

# Skip confirmation
gwtt finish "my-task" --cleanup --yes
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--target` | Target branch (default: current branch) |
| `--cleanup` | Remove worktree and branch after merge |
| `--remove-worktree` | Remove only the worktree after merge |
| `--remove-branch` | Remove only the branch after merge |
| `--force-branch` | Force delete branch (`-D` instead of `-d`) |
| `--no-ff` | Use `--no-ff` merge |
| `--squash` | Use `--squash` merge |
| `--rebase` | Rebase task branch onto target first |
| `--yes` | Skip confirmation prompts |
| `--dry-run` | Show git commands without executing |

### Applying Changes (Codex Mode)

```bash
# Non-destructive apply (default direction: worktree -> local)
gwtt --mode codex apply <opaque-id>

# Reverse non-destructive apply (local -> worktree)
gwtt --mode codex apply <opaque-id> --to worktree

# Destructive overwrite (requires confirmation unless --yes)
gwtt --mode codex overwrite <opaque-id> --to local
gwtt --mode codex overwrite <opaque-id> --to worktree --yes

# Compatibility alias for overwrite
gwtt --mode codex apply <opaque-id> --to worktree --force --yes

# Preview with structured plan + command echo
gwtt --mode codex apply <opaque-id> --dry-run
```

**Notes:**

- In codex mode, `<opaque-id>` is the directory directly under `$CODEX_HOME/worktrees`.
- `apply` is non-destructive and will not switch direction automatically on conflict.
- On conflict, `apply` exits with a next-step hint for `overwrite --to ...`.
- `overwrite` resets/cleans the destination before transfer and is destructive by design.
- `--dry-run` prints `plan`, `preflight`, and `actions` sections, then echoes the underlying git/copy operations.

### Cleanup

```bash
# Remove worktree and branch (with confirmation)
gwtt cleanup "my-task"

# Remove only the worktree (keep branch)
gwtt cleanup "my-task" --worktree-only

# Force delete branch
gwtt cleanup "my-task" --force-branch

# Skip confirmation
gwtt cleanup "my-task" --yes

# Preview without executing
gwtt cleanup "my-task" --dry-run

# Codex mode: remove a Codex-managed worktree by opaque id
gwtt --mode codex cleanup <opaque-id>
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--remove-worktree` | Remove the task worktree (default: true) |
| `--remove-branch` | Remove the task branch (default: true) |
| `--worktree-only` | Remove only worktree, keep branch |
| `--force-branch` | Force delete branch (`-D`) |
| `--yes` | Skip confirmation prompts |
| `--dry-run` | Show git commands without executing |

---

## Output Formats & Piping

The `--output` (`-o`) and `--field` (`-f`) flags enable powerful shell integrations.

### Output Formats

| Format  | Description                    | Available In     |
| ------- | ------------------------------ | ---------------- |
| `table` | Human-readable table (default) | `list`, `status` |
| `json`  | JSON array                     | `list`, `status` |
| `csv`   | CSV with headers               | `list`, `status` |
| `raw`   | Single value, no decoration    | `create`, `list` |
| `text`  | Styled text output (default)   | `create`         |

### Field Selection (for `--output raw`)

When using `--output raw` with `list`, specify which field to output:

| Field    | Description             |
| -------- | ----------------------- |
| `path`   | Worktree path (default) |
| `task`   | Task name               |
| `branch` | Branch name             |

### Piping Examples

#### Navigate to a worktree

```bash
# Change to task worktree directory
cd "$(gwtt list my-task -o raw)"

# Or using create
cd "$(gwtt create my-task -o raw)"
```

#### Copy to clipboard

```bash
# Copy worktree path
gwtt list my-task -o raw | pbcopy                    # macOS
gwtt list my-task -o raw | xclip -selection clipboard # Linux

# Copy task name
gwtt list my-task -o raw -f task | pbcopy

# Copy branch name
gwtt list my-task -o raw -f branch | pbcopy
```

#### Open in editor

```bash
# Open worktree in VS Code
code "$(gwtt list my-task -o raw)"

# Open Worktree in Zed
zed "$(gwtt list my-task -o raw)"

# Open in Cursor
cursor "$(gwtt list my-task -o raw)"
```

#### Scripting workflows

```bash
# Create and open in one command
code "$(gwtt create my-feature -o raw)"

# List all branches as plain text
gwtt list -o json | jq -r '.[].branch'

# Get paths for all worktrees
gwtt list -o json | jq -r '.[].path'

# Filter dirty worktrees
gwtt status -o json | jq '.[] | select(.dirty == true)'

# Count worktrees ahead of target
gwtt status -o json | jq '[.[] | select(.ahead > 0)] | length'
```

#### Shell function examples

```bash
# Fish: Create and cd to worktree
function gwtt-new
    set path (gwtt create $argv[1] -o raw)
    and cd $path
end

# Bash/Zsh: Create and cd to worktree
gwtt-new() {
    local path
    path=$(gwtt create "$1" -o raw) && cd "$path"
}
```

### Raw Output Fallback

When using `--output raw` with `list`:

- If no matching worktree exists but the branch does, returns the main worktree path
- Requires either a task filter or `--branch` flag

```bash
# Returns path even if no worktree exists (fallback to main repo)
gwtt list feature-branch -o raw
```

---

## Development

### Build Targets

```bash
make build        # Build binaries to dist/
make install      # Install gwtt from latest release
make uninstall    # Remove release install
make go-install   # Install to $GOPATH/bin
make go-uninstall # Remove installed binaries
make clean        # Remove dist/
make help         # Show all targets
```

### Testing and Linting

```bash
go test ./...
golangci-lint run
```

### Project Structure

```
├── main.go           # Entry point
├── cli/              # CLI command definitions
├── internal/         # Internal packages (config, git, worktree)
├── ui/               # UI/styling utilities
├── tui/              # Terminal UI components (preview)
├── examples/         # Example configs and shell functions
├── scripts/          # Installation scripts
├── docs/             # Documentation and plans
├── Makefile          # Build targets
└── go.mod            # Go module definition
```

---

## Troubleshooting

### `$GOPATH/bin` not in `$PATH`

```bash
# Check if it's in PATH
echo $PATH | grep $(go env GOPATH)/bin

# Add to shell config if missing
export PATH="$(go env GOPATH)/bin:$PATH"
```

### Permission Denied

```bash
# Use custom directory
./scripts/install.sh $HOME/.local/bin
```

### Shell Alias Not Working

Reload your shell after adding the alias:

```bash
source ~/.bashrc   # Bash
source ~/.zshrc    # Zsh
exec fish          # Fish
```

---

## License

This project is licensed under the MIT License — see the [LICENSE](https://github.com/pi2pie/git-worktree-tasks/blob/main/LICENSE) file for details.

---

## Notes

- Default worktree path: `../<repo>_<task>`
- Task names are slugified (lowercase, hyphens replace spaces)
- Paths are relative by default; use `--abs` for absolute
- Use `--dry-run` to preview git commands
- Global flags: `--mode` (`-m`), `--theme`, `--nocolor`, `--themes`
