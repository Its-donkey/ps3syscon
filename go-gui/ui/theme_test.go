package ui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func TestPS3ThemeImplementsInterface(t *testing.T) {
	var _ fyne.Theme = (*PS3Theme)(nil)
}

func TestPS3ThemeColor(t *testing.T) {
	ps3Theme := &PS3Theme{}

	tests := []struct {
		name     fyne.ThemeColorName
		expected color.Color
	}{
		{theme.ColorNameBackground, ColorBackground},
		{theme.ColorNameButton, ColorSurfaceElevated},
		{theme.ColorNameDisabledButton, color.RGBA{40, 40, 50, 255}},
		{theme.ColorNameDisabled, ColorTextMuted},
		{theme.ColorNameError, ColorError},
		{theme.ColorNameFocus, ColorPrimary},
		{theme.ColorNameForeground, ColorText},
		{theme.ColorNameForegroundOnError, ColorText},
		{theme.ColorNameForegroundOnPrimary, ColorBackground},
		{theme.ColorNameForegroundOnSuccess, ColorBackground},
		{theme.ColorNameForegroundOnWarning, ColorBackground},
		{theme.ColorNameHeaderBackground, ColorSurface},
		{theme.ColorNameHover, color.RGBA{50, 50, 70, 255}},
		{theme.ColorNameHyperlink, ColorPrimary},
		{theme.ColorNameInputBackground, ColorInputBg},
		{theme.ColorNameInputBorder, ColorBorder},
		{theme.ColorNameMenuBackground, ColorSurface},
		{theme.ColorNameOverlayBackground, color.RGBA{0, 0, 0, 200}},
		{theme.ColorNamePlaceHolder, ColorTextMuted},
		{theme.ColorNamePressed, ColorPrimaryDark},
		{theme.ColorNamePrimary, ColorPrimary},
		{theme.ColorNameScrollBar, ColorBorder},
		{theme.ColorNameSelection, color.RGBA{0, 212, 255, 80}},
		{theme.ColorNameSeparator, ColorBorder},
		{theme.ColorNameShadow, color.RGBA{0, 0, 0, 100}},
		{theme.ColorNameSuccess, ColorSuccess},
		{theme.ColorNameWarning, ColorWarning},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			result := ps3Theme.Color(tt.name, theme.VariantDark)
			if result != tt.expected {
				t.Errorf("Color(%s) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestPS3ThemeColorUnknown(t *testing.T) {
	ps3Theme := &PS3Theme{}
	result := ps3Theme.Color("unknown_color_name", theme.VariantDark)
	if result == nil {
		t.Error("Unknown color name should return a default color, not nil")
	}
}

func TestPS3ThemeColorVariants(t *testing.T) {
	ps3Theme := &PS3Theme{}

	// Test that our theme returns the same color regardless of variant
	variants := []fyne.ThemeVariant{theme.VariantDark, theme.VariantLight}

	for _, variant := range variants {
		result := ps3Theme.Color(theme.ColorNamePrimary, variant)
		if result != ColorPrimary {
			t.Errorf("Color(Primary, %v) = %v, want %v", variant, result, ColorPrimary)
		}
	}
}

func TestPS3ThemeFont(t *testing.T) {
	ps3Theme := &PS3Theme{}

	styles := []fyne.TextStyle{
		{},
		{Bold: true},
		{Italic: true},
		{Monospace: true},
		{Bold: true, Italic: true},
	}

	for _, style := range styles {
		result := ps3Theme.Font(style)
		expected := theme.DefaultTheme().Font(style)
		if result != expected {
			t.Errorf("Font(%v) should delegate to default theme", style)
		}
	}
}

func TestPS3ThemeIcon(t *testing.T) {
	ps3Theme := &PS3Theme{}

	iconNames := []fyne.ThemeIconName{
		theme.IconNameCancel,
		theme.IconNameConfirm,
		theme.IconNameDelete,
		theme.IconNameInfo,
		theme.IconNameError,
		theme.IconNameWarning,
	}

	for _, name := range iconNames {
		result := ps3Theme.Icon(name)
		expected := theme.DefaultTheme().Icon(name)
		if result != expected {
			t.Errorf("Icon(%s) should delegate to default theme", name)
		}
	}
}

func TestPS3ThemeSize(t *testing.T) {
	ps3Theme := &PS3Theme{}

	tests := []struct {
		name     fyne.ThemeSizeName
		expected float32
	}{
		{theme.SizeNamePadding, 6},
		{theme.SizeNameInnerPadding, 8},
		{theme.SizeNameScrollBar, 12},
		{theme.SizeNameScrollBarSmall, 4},
		{theme.SizeNameText, 14},
		{theme.SizeNameHeadingText, 20},
		{theme.SizeNameSubHeadingText, 16},
		{theme.SizeNameCaptionText, 12},
		{theme.SizeNameInputBorder, 2},
		{theme.SizeNameInputRadius, 6},
		{theme.SizeNameSelectionRadius, 4},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			result := ps3Theme.Size(tt.name)
			if result != tt.expected {
				t.Errorf("Size(%s) = %f, want %f", tt.name, result, tt.expected)
			}
		})
	}
}

func TestPS3ThemeSizeUnknown(t *testing.T) {
	ps3Theme := &PS3Theme{}
	// Unknown size should delegate to default theme
	result := ps3Theme.Size("unknown_size_name")
	expected := theme.DefaultTheme().Size("unknown_size_name")
	if result != expected {
		t.Errorf("Unknown size name should delegate to default theme: got %f, want %f", result, expected)
	}
}
