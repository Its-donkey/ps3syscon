package main

import (
	"ps3syscon-gui/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestFilterOptions(t *testing.T) {
	options := []string{"AUTH1", "AUTH2", "VER", "EEP", "ERRLOG", "FAN", "HALT"}

	tests := []struct {
		name     string
		search   string
		expected []string
	}{
		{
			name:     "empty search returns all",
			search:   "",
			expected: options,
		},
		{
			name:     "exact match",
			search:   "VER",
			expected: []string{"VER"},
		},
		{
			name:     "lowercase search",
			search:   "ver",
			expected: []string{"VER"},
		},
		{
			name:     "partial match",
			search:   "AUTH",
			expected: []string{"AUTH1", "AUTH2"},
		},
		{
			name:     "no match",
			search:   "NOTFOUND",
			expected: nil,
		},
		{
			name:     "single char match",
			search:   "E",
			expected: []string{"VER", "EEP", "ERRLOG"},
		},
		{
			name:     "mixed case search",
			search:   "ErR",
			expected: []string{"ERRLOG"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.FilterOptions(options, tt.search)
			if len(result) != len(tt.expected) {
				t.Errorf("FilterOptions(%v, %q) returned %d items, want %d",
					options, tt.search, len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("FilterOptions result[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestFilterOptionsEmptyInput(t *testing.T) {
	var emptyOptions []string
	result := ui.FilterOptions(emptyOptions, "test")
	if len(result) != 0 {
		t.Errorf("FilterOptions with empty input returned %d items, want 0", len(result))
	}
}

func TestFilterOptionsNilInput(t *testing.T) {
	result := ui.FilterOptions(nil, "test")
	if result != nil {
		t.Errorf("FilterOptions with nil input returned %v, want nil", result)
	}
}

func TestCreateHeader(t *testing.T) {
	header := ui.CreateHeader(LogoResource)
	if header == nil {
		t.Fatal("CreateHeader returned nil")
	}

	// Verify it's a container
	_, ok := header.(fyne.CanvasObject)
	if !ok {
		t.Error("Header is not a CanvasObject")
	}

	// Check minimum size is reasonable
	size := header.MinSize()
	if size.Width == 0 || size.Height == 0 {
		t.Errorf("Header has zero size: %v", size)
	}
}

func testWindowDeps() ui.WindowDeps {
	return ui.WindowDeps{
		LogoResource:        LogoResource,
		GetSerialPorts:      getSerialPorts,
		GetCommandNames:     GetCommandNames,
		GetCXRFCommandNames: GetCXRFCommandNames,
		GetCommand:          adaptCommand,
		GetCXRFCommand:      adaptCXRFCommand,
		SendCommand:         sendCommand,
		Authenticate:        authenticate,
		OpenSerialMonitor:   openSerialMonitor,
		ShowGuideWindow:     showGuideWindow,
	}
}

func TestCreateMainWindow(t *testing.T) {
	// Create a test app and window
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Resize(fyne.NewSize(600, 500))

	content := ui.CreateMainWindow(app, window, testWindowDeps())
	if content == nil {
		t.Fatal("CreateMainWindow returned nil")
	}

	// Set content and verify it renders without panic
	window.SetContent(content)

	// Verify minimum size
	size := content.MinSize()
	if size.Width == 0 || size.Height == 0 {
		t.Errorf("Content has zero size: %v", size)
	}
}

func TestOpenSerialMonitor(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// Test with empty port
	openSerialMonitor(app, "", "CXR")

	// Test with port set
	openSerialMonitor(app, "/dev/ttyUSB0", "CXR")

	// Test with CXRF mode
	openSerialMonitor(app, "/dev/ttyUSB0", "CXRF")
}

func TestCreateMainWindowInteraction(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Resize(fyne.NewSize(600, 500))

	content := ui.CreateMainWindow(app, window, testWindowDeps())
	window.SetContent(content)

	// Force a layout/render cycle
	window.Canvas().Refresh(content)
}

func TestCreateMainWindowSCTypeSwitch(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Resize(fyne.NewSize(600, 500))

	content := ui.CreateMainWindow(app, window, testWindowDeps())
	window.SetContent(content)

	// The content should be rendered without errors
	// We can't easily access internal widgets in tests, but we verify no panic
}

func TestFilterOptionsSpecialCharacters(t *testing.T) {
	options := []string{"R8", "R16", "R32", "W8", "W16", "W32"}

	tests := []struct {
		name     string
		search   string
		expected []string
	}{
		{
			name:     "numeric search",
			search:   "8",
			expected: []string{"R8", "W8"},
		},
		{
			name:     "R commands",
			search:   "R",
			expected: []string{"R8", "R16", "R32"},
		},
		{
			name:     "W commands",
			search:   "W",
			expected: []string{"W8", "W16", "W32"},
		},
		{
			name:     "16 bit commands",
			search:   "16",
			expected: []string{"R16", "W16"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.FilterOptions(options, tt.search)
			if len(result) != len(tt.expected) {
				t.Errorf("FilterOptions(%v, %q) returned %d items, want %d",
					options, tt.search, len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("FilterOptions result[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestCreateHeaderContent(t *testing.T) {
	header := ui.CreateHeader(LogoResource)
	if header == nil {
		t.Fatal("CreateHeader returned nil")
	}

	// Verify it renders without panic
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.SetContent(header)
	window.Canvas().Refresh(header)
}

func TestOpenSerialMonitorWithDifferentBaudRates(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// Test with different SC types which set different default baud rates
	testCases := []struct {
		scType string
	}{
		{"CXR"},
		{"CXRF"},
		{"SW"},
	}

	for _, tc := range testCases {
		t.Run(tc.scType, func(t *testing.T) {
			openSerialMonitor(app, "", tc.scType)
		})
	}
}

func TestCreateMainWindowWithAllSCTypes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Resize(fyne.NewSize(600, 500))

	content := ui.CreateMainWindow(app, window, testWindowDeps())
	window.SetContent(content)

	// Trigger layout and refresh
	window.Canvas().Refresh(content)
}

func TestFilterOptionsEdgeCases(t *testing.T) {
	options := []string{"TEST1", "TEST2", "OTHER"}

	// Test with whitespace
	result := ui.FilterOptions(options, "  ")
	if len(result) != 0 {
		t.Errorf("Expected 0 results for whitespace search, got %d", len(result))
	}

	// Test case insensitivity
	result = ui.FilterOptions(options, "test")
	if len(result) != 2 {
		t.Errorf("Expected 2 results for 'test' search, got %d", len(result))
	}
}

func TestBuildCXRCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		subCmd   string
		args     string
		expected string
	}{
		{
			name:     "command only",
			cmd:      "VER",
			subCmd:   "",
			args:     "",
			expected: "VER",
		},
		{
			name:     "command with subcommand",
			cmd:      "EEP",
			subCmd:   "GET",
			args:     "",
			expected: "EEP GET",
		},
		{
			name:     "command with subcommand and args",
			cmd:      "EEP",
			subCmd:   "GET",
			args:     "00",
			expected: "EEP GET 00",
		},
		{
			name:     "command with args only",
			cmd:      "AUTH1",
			subCmd:   "",
			args:     "10000000",
			expected: "AUTH1 10000000",
		},
		{
			name:     "empty command",
			cmd:      "",
			subCmd:   "GET",
			args:     "00",
			expected: "",
		},
		{
			name:     "whitespace command",
			cmd:      "   ",
			subCmd:   "GET",
			args:     "00",
			expected: "",
		},
		{
			name:     "command with whitespace args",
			cmd:      "EEP",
			subCmd:   "SET",
			args:     "  00 FF  ",
			expected: "EEP SET 00 FF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.BuildCXRCommand(tt.cmd, tt.subCmd, tt.args)
			if result != tt.expected {
				t.Errorf("BuildCXRCommand(%q, %q, %q) = %q, want %q",
					tt.cmd, tt.subCmd, tt.args, result, tt.expected)
			}
		})
	}
}

func TestGetSerialSpeed(t *testing.T) {
	tests := []struct {
		scType   string
		expected int
	}{
		{"CXR", 57600},
		{"SW", 57600},
		{"CXRF", 115200},
		{"UNKNOWN", 57600},
		{"", 57600},
	}

	for _, tt := range tests {
		t.Run(tt.scType, func(t *testing.T) {
			result := ui.GetSerialSpeed(tt.scType)
			if result != tt.expected {
				t.Errorf("GetSerialSpeed(%q) = %d, want %d", tt.scType, result, tt.expected)
			}
		})
	}
}

func TestFormatCommandOutput(t *testing.T) {
	tests := []struct {
		name     string
		scType   string
		result   ui.CommandResult
		expected string
	}{
		{
			name:     "CXR mode",
			scType:   "CXR",
			result:   ui.CommandResult{Code: 0, Data: []string{"DATA1", "DATA2"}},
			expected: "00000000 DATA1 DATA2",
		},
		{
			name:     "CXR mode with code",
			scType:   "CXR",
			result:   ui.CommandResult{Code: 0x12345678, Data: []string{"VALUE"}},
			expected: "12345678 VALUE",
		},
		{
			name:     "SW mode single line",
			scType:   "SW",
			result:   ui.CommandResult{Code: 0, Data: []string{"DATA"}},
			expected: "00000000 DATA",
		},
		{
			name:     "SW mode multiline",
			scType:   "SW",
			result:   ui.CommandResult{Code: 0, Data: []string{"LINE1\nLINE2"}},
			expected: "00000000\nLINE1\nLINE2",
		},
		{
			name:     "CXRF mode with data",
			scType:   "CXRF",
			result:   ui.CommandResult{Code: 0, Data: []string{"SC_READY"}},
			expected: "SC_READY",
		},
		{
			name:     "CXRF mode empty data",
			scType:   "CXRF",
			result:   ui.CommandResult{Code: 0, Data: []string{}},
			expected: "",
		},
		{
			name:     "Unknown mode with data",
			scType:   "UNKNOWN",
			result:   ui.CommandResult{Code: 0, Data: []string{"RESPONSE"}},
			expected: "RESPONSE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.FormatCommandOutput(tt.scType, tt.result)
			if result != tt.expected {
				t.Errorf("FormatCommandOutput(%q, %v) = %q, want %q",
					tt.scType, tt.result, result, tt.expected)
			}
		})
	}
}
