// Package ui provides serial port monitoring with UI integration.
package ui

import (
	"context"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// SerialPort abstracts serial port operations.
type SerialPort interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	SetReadTimeout(time.Duration) error
}

// PortOpener is a function that opens a serial port.
type PortOpener func(portName string, baudRate int) (SerialPort, error)

// SerialMonitor manages serial port monitoring with proper lifecycle control.
type SerialMonitor struct {
	mu         sync.Mutex
	port       SerialPort
	cancel     context.CancelFunc
	running    bool
	outputText *widget.Entry
	openPort   PortOpener
}

// NewSerialMonitor creates a new serial monitor instance.
func NewSerialMonitor(outputText *widget.Entry, openPort PortOpener) *SerialMonitor {
	return &SerialMonitor{
		outputText: outputText,
		openPort:   openPort,
	}
}

// Start begins monitoring the serial port.
// Returns an error if the port cannot be opened.
func (m *SerialMonitor) Start(ctx context.Context, portName string, baudRate int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return nil
	}

	port, err := m.openPort(portName, baudRate)
	if err != nil {
		return err
	}

	if err := port.SetReadTimeout(100 * time.Millisecond); err != nil {
		port.Close()
		return err
	}

	m.port = port
	m.running = true

	// Create cancellable context for this monitoring session
	monitorCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	go m.readLoop(monitorCtx)

	return nil
}

// Stop stops monitoring and closes the serial port.
func (m *SerialMonitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}

	if m.port != nil {
		m.port.Close()
		m.port = nil
	}

	m.running = false
}

// IsRunning returns whether the monitor is actively reading.
func (m *SerialMonitor) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func (m *SerialMonitor) readLoop(ctx context.Context) {
	buf := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			m.mu.Lock()
			port := m.port
			m.mu.Unlock()

			if port == nil {
				return
			}

			n, err := port.Read(buf)
			if err == nil && n > 0 {
				text := string(buf[:n])
				fyne.Do(func() {
					m.outputText.SetText(m.outputText.Text + text)
				})
			}
		}
	}
}
