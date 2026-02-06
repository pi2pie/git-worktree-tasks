---
title: "Homebrew integration strategies for git-worktree-tasks"
date: 2026-01-13
modified: 2026-01-18
status: draft
agent: Zed Agent
---

## Goal

Explore integration options with Homebrew package manager to enable users to install `git-worktree-tasks` via `brew install`. Document both distribution strategies and their trade-offs.

## Background

Homebrew is a popular package manager for macOS and Linux. It uses "formulas" (Ruby files) to define how software is built and installed. There are two primary approaches to distribute a project via Homebrew:

1. **Global Homebrew PR** — Submit formula to official `homebrew-core` repository
2. **Custom Tap** — Maintain a separate Homebrew tap repository (e.g., `homebrew-git-worktree-tasks`)

## Key Findings

### Approach 1: Global Homebrew PR (homebrew-core)

**Process:**

- Submit `git-worktree-tasks.rb` formula to `homebrew/homebrew-core` GitHub repository
- Homebrew maintainers review and merge
- Users install directly: `brew install git-worktree-tasks`

**Advantages:**

- Single, standard installation method
- Official Homebrew distribution
- No additional repository to maintain
- Higher discoverability

**Disadvantages:**

- Longer review process (days to weeks)
- Strict Homebrew guidelines and requirements
- Less control over updates and versioning
- Requires stable, tagged releases in main repository

**Requirements:**

- Project must be stable and actively maintained
- Clear, tagged releases (e.g., `v0.0.5`)
- Open source with compatible license
- Public GitHub repository

### Approach 2: Custom Tap Repository

**Process:**

- Create separate GitHub repository: `pi2pie/homebrew-git-worktree-tasks`
- Maintain `Formula/git-worktree-tasks.rb` in the tap repository
- Users add tap first: `brew tap pi2pie/git-worktree-tasks`
- Then install: `brew install git-worktree-tasks`

**Advantages:**

- Complete control over formula and updates
- Faster deployment (no review process)
- Can update formula independently of main project
- Easier for private or pre-release testing
- Can include additional formulas if needed

**Disadvantages:**

- Users must remember to tap first (extra step)
- Requires maintaining additional repository
- Lower discoverability compared to homebrew-core
- User must manage tap updates separately

**Requirements:**

- Separate GitHub repository
- Proper repository naming convention
- CI/CD to automate formula updates

## Formula Structure

### Basic Formula (Build from Source)

```ruby
class GitWorktreeTasks < Formula
  desc "Task-based git worktree helper"
  homepage "https://github.com/pi2pie/git-worktree-tasks"
  url "https://github.com/pi2pie/git-worktree-tasks/archive/refs/tags/v#{version}.tar.gz"
  sha256 "<checksum>"
  version "<version>"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"git-worktree-tasks", "."
    bin.install_symlink bin/"git-worktree-tasks" => "gwtt"
  end

  test do
    system "#{bin}/git-worktree-tasks", "--version"
  end
end
```

### Optimized Formula (Pre-built Binaries)

Requires goreleaser or similar CI/CD to generate platform-specific binaries.

```ruby
class GitWorktreeTasks < Formula
  desc "Task-based git worktree helper"
  homepage "https://github.com/pi2pie/git-worktree-tasks"
  version "<version>"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/pi2pie/git-worktree-tasks/releases/download/v#{version}/git-worktree-tasks_#{version}_darwin_arm64.tar.gz"
      sha256 "<arm64_checksum>"
    else
      url "https://github.com/pi2pie/git-worktree-tasks/releases/download/v#{version}/git-worktree-tasks_#{version}_darwin_amd64.tar.gz"
      sha256 "<amd64_checksum>"
    end
  end

  on_linux do
    url "https://github.com/pi2pie/git-worktree-tasks/releases/download/v#{version}/git-worktree-tasks_#{version}_linux_amd64.tar.gz"
    sha256 "<linux_checksum>"
  end

  def install
    bin.install "git-worktree-tasks"
    bin.install_symlink bin/"git-worktree-tasks" => "gwtt"
  end

  test do
    system "#{bin}/git-worktree-tasks", "--version"
  end
end
```

## Installation Methods

### Custom Tap (Recommended for Current Stage)

```bash
brew tap pi2pie/git-worktree-tasks
brew install git-worktree-tasks
```

One-liner variant:

```bash
brew install pi2pie/git-worktree-tasks/git-worktree-tasks
```

### Global Homebrew (Future, After Stabilization)

```bash
brew install git-worktree-tasks
```

## Implementation Roadmap

### Phase 1: Foundation (Current)

- ✅ Establish build process (Makefile)
- ⏳ Set up stable CI/CD with goreleaser

### Phase 2: Custom Tap

- Create `homebrew-git-worktree-tasks` repository
- Write and test `git-worktree-tasks.rb` formula
- Document custom tap installation in README
- Automate formula updates via CI/CD

### Phase 3: Homebrew Core (Optional, Future)

- Prepare project for homebrew-core submission
- Follow Homebrew guidelines and standards
- Submit PR to `homebrew/homebrew-core`
- Maintain formula after acceptance

## Technical Considerations

### Naming Conventions

- **Formula name**: `git-worktree-tasks` (kebab-case in formula file)
- **Formula class**: `GitWorktreeTasks` (CamelCase in Ruby)
- **Tap repository**: `homebrew-git-worktree-tasks`
- **Both commands**: `git-worktree-tasks` and `gwtt` available after install

### SHA256 Checksums

- **Build from source**: Use GitHub archive tarball SHA256
- **Pre-built binaries**: Calculate SHA256 for each platform binary
- Can be automated in release process

### Dependencies

- Go build requires `go` as build dependency
- Pre-built binaries have no dependencies
- Both approaches should test installation

## Implications and Recommendations

1. **Start with Custom Tap** — Lower barrier to entry, full control, faster iteration
2. **Automate Formula Updates** — Use CI/CD to update formula when releases are tagged
3. **Use Pre-built Binaries** — Faster installation for users, requires goreleaser setup
4. **Plan for homebrew-core** — Design with Homebrew standards in mind from the start
5. **Symmetric Command Names** — Ensure both `git-worktree-tasks` and `gwtt` are available via symlink

## Open Questions

- Should we target homebrew-core submission eventually?
- What's the versioning strategy (semver, etc.)?
- How frequently will releases be made?
- Should we support other package managers (apt, pacman, etc.) in future?

## References

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Tap Documentation](https://docs.brew.sh/Taps)
- [Homebrew Core Contributing Guidelines](https://github.com/Homebrew/homebrew-core/blob/master/CONTRIBUTING.md)

## Related Plans

- [Build and distribution flow for gwtt alias support](plans/plan-2026-01-13-build-and-distribution-flow.md) — Parent plan for overall distribution strategy
