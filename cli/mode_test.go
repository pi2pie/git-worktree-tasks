package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestModePrecedence(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_MODE", "")

		got, err := runModeCommand(t, project)
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got != modeClassic {
			t.Fatalf("mode = %q, want %q", got, modeClassic)
		}
	})

	t.Run("config", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_MODE", "")
		writeConfig(t, project, `mode = "codex"`)

		got, err := runModeCommand(t, project)
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got != modeCodex {
			t.Fatalf("mode = %q, want %q", got, modeCodex)
		}
	})

	t.Run("env_over_config", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		writeConfig(t, project, `mode = "classic"`)
		t.Setenv("GWTT_MODE", "codex")

		got, err := runModeCommand(t, project)
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got != modeCodex {
			t.Fatalf("mode = %q, want %q", got, modeCodex)
		}
	})

	t.Run("flag_over_env", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_MODE", "codex")

		got, err := runModeCommand(t, project, "--mode", "classic")
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got != modeClassic {
			t.Fatalf("mode = %q, want %q", got, modeClassic)
		}
	})
}

func TestModeValidation(t *testing.T) {
	project := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("GWTT_MODE", "nope")

	if _, err := runModeCommand(t, project); err == nil || !strings.Contains(err.Error(), "unsupported mode") {
		t.Fatalf("expected unsupported mode error, got %v", err)
	}

	t.Setenv("GWTT_MODE", "")
	if _, err := runModeCommand(t, project, "--mode", "nope"); err == nil || !strings.Contains(err.Error(), "unsupported mode") {
		t.Fatalf("expected unsupported mode error, got %v", err)
	}
}

func runModeCommand(t *testing.T, cwd string, args ...string) (string, error) {
	t.Helper()
	cmd, _ := gitWorkTreeCommand()
	var got string
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.AddCommand(&cobra.Command{
		Use: "inspect",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, ok := configFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("config missing from context")
			}
			got = cfg.Mode
			return nil
		},
	})
	cmd.SetArgs(append(args, "inspect"))

	restore := chdir(t, cwd)
	defer restore()

	err := cmd.Execute()
	return got, err
}

func writeConfig(t *testing.T, dir, content string) {
	t.Helper()
	path := filepath.Join(dir, "gwtt.config.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func chdir(t *testing.T, dir string) func() {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore chdir: %v", err)
		}
	}
}
