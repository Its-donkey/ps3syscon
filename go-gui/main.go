package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	myApp := app.New()

	// Apply custom dark theme
	myApp.Settings().SetTheme(&PS3Theme{})

	// Set application icon
	myApp.SetIcon(IconResource)

	myWindow := myApp.NewWindow("PS3 Syscon UART Tool")
	myWindow.Resize(fyne.NewSize(700, 650))
	myWindow.SetIcon(IconResource)

	// Show empty window first, then disclaimer
	myWindow.SetContent(container.NewCenter())
	myWindow.Show()

	// Show disclaimer and exit if not accepted
	go func() {
		if !showDisclaimer(myWindow) {
			myApp.Quit()
			return
		}
		// Disclaimer accepted, show main content
		fyne.Do(func() {
			myWindow.SetContent(createMainWindow(myApp, myWindow))
		})
	}()

	myApp.Run()
}
