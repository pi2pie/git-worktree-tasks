package cli

import (
	"fmt"
	"io"
	"testing"

	"github.com/spf13/cobra"
)

func TestMaskSensitivePathsPrecedence(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_DRY_RUN_MASK_SENSITIVE_PATHS", "")

		got, err := runMaskCommand(t, project)
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if !got {
			t.Fatalf("mask_sensitive_paths = false, want true")
		}
	})

	t.Run("config", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_DRY_RUN_MASK_SENSITIVE_PATHS", "")
		writeConfig(t, project, "[dry_run]\nmask_sensitive_paths = false\n")

		got, err := runMaskCommand(t, project)
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got {
			t.Fatalf("mask_sensitive_paths = true, want false")
		}
	})

	t.Run("env_over_config", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		writeConfig(t, project, "[dry_run]\nmask_sensitive_paths = true\n")
		t.Setenv("GWTT_DRY_RUN_MASK_SENSITIVE_PATHS", "false")

		got, err := runMaskCommand(t, project)
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got {
			t.Fatalf("mask_sensitive_paths = true, want false")
		}
	})

	t.Run("flag_over_env", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_DRY_RUN_MASK_SENSITIVE_PATHS", "true")

		got, err := runMaskCommand(t, project, "--mask-sensitive-paths=false")
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got {
			t.Fatalf("mask_sensitive_paths = true, want false")
		}
	})

	t.Run("flag_true_over_env_false", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_DRY_RUN_MASK_SENSITIVE_PATHS", "false")

		got, err := runMaskCommand(t, project, "--mask-sensitive-paths")
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if !got {
			t.Fatalf("mask_sensitive_paths = false, want true")
		}
	})

	t.Run("no_mask_flag_disables", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_DRY_RUN_MASK_SENSITIVE_PATHS", "true")

		got, err := runMaskCommand(t, project, "--no-mask-sensitive-paths")
		if err != nil {
			t.Fatalf("run command: %v", err)
		}
		if got {
			t.Fatalf("mask_sensitive_paths = true, want false")
		}
	})

	t.Run("mask_and_no_mask_conflict", func(t *testing.T) {
		project := t.TempDir()
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GWTT_DRY_RUN_MASK_SENSITIVE_PATHS", "")

		_, err := runMaskCommand(t, project, "--mask-sensitive-paths=false", "--no-mask-sensitive-paths")
		if err == nil {
			t.Fatalf("expected conflict error")
		}
		if err.Error() != "cannot use both --mask-sensitive-paths and --no-mask-sensitive-paths" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func runMaskCommand(t *testing.T, cwd string, args ...string) (bool, error) {
	t.Helper()
	cmd, _ := gitWorkTreeCommand()
	var got bool
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.AddCommand(&cobra.Command{
		Use: "inspect-mask",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, ok := configFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("config missing from context")
			}
			got = cfg.DryRun.MaskSensitivePaths
			return nil
		},
	})
	cmd.SetArgs(append(args, "inspect-mask"))

	restore := chdir(t, cwd)
	defer restore()

	err := cmd.Execute()
	return got, err
}
