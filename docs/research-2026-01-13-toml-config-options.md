---
title: "TOML config options for gwtt"
date: 2026-01-13
status: completed
agent: codex
---

## Goal

Survey Go TOML/config library options and document pros/cons to select an implementation for theme configuration.

## Key Findings

- **Direct TOML parsers**: `BurntSushi/toml` and `pelletier/go-toml` provide TOML decoding with struct mapping; go-toml v2 includes CLI tools like `tomljson`, `jsontoml`, and `tomll`. [^burntsushi] [^go-toml] [^go-toml-tools]
- **Full config frameworks**: `spf13/viper` supports TOML alongside env vars, flags, and multiple config sources/formats. [^viper-pkg]
- **Lightweight config manager**: `knadh/koanf` is a modular, lightweight alternative to viper that composes providers (file/env/flags) and parsers (including TOML). [^koanf]

## Pros and Cons

### BurntSushi/toml

- Pros: Small, focused TOML parser; standard struct decoding; supports current TOML versions; MIT license. [^burntsushi] [^burntsushi-pkg]
- Cons: No built-in config layering/precedence; app must handle file discovery and env overrides (inference based on scope). [^burntsushi]

### pelletier/go-toml (v2)

- Pros: TOML v1.0 support; v2 provides tools (`tomljson`, `jsontoml`, `tomll`) and published benchmarks. [^go-toml] [^go-toml-tools]
- Cons: Includes an "unstable" AST API; may be more than needed for a small config (inference). [^go-toml]

### spf13/viper

- Pros: Full config stack: files + env + flags + defaults; supports TOML and multiple formats. [^viper-pkg]
- Cons: Larger dependency surface and broader feature set than required for a single `theme` setting (inference). [^viper-pkg]

### knadh/koanf

- Pros: Modular providers/parsers; lightweight alternative to viper; explicit merge order control; TOML parser available. [^koanf]
- Cons: Requires selecting providers/parsers explicitly; more setup than a simple TOML unmarshal for one setting (inference). [^koanf]

## Implications or Recommendations

- For **minimal config (single theme)**, prefer a direct TOML parser and implement explicit precedence in our code. BurntSushi/toml or go-toml v2 are both viable; BurntSushi is simplest if we only need decoding. [^burntsushi] [^go-toml]
- If future scope grows to include more config sources (flags/env/remote), evaluate koanf or viper then. [^koanf] [^viper-pkg]

## Related Plans

- `docs/plans/plan-2026-01-13-theme-config-and-env.md`

## References

[^burntsushi]: BurntSushi/toml repository and documentation. https://github.com/BurntSushi/toml

[^burntsushi-pkg]: BurntSushi/toml package docs (compatibility notes). https://pkg.go.dev/github.com/BurntSushi/toml

[^go-toml]: pelletier/go-toml repository (v2) and tools/benchmarks notes. https://github.com/pelletier/go-toml

[^go-toml-tools]: go-toml CLI tools docs. https://github.com/pelletier/go-toml/tree/master/cmd

[^viper-pkg]: Viper package docs (format support and config sources). https://pkg.go.dev/github.com/spf13/viper

[^koanf]: knadh/koanf repository (lightweight config, providers/parsers). https://github.com/knadh/koanf
