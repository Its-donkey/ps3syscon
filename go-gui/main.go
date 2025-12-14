package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("PS3UART GUI")
	myWindow.Resize(fyne.NewSize(600, 500))

	myWindow.SetContent(createMainWindow(myApp, myWindow))
	myWindow.ShowAndRun()
}
