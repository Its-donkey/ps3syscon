// Package main provides the guide window for PS3 Syscon UART documentation.
package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed guide.txt
var guideContent string

// showGuideWindow opens a new window displaying the full PS3 UART Guide.
func ShowGuideWindow(myApp fyne.App) {
	guideWindow := myApp.NewWindow("PS3 UART Guide")
	guideWindow.Resize(fyne.NewSize(700, 600))

	// Title
	title := canvas.NewText("PS3 Syscon UART Guide", theme.Color(FTPrimary))
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}

	// motherbaords := canvas.NewText("Supported PS3 Motherbaords", theme.Color(FTText))
	// motherbaords.TextSize = 14

	// Guide text in a scrollable container
	guideText := widget.NewLabel(guideContent)
	guideText.Wrapping = fyne.TextWrapWord

	scroll := container.NewVScroll(guideText)

	content := container.NewBorder(
		container.NewVBox(
			container.NewCenter(title),
			widget.NewSeparator(),
			// motherbaords,
		),
		nil, nil, nil,
		scroll,
	)

	bg := canvas.NewRectangle(theme.Color(FTBackground))
	guideWindow.SetContent(container.NewStack(bg, container.NewPadded(content)))
	guideWindow.Show()
}
