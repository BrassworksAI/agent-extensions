package detect

import (
	"os"
	"strings"
)

type TerminalTheme string

const (
	ThemeDark  TerminalTheme = "dark"
	ThemeLight TerminalTheme = "light"
)

func DetectTerminalTheme() TerminalTheme {
	// Check COLORFGBG environment variable (format: "fg;bg")
	if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) >= 2 {
			bg := parts[len(parts)-1]
			// Background values 0-6 or 8 are typically dark
			// Background values 7 or 15 are typically light
			if bg == "7" || bg == "15" {
				return ThemeLight
			}
			return ThemeDark
		}
	}

	// Check for common dark mode indicators
	if colorterm := os.Getenv("COLORTERM"); colorterm != "" {
		// Most modern terminals with truecolor support default to dark
		if colorterm == "truecolor" || colorterm == "24bit" {
			return ThemeDark
		}
	}

	// Check macOS dark mode via AppleInterfaceStyle (if available via env)
	if style := os.Getenv("AppleInterfaceStyle"); strings.ToLower(style) == "dark" {
		return ThemeDark
	}

	// Default to dark theme
	return ThemeDark
}
