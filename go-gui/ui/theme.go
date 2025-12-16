// Package main provides a custom dark theme for the PS3 Syscon UART tool.
package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// PS3Theme implements fyne.Theme for a custom dark theme.
type PS3Theme struct{}

// Ensure PS3Theme implements fyne.Theme.
var _ fyne.Theme = (*PS3Theme)(nil)

// Color returns the color for the specified theme color name.
func (t *PS3Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return ColorBackground
	case theme.ColorNameButton:
		return ColorSurfaceElevated
	case theme.ColorNameDisabledButton:
		return color.RGBA{40, 40, 50, 255}
	case theme.ColorNameDisabled:
		return ColorTextMuted
	case theme.ColorNameError:
		return ColorError
	case theme.ColorNameFocus:
		return ColorPrimary
	case theme.ColorNameForeground:
		return ColorText
	case theme.ColorNameForegroundOnError:
		return ColorText
	case theme.ColorNameForegroundOnPrimary:
		return ColorBackground
	case theme.ColorNameForegroundOnSuccess:
		return ColorBackground
	case theme.ColorNameForegroundOnWarning:
		return ColorBackground
	case theme.ColorNameHeaderBackground:
		return ColorSurface
	case theme.ColorNameHover:
		return color.RGBA{50, 50, 70, 255}
	case theme.ColorNameHyperlink:
		return ColorPrimary
	case theme.ColorNameInputBackground:
		return ColorInputBg
	case theme.ColorNameInputBorder:
		return ColorBorder
	case theme.ColorNameMenuBackground:
		return ColorSurface
	case theme.ColorNameOverlayBackground:
		return color.RGBA{0, 0, 0, 200}
	case theme.ColorNamePlaceHolder:
		return ColorTextMuted
	case theme.ColorNamePressed:
		return ColorPrimaryDark
	case theme.ColorNamePrimary:
		return ColorPrimary
	case theme.ColorNameScrollBar:
		return ColorBorder
	case theme.ColorNameSelection:
		return color.RGBA{0, 212, 255, 80}
	case theme.ColorNameSeparator:
		return ColorBorder
	case theme.ColorNameShadow:
		return color.RGBA{0, 0, 0, 100}
	case theme.ColorNameSuccess:
		return ColorSuccess
	case theme.ColorNameWarning:
		return ColorWarning
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font returns the font resource for the specified text style.
func (t *PS3Theme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the specified icon name.
func (t *PS3Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the specified theme size name.
func (t *PS3Theme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 4
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 20
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInputBorder:
		return 2
	case theme.SizeNameInputRadius:
		return 6
	case theme.SizeNameSelectionRadius:
		return 4
	default:
		return theme.DefaultTheme().Size(name)
	}
}
