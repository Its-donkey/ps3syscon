// Package ui provides the application header component.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// CreateHeader creates the modern header with branding.
func CreateHeader(logoResource fyne.Resource) fyne.CanvasObject {
	// Logo image
	logo := canvas.NewImageFromResource(logoResource)
	logo.SetMinSize(fyne.NewSize(50, 50))
	logo.FillMode = canvas.ImageFillContain

	// Main title
	title := canvas.NewText("PS3 SYSCON", ColorPrimary)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Subtitle
	subtitle := canvas.NewText("UART INTERFACE", ColorTextMuted)
	subtitle.TextSize = 12

	// Decorative line
	line := canvas.NewRectangle(ColorPrimary)
	line.SetMinSize(fyne.NewSize(0, 2))

	titleStack := container.NewVBox(
		container.NewHBox(logo, container.NewVBox(title, subtitle), layout.NewSpacer()),
		container.NewPadded(line),
	)

	return container.NewPadded(titleStack)
}
