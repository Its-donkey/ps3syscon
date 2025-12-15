package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"go.bug.st/serial"
)

// SerialPort interface abstracts serial port operations for testing.
type SerialPort interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	SetReadTimeout(time.Duration) error
}

// SerialPortOpener is a function type for opening serial ports.
type SerialPortOpener func(portName string, mode *serial.Mode) (SerialPort, error)

// DefaultSerialPortOpener opens a real serial port.
var DefaultSerialPortOpener SerialPortOpener = func(portName string, mode *serial.Mode) (SerialPort, error) {
	return serial.Open(portName, mode)
}

// PS3UART handles serial communication with PS3 Syscon.
type PS3UART struct {
	port        SerialPort
	scType      string
	serialSpeed int
}

// CommandResult holds the result of a command execution.
type CommandResult struct {
	Code uint32
	Data []string
}

// NewPS3UART creates a new PS3UART connection.
func NewPS3UART(portName, scType string, serialSpeed int) (*PS3UART, error) {
	return NewPS3UARTWithOpener(portName, scType, serialSpeed, DefaultSerialPortOpener)
}

// NewPS3UARTWithOpener creates a new PS3UART connection with a custom port opener.
func NewPS3UARTWithOpener(portName, scType string, serialSpeed int, opener SerialPortOpener) (*PS3UART, error) {
	mode := &serial.Mode{
		BaudRate: serialSpeed,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := opener(portName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}

	port.SetReadTimeout(100 * time.Millisecond)

	return &PS3UART{
		port:        port,
		scType:      scType,
		serialSpeed: serialSpeed,
	}, nil
}

// NewPS3UARTWithPort creates a new PS3UART with an existing port (for testing).
func NewPS3UARTWithPort(port SerialPort, scType string, serialSpeed int) *PS3UART {
	return &PS3UART{
		port:        port,
		scType:      scType,
		serialSpeed: serialSpeed,
	}
}

// Close closes the serial connection.
func (p *PS3UART) Close() error {
	if p.port != nil {
		return p.port.Close()
	}
	return nil
}

// send writes ASCII data to serial port.
func (p *PS3UART) send(data string) error {
	_, err := p.port.Write([]byte(data))
	return err
}

// receive reads available data from serial port.
func (p *PS3UART) receive() (string, error) {
	buf := make([]byte, 4096)
	var result []byte

	for {
		n, err := p.port.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		result = append(result, buf[:n]...)
	}

	return string(result), nil
}

// Command sends a command and returns the result.
func (p *PS3UART) Command(cmd string, waitSec float64) CommandResult {
	switch p.scType {
	case "CXR":
		return p.commandCXR(cmd, waitSec)
	case "SW":
		return p.commandSW(cmd, waitSec)
	default:
		return p.commandCXRF(cmd, waitSec)
	}
}

func (p *PS3UART) commandCXR(cmd string, waitSec float64) CommandResult {
	length := len(cmd)
	checksum := 0
	for _, c := range cmd {
		checksum += int(c)
	}
	checksum %= 0x100

	if length <= 10 {
		p.send(fmt.Sprintf("C:%02X:%s\r\n", checksum, cmd))
	} else {
		j := 10
		p.send(fmt.Sprintf("C:%02X:%s", checksum, cmd[0:j]))
		for i := length - j; i > 15; i -= 15 {
			p.send(cmd[j : j+15])
			j += 15
		}
		p.send(cmd[j:] + "\r\n")
	}

	time.Sleep(time.Duration(waitSec * float64(time.Second)))
	answer, _ := p.receive()
	answer = strings.TrimSpace(answer)

	parts := strings.Split(answer, ":")
	if len(parts) != 3 {
		return CommandResult{Code: 0xFFFFFFFF, Data: []string{"Answer length"}}
	}

	checksum = 0
	for _, c := range parts[2] {
		checksum += int(c)
	}
	checksum %= 0x100

	if parts[0] != "R" && parts[0] != "E" {
		return CommandResult{Code: 0xFFFFFFFF, Data: []string{"Magic"}}
	}
	if parts[1] != fmt.Sprintf("%02X", checksum) {
		return CommandResult{Code: 0xFFFFFFFF, Data: []string{"Checksum"}}
	}

	data := strings.Split(parts[2], " ")
	if (parts[0] == "R" && len(data) < 2) || (parts[0] == "E" && len(data) != 2) {
		return CommandResult{Code: 0xFFFFFFFF, Data: []string{"Data length"}}
	}

	code := parseHexUint32(data[1])
	if data[0] != "OK" || len(data) < 2 {
		return CommandResult{Code: code, Data: []string{}}
	}
	return CommandResult{Code: code, Data: data[2:]}
}

func (p *PS3UART) commandSW(cmd string, waitSec float64) CommandResult {
	length := len(cmd)
	if length >= 0x40 {
		result := p.Command("SETCMDLONG FF FF", 1)
		if result.Code != 0 {
			return CommandResult{Code: 0xFFFFFFFF, Data: []string{"Setcmdlong"}}
		}
	}

	checksum := 0
	for _, c := range cmd {
		checksum += int(c)
	}
	checksum %= 0x100

	p.send(fmt.Sprintf("%s:%02X\r\n", cmd, checksum))

	time.Sleep(time.Duration(waitSec * float64(time.Second)))
	answer, _ := p.receive()
	answer = strings.TrimSpace(answer)

	lines := strings.Split(answer, "\n")
	for i, line := range lines {
		line = strings.ReplaceAll(line, "\n", "")
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return CommandResult{Code: 0xFFFFFFFF, Data: []string{"Answer length"}}
		}

		checksum = 0
		for _, c := range parts[0] {
			checksum += int(c)
		}
		checksum %= 0x100

		if parts[1] != fmt.Sprintf("%02X", checksum) {
			return CommandResult{Code: 0xFFFFFFFF, Data: []string{"Checksum"}}
		}
		lines[i] = parts[0] + "\n"
	}

	ret := strings.Split(strings.ReplaceAll(lines[len(lines)-1], "\n", ""), " ")
	if len(ret) < 2 || len(ret[1]) != 8 {
		return CommandResult{Code: 0, Data: lines}
	} else if len(lines) == 1 {
		return CommandResult{Code: parseHexUint32(ret[1]), Data: ret[2:]}
	}
	return CommandResult{Code: parseHexUint32(ret[1]), Data: lines[:len(lines)-1]}
}

func (p *PS3UART) commandCXRF(cmd string, waitSec float64) CommandResult {
	p.send(cmd + "\r\n")
	time.Sleep(time.Duration(waitSec * float64(time.Second)))
	answer, _ := p.receive()
	answer = strings.TrimSpace(answer)
	return CommandResult{Code: 0, Data: []string{answer}}
}

func parseHexUint32(s string) uint32 {
	var val uint32
	fmt.Sscanf(s, "%x", &val)
	return val
}

// Auth performs authentication with the Syscon.
func (p *PS3UART) Auth() string {
	if p.scType == "CXR" || p.scType == "SW" {
		auth1r := p.Command("AUTH1 10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 1)
		if auth1r.Code == 0 && len(auth1r.Data) > 0 {
			auth1rBytes, err := hex.DecodeString(auth1r.Data[0])
			if err != nil {
				return "Auth1 response decode error"
			}

			if len(auth1rBytes) >= 0x40 && bytesEqual(auth1rBytes[0:0x10], auth1rHdr) {
				data, err := aesDecryptCBC(sc2tb, zeroIV, auth1rBytes[0x10:0x40])
				if err != nil {
					return "Decryption error"
				}

				if bytesEqual(data[0x8:0x10], zeroIV[0:0x8]) &&
					bytesEqual(data[0x10:0x20], authValue) &&
					bytesEqual(data[0x20:0x30], zeroIV) {

					newData := make([]byte, 0x30)
					copy(newData[0:0x8], data[0x8:0x10])
					copy(newData[0x8:0x10], data[0x0:0x8])
					// Rest is zeros

					auth2Body, err := aesEncryptCBC(tb2sc, zeroIV, newData)
					if err != nil {
						return "Encryption error"
					}

					auth2Cmd := "AUTH2 " + strings.ToUpper(hex.EncodeToString(append(auth2Header, auth2Body...)))
					auth2r := p.Command(auth2Cmd, 1)
					if auth2r.Code == 0 {
						return "Auth successful"
					}
					return "Auth failed"
				}
				return "Auth1 response body invalid"
			}
			return "Auth1 response header invalid"
		}
		return "Auth1 response invalid"
	}

	// CXRF mode
	scopen := p.Command("scopen", 1)
	if len(scopen.Data) > 0 && strings.Contains(scopen.Data[0], "SC_READY") {
		auth1r := p.Command("10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 1)
		if len(auth1r.Data) > 0 {
			parts := strings.Split(auth1r.Data[0], "\r")
			if len(parts) > 1 && len(parts[1]) > 1 {
				auth1rHex := parts[1][1:]
				if len(auth1rHex) == 128 {
					auth1rBytes, _ := hex.DecodeString(auth1rHex)
					if bytesEqual(auth1rBytes[0:0x10], auth1rHdr) {
						data, _ := aesDecryptCBC(sc2tb, zeroIV, auth1rBytes[0x10:0x40])
						if bytesEqual(data[0x8:0x10], zeroIV[0:0x8]) &&
							bytesEqual(data[0x10:0x20], authValue) &&
							bytesEqual(data[0x20:0x30], zeroIV) {

							newData := make([]byte, 0x30)
							copy(newData[0:0x8], data[0x8:0x10])
							copy(newData[0x8:0x10], data[0x0:0x8])

							auth2Body, _ := aesEncryptCBC(tb2sc, zeroIV, newData)
							auth2Cmd := strings.ToUpper(hex.EncodeToString(append(auth2Header, auth2Body...)))
							auth2r := p.Command(auth2Cmd, 1)
							if len(auth2r.Data) > 0 && strings.Contains(auth2r.Data[0], "SC_SUCCESS") {
								return "Auth successful"
							}
							return "Auth failed"
						}
						return "Auth1 response body invalid"
					}
					return "Auth1 response header invalid"
				}
				return "Auth1 response invalid"
			}
		}
	}
	return "scopen response invalid"
}

// getSerialPorts returns a list of available serial ports.
func getSerialPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		return []string{}
	}
	return ports
}
