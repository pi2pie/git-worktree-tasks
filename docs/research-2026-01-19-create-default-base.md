---
title: "Create default base selection"
date: 2026-01-19
status: in-progress
agent: Codex
---

## Goal
Define how `create` selects a base branch by default when invoked from different contexts, including older Git default branches (`master`) and detached HEAD cases.

## Key Findings
- Current behavior: `create` defaults to base `main` unconditionally.
- The command currently requires a named branch (errors on detached HEAD) before any worktree creation occurs.
- Worktrees always share the same Git common directory; creating a worktree from another worktree does not create a new repo, it creates another worktree in the same repo.

## Implications or Recommendations
- Prefer base = current branch when on a named branch.
- If detached HEAD, either:
  - Error with a clear message and require `--base`, or
  - Fall back to the repository default branch.
- If current branch cannot be determined, determine the repository default branch (prefer `main`, then `master`, or use Git’s symbolic ref to origin/HEAD when available).
- `--base` always overrides defaults.

## Open Questions
- Should detached HEAD be allowed with an explicit `--base`, or should we require checkout of a named branch?
- Should we follow local default branch or remote `origin/HEAD` when deciding the fallback?

## References
- (To be added; consider `git symbolic-ref refs/remotes/origin/HEAD` and `git config --get init.defaultBranch` behavior.)

## Related Plans
- docs/plans/plan-2026-01-18-open-source-readiness.md
