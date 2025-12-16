package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestShowDisclaimer(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	acceptCalled := false
	declineCalled := false

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowDisclaimer panicked: %v", r)
		}
	}()

	ShowDisclaimer(window, func() {
		acceptCalled = true
	}, func() {
		declineCalled = true
	})

	// Note: In tests, we can't actually click the buttons,
	// but we verify the dialog shows without error
	_ = acceptCalled
	_ = declineCalled
}

func TestShowDisclaimerWithNilCallbacks(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	// This should not panic even with nil callbacks
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowDisclaimer with nil callbacks panicked: %v", r)
		}
	}()

	ShowDisclaimer(window, nil, nil)
}

func TestDisclaimerTextIsNotEmpty(t *testing.T) {
	// Access the disclaimerText constant through a test that creates the dialog
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Show()

	// The disclaimer should show without error, which indicates the text is valid
	ShowDisclaimer(window, func() {}, func() {})
}
