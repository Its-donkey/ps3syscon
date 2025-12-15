package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

	myWindow.SetContent(createMainWindow(myApp, myWindow))
	myWindow.ShowAndRun()
}
