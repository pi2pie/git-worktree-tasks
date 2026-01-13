package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	envThemeName           = "GWTT_THEME"
	userConfigRelativePath = ".config/gwtt/config.toml"
	projectConfigPrimary   = "gwtt.config.toml"
	projectConfigFallback  = "gwtt.toml"
)

type themeConfig struct {
	Name string `toml:"name"`
}

type configFile struct {
	Theme themeConfig `toml:"theme"`
}

// ResolveThemeName returns the theme name from env or config files.
// Empty string means no theme configured.
func ResolveThemeName() (string, error) {
	if name := themeFromEnv(); name != "" {
		return name, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if name, ok, err := themeFromFile(filepath.Join(cwd, projectConfigPrimary)); err != nil {
		return "", err
	} else if ok {
		return name, nil
	}

	if name, ok, err := themeFromFile(filepath.Join(cwd, projectConfigFallback)); err != nil {
		return "", err
	} else if ok {
		return name, nil
	}

	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "", nil
	}

	if name, ok, err := themeFromFile(filepath.Join(home, userConfigRelativePath)); err != nil {
		return "", err
	} else if ok {
		return name, nil
	}

	return "", nil
}

func themeFromEnv() string {
	value, ok := os.LookupEnv(envThemeName)
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}

func themeFromFile(path string) (string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}

	var cfg configFile
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return "", false, err
	}

	name := strings.TrimSpace(cfg.Theme.Name)
	if name == "" {
		return "", false, nil
	}

	return name, true, nil
}
