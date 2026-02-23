package detect

import (
	"os"
	"testing"
)

func TestDetectTerminalTheme_COLORFGBG_Dark(t *testing.T) {
	original := os.Getenv("COLORFGBG")
	defer func() {
		os.Setenv("COLORFGBG", original)
	}()

	// Dark backgrounds (0-6, 8)
	darkValues := []string{"15;0", "15;1", "15;8", "7;0"}
	for _, val := range darkValues {
		os.Setenv("COLORFGBG", val)
		theme := DetectTerminalTheme()
		if theme != ThemeDark {
			t.Errorf("COLORFGBG=%q should be dark, got %v", val, theme)
		}
	}
}

func TestDetectTerminalTheme_COLORFGBG_Light(t *testing.T) {
	original := os.Getenv("COLORFGBG")
	defer func() {
		os.Setenv("COLORFGBG", original)
	}()

	// Light backgrounds (7, 15)
	lightValues := []string{"0;7", "0;15", "8;15"}
	for _, val := range lightValues {
		os.Setenv("COLORFGBG", val)
		theme := DetectTerminalTheme()
		if theme != ThemeLight {
			t.Errorf("COLORFGBG=%q should be light, got %v", val, theme)
		}
	}
}

func TestDetectTerminalTheme_COLORTERM(t *testing.T) {
	// Save original env
	origCOLORFGBG := os.Getenv("COLORFGBG")
	origCOLORTERM := os.Getenv("COLORTERM")
	defer func() {
		os.Setenv("COLORFGBG", origCOLORFGBG)
		os.Setenv("COLORTERM", origCOLORTERM)
	}()

	// Clear COLORFGBG so COLORTERM is used
	os.Setenv("COLORFGBG", "")

	tests := []struct {
		colorterm string
		expected  TerminalTheme
	}{
		{"truecolor", ThemeDark},
		{"24bit", ThemeDark},
	}

	for _, tt := range tests {
		os.Setenv("COLORTERM", tt.colorterm)
		theme := DetectTerminalTheme()
		if theme != tt.expected {
			t.Errorf("COLORTERM=%q: expected %v, got %v", tt.colorterm, tt.expected, theme)
		}
	}
}

func TestDetectTerminalTheme_Default(t *testing.T) {
	// Save and clear relevant env vars
	origCOLORFGBG := os.Getenv("COLORFGBG")
	origCOLORTERM := os.Getenv("COLORTERM")
	origApple := os.Getenv("AppleInterfaceStyle")
	defer func() {
		os.Setenv("COLORFGBG", origCOLORFGBG)
		os.Setenv("COLORTERM", origCOLORTERM)
		os.Setenv("AppleInterfaceStyle", origApple)
	}()

	os.Setenv("COLORFGBG", "")
	os.Setenv("COLORTERM", "")
	os.Setenv("AppleInterfaceStyle", "")

	theme := DetectTerminalTheme()
	if theme != ThemeDark {
		t.Errorf("default theme should be dark, got %v", theme)
	}
}

func TestThemeConstants(t *testing.T) {
	if ThemeDark != "dark" {
		t.Errorf("ThemeDark should be 'dark', got %q", ThemeDark)
	}
	if ThemeLight != "light" {
		t.Errorf("ThemeLight should be 'light', got %q", ThemeLight)
	}
}
