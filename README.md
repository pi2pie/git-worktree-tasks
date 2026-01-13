# git-worktree-tasks

A small CLI to manage task-based Git worktrees with predictable naming and cleanup flows.

## Installation

### Quick Start (For Go Developers)

If you have Go installed and want both `git-worktree-tasks` and `gwtt` commands:

```bash
make go-install
```

Then add a shell alias (choose one):

**Bash/Zsh** — Add to `~/.bashrc` or `~/.zshrc`:
```bash
alias gwtt="git-worktree-tasks"
```

**Fish** — Add to `~/.config/fish/config.fish`:
```fish
alias gwtt git-worktree-tasks
```

### Installation Options

#### Option 1: Using Makefile (Recommended for developers)

Clone the repo and use the Makefile target:

```bash
git clone https://github.com/dev-pi2pie/git-worktree-tasks
cd git-worktree-tasks
make go-install
```

This builds and installs both `git-worktree-tasks` and `gwtt` binaries to `$GOPATH/bin`.

#### Option 2: Using Installation Script Directly

```bash
./scripts/go-install.sh
# Or with custom installation directory:
./scripts/go-install.sh /usr/local/bin
```

#### Option 3: Standard Go Install (Basic)

For just the `git-worktree-tasks` binary:

```bash
# From source (local clone)
git clone https://github.com/dev-pi2pie/git-worktree-tasks
cd git-worktree-tasks
go install ./

# Or from GitHub directly
go install github.com/dev-pi2pie/git-worktree-tasks@latest
```

This creates only the `git-worktree-tasks` binary. To use `gwtt`, see [Shell Configuration](#shell-configuration) below.

#### Option 4: Build Locally

```bash
git clone https://github.com/dev-pi2pie/git-worktree-tasks
cd git-worktree-tasks
make build
# or: go build -o dist/git-worktree-tasks ./
```

The binaries will be in the `dist/` directory.

### Requirements

- **Go 1.25.5 or higher** — Required for building from source
- **`$GOPATH/bin` in your `$PATH`** — Required for `go-install` targets (standard Go setup)

## Shell Configuration

After installation, you can configure your shell to use the `gwtt` shorthand.

### Option A: Shell Alias (Recommended)

Simple, easy to add/remove, works in any shell.

**Bash** — Add to `~/.bashrc`:
```bash
alias gwtt="git-worktree-tasks"
```

Then reload:
```bash
source ~/.bashrc
```

**Zsh** — Add to `~/.zshrc`:
```bash
alias gwtt="git-worktree-tasks"
```

Then reload:
```bash
source ~/.zshrc
```

**Fish** — Add to `~/.config/fish/config.fish`:
```fish
alias gwtt git-worktree-tasks
```

Note: Fish uses different syntax (no `=` operator).

Then reload:
```bash
exec fish
```

### Option B: Manual Symlink

Creates a binary-level symlink (works in all contexts, including scripts and IDEs).

```bash
ln -s $(which git-worktree-tasks) $(dirname $(which git-worktree-tasks))/gwtt
```

### Removing the `gwtt` Shorthand

**If using an alias:**
- Delete the alias line from your shell config file
- Reload your shell or restart your terminal

**If using a symlink:**
```bash
rm $(which gwtt)
```

**If using uninstall script:**
```bash
make go-uninstall
# or: ./scripts/go-uninstall.sh
```

## Verification

Verify your installation:

```bash
# Check git-worktree-tasks
git-worktree-tasks --version

# Check gwtt (if configured)
gwtt --version
```

## Configuration

### Theme Selection

You can set a default theme using either an environment variable or a TOML config file.

**Precedence (highest to lowest):**
1. `--theme` flag
2. `GWTT_THEME` environment variable
3. Project config (`gwtt.config.toml`, then `gwtt.toml`)
4. User config (`$HOME/.config/gwtt/config.toml`)
5. Built-in default theme

**Environment variable:**
```bash
export GWTT_THEME=nord
```

**Project config (repo root):** `gwtt.config.toml` or `gwtt.toml`
```toml
[theme]
name = "nord"
```

Example file: `examples/gwtt.config.toml`

**User config:** `$HOME/.config/gwtt/config.toml`
```toml
[theme]
name = "nord"
```

## Usage

Create a worktree for a task:

```bash
git-worktree-tasks create "my-task" --base main
```

Create in a custom location (relative to repo root or absolute path):

```bash
git-worktree-tasks create "my-task" --path ../custom-location
```

Copy a ready-to-run `cd` command after creation:

```bash
git-worktree-tasks create "my-task" --base main --copy-cd
```

Output only the worktree path (raw mode, easy to pipe; `-o` alias):

```bash
cd "$(git-worktree-tasks create \"my-task\" --base main --output raw)"
```

List worktrees (relative paths by default):

```bash
git-worktree-tasks list
```

Show detailed status:

```bash
git-worktree-tasks status
```

Show absolute paths when needed:

```bash
git-worktree-tasks list --absolute-path
git-worktree-tasks status --absolute-path
```

Finish a task (merge into target and cleanup):

```bash
git-worktree-tasks finish "my-task" --target main
```

Cleanup without merge:

```bash
git-worktree-tasks cleanup "my-task"
```

Cleanup only the worktree (keep the branch):

```bash
git-worktree-tasks cleanup "my-task" --worktree-only
```

## Development

### Build Targets

```bash
# Build both binaries to dist/ directory
make build

# Install both binaries to $GOPATH/bin (for local development)
make go-install

# Remove installed binaries
make go-uninstall

# Clean up build artifacts (removes dist/ directory)
make clean

# Show all available targets
make help
```

### Project Structure

```
./
├── main.go                 # Entry point
├── cli/                    # CLI command definitions
├── internal/               # Internal packages
├── ui/                     # UI/styling utilities
├── tui/                    # Terminal UI components
├── examples/               # Usage examples
├── scripts/                # Installation/build scripts
│   ├── go-install.sh       # Install both binaries to $GOPATH/bin
│   └── go-uninstall.sh     # Remove installed binaries
├── dist/                   # Build output directory (created by make build)
├── docs/                   # Documentation and plans
├── Makefile                # Build and installation targets
├── go.mod                  # Go module definition
└── README.md               # This file
```

**Note:** The `dist/` directory is created during `make build` and should not be committed to git (see `.gitignore`).

## Notes

- Default worktree path uses the pattern `../<repo>_<task>`.
- Create output shows relative paths by default.
- Create `--output raw` prints only the worktree path (no extra text).
- Create will report an existing task worktree and return the path instead of failing.
- List/status include the matching branch column in table and JSON output.
- Use `--output json` on list/status for machine-readable output.
- Cleanup defaults to removing both the worktree and the task branch (with separate confirmations).

## Troubleshooting

### `$GOPATH/bin` not in `$PATH`

If `make go-install` succeeds but commands don't work, ensure `$GOPATH/bin` is in your `$PATH`:

```bash
# Check if it's in PATH
echo $PATH | grep $(go env GOPATH)/bin

# If not, add to your shell config (~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish):
export PATH="$(go env GOPATH)/bin:$PATH"
```

### Permission Denied on Install

If you get permission errors during `make go-install`:

```bash
# Check if directory is writable
ls -ld $(go env GOBIN)

# Use a custom installation directory
make go-install INSTALL_DIR=$HOME/.local/bin
# or
./scripts/go-install.sh $HOME/.local/bin
```

### Shell Alias Not Working

After adding an alias, reload your shell:

```bash
# Bash
source ~/.bashrc

# Zsh
source ~/.zshrc

# Fish
exec fish
```

Or restart your terminal.
