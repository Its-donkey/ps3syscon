package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestShowGuideWindow(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ShowGuideWindow panicked: %v", r)
		}
	}()

	ShowGuideWindow(app)
}

func TestShowGuideWindowMultipleTimes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// Opening multiple guide windows should not cause issues
	ShowGuideWindow(app)
	ShowGuideWindow(app)
	ShowGuideWindow(app)
}

func TestGuideContentNotEmpty(t *testing.T) {
	if guideContent == "" {
		t.Error("guideContent should not be empty")
	}
}

func TestGuideContentHasContent(t *testing.T) {
	// The guide should have reasonable content
	if len(guideContent) < 100 {
		t.Errorf("guideContent is too short: %d characters", len(guideContent))
	}
}
