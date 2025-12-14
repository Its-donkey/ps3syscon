package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
)

// createBrandingHeader creates the header with logo and text.
func createBrandingHeader() fyne.CanvasObject {
	brandingText := canvas.NewText("BARNEY'S BETA", color.RGBA{100, 100, 100, 255})
	brandingText.TextSize = 20
	brandingText.Alignment = fyne.TextAlignCenter

	header := container.NewVBox(
		container.NewCenter(brandingText),
		widget.NewSeparator(),
	)
	return header
}

// createMainWindow builds the main application window content.
func createMainWindow(myApp fyne.App, myWindow fyne.Window) fyne.CanvasObject {
	header := createBrandingHeader()

	portSelect := widget.NewSelect(getSerialPorts(), nil)
	portSelect.PlaceHolder = "Select serial port"

	refreshBtn := widget.NewButton("Refresh", func() {
		portSelect.Options = getSerialPorts()
		portSelect.Refresh()
	})

	portRow := container.NewBorder(nil, nil, nil, refreshBtn, portSelect)

	scTypeSelect := widget.NewSelect([]string{"CXR", "CXRF", "SW"}, nil)
	scTypeSelect.SetSelected("CXR")

	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("Enter command (e.g., EEP GET 00)")

	outputText := widget.NewMultiLineEntry()
	outputText.SetMinRowsVisible(10)
	outputText.Wrapping = fyne.TextWrapWord

	sendCmd := func() {
		if portSelect.Selected == "" || scTypeSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("please select serial port and SC type"), myWindow)
			return
		}

		serialSpeed := 57600
		if scTypeSelect.Selected == "CXRF" {
			serialSpeed = 115200
		}

		ps3, err := NewPS3UART(portSelect.Selected, scTypeSelect.Selected, serialSpeed)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		defer ps3.Close()

		result := ps3.Command(commandEntry.Text, 1)
		if result.Code == 0xFFFFFFFF {
			dialog.ShowError(fmt.Errorf("command failed: %s", result.Data[0]), myWindow)
			return
		}

		var output string
		switch scTypeSelect.Selected {
		case "CXR":
			output = fmt.Sprintf("%08X %s", result.Code, strings.Join(result.Data, " "))
		case "SW":
			if len(result.Data) > 0 && !strings.Contains(result.Data[0], "\n") {
				output = fmt.Sprintf("%08X %s", result.Code, strings.Join(result.Data, " "))
			} else {
				output = fmt.Sprintf("%08X\n%s", result.Code, strings.Join(result.Data, ""))
			}
		default:
			if len(result.Data) > 0 {
				output = result.Data[0]
			}
		}

		outputText.SetText(outputText.Text + output + "\n")
	}

	authCmd := func() {
		if portSelect.Selected == "" || scTypeSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("please select serial port and SC type"), myWindow)
			return
		}

		serialSpeed := 57600
		if scTypeSelect.Selected == "CXRF" {
			serialSpeed = 115200
		}

		ps3, err := NewPS3UART(portSelect.Selected, scTypeSelect.Selected, serialSpeed)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		defer ps3.Close()

		result := ps3.Auth()
		dialog.ShowInformation("Authentication Result", result, myWindow)
	}

	showHelp := func() {
		helpText := `External mode:
EEP GET (get EEPROM address)
EEP SET (set EEPROM address value)
ERRLOG GET 00 (get errorlog from code 0 - repeat until 1F)

Internal mode:
eepcsum (check EEPROM checksum)
errlog (get errlog)
clearerrlog (clear errorlog)
r (read from eeprom address)
w (write to eeprom address)
fantbl (get/set/getini/setini/gettable/settable)
patchvereep (get patched version)

Read the PS3-Uart-Guide-V2.pdf for further information`

		dialog.ShowInformation("Available Commands", helpText, myWindow)
	}

	sendBtn := widget.NewButton("Send Command", sendCmd)
	authBtn := widget.NewButton("Auth", authCmd)
	helpBtn := widget.NewButton("Help", showHelp)

	serialMonitorBtn := widget.NewButton("Serial Monitor", func() {
		openSerialMonitor(myApp, portSelect.Selected, scTypeSelect.Selected)
	})

	btnRow1 := container.NewGridWithColumns(2, sendBtn, authBtn)
	btnRow2 := container.NewGridWithColumns(2, helpBtn, serialMonitorBtn)

	formContent := container.NewVBox(
		widget.NewLabel("Serial Port:"),
		portRow,
		widget.NewLabel("SC Type:"),
		scTypeSelect,
		widget.NewLabel("Command:"),
		commandEntry,
		widget.NewLabel("Output:"),
		outputText,
		btnRow1,
		btnRow2,
	)

	content := container.NewBorder(header, nil, nil, nil,
		container.NewPadded(formContent),
	)

	return content
}

// openSerialMonitor opens the serial monitor window.
func openSerialMonitor(myApp fyne.App, defaultPort, scType string) {
	monitorWindow := myApp.NewWindow("Serial Monitor")
	monitorWindow.Resize(fyne.NewSize(500, 400))

	portSelect := widget.NewSelect(getSerialPorts(), nil)
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
	outputText.SetMinRowsVisible(15)
	outputText.Wrapping = fyne.TextWrapWord

	var port serial.Port
	var monitoring bool
	var stopChan chan struct{}

	startBtn := widget.NewButton("Start", nil)
	stopBtn := widget.NewButton("Stop", nil)
	stopBtn.Disable()

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
						outputText.SetText(outputText.Text + text)
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
	}

	monitorWindow.SetOnClosed(func() {
		if monitoring {
			stopBtn.OnTapped()
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Serial Port:"),
		portSelect,
		widget.NewLabel("Baud Rate:"),
		baudSelect,
		container.NewGridWithColumns(2, startBtn, stopBtn),
		widget.NewLabel("Output:"),
		outputText,
	)

	monitorWindow.SetContent(content)
	monitorWindow.Show()
}
