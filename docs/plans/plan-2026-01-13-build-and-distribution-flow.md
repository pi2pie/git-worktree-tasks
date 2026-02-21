---
title: "Build flow and installation strategy for gwtt command"
created-date: 2026-01-13
status: completed
agent: Zed Agent
---

## Problem

Currently, the CLI tool has two command names:

- `git-worktree-tasks` (primary name in `go.mod`)
- `gwtt` (alias defined in Cobra)

However, the Cobra `Aliases` field only works for **subcommands**, not the binary name itself. When users install via `go install`, only the `git-worktree-tasks` binary is created, and the `gwtt` shorthand doesn't work as expected.

## Goal

Establish a simple, local build and installation workflow that produces both `git-worktree-tasks` and `gwtt` binaries, and provide clear guidance for users to configure their shell environment. This is specifically for Go developers who have Go installed.

## Key Findings

1. **Cobra Aliases Limitation**: The `Aliases` field in the root command only applies to subcommand invocation, not the binary executable name.

2. **Standard Go Installation**: When users run `go install`, the binary is named after the module path and placed in `$GOPATH/bin`.

3. **Symlink Complexity**: While symlinks work, they create cleanup/management burden during uninstallation. Users need clear removal instructions.

4. **Better Approach**: Build both binaries and place them directly in `$GOPATH/bin`. Let users customize their shell environment (alias or symlink) based on their preference.

5. **Go Developers Only**: This workflow is for users with Go installed. For broader distribution, see future plans for package managers (Homebrew, apt, etc.).

6. **No CICD Required Now**: Focus on local development workflow. CICD automation can be addressed in a separate future plan.

## Proposed Implementation

### 1. Project Structure

Create a dedicated scripts directory for Go-specific workflows:

```
./scripts/
├── go-install.sh       # Build and install both binaries to $GOPATH/bin (Go developers only)
└── go-uninstall.sh     # Remove both binaries from $GOPATH/bin (Go developers only)
```

### 2. Build and Installation

**Makefile Targets:**

- `make build` — Build both `git-worktree-tasks` and `gwtt` binaries locally
- `make go-install` — Build and install both binaries to `$GOPATH/bin` (requires Go)
- `make go-uninstall` — Remove both binaries from `$GOPATH/bin` (requires Go)
- `make clean` — Clean up local build artifacts

**Shell Script (`./scripts/go-install.sh`):**

- Requires Go to be installed
- Builds both binaries
- Installs to `$GOPATH/bin` (or custom path if provided)
- Provides clear feedback on installed location
- Does NOT create symlinks (user's choice)

### 3. User Shell Configuration (Documentation)

Document three approaches for users across different shells:

#### Option A: Shell Alias (Recommended for Simplicity)

**For Bash (`~/.bashrc`):**

```bash
alias gwtt="git-worktree-tasks"
```

**For Zsh (`~/.zshrc`):**

```bash
alias gwtt="git-worktree-tasks"
```

**For Fish (`~/.config/fish/config.fish`):**

```fish
alias gwtt git-worktree-tasks
```

**Pros:**

- Easy to add and remove
- No filesystem clutter
- Works across systems
- Can be disabled quickly

**Cons:**

- Must be added to each shell profile
- Not available to non-shell tools

**Removal:** Delete the alias line from the respective shell config file

#### Option B: Symlink (For Direct Command Availability)

Create symlink manually:

```bash
ln -s $(which git-worktree-tasks) $(dirname $(which git-worktree-tasks))/gwtt
```

**Pros:**

- Works everywhere (shell, scripts, IDEs)
- Binary-level shortcut
- No shell-specific configuration needed

**Cons:**

- Requires cleanup
- May conflict with other installations
- Requires manual removal

**Removal:** Delete the symlink

```bash
rm $(dirname $(which git-worktree-tasks))/gwtt
```

### 4. Installation Instructions

Users have three paths:

**Path 1: Using Makefile (Easiest for developers)**

```bash
make go-install
# Then add alias to shell config (Bash, Zsh, or Fish), or create symlink manually
```

**Path 2: Using Script Directly**

```bash
./scripts/go-install.sh
# Then add alias to shell config (Bash, Zsh, or Fish), or create symlink manually
```

**Path 3: Standard go install (Basic, Go-native)**

```bash
go install ./
# Only git-worktree-tasks available; use alias or symlink for gwtt
```

### 5. Uninstallation Instructions

**Using Makefile:**

```bash
make go-uninstall
# Also remove alias from shell config or delete symlink if created
```

**Using Script Directly:**

```bash
./scripts/go-uninstall.sh
# Also remove alias from shell config or delete symlink if created
```

**Manual Cleanup:**

```bash
rm $(which git-worktree-tasks)
rm $(which gwtt)  # If symlink was created
```

## Next Steps

1. Create `./scripts/go-install.sh` for building and installing both binaries
2. Create `./scripts/go-uninstall.sh` for cleanup
3. Create/update Makefile with `build`, `go-install`, `go-uninstall`, `clean` targets
4. Document in README:
   - Installation methods (Makefile, script, go install)
   - Shell alias setup instructions (Bash, Zsh, Fish)
   - Manual symlink setup instructions
   - How to verify installation
   - How to uninstall

## Scope Limitations

- **Go Developers Only**: Requires Go to be installed and `$GOPATH/bin` in `$PATH`
- **No CICD automation**: This plan focuses on local workflow only
- **No pre-built binaries**: Users build from source
- **No package managers**: Homebrew, apt, etc. are separate future plans
- **No curl-to-install pattern**: Not suitable without CI/CD and hosted releases

## Success Criteria

- ✅ Makefile with clean `go-install` target for Go developers
- ✅ Both `git-worktree-tasks` and `gwtt` binaries available in `$GOPATH/bin` after `make go-install`
- ✅ Clear documentation for shell alias setup (Bash, Zsh, Fish)
- ✅ Clear documentation for manual symlink setup
- ✅ Clear documentation for uninstallation
- ✅ Users understand the different setup options and trade-offs
- ✅ Shell configuration examples work correctly for all documented shells

## Future Considerations

- **CI/CD Plan**: When ready, will add automated release builds and versioning
- **Distribution Plan**: Will handle Homebrew taps, package managers, etc.
- **Release Artifacts**: Will build pre-compiled binaries for multiple platforms (for non-Go users)

## Related Documents

- `docs/researches/research-2026-01-13-homebrew-integration.md` — Future distribution approach for broader user base
- [Verify root command configuration for gwtt alias](jobs/2026-01-13-verify-root-command-config.md) — Initial investigation

## Notes

- Users should use alias or symlink based on their preference, not forced approach
- `$GOPATH/bin` is standard location; script can accept custom paths as argument
- Document both setup and removal procedures clearly for user convenience
- This plan sets foundation for future automated distribution methods
- Go-specific targets (`go-install`, `go-uninstall`) make it clear this requires Go
- Fish shell uses different syntax for aliases (no `=` operator)
