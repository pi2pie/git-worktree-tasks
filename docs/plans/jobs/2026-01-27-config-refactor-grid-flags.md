---
title: "Config refactor and grid flags preservation"
created-date: 2026-01-27
status: completed
agent: github-copilot
---

# Config Refactor and Grid Flags Preservation

## Summary

Refactored `internal/config/config.go` to improve code quality and restore correct grid cascade behavior that preserves explicit user settings.

## Changes

### 1. DRY Refactor for `applyProjectConfig`

Replaced duplicated config loading pattern with a loop:

```go
paths := []string{
    filepath.Join(cwd, projectConfigPrimary),
    filepath.Join(cwd, projectConfigFallback),
}

for _, path := range paths {
    file, ok, err := loadConfigFile(path)
    if err != nil {
        return err
    }
    if ok {
        applyConfig(cfg, flags, file)
        return nil
    }
}
```

Benefits:

- Eliminates code duplication
- Easier to add more fallback paths in the future
- Maintains "first match wins" semantics

### 2. Error Context for `os.Getwd()`

Added error wrapping:

```go
if err != nil {
    return fmt.Errorf("get working directory: %w", err)
}
```

### 3. Removed Unused Parameter

Removed `flags *gridFlags` parameter from `applyEnvConfig` since env config doesn't control grid settings.

### 4. Restored Global Grid Flags Tracking

Fixed grid cascade behavior to preserve explicit settings across config layers:

- `gridFlags` struct tracks whether `list.grid` or `status.grid` was explicitly set
- Cascade from `table.grid` only happens in `Load()` after all config layers are applied
- User's explicit `list.grid = true` is preserved even if project sets `table.grid = false`

## Test Coverage

Added `TestLoadConfigExplicitGridPreserved` to verify:

- User config sets `list.grid = true`
- Project config sets `table.grid = false` (no explicit `list.grid`)
- Result: `list.grid` stays `true`, `status.grid` cascades to `false`

## Files Modified

- `internal/config/config.go`
- `internal/config/config_test.go`
