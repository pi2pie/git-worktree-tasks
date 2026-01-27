package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	envColorEnabled = "GWTT_COLOR"
)

type Config struct {
	Theme   ThemeConfig
	UI      UIConfig
	Table   TableConfig
	Create  CreateConfig
	List    ListConfig
	Status  StatusConfig
	Finish  FinishConfig
	Cleanup CleanupConfig
}

type ThemeConfig struct {
	Name string
}

type UIConfig struct {
	ColorEnabled bool
}

type TableConfig struct {
	Grid bool
}

type CreateConfig struct {
	Output       string
	SkipExisting bool
	Path         CreatePathConfig
}

type CreatePathConfig struct {
	Root   string
	Format string
}

type ListConfig struct {
	Output       string
	Field        string
	AbsolutePath bool
	Grid         bool
	Strict       bool
}

type StatusConfig struct {
	Output       string
	AbsolutePath bool
	Grid         bool
	Strict       bool
}

type FinishConfig struct {
	Cleanup        bool
	RemoveWorktree bool
	RemoveBranch   bool
	ForceBranch    bool
	MergeMode      string
	Confirm        bool
}

type CleanupConfig struct {
	RemoveWorktree bool
	RemoveBranch   bool
	WorktreeOnly   bool
	ForceBranch    bool
	Confirm        bool
}

func DefaultConfig() Config {
	return Config{
		Theme: ThemeConfig{
			Name: "default",
		},
		UI: UIConfig{
			ColorEnabled: true,
		},
		Table: TableConfig{
			Grid: false,
		},
		Create: CreateConfig{
			Output:       "text",
			SkipExisting: false,
			Path: CreatePathConfig{
				Root:   "../",
				Format: "{repo}_{task}",
			},
		},
		List: ListConfig{
			Output:       "table",
			Field:        "path",
			AbsolutePath: false,
			Grid:         false,
			Strict:       false,
		},
		Status: StatusConfig{
			Output:       "table",
			AbsolutePath: false,
			Grid:         false,
			Strict:       false,
		},
		Finish: FinishConfig{
			Cleanup:        false,
			RemoveWorktree: false,
			RemoveBranch:   false,
			ForceBranch:    false,
			MergeMode:      "ff",
			Confirm:        true,
		},
		Cleanup: CleanupConfig{
			RemoveWorktree: true,
			RemoveBranch:   true,
			WorktreeOnly:   false,
			ForceBranch:    false,
			Confirm:        true,
		},
	}
}

type loadedConfigFile struct {
	Theme   themeConfigFile   `toml:"theme"`
	UI      uiConfigFile      `toml:"ui"`
	Table   tableConfigFile   `toml:"table"`
	Create  createConfigFile  `toml:"create"`
	List    listConfigFile    `toml:"list"`
	Status  statusConfigFile  `toml:"status"`
	Finish  finishConfigFile  `toml:"finish"`
	Cleanup cleanupConfigFile `toml:"cleanup"`
}

type themeConfigFile struct {
	Name *string `toml:"name"`
}

type uiConfigFile struct {
	ColorEnabled *bool `toml:"color_enabled"`
}

type tableConfigFile struct {
	Grid *bool `toml:"grid"`
}

type createConfigFile struct {
	Output       *string        `toml:"output"`
	SkipExisting *bool          `toml:"skip_existing"`
	Path         createPathFile `toml:"path"`
}

type createPathFile struct {
	Root   *string `toml:"root"`
	Format *string `toml:"format"`
}

type listConfigFile struct {
	Output       *string `toml:"output"`
	Field        *string `toml:"field"`
	AbsolutePath *bool   `toml:"absolute_path"`
	Grid         *bool   `toml:"grid"`
	Strict       *bool   `toml:"strict"`
}

type statusConfigFile struct {
	Output       *string `toml:"output"`
	AbsolutePath *bool   `toml:"absolute_path"`
	Grid         *bool   `toml:"grid"`
	Strict       *bool   `toml:"strict"`
}

type finishConfigFile struct {
	Cleanup        *bool   `toml:"cleanup"`
	RemoveWorktree *bool   `toml:"remove_worktree"`
	RemoveBranch   *bool   `toml:"remove_branch"`
	ForceBranch    *bool   `toml:"force_branch"`
	MergeMode      *string `toml:"merge_mode"`
	Confirm        *bool   `toml:"confirm"`
}

type cleanupConfigFile struct {
	RemoveWorktree *bool `toml:"remove_worktree"`
	RemoveBranch   *bool `toml:"remove_branch"`
	WorktreeOnly   *bool `toml:"worktree_only"`
	ForceBranch    *bool `toml:"force_branch"`
	Confirm        *bool `toml:"confirm"`
}

func Load() (Config, error) {
	cfg := DefaultConfig()

	if err := applyUserConfig(&cfg); err != nil {
		return cfg, err
	}
	if err := applyProjectConfig(&cfg); err != nil {
		return cfg, err
	}
	if err := applyEnvConfig(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func applyUserConfig(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return nil
	}
	path := filepath.Join(home, userConfigRelativePath)
	file, ok, err := loadConfigFile(path)
	if err != nil {
		return err
	}
	if ok {
		applyConfig(cfg, file)
	}
	return nil
}

func applyProjectConfig(cfg *Config) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

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
			applyConfig(cfg, file)
			return nil
		}
	}
	return nil
}

func loadConfigFile(path string) (loadedConfigFile, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return loadedConfigFile{}, false, nil
		}
		return loadedConfigFile{}, false, err
	}
	var cfg loadedConfigFile
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return loadedConfigFile{}, false, err
	}
	return cfg, true, nil
}

func applyEnvConfig(cfg *Config) error {
	if name, ok := envString(envThemeName); ok {
		cfg.Theme.Name = name
	}
	if enabled, ok, err := envBool(envColorEnabled); err != nil {
		return err
	} else if ok {
		cfg.UI.ColorEnabled = enabled
	}
	return nil
}

func envString(key string) (string, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", false
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}
	return trimmed, true
}

func envBool(key string) (bool, bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false, false, nil
	}
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return false, false, nil
	}
	switch trimmed {
	case "1", "true", "yes", "on":
		return true, true, nil
	case "0", "false", "no", "off":
		return false, true, nil
	default:
		return false, false, fmt.Errorf("invalid %s value %q", key, value)
	}
}

func applyConfig(cfg *Config, file loadedConfigFile) {
	tableGridSet := file.Table.Grid != nil
	listGridSet := file.List.Grid != nil
	statusGridSet := file.Status.Grid != nil

	if name, ok := trimString(file.Theme.Name); ok {
		cfg.Theme.Name = name
	}
	if file.UI.ColorEnabled != nil {
		cfg.UI.ColorEnabled = *file.UI.ColorEnabled
	}
	if file.Table.Grid != nil {
		cfg.Table.Grid = *file.Table.Grid
	}
	if output, ok := trimString(file.Create.Output); ok {
		cfg.Create.Output = output
	}
	if file.Create.SkipExisting != nil {
		cfg.Create.SkipExisting = *file.Create.SkipExisting
	}
	if root, ok := trimString(file.Create.Path.Root); ok {
		cfg.Create.Path.Root = root
	}
	if format, ok := trimString(file.Create.Path.Format); ok {
		cfg.Create.Path.Format = format
	}
	if output, ok := trimString(file.List.Output); ok {
		cfg.List.Output = output
	}
	if field, ok := trimString(file.List.Field); ok {
		cfg.List.Field = field
	}
	if file.List.AbsolutePath != nil {
		cfg.List.AbsolutePath = *file.List.AbsolutePath
	}
	if file.List.Grid != nil {
		cfg.List.Grid = *file.List.Grid
	}
	if file.List.Strict != nil {
		cfg.List.Strict = *file.List.Strict
	}
	if output, ok := trimString(file.Status.Output); ok {
		cfg.Status.Output = output
	}
	if file.Status.AbsolutePath != nil {
		cfg.Status.AbsolutePath = *file.Status.AbsolutePath
	}
	if file.Status.Grid != nil {
		cfg.Status.Grid = *file.Status.Grid
	}
	if file.Status.Strict != nil {
		cfg.Status.Strict = *file.Status.Strict
	}
	if file.Finish.Cleanup != nil {
		cfg.Finish.Cleanup = *file.Finish.Cleanup
	}
	if file.Finish.RemoveWorktree != nil {
		cfg.Finish.RemoveWorktree = *file.Finish.RemoveWorktree
	}
	if file.Finish.RemoveBranch != nil {
		cfg.Finish.RemoveBranch = *file.Finish.RemoveBranch
	}
	if file.Finish.ForceBranch != nil {
		cfg.Finish.ForceBranch = *file.Finish.ForceBranch
	}
	if mode, ok := trimString(file.Finish.MergeMode); ok {
		cfg.Finish.MergeMode = mode
	}
	if file.Finish.Confirm != nil {
		cfg.Finish.Confirm = *file.Finish.Confirm
	}
	if file.Cleanup.RemoveWorktree != nil {
		cfg.Cleanup.RemoveWorktree = *file.Cleanup.RemoveWorktree
	}
	if file.Cleanup.RemoveBranch != nil {
		cfg.Cleanup.RemoveBranch = *file.Cleanup.RemoveBranch
	}
	if file.Cleanup.WorktreeOnly != nil {
		cfg.Cleanup.WorktreeOnly = *file.Cleanup.WorktreeOnly
	}
	if file.Cleanup.ForceBranch != nil {
		cfg.Cleanup.ForceBranch = *file.Cleanup.ForceBranch
	}
	if file.Cleanup.Confirm != nil {
		cfg.Cleanup.Confirm = *file.Cleanup.Confirm
	}
	if tableGridSet {
		if !listGridSet {
			cfg.List.Grid = cfg.Table.Grid
		}
		if !statusGridSet {
			cfg.Status.Grid = cfg.Table.Grid
		}
	}
}

func trimString(value *string) (string, bool) {
	if value == nil {
		return "", false
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return "", false
	}
	return trimmed, true
}
