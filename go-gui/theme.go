package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Custom color palette
var (
	colorBackground      = color.RGBA{18, 18, 24, 255}      // Deep dark blue-black
	colorSurface         = color.RGBA{28, 28, 38, 255}      // Slightly lighter surface
	colorSurfaceElevated = color.RGBA{38, 38, 52, 255}      // Cards and elevated surfaces
	colorPrimary         = color.RGBA{0, 212, 255, 255}     // Cyan accent
	colorPrimaryDark     = color.RGBA{0, 170, 204, 255}     // Darker cyan
	colorSecondary       = color.RGBA{138, 43, 226, 255}    // Purple accent
	colorSuccess         = color.RGBA{0, 255, 136, 255}     // Green for success
	colorWarning         = color.RGBA{255, 170, 0, 255}     // Orange for warnings
	colorError           = color.RGBA{255, 82, 82, 255}     // Red for errors
	colorText            = color.RGBA{240, 240, 250, 255}   // Light text
	colorTextMuted       = color.RGBA{140, 140, 160, 255}   // Muted text
	colorBorder          = color.RGBA{60, 60, 80, 255}      // Border color
	colorInputBg         = color.RGBA{22, 22, 32, 255}      // Input background
)

// PS3Theme implements a custom dark theme for the PS3 UART tool.
type PS3Theme struct{}

var _ fyne.Theme = (*PS3Theme)(nil)

// Color returns the color for the specified theme color name.
func (t *PS3Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colorBackground
	case theme.ColorNameButton:
		return colorSurfaceElevated
	case theme.ColorNameDisabledButton:
		return color.RGBA{40, 40, 50, 255}
	case theme.ColorNameDisabled:
		return colorTextMuted
	case theme.ColorNameError:
		return colorError
	case theme.ColorNameFocus:
		return colorPrimary
	case theme.ColorNameForeground:
		return colorText
	case theme.ColorNameForegroundOnError:
		return colorText
	case theme.ColorNameForegroundOnPrimary:
		return colorBackground
	case theme.ColorNameForegroundOnSuccess:
		return colorBackground
	case theme.ColorNameForegroundOnWarning:
		return colorBackground
	case theme.ColorNameHeaderBackground:
		return colorSurface
	case theme.ColorNameHover:
		return color.RGBA{50, 50, 70, 255}
	case theme.ColorNameHyperlink:
		return colorPrimary
	case theme.ColorNameInputBackground:
		return colorInputBg
	case theme.ColorNameInputBorder:
		return colorBorder
	case theme.ColorNameMenuBackground:
		return colorSurface
	case theme.ColorNameOverlayBackground:
		return color.RGBA{0, 0, 0, 200}
	case theme.ColorNamePlaceHolder:
		return colorTextMuted
	case theme.ColorNamePressed:
		return colorPrimaryDark
	case theme.ColorNamePrimary:
		return colorPrimary
	case theme.ColorNameScrollBar:
		return colorBorder
	case theme.ColorNameSelection:
		return color.RGBA{0, 212, 255, 80}
	case theme.ColorNameSeparator:
		return colorBorder
	case theme.ColorNameShadow:
		return color.RGBA{0, 0, 0, 100}
	case theme.ColorNameSuccess:
		return colorSuccess
	case theme.ColorNameWarning:
		return colorWarning
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
