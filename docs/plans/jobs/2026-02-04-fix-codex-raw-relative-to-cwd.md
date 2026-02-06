---
title: "Fix codex raw paths to be relative to cwd"
date: 2026-02-04
status: completed
agent: codex
---

## Summary

Adjusted codex-mode `list --output raw` to return paths relative to the current working directory when `--abs` is not set.
