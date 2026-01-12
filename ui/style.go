package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const defaultThemeName = "default"

type Theme struct {
	Name    string
	Title   string
	Muted   string
	Success string
	Warning string
	Error   string
	Header  string
	Prompt  string
	Accent  string
	Border  string
}

var (
	TitleStyle   lipgloss.Style
	MutedStyle   lipgloss.Style
	BorderStyle  lipgloss.Style
	SuccessStyle lipgloss.Style
	WarningStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
	HeaderStyle  lipgloss.Style
	PromptStyle  lipgloss.Style
	AccentStyle  lipgloss.Style
)

var (
	colorEnabled = true
	themes       = map[string]Theme{
		"default": {
			Name:    "default",
			Title:   "63",
			Muted:   "244",
			Success: "35",
			Warning: "214",
			Error:   "203",
			Header:  "81",
			Prompt:  "75",
			Accent:  "213",
			Border:  "239",
		},
		"nord": {
			Name:    "nord",
			Title:   "111",
			Muted:   "245",
			Success: "108",
			Warning: "179",
			Error:   "167",
			Header:  "110",
			Prompt:  "109",
			Accent:  "117",
			Border:  "238",
		},
		"dracula": {
			Name:    "dracula",
			Title:   "141",
			Muted:   "245",
			Success: "84",
			Warning: "220",
			Error:   "203",
			Header:  "147",
			Prompt:  "75",
			Accent:  "213",
			Border:  "239",
		},
		"solarized": {
			Name:    "solarized",
			Title:   "136",
			Muted:   "245",
			Success: "37",
			Warning: "166",
			Error:   "160",
			Header:  "33",
			Prompt:  "32",
			Accent:  "135",
			Border:  "240",
		},
		"gruvbox": {
			Name:    "gruvbox",
			Title:   "208",
			Muted:   "245",
			Success: "142",
			Warning: "214",
			Error:   "167",
			Header:  "223",
			Prompt:  "214",
			Accent:  "175",
			Border:  "239",
		},
	}
	activeTheme = themes[defaultThemeName]
)

func init() {
	SetColorEnabled(true)
}

func SetColorEnabled(enabled bool) {
	colorEnabled = enabled
	applyStyles()
}

func SetTheme(name string) error {
	if strings.TrimSpace(name) == "" {
		name = defaultThemeName
	}
	theme, ok := themes[strings.ToLower(name)]
	if !ok {
		return fmt.Errorf("unknown theme %q (available: %s)", name, strings.Join(ThemeNames(), ", "))
	}
	activeTheme = theme
	applyStyles()
	return nil
}

func ThemeNames() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func DefaultThemeName() string {
	return defaultThemeName
}

func applyStyles() {
	TitleStyle = lipgloss.NewStyle().Bold(true)
	MutedStyle = lipgloss.NewStyle()
	BorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
	SuccessStyle = lipgloss.NewStyle()
	WarningStyle = lipgloss.NewStyle()
	ErrorStyle = lipgloss.NewStyle().Bold(true)
	HeaderStyle = lipgloss.NewStyle().Bold(true)
	PromptStyle = lipgloss.NewStyle().Bold(true)
	AccentStyle = lipgloss.NewStyle().Bold(true)

	if !colorEnabled {
		return
	}

	TitleStyle = TitleStyle.Foreground(lipgloss.Color(activeTheme.Title))
	MutedStyle = MutedStyle.Foreground(lipgloss.Color(activeTheme.Muted))
	SuccessStyle = SuccessStyle.Foreground(lipgloss.Color(activeTheme.Success))
	WarningStyle = WarningStyle.Foreground(lipgloss.Color(activeTheme.Warning))
	ErrorStyle = ErrorStyle.Foreground(lipgloss.Color(activeTheme.Error))
	HeaderStyle = HeaderStyle.Foreground(lipgloss.Color(activeTheme.Header))
	PromptStyle = PromptStyle.Foreground(lipgloss.Color(activeTheme.Prompt))
	AccentStyle = AccentStyle.Foreground(lipgloss.Color(activeTheme.Accent))
	if activeTheme.Border != "" {
		BorderStyle = BorderStyle.BorderForeground(lipgloss.Color(activeTheme.Border))
	}
}
