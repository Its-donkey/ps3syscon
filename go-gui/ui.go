package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
)

// createCard creates a card-like container with a title and content.
func createCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	titleLabel := canvas.NewText(title, colorPrimary)
	titleLabel.TextSize = 12
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	bg := canvas.NewRectangle(colorSurfaceElevated)
	bg.CornerRadius = 8

	cardContent := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		content,
	)

	return container.NewStack(bg, container.NewPadded(cardContent))
}

// createHeader creates the modern header with branding.
func createHeader() fyne.CanvasObject {
	// Logo image
	logo := canvas.NewImageFromResource(LogoResource)
	logo.SetMinSize(fyne.NewSize(50, 50))
	logo.FillMode = canvas.ImageFillContain

	// Main title
	title := canvas.NewText("PS3 SYSCON", colorPrimary)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Subtitle
	subtitle := canvas.NewText("UART INTERFACE", colorTextMuted)
	subtitle.TextSize = 12

	// Decorative line
	line := canvas.NewRectangle(colorPrimary)
	line.SetMinSize(fyne.NewSize(0, 2))

	titleStack := container.NewVBox(
		container.NewHBox(logo, container.NewVBox(title, subtitle), layout.NewSpacer()),
		container.NewPadded(line),
	)

	return container.NewPadded(titleStack)
}

// filterOptions filters a list of options based on a search string.
func filterOptions(options []string, search string) []string {
	if search == "" {
		return options
	}
	search = strings.ToUpper(search)
	var filtered []string
	for _, opt := range options {
		if strings.Contains(strings.ToUpper(opt), search) {
			filtered = append(filtered, opt)
		}
	}
	return filtered
}

// buildCXRCommand builds a command string from parts for CXR mode.
func buildCXRCommand(cmd, subCmd, args string) string {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}
	parts := []string{cmd}
	if subCmd != "" {
		parts = append(parts, subCmd)
	}
	args = strings.TrimSpace(args)
	if args != "" {
		parts = append(parts, args)
	}
	return strings.Join(parts, " ")
}

// getSerialSpeed returns the appropriate baud rate for the SC type.
func getSerialSpeed(scType string) int {
	if scType == "CXRF" {
		return 115200
	}
	return 57600
}

// formatCommandOutput formats the command result for display.
func formatCommandOutput(scType string, result CommandResult) string {
	switch scType {
	case "CXR":
		return fmt.Sprintf("%08X %s", result.Code, strings.Join(result.Data, " "))
	case "SW":
		if len(result.Data) > 0 && !strings.Contains(result.Data[0], "\n") {
			return fmt.Sprintf("%08X %s", result.Code, strings.Join(result.Data, " "))
		}
		return fmt.Sprintf("%08X\n%s", result.Code, strings.Join(result.Data, ""))
	default:
		if len(result.Data) > 0 {
			return result.Data[0]
		}
		return ""
	}
}

// createMainWindow builds the main application window content.
func createMainWindow(myApp fyne.App, myWindow fyne.Window) fyne.CanvasObject {
	header := createHeader()

	// Connection section
	portSelect := widget.NewSelect(getSerialPorts(), nil)
	portSelect.PlaceHolder = "Select serial port..."

	refreshBtn := widget.NewButtonWithIcon("", fyne.NewStaticResource("refresh", nil), func() {
		portSelect.Options = getSerialPorts()
		portSelect.Refresh()
	})
	refreshBtn.Importance = widget.LowImportance

	// Use text button instead of icon since we may not have the icon resource
	refreshBtn = widget.NewButton("Refresh", func() {
		portSelect.Options = getSerialPorts()
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

	connectionCard := createCard("CONNECTION", connectionContent)

	// CXR command selection widgets
	cmdSelectEntry := widget.NewSelectEntry(GetCommandNames())
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
	cxrfCmdSelectEntry := widget.NewSelectEntry(GetCXRFCommandNames())
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
	commandCard := createCard("COMMAND", commandSection)

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
		canvas.NewText("TERMINAL OUTPUT", colorPrimary),
		clearBtn,
	)

	terminalBg := canvas.NewRectangle(colorInputBg)
	terminalBg.CornerRadius = 6

	terminalContent := container.NewBorder(
		terminalHeader,
		nil, nil, nil,
		container.NewStack(terminalBg, container.NewPadded(outputText)),
	)

	terminalCard := createCard("OUTPUT", terminalContent)

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

		cmd := GetCommand(cmdName)
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

		cmd := GetCXRFCommand(cmdName)
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
			return buildCXRCommand(cmdSelectEntry.Text, subCmdSelect.Selected, argsEntry.Text)
		case "CXRF":
			return buildCXRCommand(cxrfCmdSelectEntry.Text, cxrfSubCmdSelect.Selected, cxrfArgsEntry.Text)
		default:
			return commandEntry.Text
		}
	}

	sendCmd := func() {
		if portSelect.Selected == "" || scTypeSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("please select serial port and mode"), myWindow)
			return
		}

		cmdText := buildCommand()
		if cmdText == "" {
			dialog.ShowError(fmt.Errorf("please enter a command"), myWindow)
			return
		}

		serialSpeed := getSerialSpeed(scTypeSelect.Selected)

		ps3, err := NewPS3UART(portSelect.Selected, scTypeSelect.Selected, serialSpeed)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		defer ps3.Close()

		result := ps3.Command(cmdText, 1)
		if result.Code == 0xFFFFFFFF {
			dialog.ShowError(fmt.Errorf("command failed: %s", result.Data[0]), myWindow)
			return
		}

		output := formatCommandOutput(scTypeSelect.Selected, result)
		timestamp := time.Now().Format("15:04:05")
		outputText.SetText(outputText.Text + fmt.Sprintf("[%s] > %s\n%s\n", timestamp, cmdText, output))
	}

	// Enter key handlers
	cmdSelectEntry.OnSubmitted = func(s string) {
		cmd := GetCommand(s)
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
		cmd := GetCXRFCommand(s)
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
		if portSelect.Selected == "" || scTypeSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("please select serial port and mode"), myWindow)
			return
		}

		serialSpeed := getSerialSpeed(scTypeSelect.Selected)

		ps3, err := NewPS3UART(portSelect.Selected, scTypeSelect.Selected, serialSpeed)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		defer ps3.Close()

		result := ps3.Auth()
		timestamp := time.Now().Format("15:04:05")
		outputText.SetText(outputText.Text + fmt.Sprintf("[%s] > AUTH\n%s\n", timestamp, result))
	}

	// Help function
	showHelp := func() {
		helpText := `QUICK START
1. Select serial port and click Authenticate
2. Use dropdowns to build commands

CXR MODE - Common Commands
  EEP GET <addr> <len>    Read EEPROM
  EEP SET <addr> <len> <val>  Write EEPROM
  ERRLOG GET <00-1F>      Read error log
  ERRLOG CLEAR            Clear errors
  VER                     Firmware version

CXRF MODE - Common Commands
  r <offset> [len]        Read from syscon
  w <offset> <value>      Write to syscon
  errlog / clearerrlog    Error log ops
  eepcsum                 Verify checksum
  version                 SC firmware ver

TIPS
• Press Enter to send command
• CXRF requires DIAG pin grounded
• Auth is required before commands`

		dialog.ShowInformation("Quick Reference", helpText, myWindow)
	}

	// Action buttons
	sendBtn := widget.NewButton("Send Command", sendCmd)
	sendBtn.Importance = widget.HighImportance

	authBtn := widget.NewButton("Authenticate", authCmd)

	helpBtn := widget.NewButton("Help", showHelp)
	helpBtn.Importance = widget.LowImportance

	monitorBtn := widget.NewButton("Serial Monitor", func() {
		openSerialMonitor(myApp, portSelect.Selected, scTypeSelect.Selected)
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
	bg := canvas.NewRectangle(colorBackground)

	return container.NewStack(bg, container.NewPadded(mainContent))
}

// openSerialMonitor opens the serial monitor window.
func openSerialMonitor(myApp fyne.App, defaultPort, scType string) {
	monitorWindow := myApp.NewWindow("Serial Monitor")
	monitorWindow.Resize(fyne.NewSize(600, 500))

	// Header
	title := canvas.NewText("SERIAL MONITOR", colorPrimary)
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}

	portSelect := widget.NewSelect(getSerialPorts(), nil)
	portSelect.PlaceHolder = "Select port..."
	if defaultPort != "" {
		portSelect.SetSelected(defaultPort)
	}

	baudSelect := widget.NewSelect([]string{"57600", "115200"}, nil)
	if scType == "CXRF" {
		baudSelect.SetSelected("115200")
	} else {
		baudSelect.SetSelected("57600")
	}

	outputText := widget.NewMultiLineEntry()
	outputText.SetMinRowsVisible(18)
	outputText.Wrapping = fyne.TextWrapWord
	outputText.TextStyle = fyne.TextStyle{Monospace: true}

	var port serial.Port
	var monitoring bool
	var stopChan chan struct{}

	// Status indicator
	statusLabel := canvas.NewText("DISCONNECTED", colorTextMuted)
	statusLabel.TextSize = 10

	startBtn := widget.NewButton("Start", nil)
	startBtn.Importance = widget.HighImportance
	stopBtn := widget.NewButton("Stop", nil)
	stopBtn.Disable()

	clearBtn := widget.NewButton("Clear", func() {
		outputText.SetText("")
	})
	clearBtn.Importance = widget.LowImportance

	startBtn.OnTapped = func() {
		if portSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("please select a serial port"), monitorWindow)
			return
		}

		baudRate := 57600
		if baudSelect.Selected == "115200" {
			baudRate = 115200
		}

		mode := &serial.Mode{
			BaudRate: baudRate,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		}

		var err error
		port, err = serial.Open(portSelect.Selected, mode)
		if err != nil {
			dialog.ShowError(err, monitorWindow)
			return
		}

		port.SetReadTimeout(100 * time.Millisecond)
		monitoring = true
		stopChan = make(chan struct{})

		startBtn.Disable()
		stopBtn.Enable()
		portSelect.Disable()
		baudSelect.Disable()
		statusLabel.Text = "CONNECTED"
		statusLabel.Color = colorSuccess
		statusLabel.Refresh()

		go func() {
			buf := make([]byte, 1024)
			for monitoring {
				select {
				case <-stopChan:
					return
				default:
					n, err := port.Read(buf)
					if err == nil && n > 0 {
						text := string(buf[:n])
						fyne.Do(func() {
							outputText.SetText(outputText.Text + text)
						})
					}
				}
			}
		}()
	}

	stopBtn.OnTapped = func() {
		monitoring = false
		if stopChan != nil {
			close(stopChan)
		}
		if port != nil {
			port.Close()
		}

		startBtn.Enable()
		stopBtn.Disable()
		portSelect.Enable()
		baudSelect.Enable()
		statusLabel.Text = "DISCONNECTED"
		statusLabel.Color = colorTextMuted
		statusLabel.Refresh()
	}

	monitorWindow.SetOnClosed(func() {
		if monitoring {
			stopBtn.OnTapped()
		}
	})

	// Layout
	configRow := container.NewGridWithColumns(2,
		container.NewVBox(widget.NewLabel("Port"), portSelect),
		container.NewVBox(widget.NewLabel("Baud Rate"), baudSelect),
	)

	buttonRow := container.NewHBox(startBtn, stopBtn, layout.NewSpacer(), statusLabel, layout.NewSpacer(), clearBtn)

	terminalBg := canvas.NewRectangle(colorInputBg)
	terminalBg.CornerRadius = 6

	content := container.NewBorder(
		container.NewVBox(
			container.NewPadded(title),
			widget.NewSeparator(),
			container.NewPadded(configRow),
			container.NewPadded(buttonRow),
		),
		nil, nil, nil,
		container.NewPadded(container.NewStack(terminalBg, container.NewPadded(outputText))),
	)

	bg := canvas.NewRectangle(colorBackground)
	monitorWindow.SetContent(container.NewStack(bg, content))
	monitorWindow.Show()
}
