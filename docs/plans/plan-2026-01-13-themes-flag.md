---
title: "Make --themes available at root"
date: 2026-01-13
status: completed
agent: codex
---

## Goal
- Allow `gwtt --themes` to list themes without requiring a subcommand.

## Plan
- Add root command handling so `--themes` prints the theme list when no subcommand is provided.
- Verify help still shows when no args and `--themes` is not set.
- Confirm README examples match the updated behavior.
