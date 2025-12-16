package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestCreateHeader(t *testing.T) {
	// Create a dummy resource for testing
	logoResource := fyne.NewStaticResource("test_logo", []byte{})

	header := CreateHeader(logoResource)
	if header == nil {
		t.Fatal("CreateHeader returned nil")
	}
}

func TestCreateHeaderMinSize(t *testing.T) {
	logoResource := fyne.NewStaticResource("test_logo", []byte{})

	header := CreateHeader(logoResource)
	size := header.MinSize()

	if size.Width == 0 || size.Height == 0 {
		t.Errorf("Header has zero size: %v", size)
	}
}

func TestCreateHeaderWithNilResource(t *testing.T) {
	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CreateHeader with nil resource panicked: %v", r)
		}
	}()

	header := CreateHeader(nil)
	if header == nil {
		t.Error("CreateHeader with nil resource returned nil")
	}
}

func TestCreateHeaderRenders(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	logoResource := fyne.NewStaticResource("test_logo", []byte{})

	header := CreateHeader(logoResource)
	window.SetContent(header)
	window.Canvas().Refresh(header)

	// Verify no panic during rendering
}

func TestCreateHeaderReasonableSize(t *testing.T) {
	logoResource := fyne.NewStaticResource("test_logo", []byte{})

	header := CreateHeader(logoResource)
	size := header.MinSize()

	// Header should have reasonable minimum dimensions
	if size.Height < 50 {
		t.Errorf("Header height %f is too small (expected >= 50)", size.Height)
	}
}
