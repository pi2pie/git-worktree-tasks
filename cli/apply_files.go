package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func writeTempPatch(contents string) (string, error) {
	tmp, err := os.CreateTemp("", "gwtt-apply-*.patch")
	if err != nil {
		return "", err
	}
	if _, err := io.WriteString(tmp, contents); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			return "", fmt.Errorf("write patch: %w (close error: %v)", err, closeErr)
		}
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func copyFile(srcRoot, dstRoot, rel string, dryRun bool, out io.Writer) (err error) {
	srcPath := filepath.Join(srcRoot, rel)
	dstPath := filepath.Join(dstRoot, rel)

	info, err := os.Lstat(srcPath)
	if err != nil {
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(srcPath)
		if err != nil {
			return err
		}
		if dryRun {
			_, err := fmt.Fprintf(out, "symlink %s -> %s (%s)\n", srcPath, dstPath, target)
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return err
		}
		if err := os.Remove(dstPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return os.Symlink(target, dstPath)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("unsupported file type for copy: %s", srcPath)
	}

	if dryRun {
		_, err := fmt.Fprintf(out, "copy %s -> %s\n", srcPath, dstPath)
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	outFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if _, err := io.Copy(outFile, in); err != nil {
		return err
	}
	return nil
}

func removeTempPatch(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
