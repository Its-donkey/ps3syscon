package ui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/theme"
)

func TestColorPaletteValues(t *testing.T) {
	tests := []struct {
		name     string
		c        color.RGBA
		expected color.RGBA
	}{
		{"ColorBackground", ColorBackground, color.RGBA{18, 18, 24, 255}},
		{"ColorSurface", ColorSurface, color.RGBA{28, 28, 38, 255}},
		{"ColorSurfaceElevated", ColorSurfaceElevated, color.RGBA{38, 38, 52, 255}},
		{"ColorPrimary", ColorPrimary, color.RGBA{0, 212, 255, 255}},
		{"ColorPrimaryDark", ColorPrimaryDark, color.RGBA{0, 170, 204, 255}},
		{"ColorSuccess", ColorSuccess, color.RGBA{0, 255, 136, 255}},
		{"ColorWarning", ColorWarning, color.RGBA{255, 170, 0, 255}},
		{"ColorError", ColorError, color.RGBA{255, 82, 82, 255}},
		{"ColorText", ColorText, color.RGBA{240, 240, 250, 255}},
		{"ColorTextMuted", ColorTextMuted, color.RGBA{140, 140, 160, 255}},
		{"ColorBorder", ColorBorder, color.RGBA{60, 60, 80, 255}},
		{"ColorInputBg", ColorInputBg, color.RGBA{22, 22, 32, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.c != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.c, tt.expected)
			}
		})
	}
}

func TestColorPaletteFullAlpha(t *testing.T) {
	colors := []color.RGBA{
		ColorBackground,
		ColorSurface,
		ColorSurfaceElevated,
		ColorPrimary,
		ColorPrimaryDark,
		ColorSuccess,
		ColorWarning,
		ColorError,
		ColorText,
		ColorTextMuted,
		ColorBorder,
		ColorInputBg,
	}

	for i, c := range colors {
		if c.A != 255 {
			t.Errorf("Color %d has alpha %d, expected 255", i, c.A)
		}
	}
}

func TestThemeColorNameAliases(t *testing.T) {
	// Verify the theme color name aliases are correctly mapped
	tests := []struct {
		alias    string
		expected string
	}{
		{string(FTBackground), string(theme.ColorNameBackground)},
		{string(FTButton), string(theme.ColorNameButton)},
		{string(FTDisabledButton), string(theme.ColorNameDisabledButton)},
		{string(FTDisabled), string(theme.ColorNameDisabled)},
		{string(FTError), string(theme.ColorNameError)},
		{string(FTFocus), string(theme.ColorNameFocus)},
		{string(FTForeground), string(theme.ColorNameForeground)},
		{string(FTText), string(theme.ColorNameForeground)},
		{string(FTPrimary), string(theme.ColorNamePrimary)},
		{string(FTSuccess), string(theme.ColorNameSuccess)},
		{string(FTWarning), string(theme.ColorNameWarning)},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			if tt.alias != tt.expected {
				t.Errorf("alias %s = %q, want %q", tt.alias, tt.alias, tt.expected)
			}
		})
	}
}

func TestColorsAreDistinct(t *testing.T) {
	colors := map[string]color.RGBA{
		"Background":      ColorBackground,
		"Surface":         ColorSurface,
		"SurfaceElevated": ColorSurfaceElevated,
		"Primary":         ColorPrimary,
		"PrimaryDark":     ColorPrimaryDark,
		"Success":         ColorSuccess,
		"Warning":         ColorWarning,
		"Error":           ColorError,
		"Text":            ColorText,
		"TextMuted":       ColorTextMuted,
		"Border":          ColorBorder,
		"InputBg":         ColorInputBg,
	}

	names := make([]string, 0, len(colors))
	for name := range colors {
		names = append(names, name)
	}

	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			c1 := colors[names[i]]
			c2 := colors[names[j]]
			if c1 == c2 {
				t.Errorf("Colors %s and %s are identical: %v", names[i], names[j], c1)
			}
		}
	}
}
