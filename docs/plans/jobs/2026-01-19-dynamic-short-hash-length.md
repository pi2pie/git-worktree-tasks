---
title: "Dynamic short hash length"
created-date: 2026-01-19
status: completed
agent: "Codex, GitHub Copilot"
---

## Summary

Updated short hash rendering to query Git directly for the appropriate abbreviation length, replacing manual threshold-based guessing.

## Problem

The original implementation used hardcoded commit-count thresholds to decide hash length (7/8/10/12 chars). This approach:

- Required guessing appropriate thresholds
- Ignored user's `core.abbrev` configuration
- Didn't match Git's actual collision-avoidance algorithm

## Solution

Query Git directly via `git rev-parse --short HEAD` and measure the returned hash length. This:

- Uses Git's native algorithm based on packed object count
- Respects user's `core.abbrev` setting if configured
- Automatically adapts as Git's algorithm evolves

## Implementation

```go
// Before: manual thresholds
git rev-list --count --all  // count commits, then apply thresholds

// After: let Git decide
git rev-parse --short HEAD  // returns abbreviated hash at correct length
```

## Research Findings

### Git's Default Behavior

- Default abbreviation is **7 characters**
- `--abbrev-commit` extends dynamically when collisions are detected
- `core.abbrev` can be set to "auto" (default) or a specific length (min 4)

### Official Documentation (git-config)

> **core.abbrev**: Set the length object names are abbreviated to. If unspecified or set to "auto", an appropriate value is computed based on the approximate number of packed objects in your repository.

### Key Insight

Git uses **packed object count**, not commit count. Object count is typically 5-10x the commit count (commits + trees + blobs + tags). This explains why a 3k commit repo might show 8-char hashes.

### Real-World Examples

| Repository     | Commits | Objects (est.) | Hash Length |
| -------------- | ------- | -------------- | ----------- |
| Small project  | < 1K    | < 10K          | 7 chars     |
| Medium project | ~3K     | ~30K           | 7-8 chars   |
| Large project  | ~30K    | ~300K          | 8-10 chars  |
| Linux kernel   | ~1M     | ~7M            | 12 chars    |

## References

- [git-config(1)](https://git-scm.com/docs/git-config) — `core.abbrev` documentation
- [git-rev-parse(1)](https://git-scm.com/docs/git-rev-parse) — `--short` option
- [gitrevisions(7)](https://git-scm.com/docs/gitrevisions) — revision selection
- [Pro Git Book - Short SHA-1](https://git-scm.com/book/en/v2/Git-Tools-Revision-Selection#Short-SHA-1)
