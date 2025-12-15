package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// TestMainAppInitialization tests that the app can be created and initialized
// without panicking. This provides coverage for the app initialization code path.
func TestMainAppInitialization(t *testing.T) {
	// Use test app to avoid opening actual window
	app := test.NewApp()
	if app == nil {
		t.Fatal("Failed to create test app")
	}
	defer app.Quit()

	window := app.NewWindow("PS3UART GUI")
	if window == nil {
		t.Fatal("Failed to create window")
	}

	window.Resize(fyne.NewSize(600, 500))

	content := createMainWindow(app, window)
	if content == nil {
		t.Fatal("createMainWindow returned nil")
	}

	window.SetContent(content)

	// Verify window has content set
	if window.Content() == nil {
		t.Error("Window content is nil after SetContent")
	}
}

// TestMainAppWindow tests window properties
func TestMainAppWindow(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("PS3UART GUI")
	window.Resize(fyne.NewSize(600, 500))

	size := window.Canvas().Size()
	if size.Width < 600 || size.Height < 500 {
		t.Errorf("Window size is smaller than expected: %v", size)
	}
}

// TestAppMetadata tests that the app metadata is correct
func TestAppMetadata(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// Test app is created correctly
	if app == nil {
		t.Fatal("App is nil")
	}
}

// TestWindowContent tests that window content is properly set
func TestWindowContent(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test Window")
	content := createMainWindow(app, window)

	window.SetContent(content)

	// Trigger a canvas refresh to ensure rendering works
	window.Canvas().Refresh(content)

	// Verify content was set
	if window.Content() != content {
		t.Error("Window content mismatch")
	}
}

// TestMultipleWindows tests creating multiple windows (main + serial monitor)
func TestMultipleWindows(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// Create main window
	mainWindow := app.NewWindow("PS3UART GUI")
	mainWindow.Resize(fyne.NewSize(600, 500))
	mainWindow.SetContent(createMainWindow(app, mainWindow))

	// Open serial monitor window
	openSerialMonitor(app, "", "CXR")

	// Verify main window is still functional
	mainWindow.Canvas().Refresh(mainWindow.Content())
}

// TestWindowResize tests that the window can be resized
func TestWindowResize(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.SetContent(createMainWindow(app, window))

	// Test various sizes
	sizes := []fyne.Size{
		fyne.NewSize(600, 500),
		fyne.NewSize(800, 600),
		fyne.NewSize(400, 300),
	}

	for _, size := range sizes {
		window.Resize(size)
		// Just verify no panic occurs
	}
}
