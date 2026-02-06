---
title: "Add project_root.go unit tests and documentation"
date: 2026-01-27
status: completed
agent: copilot
---

# Add project_root.go unit tests and documentation

## Summary

Follow-up to the project config root resolution fix. Adds:

- Direct unit tests for `project_root.go` edge cases
- Documentation comments clarifying `.git` file handling for worktrees/submodules

## Why

The initial fix had coverage only via integration tests in `config_test.go`. Direct unit tests improve:

- Edge case coverage (no `.git` found, `.git` as file)
- Documentation of intentional behavior
- Faster feedback during refactoring

## Files

- `internal/config/project_root.go` — added doc comments
- `internal/config/project_root_test.go` — new test file

## Test Cases

| Test                            | Scenario                                 |
| ------------------------------- | ---------------------------------------- |
| `TestFindRepoRoot_Found`        | Nested directory with `.git` at ancestor |
| `TestFindRepoRoot_NotFound`     | No `.git` in hierarchy                   |
| `TestFindRepoRoot_GitFile`      | `.git` is a file (worktree scenario)     |
| `TestHasGitDir_Directory`       | `.git` directory exists                  |
| `TestHasGitDir_File`            | `.git` file exists                       |
| `TestHasGitDir_NotExists`       | No `.git` present                        |
| `TestProjectConfigRoot_WithGit` | Returns repo root from subdir            |
| `TestProjectConfigRoot_NoGit`   | Falls back to cwd                        |

## Related Plans

- [plan-2026-01-27-fix-project-config-root-resolution.md](../plan-2026-01-27-fix-project-config-root-resolution.md)
