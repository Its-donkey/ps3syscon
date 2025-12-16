// Package ui provides card container components.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateCard creates a card-like container with a title and content.
func CreateCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	titleLabel := canvas.NewText(title, ColorPrimary)
	titleLabel.TextSize = 12
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	bg := canvas.NewRectangle(ColorSurfaceElevated)
	bg.CornerRadius = 8

	cardContent := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		content,
	)

	return container.NewStack(bg, container.NewPadded(cardContent))
}
