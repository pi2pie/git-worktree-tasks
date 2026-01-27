package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigPrecedence(t *testing.T) {
	home := t.TempDir()
	project := t.TempDir()

	userConfigPath := filepath.Join(home, userConfigRelativePath)
	if err := os.MkdirAll(filepath.Dir(userConfigPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	writeFile(t, userConfigPath, `
[theme]
name = "user"

[ui]
color_enabled = true

[table]
grid = true

[list]
output = "json"
`)

	writeFile(t, filepath.Join(project, projectConfigPrimary), `
[theme]
name = "project"

[table]
grid = false

[list]
output = "csv"
`)

	restore := chdir(t, project)
	defer restore()

	t.Setenv("HOME", home)
	t.Setenv(envThemeName, "env")
	t.Setenv(envColorEnabled, "0")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Theme.Name != "env" {
		t.Fatalf("Theme.Name = %q, want %q", cfg.Theme.Name, "env")
	}
	if cfg.UI.ColorEnabled {
		t.Fatalf("UI.ColorEnabled = true, want false")
	}
	if cfg.List.Output != "csv" {
		t.Fatalf("List.Output = %q, want %q", cfg.List.Output, "csv")
	}
	if cfg.List.Grid {
		t.Fatalf("List.Grid = true, want false")
	}
	if cfg.Status.Grid {
		t.Fatalf("Status.Grid = true, want false")
	}
}

func TestLoadConfigTableGridFallback(t *testing.T) {
	project := t.TempDir()
	writeFile(t, filepath.Join(project, projectConfigPrimary), `
[table]
grid = true

[list]
grid = false
`)

	restore := chdir(t, project)
	defer restore()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.List.Grid {
		t.Fatalf("List.Grid = true, want false")
	}
	if !cfg.Status.Grid {
		t.Fatalf("Status.Grid = false, want true")
	}
}

func TestLoadEnvColorInvalid(t *testing.T) {
	t.Setenv(envColorEnabled, "maybe")

	_, err := Load()
	if err == nil {
		t.Fatalf("Load() expected error")
	}
}

func TestLoadConfigExplicitGridPreserved(t *testing.T) {
	home := t.TempDir()
	project := t.TempDir()

	userConfigPath := filepath.Join(home, userConfigRelativePath)
	if err := os.MkdirAll(filepath.Dir(userConfigPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	writeFile(t, userConfigPath, `
[list]
grid = true
`)
	writeFile(t, filepath.Join(project, projectConfigPrimary), `
[table]
grid = false
`)

	restore := chdir(t, project)
	defer restore()
	t.Setenv("HOME", home)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.List.Grid {
		t.Fatalf("List.Grid = false, want true (user's explicit setting should be preserved)")
	}
	if cfg.Status.Grid {
		t.Fatalf("Status.Grid = true, want false (should cascade from table.grid)")
	}
}
