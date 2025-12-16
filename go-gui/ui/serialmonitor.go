// Package ui provides the serial monitor window component.
package ui

import (
	"context"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// MonitorDeps contains dependencies for the serial monitor window.
type MonitorDeps struct {
	GetSerialPorts func() []string
	OpenPort       PortOpener
}

// OpenSerialMonitor opens the serial monitor window.
func OpenSerialMonitor(myApp fyne.App, defaultPort, scType string, deps MonitorDeps) {
	monitorWindow := myApp.NewWindow("Serial Monitor")
	monitorWindow.Resize(fyne.NewSize(600, 500))

	// Header
	title := canvas.NewText("SERIAL MONITOR", ColorPrimary)
	title.TextSize = 18
	title.TextStyle = fyne.TextStyle{Bold: true}

	portSelect := widget.NewSelect(deps.GetSerialPorts(), nil)
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

	monitor := NewSerialMonitor(outputText, deps.OpenPort)

	// Status indicator
	statusLabel := canvas.NewText("DISCONNECTED", ColorTextMuted)
	statusLabel.TextSize = 10

	startBtn := widget.NewButton("Start", nil)
	startBtn.Importance = widget.HighImportance
	stopBtn := widget.NewButton("Stop", nil)
	stopBtn.Disable()

	clearBtn := widget.NewButton("Clear", func() {
		outputText.SetText("")
	})
	clearBtn.Importance = widget.LowImportance

	updateStatus := func(connected bool) {
		if connected {
			statusLabel.Text = "CONNECTED"
			statusLabel.Color = ColorSuccess
			startBtn.Disable()
			stopBtn.Enable()
			portSelect.Disable()
			baudSelect.Disable()
		} else {
			statusLabel.Text = "DISCONNECTED"
			statusLabel.Color = ColorTextMuted
			startBtn.Enable()
			stopBtn.Disable()
			portSelect.Enable()
			baudSelect.Enable()
		}
		statusLabel.Refresh()
	}

	startBtn.OnTapped = func() {
		if portSelect.Selected == "" {
			dialog.ShowError(errors.New("mode not selected"), monitorWindow)
			return
		}

		baudRate := 57600
		if baudSelect.Selected == "115200" {
			baudRate = 115200
		}

		if err := monitor.Start(context.Background(), portSelect.Selected, baudRate); err != nil {
			dialog.ShowError(err, monitorWindow)
			return
		}

		updateStatus(true)
	}

	stopBtn.OnTapped = func() {
		monitor.Stop()
		updateStatus(false)
	}

	monitorWindow.SetOnClosed(func() {
		monitor.Stop()
	})

	// Layout
	configRow := container.NewGridWithColumns(2,
		container.NewVBox(widget.NewLabel("Port"), portSelect),
		container.NewVBox(widget.NewLabel("Baud Rate"), baudSelect),
	)

	buttonRow := container.NewHBox(startBtn, stopBtn, layout.NewSpacer(), statusLabel, layout.NewSpacer(), clearBtn)

	terminalBg := canvas.NewRectangle(ColorInputBg)
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

	bg := canvas.NewRectangle(ColorBackground)
	monitorWindow.SetContent(container.NewStack(bg, content))
	monitorWindow.Show()
}
