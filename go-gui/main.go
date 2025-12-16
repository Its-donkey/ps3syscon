// Package main provides the entry point for the PS3 Syscon UART GUI tool.
package main

import (
	"ps3syscon-gui/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"go.bug.st/serial"
)

func main() {
	myApp := app.New()

	// Apply custom dark theme
	myApp.Settings().SetTheme(&ui.PS3Theme{})

	// Set application icon
	myApp.SetIcon(ui.IconResource)

	myWindow := myApp.NewWindow("PS3 Syscon UART Tool")
	myWindow.Resize(fyne.NewSize(700, 650))
	myWindow.SetIcon(ui.IconResource)

	// Set up the application menu
	ui.SetupMainMenu(myApp, myWindow)

	myWindow.Show()

	// Show disclaimer first - callbacks will handle accept/decline
	ui.ShowDisclaimer(myWindow,
		func() {
			// On accept - show main window
			deps := ui.WindowDeps{
				LogoResource:        ui.LogoResource,
				GetSerialPorts:      getSerialPorts,
				GetCommandNames:     GetCommandNames,
				GetCXRFCommandNames: GetCXRFCommandNames,
				GetCommand:          adaptCommand,
				GetCXRFCommand:      adaptCXRFCommand,
				SendCommand:         sendCommand,
				Authenticate:        authenticate,
				OpenSerialMonitor:   openSerialMonitor,
				ShowGuideWindow:     ui.ShowGuideWindow,
			}
			myWindow.SetContent(ui.CreateMainWindow(myApp, myWindow, deps))
		},
		func() {
			// On decline - quit
			myApp.Quit()
		},
	)

	myApp.Run()
}

// adaptCommand adapts the Command type from commands.go to ui.Command.
func adaptCommand(name string) *ui.Command {
	cmd := GetCommand(name)
	if cmd == nil {
		return nil
	}
	return &ui.Command{
		Name:        cmd.Name,
		Subcommands: cmd.Subcommands,
		Description: cmd.Description,
	}
}

// adaptCXRFCommand adapts the Command type for CXRF commands.
func adaptCXRFCommand(name string) *ui.Command {
	cmd := GetCXRFCommand(name)
	if cmd == nil {
		return nil
	}
	return &ui.Command{
		Name:        cmd.Name,
		Subcommands: cmd.Subcommands,
		Description: cmd.Description,
	}
}

// sendCommand wraps the serial command execution.
func sendCommand(port, scType, cmd string, speed int) (ui.CommandResult, error) {
	ps3, err := NewPS3UART(port, scType, speed)
	if err != nil {
		return ui.CommandResult{}, err
	}
	defer ps3.Close()

	result := ps3.Command(cmd, 1)
	return ui.CommandResult{
		Code: result.Code,
		Data: result.Data,
	}, nil
}

// authenticate wraps the serial authentication.
func authenticate(port, scType string, speed int) error {
	ps3, err := NewPS3UART(port, scType, speed)
	if err != nil {
		return err
	}
	defer ps3.Close()

	return ps3.Auth()
}

// openSerialMonitor wraps the ui.OpenSerialMonitor with dependencies.
func openSerialMonitor(myApp fyne.App, port, scType string) {
	deps := ui.MonitorDeps{
		GetSerialPorts: getSerialPorts,
		OpenPort:       openSerialPort,
	}
	ui.OpenSerialMonitor(myApp, port, scType, deps)
}

// openSerialPort opens a serial port with the given settings.
func openSerialPort(portName string, baudRate int) (ui.SerialPort, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	return serial.Open(portName, mode)
}
