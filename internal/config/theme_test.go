package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveThemeName_EnvOverride(t *testing.T) {
	t.Setenv(envThemeName, "nord")

	name, err := ResolveThemeName()
	if err != nil {
		t.Fatalf("ResolveThemeName() error = %v", err)
	}
	if name != "nord" {
		t.Fatalf("ResolveThemeName() = %q, want %q", name, "nord")
	}
}

func TestResolveThemeName_ProjectConfigOrder(t *testing.T) {
	t.Setenv(envThemeName, "")

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, projectConfigFallback), "[theme]\nname = \"gruvbox\"\n")
	writeFile(t, filepath.Join(dir, projectConfigPrimary), "[theme]\nname = \"dracula\"\n")

	restore := chdir(t, dir)
	defer restore()

	name, err := ResolveThemeName()
	if err != nil {
		t.Fatalf("ResolveThemeName() error = %v", err)
	}
	if name != "dracula" {
		t.Fatalf("ResolveThemeName() = %q, want %q", name, "dracula")
	}
}

func TestResolveThemeName_UserConfigFallback(t *testing.T) {
	t.Setenv(envThemeName, "")

	home := t.TempDir()
	t.Setenv("HOME", home)

	configPath := filepath.Join(home, userConfigRelativePath)
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	writeFile(t, configPath, "[theme]\nname = \"solarized\"\n")

	restore := chdir(t, t.TempDir())
	defer restore()

	name, err := ResolveThemeName()
	if err != nil {
		t.Fatalf("ResolveThemeName() error = %v", err)
	}
	if name != "solarized" {
		t.Fatalf("ResolveThemeName() = %q, want %q", name, "solarized")
	}
}

func TestResolveThemeName_InvalidToml(t *testing.T) {
	t.Setenv(envThemeName, "")

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, projectConfigPrimary), "theme =\n")

	restore := chdir(t, dir)
	defer restore()

	_, err := ResolveThemeName()
	if err == nil {
		t.Fatalf("ResolveThemeName() expected error")
	}
	if !strings.Contains(err.Error(), "toml") {
		t.Fatalf("ResolveThemeName() error = %v, want TOML parse error", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}

func chdir(t *testing.T, dir string) func() {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%s) error = %v", dir, err)
	}
	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("Chdir restore error = %v", err)
		}
	}
}
