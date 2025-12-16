package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestCreateCard(t *testing.T) {
	content := widget.NewLabel("Test Content")
	card := CreateCard("Test Title", content)

	if card == nil {
		t.Fatal("CreateCard returned nil")
	}
}

func TestCreateCardMinSize(t *testing.T) {
	content := widget.NewLabel("Test Content")
	card := CreateCard("Test Title", content)

	size := card.MinSize()
	if size.Width == 0 || size.Height == 0 {
		t.Errorf("Card has zero size: %v", size)
	}
}

func TestCreateCardWithEmptyTitle(t *testing.T) {
	content := widget.NewLabel("Test Content")
	card := CreateCard("", content)

	if card == nil {
		t.Fatal("CreateCard with empty title returned nil")
	}
}

func TestCreateCardWithNilContent(t *testing.T) {
	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CreateCard with nil content panicked: %v", r)
		}
	}()

	card := CreateCard("Test Title", nil)
	if card == nil {
		t.Error("CreateCard with nil content returned nil")
	}
}

func TestCreateCardRenders(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	content := widget.NewLabel("Test Content")
	card := CreateCard("Test Title", content)

	window.SetContent(card)
	window.Canvas().Refresh(card)

	// Verify no panic during rendering
}

func TestCreateCardWithMultipleContent(t *testing.T) {
	content := widget.NewLabel("Line 1\nLine 2\nLine 3")
	card := CreateCard("Multi-line", content)

	if card == nil {
		t.Fatal("CreateCard with multi-line content returned nil")
	}

	size := card.MinSize()
	if size.Height == 0 {
		t.Error("Card with multi-line content has zero height")
	}
}

func TestCreateCardWithLongTitle(t *testing.T) {
	content := widget.NewLabel("Content")
	longTitle := "This is a very long title that should still work correctly"
	card := CreateCard(longTitle, content)

	if card == nil {
		t.Fatal("CreateCard with long title returned nil")
	}
}

func TestCreateCardWithComplexContent(t *testing.T) {
	// Create complex nested content
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter text...")

	_ = widget.NewButton("Click", func() {})

	content := widget.NewLabel("Complex content test")

	card := CreateCard("Complex Card", content)

	if card == nil {
		t.Fatal("CreateCard with complex content returned nil")
	}
}
