// Package ui provides the about dialog and menu setup.
package ui

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SetupMainMenu configures the application menu with About dialog.
func SetupMainMenu(myApp fyne.App, myWindow fyne.Window) {
	aboutItem := fyne.NewMenuItem("About", func() {
		ShowAboutDialog(myWindow)
	})

	helpMenu := fyne.NewMenu("Help", aboutItem)
	mainMenu := fyne.NewMainMenu(helpMenu)
	myWindow.SetMainMenu(mainMenu)
}

// ShowAboutDialog displays the About dialog with version, author, and links.
func ShowAboutDialog(myWindow fyne.Window) {
	versionText := widget.NewLabel(fmt.Sprintf("Version %s", AppVersion))
	versionText.Alignment = fyne.TextAlignCenter

	authorText := widget.NewLabel(fmt.Sprintf("Author: %s", AppAuthor))
	authorText.Alignment = fyne.TextAlignCenter

	var emailHyperlink *widget.Hyperlink
	if emailLink, err := url.Parse("mailto:" + AppEmail); err == nil {
		emailHyperlink = widget.NewHyperlink(AppEmail, emailLink)
	} else {
		emailHyperlink = widget.NewHyperlink(AppEmail, nil)
	}

	var repoHyperlink *widget.Hyperlink
	if repoLink, err := url.Parse(AppRepo); err == nil {
		repoHyperlink = widget.NewHyperlink("GitHub Repository", repoLink)
	} else {
		repoHyperlink = widget.NewHyperlink("GitHub Repository", nil)
	}

	descText := widget.NewLabel("A cross-platform GUI tool for communicating with\nPS3 Syscon via UART interface.")
	descText.Alignment = fyne.TextAlignCenter
	descText.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		versionText,
		widget.NewSeparator(),
		authorText,
		container.NewCenter(emailHyperlink),
		container.NewCenter(repoHyperlink),
		widget.NewSeparator(),
		descText,
	)

	dialog.ShowCustom("About PS3 Syscon UART Tool", "Close", container.NewPadded(content), myWindow)
}
