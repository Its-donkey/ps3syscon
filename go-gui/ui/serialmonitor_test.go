package ui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func testMonitorDeps() MonitorDeps {
	return MonitorDeps{
		GetSerialPorts: func() []string {
			return []string{"/dev/ttyUSB0", "/dev/ttyUSB1"}
		},
		OpenPort: func(portName string, baudRate int) (SerialPort, error) {
			return &mockSerialPort{}, nil
		},
	}
}

func TestOpenSerialMonitor(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OpenSerialMonitor panicked: %v", r)
		}
	}()

	OpenSerialMonitor(app, "", "CXR", testMonitorDeps())
}

func TestOpenSerialMonitorWithPort(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	OpenSerialMonitor(app, "/dev/ttyUSB0", "CXR", testMonitorDeps())
}

func TestOpenSerialMonitorCXRF(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	OpenSerialMonitor(app, "/dev/ttyUSB0", "CXRF", testMonitorDeps())
}

func TestOpenSerialMonitorSW(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	OpenSerialMonitor(app, "", "SW", testMonitorDeps())
}

func TestOpenSerialMonitorAllModes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	modes := []string{"CXR", "CXRF", "SW", ""}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			OpenSerialMonitor(app, "", mode, testMonitorDeps())
		})
	}
}

func TestOpenSerialMonitorWithEmptyPorts(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	deps := MonitorDeps{
		GetSerialPorts: func() []string { return []string{} },
		OpenPort: func(portName string, baudRate int) (SerialPort, error) {
			return &mockSerialPort{}, nil
		},
	}

	OpenSerialMonitor(app, "", "CXR", deps)
}

func TestMonitorDepsStruct(t *testing.T) {
	deps := MonitorDeps{
		GetSerialPorts: func() []string { return []string{"port1"} },
		OpenPort: func(portName string, baudRate int) (SerialPort, error) {
			return nil, nil
		},
	}

	ports := deps.GetSerialPorts()
	if len(ports) != 1 || ports[0] != "port1" {
		t.Errorf("GetSerialPorts returned unexpected result: %v", ports)
	}

	port, err := deps.OpenPort("test", 57600)
	if port != nil || err != nil {
		t.Error("OpenPort should return nil, nil")
	}
}
