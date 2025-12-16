// Package ui provides the user interface components for the PS3 Syscon UART tool.
package ui

import (
	"image/color"

	"fyne.io/fyne/v2/theme"
)

// Theme color name aliases for shorter syntax.
const (
	FTBackground          = theme.ColorNameBackground
	FTButton              = theme.ColorNameButton
	FTDisabledButton      = theme.ColorNameDisabledButton
	FTDisabled            = theme.ColorNameDisabled
	FTError               = theme.ColorNameError
	FTFocus               = theme.ColorNameFocus
	FTForeground          = theme.ColorNameForeground
	FTText		          = theme.ColorNameForeground
	FTForegroundOnError   = theme.ColorNameForegroundOnError
	FTForegroundOnPrimary = theme.ColorNameForegroundOnPrimary
	FTForegroundOnSuccess = theme.ColorNameForegroundOnSuccess
	FTForegroundOnWarning = theme.ColorNameForegroundOnWarning
	FTHeaderBackground    = theme.ColorNameHeaderBackground
	FTHover               = theme.ColorNameHover
	FTHyperlink           = theme.ColorNameHyperlink
	FTInputBackground     = theme.ColorNameInputBackground
	FTInputBorder         = theme.ColorNameInputBorder
	FTMenuBackground      = theme.ColorNameMenuBackground
	FTOverlayBackground   = theme.ColorNameOverlayBackground
	FTPlaceHolder         = theme.ColorNamePlaceHolder
	FTPressed             = theme.ColorNamePressed
	FTPrimary             = theme.ColorNamePrimary
	FTScrollBar           = theme.ColorNameScrollBar
	FTSelection           = theme.ColorNameSelection
	FTSeparator           = theme.ColorNameSeparator
	FTShadow              = theme.ColorNameShadow
	FTSuccess             = theme.ColorNameSuccess
	FTWarning             = theme.ColorNameWarning
)

// Color palette for the PS3 Syscon UI theme.
var (
	ColorBackground      = color.RGBA{18, 18, 24, 255}    // Deep dark blue-black
	ColorSurface         = color.RGBA{28, 28, 38, 255}    // Slightly lighter surface
	ColorSurfaceElevated = color.RGBA{38, 38, 52, 255}    // Cards and elevated surfaces
	ColorPrimary         = color.RGBA{0, 212, 255, 255}   // Cyan accent
	ColorPrimaryDark     = color.RGBA{0, 170, 204, 255}   // Darker cyan for pressed state
	ColorSuccess         = color.RGBA{0, 255, 136, 255}   // Green for success
	ColorWarning         = color.RGBA{255, 170, 0, 255}   // Orange for warnings
	ColorError           = color.RGBA{255, 82, 82, 255}   // Red for errors
	ColorText            = color.RGBA{240, 240, 250, 255} // Light text
	ColorTextMuted       = color.RGBA{140, 140, 160, 255} // Muted text
	ColorBorder          = color.RGBA{60, 60, 80, 255}    // Border color
	ColorInputBg         = color.RGBA{22, 22, 32, 255}    // Input background
)
