# git-worktree-tasks

A small CLI to manage task-based Git worktrees with predictable naming and cleanup flows.

## Quick Start

```bash
# Install (requires Go 1.25.5+)
make go-install

# Create a task worktree
gwtt create "my-feature" --base main

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
- [Shell Configuration](#shell-configuration)
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

---

## Installation

### Option 1: Using Makefile (Recommended)

```bash
git clone https://github.com/dev-pi2pie/git-worktree-tasks
cd git-worktree-tasks
make go-install
```

This installs both `git-worktree-tasks` and `gwtt` binaries to `$GOPATH/bin`.

### Option 2: Using Installation Script

```bash
./scripts/go-install.sh
# Or with custom directory:
./scripts/go-install.sh /usr/local/bin
```

### Option 3: Standard Go Install

```bash
go install github.com/dev-pi2pie/git-worktree-tasks@latest
```

> **Note:** This creates only `git-worktree-tasks`. For `gwtt`, add a shell alias (see below).

### Option 4: Build Locally

```bash
make build
# Binaries in dist/
```

### Requirements

- **Go 1.25.5+** for building
- **`$GOPATH/bin` in `$PATH`** for `go-install` targets

---

## Shell Configuration

Set up the `gwtt` alias for convenience:

| Shell | Config File | Alias Syntax |
|-------|-------------|--------------|
| Bash  | `~/.bashrc` | `alias gwtt="git-worktree-tasks"` |
| Zsh   | `~/.zshrc`  | `alias gwtt="git-worktree-tasks"` |
| Fish  | `~/.config/fish/config.fish` | `alias gwtt git-worktree-tasks` |

After adding, reload your shell (`source ~/.bashrc`, `source ~/.zshrc`, or `exec fish`).

**Alternative:** Create a symlink:
```bash
ln -s $(which git-worktree-tasks) $(dirname $(which git-worktree-tasks))/gwtt
```

---

## Configuration

### Theme Selection

Set a color theme using (highest precedence first):

1. `--theme` flag
2. `GWTT_THEME` environment variable
3. Project config (`gwtt.config.toml` or `gwtt.toml` in repo root)
4. User config (`$HOME/.config/gwtt/config.toml`)
5. Built-in default

```bash
# Environment variable
export GWTT_THEME=nord

# List available themes
gwtt --themes
```

**Config file format:**
```toml
[theme]
name = "nord"
```

---

## Usage Guide

### Commands Overview

| Command   | Alias | Description |
|-----------|-------|-------------|
| `create`  |       | Create a worktree and branch for a task |
| `list`    | `ls`  | List task worktrees |
| `status`  |       | Show detailed worktree status |
| `finish`  |       | Merge a task branch into target |
| `cleanup` | `rm`  | Remove a task worktree and/or branch |

### Creating Worktrees

```bash
# Basic usage
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
| `--base` | | Base branch to create from (default: `main`) |
| `--path` | `-p` | Override worktree path |
| `--output` | `-o` | Output format: `text`, `raw` |
| `--skip-existing` | `--skip` | Reuse existing worktree |
| `--dry-run` | | Show git commands without executing |

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
```

**Status columns:** Task, Branch, Path, Base, Target, Last Commit, Dirty, Ahead, Behind

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

| Format | Description | Available In |
|--------|-------------|--------------|
| `table` | Human-readable table (default) | `list`, `status` |
| `json` | JSON array | `list`, `status` |
| `csv` | CSV with headers | `list`, `status` |
| `raw` | Single value, no decoration | `create`, `list` |
| `text` | Styled text output (default) | `create` |

### Field Selection (for `--output raw`)

When using `--output raw` with `list`, specify which field to output:

| Field | Description |
|-------|-------------|
| `path` | Worktree path (default) |
| `task` | Task name |
| `branch` | Branch name |

### Piping Examples

#### Navigate to a worktree

```bash
# Change to task worktree directory
cd "$(gwtt list my-task -o raw)"

# Or using create
cd "$(gwtt create my-task --base main -o raw)"
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

# Open in Cursor
cursor "$(gwtt list my-task -o raw)"
```

#### Scripting workflows

```bash
# Create and open in one command
code "$(gwtt create my-feature --base main -o raw)"

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
    set path (gwtt create $argv[1] --base main -o raw)
    and cd $path
end

# Bash/Zsh: Create and cd to worktree
gwtt-new() {
    local path
    path=$(gwtt create "$1" --base main -o raw) && cd "$path"
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
make go-install   # Install to $GOPATH/bin
make go-uninstall # Remove installed binaries
make clean        # Remove dist/
make help         # Show all targets
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
./scripts/go-install.sh $HOME/.local/bin
```

### Shell Alias Not Working

Reload your shell after adding the alias:
```bash
source ~/.bashrc   # Bash
source ~/.zshrc    # Zsh
exec fish          # Fish
```

---

## Notes

- Default worktree path: `../<repo>_<task>`
- Task names are slugified (lowercase, hyphens replace spaces)
- Paths are relative by default; use `--abs` for absolute
- Use `--dry-run` to preview git commands
- Global flags: `--theme`, `--nocolor`, `--themes`
