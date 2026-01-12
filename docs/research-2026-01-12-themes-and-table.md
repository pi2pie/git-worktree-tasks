---
title: "Theme candidates and responsive table approach"
date: 2026-01-12
status: completed
agent: codex
---

## Goal
Identify 3–5 UI themes suitable for CLI/TUI output and propose a responsive table strategy for narrow terminals.

## Key Findings
- Theme candidates (5 total including default):
  - Default: balanced cyan/magenta highlights for general usage.
  - Nord: cool blue/ice palette for low-glare terminals.
  - Dracula: high-saturation accent for dark backgrounds.
  - Solarized: warm/cool contrast with restrained highlights.
  - Gruvbox: earthy warm palette for readability on dim displays.
- Consistent roles across themes keep UX predictable: Title, Muted, Success, Warning, Error, Header, Prompt, Accent, Border.
- Responsive CLI tables should shrink/truncate flexible columns (PATH, BRANCH) based on terminal width, while keeping status columns readable (PRESENT, HEAD, AHEAD/BEHIND).
- Bubble Tea `bubbles/table` is best used in TUI views with fixed column roles and width adjustments on `tea.WindowSizeMsg`.

## Implications or Recommendations
- Implement `--theme` for CLI/TUI with a small set of named palettes and a shared role-based style map.
- Add width-aware truncation in CLI tables so narrow terminals remain usable without dropping essential columns.
- For TUI, model the list view using `bubbles/table` and update column widths when the terminal resizes.

## Open Questions
- Should the theme list allow user-defined palettes (future config file)?
- Which columns should be hidden when terminal width is extremely small (if truncation is insufficient)?

## References
- None (internal reasoning)

## Related Plans
- plans/plan-2026-01-12-themes-and-rwd.md
