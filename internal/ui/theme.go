package ui

import "github.com/shanepadgett/agent-extensions/internal/detect"

type Theme struct {
	Primary   string
	Secondary string
	Success   string
	Warning   string
	Error     string
	Muted     string
}

var DarkTheme = Theme{
	Primary:   "#f5a97f", // peach
	Secondary: "#c6a0f6", // mauve
	Success:   "#a6da95", // green
	Warning:   "#eed49f", // yellow
	Error:     "#ed8796", // red
	Muted:     "#6e738d", // overlay
}

var LightTheme = Theme{
	Primary:   "#d35d35", // burnt orange
	Secondary: "#7c3aed", // violet
	Success:   "#16a34a", // green
	Warning:   "#ca8a04", // amber
	Error:     "#dc2626", // red
	Muted:     "#71717a", // zinc
}

func GetTheme() Theme {
	if detect.DetectTerminalTheme() == detect.ThemeLight {
		return LightTheme
	}
	return DarkTheme
}
