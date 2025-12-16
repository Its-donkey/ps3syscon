package main

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
		{theme.ColorNameBackground, colorBackground},
		{theme.ColorNameButton, colorSurfaceElevated},
		{theme.ColorNameDisabledButton, color.RGBA{40, 40, 50, 255}},
		{theme.ColorNameDisabled, colorTextMuted},
		{theme.ColorNameError, colorError},
		{theme.ColorNameFocus, colorPrimary},
		{theme.ColorNameForeground, colorText},
		{theme.ColorNamePrimary, colorPrimary},
		{theme.ColorNameSuccess, colorSuccess},
		{theme.ColorNameWarning, colorWarning},
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

func TestPS3ThemeFont(t *testing.T) {
	ps3Theme := &PS3Theme{}

	styles := []fyne.TextStyle{
		{},
		{Bold: true},
		{Italic: true},
		{Monospace: true},
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
		{theme.SizeNameText, 14},
		{theme.SizeNameHeadingText, 20},
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
	result := ps3Theme.Size("unknown_size_name")
	expected := theme.DefaultTheme().Size("unknown_size_name")
	if result != expected {
		t.Errorf("Unknown size should delegate to default theme: got %f, want %f", result, expected)
	}
}
