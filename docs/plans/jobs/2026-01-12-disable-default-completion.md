---
title: "Disable default completion subcommand"
created-date: 2026-01-12
status: completed
agent: codex
---

## Summary

- Disabled Cobra's auto-registered completion subcommand so help output matches intended commands.

## Changes

- Set `cmd.CompletionOptions.DisableDefaultCmd = true` in `cli/root.go`.

## Rationale

- The tool does not define a completion command, so the default Cobra command was misleading.
