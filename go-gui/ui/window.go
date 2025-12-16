// Package ui provides the main application window.
package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Command represents a PS3 Syscon command with optional subcommands.
type Command struct {
	Name        string
	Subcommands []string
	Description string
}

// HasSubcommands returns true if the command has subcommands.
func (c *Command) HasSubcommands() bool {
	return len(c.Subcommands) > 0
}

// WindowDeps contains dependencies for the main window.
type WindowDeps struct {
	LogoResource       fyne.Resource
	GetSerialPorts     func() []string
	GetCommandNames    func() []string
	GetCXRFCommandNames func() []string
	GetCommand         func(name string) *Command
	GetCXRFCommand     func(name string) *Command
	SendCommand        func(port, scType, cmd string, speed int) (CommandResult, error)
	Authenticate       func(port, scType string, speed int) error
	OpenSerialMonitor  func(myApp fyne.App, port, scType string)
	ShowGuideWindow    func(myApp fyne.App)
}

// CreateMainWindow builds the main application window content.
func CreateMainWindow(myApp fyne.App, myWindow fyne.Window, deps WindowDeps) fyne.CanvasObject {
	header := CreateHeader(deps.LogoResource)

	// Connection section
	portSelect := widget.NewSelect(deps.GetSerialPorts(), nil)
	portSelect.PlaceHolder = "Select serial port..."

	refreshBtn := widget.NewButton("Refresh", func() {
		portSelect.Options = deps.GetSerialPorts()
		portSelect.Refresh()
	})
	refreshBtn.Importance = widget.LowImportance

	scTypeSelect := widget.NewSelect([]string{"CXR", "CXRF", "SW"}, nil)
	scTypeSelect.SetSelected("CXR")

	// Mode descriptions
	modeDesc := widget.NewLabel("External commands via UART")
	modeDesc.TextStyle = fyne.TextStyle{Italic: true}

	connectionContent := container.NewVBox(
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabel("Serial Port"),
				container.NewBorder(nil, nil, nil, refreshBtn, portSelect),
			),
			container.NewVBox(
				widget.NewLabel("Mode"),
				scTypeSelect,
			),
		),
		modeDesc,
	)

	connectionCard := CreateCard("CONNECTION", connectionContent)

	// CXR command selection widgets
	cmdSelectEntry := widget.NewSelectEntry(deps.GetCommandNames())
	cmdSelectEntry.PlaceHolder = "Select command..."

	subCmdSelect := widget.NewSelect([]string{}, nil)
	subCmdSelect.PlaceHolder = "Subcommand"
	subCmdSelect.Disable()

	argsEntry := widget.NewEntry()
	argsEntry.SetPlaceHolder("Arguments")

	cxrCommandContent := container.NewGridWithColumns(3,
		container.NewVBox(widget.NewLabel("Command"), cmdSelectEntry),
		container.NewVBox(widget.NewLabel("Subcommand"), subCmdSelect),
		container.NewVBox(widget.NewLabel("Arguments"), argsEntry),
	)

	// CXRF command selection widgets
	cxrfCmdSelectEntry := widget.NewSelectEntry(deps.GetCXRFCommandNames())
	cxrfCmdSelectEntry.PlaceHolder = "Select command..."

	cxrfSubCmdSelect := widget.NewSelect([]string{}, nil)
	cxrfSubCmdSelect.PlaceHolder = "Subcommand"
	cxrfSubCmdSelect.Disable()

	cxrfArgsEntry := widget.NewEntry()
	cxrfArgsEntry.SetPlaceHolder("Arguments")

	cxrfDescLabel := widget.NewLabel("")
	cxrfDescLabel.Wrapping = fyne.TextWrapWord
	cxrfDescLabel.TextStyle = fyne.TextStyle{Italic: true}

	cxrfCommandContent := container.NewVBox(
		container.NewGridWithColumns(3,
			container.NewVBox(widget.NewLabel("Command"), cxrfCmdSelectEntry),
			container.NewVBox(widget.NewLabel("Subcommand"), cxrfSubCmdSelect),
			container.NewVBox(widget.NewLabel("Arguments"), cxrfArgsEntry),
		),
		cxrfDescLabel,
	)

	// Legacy command entry for SW mode
	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("Enter raw command...")

	legacyCommandContent := container.NewVBox(
		widget.NewLabel("Raw Command"),
		commandEntry,
	)

	// Command sections
	cxrSection := cxrCommandContent
	cxrfSection := cxrfCommandContent
	legacySection := legacyCommandContent

	commandSection := container.NewStack(cxrSection)
	commandCard := CreateCard("COMMAND", commandSection)

	// Output terminal
	outputText := widget.NewMultiLineEntry()
	outputText.SetMinRowsVisible(12)
	outputText.Wrapping = fyne.TextWrapWord
	outputText.TextStyle = fyne.TextStyle{Monospace: true}

	// Terminal header with clear button
	clearBtn := widget.NewButton("Clear", func() {
		outputText.SetText("")
	})
	clearBtn.Importance = widget.LowImportance

	terminalHeader := container.NewBorder(nil, nil,
		canvas.NewText("TERMINAL OUTPUT", ColorPrimary),
		clearBtn,
	)

	terminalBg := canvas.NewRectangle(ColorInputBg)
	terminalBg.CornerRadius = 6

	terminalContent := container.NewBorder(
		terminalHeader,
		nil, nil, nil,
		container.NewStack(terminalBg, container.NewPadded(outputText)),
	)

	terminalCard := CreateCard("OUTPUT", terminalContent)

	// Update mode description
	updateModeDesc := func(mode string) {
		switch mode {
		case "CXR":
			modeDesc.SetText("External commands via UART (57600 baud)")
		case "CXRF":
			modeDesc.SetText("Internal commands - DIAG mode (115200 baud)")
		case "SW":
			modeDesc.SetText("Legacy raw command mode")
		}
	}

	// Update subcommand dropdown when command changes
	cmdSelectEntry.OnChanged = func(cmdName string) {
		argsEntry.SetText("")
		subCmdSelect.ClearSelected()

		cmd := deps.GetCommand(cmdName)
		if cmd != nil && cmd.HasSubcommands() {
			subCmdSelect.Options = cmd.Subcommands
			subCmdSelect.PlaceHolder = "Select..."
			subCmdSelect.Enable()
		} else {
			subCmdSelect.Options = []string{}
			subCmdSelect.PlaceHolder = "N/A"
			subCmdSelect.Disable()
		}
		subCmdSelect.Refresh()
	}

	// Update CXRF subcommand dropdown and description when command changes
	cxrfCmdSelectEntry.OnChanged = func(cmdName string) {
		cxrfArgsEntry.SetText("")
		cxrfSubCmdSelect.ClearSelected()

		cmd := deps.GetCXRFCommand(cmdName)
		if cmd != nil {
			cxrfDescLabel.SetText(cmd.Description)
			if cmd.HasSubcommands() {
				cxrfSubCmdSelect.Options = cmd.Subcommands
				cxrfSubCmdSelect.PlaceHolder = "Select..."
				cxrfSubCmdSelect.Enable()
			} else {
				cxrfSubCmdSelect.Options = []string{}
				cxrfSubCmdSelect.PlaceHolder = "N/A"
				cxrfSubCmdSelect.Disable()
			}
		} else {
			cxrfDescLabel.SetText("")
			cxrfSubCmdSelect.Options = []string{}
			cxrfSubCmdSelect.PlaceHolder = "N/A"
			cxrfSubCmdSelect.Disable()
		}
		cxrfSubCmdSelect.Refresh()
	}

	// Build command string based on mode
	buildCommand := func() string {
		switch scTypeSelect.Selected {
		case "CXR":
			return BuildCXRCommand(cmdSelectEntry.Text, subCmdSelect.Selected, argsEntry.Text)
		case "CXRF":
			return BuildCXRCommand(cxrfCmdSelectEntry.Text, cxrfSubCmdSelect.Selected, cxrfArgsEntry.Text)
		default:
			return commandEntry.Text
		}
	}

	sendCmd := func() {
		if portSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("serial port not selected"), myWindow)
			return
		}
		if scTypeSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("mode not selected"), myWindow)
			return
		}

		cmdText := buildCommand()
		if cmdText == "" {
			dialog.ShowError(fmt.Errorf("command is empty"), myWindow)
			return
		}

		serialSpeed := GetSerialSpeed(scTypeSelect.Selected)

		result, err := deps.SendCommand(portSelect.Selected, scTypeSelect.Selected, cmdText, serialSpeed)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		if result.Code == 0xFFFFFFFF {
			errMsg := "unknown error"
			if len(result.Data) > 0 {
				errMsg = result.Data[0]
			}
			dialog.ShowError(fmt.Errorf("command failed: %s", errMsg), myWindow)
			return
		}

		output := FormatCommandOutput(scTypeSelect.Selected, result)
		timestamp := time.Now().Format("15:04:05")
		outputText.SetText(outputText.Text + fmt.Sprintf("[%s] > %s\n%s\n", timestamp, cmdText, output))
	}

	// Enter key handlers
	cmdSelectEntry.OnSubmitted = func(s string) {
		cmd := deps.GetCommand(s)
		if cmd == nil || !cmd.HasSubcommands() {
			sendCmd()
		}
	}

	argsEntry.OnSubmitted = func(s string) { sendCmd() }

	subCmdSelect.OnChanged = func(s string) {
		if s != "" {
			myWindow.Canvas().Focus(argsEntry)
		}
	}

	cxrfCmdSelectEntry.OnSubmitted = func(s string) {
		cmd := deps.GetCXRFCommand(s)
		if cmd == nil || !cmd.HasSubcommands() {
			sendCmd()
		}
	}

	cxrfArgsEntry.OnSubmitted = func(s string) { sendCmd() }

	cxrfSubCmdSelect.OnChanged = func(s string) {
		if s != "" {
			myWindow.Canvas().Focus(cxrfArgsEntry)
		}
	}

	commandEntry.OnSubmitted = func(s string) { sendCmd() }

	// Auth function
	authCmd := func() {
		if portSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("serial port not selected"), myWindow)
			return
		}
		if scTypeSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("mode not selected"), myWindow)
			return
		}

		serialSpeed := GetSerialSpeed(scTypeSelect.Selected)
		timestamp := time.Now().Format("15:04:05")

		if err := deps.Authenticate(portSelect.Selected, scTypeSelect.Selected, serialSpeed); err != nil {
			outputText.SetText(outputText.Text + fmt.Sprintf("[%s] > AUTH\nFailed: %v\n", timestamp, err))
			dialog.ShowError(err, myWindow)
			return
		}

		outputText.SetText(outputText.Text + fmt.Sprintf("[%s] > AUTH\nAuth successful\n", timestamp))
	}

	// Action buttons
	sendBtn := widget.NewButton("Send Command", sendCmd)
	sendBtn.Importance = widget.HighImportance

	authBtn := widget.NewButton("Authenticate", authCmd)

	helpBtn := widget.NewButton("Help", func() {
		ShowHelpDialog(myApp, myWindow, func() {
			deps.ShowGuideWindow(myApp)
		})
	})
	helpBtn.Importance = widget.LowImportance

	monitorBtn := widget.NewButton("Serial Monitor", func() {
		deps.OpenSerialMonitor(myApp, portSelect.Selected, scTypeSelect.Selected)
	})
	monitorBtn.Importance = widget.LowImportance

	// Button layout
	actionButtons := container.NewGridWithColumns(4, sendBtn, authBtn, monitorBtn, helpBtn)

	// Toggle visibility based on SC type
	scTypeSelect.OnChanged = func(scType string) {
		updateModeDesc(scType)
		commandSection.RemoveAll()
		switch scType {
		case "CXR":
			commandSection.Add(cxrSection)
		case "CXRF":
			commandSection.Add(cxrfSection)
		case "SW":
			commandSection.Add(legacySection)
		}
		commandSection.Refresh()
	}

	// Main layout
	leftColumn := container.NewVBox(
		connectionCard,
		commandCard,
		actionButtons,
	)

	// Use border layout for main content
	mainContent := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		container.NewBorder(
			leftColumn,
			nil, nil, nil,
			terminalCard,
		),
	)

	// Background
	bg := canvas.NewRectangle(ColorBackground)

	return container.NewStack(bg, container.NewPadded(mainContent))
}
