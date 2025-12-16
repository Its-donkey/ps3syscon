package ui

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// mockSerialPort implements SerialPort for testing
type mockSerialPort struct {
	mu          sync.Mutex
	readData    []byte
	readErr     error
	closed      bool
	readTimeout time.Duration
}

func (m *mockSerialPort) Read(buf []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.readErr != nil {
		return 0, m.readErr
	}
	if len(m.readData) > 0 {
		n := copy(buf, m.readData)
		m.readData = m.readData[n:]
		return n, nil
	}
	return 0, nil
}

func (m *mockSerialPort) Write(buf []byte) (int, error) {
	return len(buf), nil
}

func (m *mockSerialPort) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (m *mockSerialPort) SetReadTimeout(d time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.readTimeout = d
	return nil
}

func (m *mockSerialPort) IsClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

func TestNewSerialMonitor(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return &mockSerialPort{}, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	if monitor == nil {
		t.Fatal("NewSerialMonitor returned nil")
	}

	if monitor.outputText != outputText {
		t.Error("outputText not set correctly")
	}

	if monitor.openPort == nil {
		t.Error("openPort not set correctly")
	}
}

func TestSerialMonitorStart(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	mockPort := &mockSerialPort{}
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return mockPort, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	err := monitor.Start(context.Background(), "/dev/test", 57600)
	if err != nil {
		t.Errorf("Start returned error: %v", err)
	}

	if !monitor.IsRunning() {
		t.Error("Monitor should be running after Start")
	}

	monitor.Stop()
}

func TestSerialMonitorStartError(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	expectedErr := errors.New("port open failed")
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return nil, expectedErr
	}

	monitor := NewSerialMonitor(outputText, openPort)

	err := monitor.Start(context.Background(), "/dev/test", 57600)
	if err != expectedErr {
		t.Errorf("Start should return port open error, got: %v", err)
	}

	if monitor.IsRunning() {
		t.Error("Monitor should not be running after failed Start")
	}
}

func TestSerialMonitorStop(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	mockPort := &mockSerialPort{}
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return mockPort, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	_ = monitor.Start(context.Background(), "/dev/test", 57600)
	monitor.Stop()

	if monitor.IsRunning() {
		t.Error("Monitor should not be running after Stop")
	}

	if !mockPort.IsClosed() {
		t.Error("Port should be closed after Stop")
	}
}

func TestSerialMonitorStopWhenNotRunning(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return &mockSerialPort{}, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	// Stop when not running should not panic
	monitor.Stop()

	if monitor.IsRunning() {
		t.Error("Monitor should not be running")
	}
}

func TestSerialMonitorIsRunning(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return &mockSerialPort{}, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	if monitor.IsRunning() {
		t.Error("New monitor should not be running")
	}

	_ = monitor.Start(context.Background(), "/dev/test", 57600)

	if !monitor.IsRunning() {
		t.Error("Monitor should be running after Start")
	}

	monitor.Stop()

	if monitor.IsRunning() {
		t.Error("Monitor should not be running after Stop")
	}
}

func TestSerialMonitorStartWhileRunning(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return &mockSerialPort{}, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	_ = monitor.Start(context.Background(), "/dev/test", 57600)

	// Starting again while running should return nil (no-op)
	err := monitor.Start(context.Background(), "/dev/test", 57600)
	if err != nil {
		t.Errorf("Start while running should return nil, got: %v", err)
	}

	monitor.Stop()
}

func TestSerialMonitorMultipleStartStop(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return &mockSerialPort{}, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	// Multiple start/stop cycles should work
	for i := 0; i < 3; i++ {
		err := monitor.Start(context.Background(), "/dev/test", 57600)
		if err != nil {
			t.Errorf("Start cycle %d failed: %v", i, err)
		}
		monitor.Stop()
	}
}

func TestSerialMonitorSetReadTimeoutError(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	outputText := widget.NewMultiLineEntry()

	// Create a mock port that returns error on SetReadTimeout
	mockPort := &mockSerialPortWithTimeoutError{}
	openPort := func(portName string, baudRate int) (SerialPort, error) {
		return mockPort, nil
	}

	monitor := NewSerialMonitor(outputText, openPort)

	err := monitor.Start(context.Background(), "/dev/test", 57600)
	if err == nil {
		t.Error("Start should return error when SetReadTimeout fails")
		monitor.Stop()
	}
}

// mockSerialPortWithTimeoutError is a mock that fails on SetReadTimeout
type mockSerialPortWithTimeoutError struct {
	mockSerialPort
}

func (m *mockSerialPortWithTimeoutError) SetReadTimeout(d time.Duration) error {
	return errors.New("timeout error")
}
