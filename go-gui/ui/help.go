// Package ui provides the help dialog component.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowHelpDialog displays the command reference help dialog.
func ShowHelpDialog(myApp fyne.App, myWindow fyne.Window, onMoreHelp func()) {
	helpTitle := canvas.NewText("Command Reference", ColorPrimary)
	helpTitle.TextSize = 18
	helpTitle.TextStyle = fyne.TextStyle{Bold: true}

	helpText := widget.NewLabel(`CXR MODE - External Commands (57600 baud)
  EEP GET/SET, ERRLOG, VER, AUTH

CXRF MODE - Internal Commands (115200 baud)
  r/w, errlog, eepcsum, version
  Requires DIAG pin grounded

TIPS
  Press Enter to send command
  Auth is required before commands`)

	// "Need more help?" link
	moreHelpBtn := widget.NewButton("Need more help? View Full Guide", onMoreHelp)
	moreHelpBtn.Importance = widget.LowImportance

	// Spacer to force minimum width
	spacer := canvas.NewRectangle(ColorBackground)
	spacer.SetMinSize(fyne.NewSize(380, 1))

	content := container.NewVBox(
		helpTitle,
		widget.NewSeparator(),
		helpText,
		widget.NewSeparator(),
		moreHelpBtn,
		spacer,
	)

	dialog.ShowCustom("Help", "Close", container.NewPadded(content), myWindow)
}
