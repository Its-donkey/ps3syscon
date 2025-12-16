package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestShowHelpDialog(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	moreHelpCalled := false

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowHelpDialog panicked: %v", r)
		}
	}()

	ShowHelpDialog(app, window, func() {
		moreHelpCalled = true
	})

	_ = moreHelpCalled
}

func TestShowHelpDialogWithNilCallback(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	// This should not panic even with nil callback
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowHelpDialog with nil callback panicked: %v", r)
		}
	}()

	ShowHelpDialog(app, window, nil)
}

func TestShowHelpDialogMultipleTimes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	// Showing dialog multiple times should not cause issues
	ShowHelpDialog(app, window, func() {})
	ShowHelpDialog(app, window, func() {})
	ShowHelpDialog(app, window, func() {})
}
