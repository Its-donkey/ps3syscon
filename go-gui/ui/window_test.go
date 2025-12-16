package ui

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func testWindowDeps() WindowDeps {
	return WindowDeps{
		LogoResource: fyne.NewStaticResource("test_logo", []byte{}),
		GetSerialPorts: func() []string {
			return []string{"/dev/ttyUSB0", "/dev/ttyUSB1"}
		},
		GetCommandNames: func() []string {
			return []string{"VER", "AUTH1", "EEP"}
		},
		GetCXRFCommandNames: func() []string {
			return []string{"version", "eepcsum", "errlog"}
		},
		GetCommand: func(name string) *Command {
			if name == "EEP" {
				return &Command{Name: "EEP", Subcommands: []string{"GET", "SET"}}
			}
			return &Command{Name: name}
		},
		GetCXRFCommand: func(name string) *Command {
			return &Command{Name: name, Description: "Test description"}
		},
		SendCommand: func(port, scType, cmd string, speed int) (CommandResult, error) {
			return CommandResult{Code: 0, Data: []string{"OK"}}, nil
		},
		Authenticate: func(port, scType string, speed int) error {
			return nil
		},
		OpenSerialMonitor: func(myApp fyne.App, port, scType string) {},
		ShowGuideWindow:   func(myApp fyne.App) {},
	}
}

func TestCreateMainWindow(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Resize(fyne.NewSize(600, 500))

	content := CreateMainWindow(app, window, testWindowDeps())
	if content == nil {
		t.Fatal("CreateMainWindow returned nil")
	}

	window.SetContent(content)

	// Verify minimum size
	size := content.MinSize()
	if size.Width == 0 || size.Height == 0 {
		t.Errorf("Content has zero size: %v", size)
	}
}

func TestCreateMainWindowRenders(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	window := app.NewWindow("Test")
	window.Resize(fyne.NewSize(600, 500))

	content := CreateMainWindow(app, window, testWindowDeps())
	window.SetContent(content)

	// Force a canvas refresh
	window.Canvas().Refresh(content)
}

func TestCreateMainWindowWithEmptyPorts(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	deps := testWindowDeps()
	deps.GetSerialPorts = func() []string { return []string{} }

	window := app.NewWindow("Test")
	content := CreateMainWindow(app, window, deps)

	if content == nil {
		t.Fatal("CreateMainWindow with empty ports returned nil")
	}
}

func TestCreateMainWindowWithEmptyCommands(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	deps := testWindowDeps()
	deps.GetCommandNames = func() []string { return []string{} }
	deps.GetCXRFCommandNames = func() []string { return []string{} }

	window := app.NewWindow("Test")
	content := CreateMainWindow(app, window, deps)

	if content == nil {
		t.Fatal("CreateMainWindow with empty commands returned nil")
	}
}

func TestWindowDepsCallbacks(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	getSerialPortsCalled := false
	getCommandNamesCalled := false
	getCXRFCommandNamesCalled := false

	deps := WindowDeps{
		LogoResource: fyne.NewStaticResource("test_logo", []byte{}),
		GetSerialPorts: func() []string {
			getSerialPortsCalled = true
			return []string{"/dev/test"}
		},
		GetCommandNames: func() []string {
			getCommandNamesCalled = true
			return []string{"VER"}
		},
		GetCXRFCommandNames: func() []string {
			getCXRFCommandNamesCalled = true
			return []string{"version"}
		},
		GetCommand:        func(name string) *Command { return nil },
		GetCXRFCommand:    func(name string) *Command { return nil },
		SendCommand:       func(port, scType, cmd string, speed int) (CommandResult, error) { return CommandResult{}, nil },
		Authenticate:      func(port, scType string, speed int) error { return nil },
		OpenSerialMonitor: func(myApp fyne.App, port, scType string) {},
		ShowGuideWindow:   func(myApp fyne.App) {},
	}

	window := app.NewWindow("Test")
	_ = CreateMainWindow(app, window, deps)

	if !getSerialPortsCalled {
		t.Error("GetSerialPorts was not called")
	}
	if !getCommandNamesCalled {
		t.Error("GetCommandNames was not called")
	}
	if !getCXRFCommandNamesCalled {
		t.Error("GetCXRFCommandNames was not called")
	}
}

func TestCommandHasSubcommands(t *testing.T) {
	tests := []struct {
		name     string
		cmd      Command
		expected bool
	}{
		{
			name:     "with subcommands",
			cmd:      Command{Name: "EEP", Subcommands: []string{"GET", "SET"}},
			expected: true,
		},
		{
			name:     "empty subcommands",
			cmd:      Command{Name: "VER", Subcommands: []string{}},
			expected: false,
		},
		{
			name:     "nil subcommands",
			cmd:      Command{Name: "AUTH"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.HasSubcommands()
			if result != tt.expected {
				t.Errorf("HasSubcommands() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCommandStruct(t *testing.T) {
	cmd := Command{
		Name:        "TEST",
		Subcommands: []string{"SUB1", "SUB2"},
		Description: "Test command",
	}

	if cmd.Name != "TEST" {
		t.Errorf("Command.Name = %q, want %q", cmd.Name, "TEST")
	}

	if len(cmd.Subcommands) != 2 {
		t.Errorf("Command.Subcommands length = %d, want 2", len(cmd.Subcommands))
	}

	if cmd.Description != "Test command" {
		t.Errorf("Command.Description = %q, want %q", cmd.Description, "Test command")
	}
}

func TestCreateMainWindowSendCommandError(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	deps := testWindowDeps()
	deps.SendCommand = func(port, scType, cmd string, speed int) (CommandResult, error) {
		return CommandResult{}, errors.New("send error")
	}

	window := app.NewWindow("Test")
	content := CreateMainWindow(app, window, deps)

	if content == nil {
		t.Fatal("CreateMainWindow returned nil")
	}
}

func TestCreateMainWindowAuthError(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	deps := testWindowDeps()
	deps.Authenticate = func(port, scType string, speed int) error {
		return errors.New("auth error")
	}

	window := app.NewWindow("Test")
	content := CreateMainWindow(app, window, deps)

	if content == nil {
		t.Fatal("CreateMainWindow returned nil")
	}
}
