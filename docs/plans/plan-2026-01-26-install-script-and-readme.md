---
title: "Install Script + README Clarity for Binary Names"
date: 2026-01-26
status: active
agent: codex
---

## Goal
Create a dedicated install/uninstall experience that detects platform, downloads the correct release asset, and clarifies the `gwtt` vs `git-worktree-tasks` command naming in the README.

## Context
- Release assets built via GoReleaser ship the binary as `gwtt` (or `gwtt.exe`).
- `go install` yields a `git-worktree-tasks` binary name in GOPATH.
- Users need a clear explanation of why both names exist and how to use them.
- Existing `scripts/go-install.sh` and `scripts/go-uninstall.sh` are wired into the Makefile and use a Go build flow, not release assets.
- New install/uninstall scripts should download the correct release asset via `curl` or `wget`, based on platform detection.
- Asset naming and packaging are controlled by `.goreleaser.yml`, which must be aligned with install logic.

## Plan
1. Review current release assets, README install wording, and existing scripts to understand gaps and naming inconsistencies.
2. Define the install/uninstall approach:
   - Platform/arch detection and asset naming rules.
   - Release endpoint selection: use GitHub `releases/latest` to avoid prereleases/drafts and pick stable assets.
   - Where binaries are installed (e.g., `~/.local/bin` or user-specified).
   - How to handle `gwtt` vs `git-worktree-tasks` (symlink, rename, or document both).
   - Whether to replace or supplement existing go-install/go-uninstall scripts and Makefile targets.
3. Implement scripts under `./scripts` (install + uninstall), and update README to:
   - Explain the naming difference.
   - Provide install methods (script + `go install`).
   - Provide uninstall steps.
   - Note that release-asset install uses `curl` or `wget`, and document prerequisites.

## Draft Install Logic (to implement)
- **Policy (explicit)**
  - Primary command: `gwtt` (short and consistent with GoReleaser release assets).
  - Release-asset install: downloads and installs `gwtt` into the **current directory by default** (or a provided path).
  - `go install` behavior remains: it produces `git-worktree-tasks`; README explains the mismatch and how to add `gwtt` or `git-worktree-tasks` via manual alias/symlink.
- **Detect platform/arch**
  - `uname -s` → `darwin|linux|windows` (treat MSYS/Cygwin/Git Bash as `windows`).
  - `uname -m` → map `x86_64|amd64` → `amd64`, `arm64|aarch64` → `arm64`.
- **Resolve stable release + asset**
  - Fetch `https://api.github.com/repos/pi2pie/git-worktree-tasks/releases/latest`.
  - Extract `tag_name` and build asset name:
    - `git-worktree-tasks_${tag}_${os}_${arch}.tar.gz` (macOS/Linux)
    - `git-worktree-tasks_${tag}_${os}_${arch}.zip` (Windows)
  - Select `assets[].browser_download_url` that matches the asset name.
- **Verify download**
  - Download `checksums.txt` from the same release and verify the archive checksum before extracting.
- **Download & install**
  - Use `curl -fsSL` or `wget -qO` to download to a temp dir.
  - Extract `gwtt` (or `gwtt.exe`) from the archive.
  - Install to target dir (default current directory, or user-specified arg), ensure executable bit on *nix.
  - Do not create symlinks or install man pages; document optional manual setup in README.

## Packaging Notes
- `.goreleaser.yml` should include man pages for manual install (even if the script does not auto-install them).
- **Uninstall**
  - Remove installed `gwtt` (or `gwtt.exe`) from the target dir (default current directory).
  - Keep uninstall script simple and path-scoped; no global cleanup beyond installed files.

## Makefile Alignment
- Add new Makefile targets: `install` and `uninstall` for the release-asset scripts.
- Keep existing Go-based install targets (e.g., `go-install`, `go-uninstall`) for developers or `go install` flow.

## Risks / Open Questions
- How to handle Windows PATH guidance for `.exe` installs.

## Success Criteria
- A single, clear install path that auto-selects platform/arch and fetches the correct asset.
- README clearly explains binary naming and offers consistent commands.
- Uninstall instructions remove installed artifacts cleanly.
