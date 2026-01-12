package cli

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func copyToClipboard(ctx context.Context, text string) error {
	switch runtime.GOOS {
	case "darwin":
		return runClipboardCommand(ctx, "pbcopy", nil, text)
	case "windows":
		return runClipboardCommand(ctx, "clip", nil, text)
	case "linux":
		if err := tryClipboardCommand(ctx, "wl-copy", nil, text); err == nil {
			return nil
		}
		return runClipboardCommand(ctx, "xclip", []string{"-selection", "clipboard"}, text)
	default:
		return fmt.Errorf("clipboard copy not supported on %s", runtime.GOOS)
	}
}

func tryClipboardCommand(ctx context.Context, name string, args []string, text string) error {
	if _, err := exec.LookPath(name); err != nil {
		return err
	}
	return runClipboardCommand(ctx, name, args, text)
}

func runClipboardCommand(ctx context.Context, name string, args []string, text string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copy to clipboard via %s: %w", name, err)
	}
	return nil
}
