package cli

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	posixHomeToken   = "$HOME"
	windowsHomeToken = "%USERPROFILE%"
)

type pathMaskContext struct {
	enabled bool
	home    string
	windows bool
}

func shouldMaskSensitivePaths(ctx context.Context) bool {
	cfg, ok := configFromContext(ctx)
	if !ok || cfg == nil {
		return true
	}
	return cfg.DryRun.MaskSensitivePaths
}

func resolvePathMaskContext(mask bool) pathMaskContext {
	if !mask {
		return pathMaskContext{}
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return pathMaskContext{}
	}
	return pathMaskContext{
		enabled: true,
		home:    home,
		windows: runtime.GOOS == "windows",
	}
}

func formatGitCommandForDryRun(args []string, mask bool) string {
	return formatGitCommandForDryRunWithContext(args, resolvePathMaskContext(mask))
}

func formatGitCommandForDryRunWithContext(args []string, maskCtx pathMaskContext) string {
	if !maskCtx.enabled {
		return formatGitCommand(args)
	}
	return "git " + formatDryRunArgs(maskGitArgs(args, maskCtx), maskCtx)
}

func formatDryRunArgs(args []string, maskCtx pathMaskContext) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, quoteDryRunArg(arg, maskCtx))
	}
	return strings.Join(parts, " ")
}

func quoteDryRunArg(arg string, maskCtx pathMaskContext) string {
	if !maskCtx.windows && strings.HasPrefix(arg, posixHomeToken) {
		return quotePosixHomeArg(arg)
	}
	return shellQuote(arg)
}

func quotePosixHomeArg(path string) string {
	rest := path[len(posixHomeToken):]
	escapedRest := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"`", "\\`",
		"$", `\$`,
	).Replace(rest)
	return `"$HOME` + escapedRest + `"`
}

func maskGitArgs(args []string, maskCtx pathMaskContext) []string {
	masked := make([]string, len(args))
	pathValue := false
	for i, arg := range args {
		value := arg
		if pathValue || looksLikeAbsolutePath(arg, maskCtx.windows) {
			value = maskHomePathWithContext(arg, maskCtx)
		}
		masked[i] = value
		pathValue = expectsPathValue(arg)
	}
	return masked
}

func expectsPathValue(arg string) bool {
	switch arg {
	case "-C", "--git-dir", "--work-tree":
		return true
	default:
		return false
	}
}

func looksLikeAbsolutePath(value string, windows bool) bool {
	if windows {
		return isWindowsAbsolutePath(value)
	}
	return filepath.IsAbs(value)
}

func isWindowsAbsolutePath(value string) bool {
	if len(value) >= 3 && isASCIILetter(value[0]) && value[1] == ':' && (value[2] == '\\' || value[2] == '/') {
		return true
	}
	if strings.HasPrefix(value, `\\`) || strings.HasPrefix(value, "//") {
		return true
	}
	return strings.HasPrefix(value, `\`)
}

func isASCIILetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func maskPathForDryRun(path string, mask bool) string {
	return maskPathForDryRunWithContext(path, resolvePathMaskContext(mask))
}

func maskPathForDryRunWithContext(path string, maskCtx pathMaskContext) string {
	return maskHomePathWithContext(path, maskCtx)
}

func maskHomePathWithContext(path string, maskCtx pathMaskContext) string {
	if !maskCtx.enabled {
		return path
	}
	return maskHomePathWith(path, maskCtx.home, maskCtx.windows)
}

func maskHomePathWith(path, home string, windows bool) string {
	if strings.TrimSpace(home) == "" {
		return path
	}
	normalizedPath := normalizePathForMask(path, windows)
	normalizedHome := normalizePathForMask(home, windows)
	if normalizedPath == "" || normalizedHome == "" {
		return path
	}

	token := posixHomeToken
	separator := "/"
	if windows {
		token = windowsHomeToken
		separator = `\`
	}

	if pathEqualsForMask(normalizedPath, normalizedHome, windows) {
		return token
	}
	if !hasPathPrefixForMask(normalizedPath, normalizedHome, windows) {
		return path
	}

	relative := trimLeadingSeparators(normalizedPath[len(normalizedHome):], windows)
	if relative == "" {
		return token
	}
	return token + separator + relative
}

func normalizePathForMask(path string, windows bool) string {
	normalized := strings.TrimSpace(path)
	if normalized == "" {
		return ""
	}
	if windows {
		normalized = strings.ReplaceAll(normalized, "/", `\`)
		for strings.HasSuffix(normalized, `\`) && !isWindowsDriveRoot(normalized) {
			normalized = strings.TrimSuffix(normalized, `\`)
		}
		return normalized
	}
	if normalized != "/" {
		normalized = strings.TrimRight(normalized, "/")
		if normalized == "" {
			return "/"
		}
	}
	return normalized
}

func isWindowsDriveRoot(path string) bool {
	return len(path) == 3 && isASCIILetter(path[0]) && path[1] == ':' && path[2] == '\\'
}

func pathEqualsForMask(left, right string, windows bool) bool {
	if windows {
		return strings.EqualFold(left, right)
	}
	return left == right
}

func hasPathPrefixForMask(path, prefix string, windows bool) bool {
	separator := "/"
	if windows {
		path = strings.ToLower(path)
		prefix = strings.ToLower(prefix)
		separator = `\`
	}
	if path == prefix {
		return true
	}
	if !strings.HasSuffix(prefix, separator) {
		prefix += separator
	}
	return strings.HasPrefix(path, prefix)
}

func trimLeadingSeparators(value string, windows bool) string {
	if windows {
		return strings.TrimLeft(value, `/\`)
	}
	return strings.TrimLeft(value, "/")
}
