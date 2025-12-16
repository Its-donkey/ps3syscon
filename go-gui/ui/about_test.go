package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestSetupMainMenu(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetupMainMenu panicked: %v", r)
		}
	}()

	SetupMainMenu(app, window)
}

func TestShowAboutDialog(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowAboutDialog panicked: %v", r)
		}
	}()

	ShowAboutDialog(window)
}

func TestShowAboutDialogContent(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	// Verify the dialog shows without error
	ShowAboutDialog(window)

	// The dialog should be shown (we can't easily verify content in tests)
}

func TestSetupMainMenuMultipleTimes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")

	// Setting up menu multiple times should not cause issues
	SetupMainMenu(app, window)
	SetupMainMenu(app, window)
	SetupMainMenu(app, window)
}
