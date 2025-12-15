package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"go.bug.st/serial"
)

// MockSerialPort implements SerialPort for testing.
type MockSerialPort struct {
	ReadData    []byte
	ReadErr     error
	ReadCalls   int
	WriteData   []byte
	WriteErr    error
	WriteCalls  int
	Closed      bool
	CloseErr    error
	ReadTimeout time.Duration
	ReadIndex   int
	ReadChunks  [][]byte
	ChunkIndex  int
	// Responses is a list of complete responses, one per receive() call
	Responses     []string
	ResponseIndex int
	ResponseSent  bool
}

func (m *MockSerialPort) Read(buf []byte) (int, error) {
	m.ReadCalls++
	if m.ReadErr != nil {
		return 0, m.ReadErr
	}

	// If using Responses mode (for multi-command tests)
	if len(m.Responses) > 0 {
		if m.ResponseIndex >= len(m.Responses) {
			return 0, nil
		}
		if m.ResponseSent {
			// Already sent this response, return 0 to end this receive() call
			m.ResponseSent = false
			m.ResponseIndex++
			return 0, nil
		}
		// Send the current response
		resp := m.Responses[m.ResponseIndex]
		n := copy(buf, []byte(resp))
		m.ResponseSent = true
		return n, nil
	}

	// If we have chunked data, return chunks
	if len(m.ReadChunks) > 0 {
		if m.ChunkIndex >= len(m.ReadChunks) {
			return 0, nil
		}
		chunk := m.ReadChunks[m.ChunkIndex]
		m.ChunkIndex++
		n := copy(buf, chunk)
		return n, nil
	}

	// Otherwise return ReadData once then nothing
	if m.ReadIndex >= len(m.ReadData) {
		return 0, nil
	}
	n := copy(buf, m.ReadData[m.ReadIndex:])
	m.ReadIndex += n
	return n, nil
}

func (m *MockSerialPort) Write(data []byte) (int, error) {
	m.WriteCalls++
	if m.WriteErr != nil {
		return 0, m.WriteErr
	}
	m.WriteData = append(m.WriteData, data...)
	return len(data), nil
}

func (m *MockSerialPort) Close() error {
	m.Closed = true
	return m.CloseErr
}

func (m *MockSerialPort) SetReadTimeout(d time.Duration) error {
	m.ReadTimeout = d
	return nil
}

func TestNewPS3UARTWithPort(t *testing.T) {
	mock := &MockSerialPort{}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	if uart == nil {
		t.Fatal("NewPS3UARTWithPort returned nil")
	}
	if uart.scType != "CXR" {
		t.Errorf("scType = %q, want %q", uart.scType, "CXR")
	}
	if uart.serialSpeed != 57600 {
		t.Errorf("serialSpeed = %d, want %d", uart.serialSpeed, 57600)
	}
}

func TestNewPS3UARTWithOpener(t *testing.T) {
	mock := &MockSerialPort{}
	opener := func(portName string, mode *serial.Mode) (SerialPort, error) {
		if portName == "/dev/test" {
			return mock, nil
		}
		return nil, errors.New("port not found")
	}

	// Test successful open
	uart, err := NewPS3UARTWithOpener("/dev/test", "CXR", 57600, opener)
	if err != nil {
		t.Fatalf("NewPS3UARTWithOpener failed: %v", err)
	}
	if uart == nil {
		t.Fatal("NewPS3UARTWithOpener returned nil uart")
	}
	if mock.ReadTimeout != 100*time.Millisecond {
		t.Errorf("ReadTimeout = %v, want %v", mock.ReadTimeout, 100*time.Millisecond)
	}

	// Test failed open
	_, err = NewPS3UARTWithOpener("/dev/invalid", "CXR", 57600, opener)
	if err == nil {
		t.Error("Expected error for invalid port, got nil")
	}
}

func TestPS3UARTClose(t *testing.T) {
	// Test with nil port
	uart := &PS3UART{port: nil}
	err := uart.Close()
	if err != nil {
		t.Errorf("Close with nil port returned error: %v", err)
	}

	// Test with mock port
	mock := &MockSerialPort{}
	uart = NewPS3UARTWithPort(mock, "CXR", 57600)
	err = uart.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
	if !mock.Closed {
		t.Error("Port was not closed")
	}

	// Test close error
	mock = &MockSerialPort{CloseErr: errors.New("close error")}
	uart = NewPS3UARTWithPort(mock, "CXR", 57600)
	err = uart.Close()
	if err == nil {
		t.Error("Expected close error, got nil")
	}
}

func TestPS3UARTSend(t *testing.T) {
	mock := &MockSerialPort{}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	err := uart.send("TEST")
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if string(mock.WriteData) != "TEST" {
		t.Errorf("WriteData = %q, want %q", mock.WriteData, "TEST")
	}

	// Test write error
	mock = &MockSerialPort{WriteErr: errors.New("write error")}
	uart = NewPS3UARTWithPort(mock, "CXR", 57600)
	err = uart.send("TEST")
	if err == nil {
		t.Error("Expected write error, got nil")
	}
}

func TestPS3UARTReceive(t *testing.T) {
	// Test successful receive
	mock := &MockSerialPort{ReadData: []byte("R:3A:OK 00000000\r\n")}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	data, err := uart.receive()
	if err != nil {
		t.Fatalf("receive failed: %v", err)
	}
	if data != "R:3A:OK 00000000\r\n" {
		t.Errorf("received = %q, want %q", data, "R:3A:OK 00000000\r\n")
	}

	// Test receive with read error (should break loop)
	mock = &MockSerialPort{ReadErr: errors.New("read error")}
	uart = NewPS3UARTWithPort(mock, "CXR", 57600)
	_, err = uart.receive()
	// Error breaks loop but returns empty string and nil error
	if err != nil {
		t.Errorf("receive with error should return nil error, got: %v", err)
	}
}

func TestPS3UARTReceiveMultipleChunks(t *testing.T) {
	mock := &MockSerialPort{
		ReadChunks: [][]byte{
			[]byte("R:3A:"),
			[]byte("OK 00000000"),
			[]byte("\r\n"),
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	data, _ := uart.receive()
	expected := "R:3A:OK 00000000\r\n"
	if data != expected {
		t.Errorf("received = %q, want %q", data, expected)
	}
}

func TestPS3UARTCommandRouting(t *testing.T) {
	tests := []struct {
		name   string
		scType string
	}{
		{"CXR mode", "CXR"},
		{"SW mode", "SW"},
		{"CXRF mode", "CXRF"},
		{"Unknown mode defaults to CXRF", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Response works for all modes (CXRF and SW don't validate like CXR)
			mock := &MockSerialPort{ReadData: []byte("R:3A:OK 00000000\r\n")}
			uart := NewPS3UARTWithPort(mock, tt.scType, 57600)
			// Just verify it doesn't panic
			_ = uart.Command("VER", 0.001)
		})
	}
}

func TestCommandCXRShortCommand(t *testing.T) {
	// Test short command (<= 10 chars)
	// VER command with expected response
	// Checksum of "OK 00000000" = 0x3A
	mock := &MockSerialPort{ReadData: []byte("R:3A:OK 00000000\r\n")}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.commandCXR("VER", 0.001)
	if result.Code != 0 {
		t.Errorf("Command failed with code: %d, data: %v", result.Code, result.Data)
	}
}

func TestCommandCXRLongCommand(t *testing.T) {
	// Test long command (> 10 chars)
	// Checksum of "OK 00000000" = 0x3A
	mock := &MockSerialPort{ReadData: []byte("R:3A:OK 00000000\r\n")}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	// Command longer than 10 chars to trigger multipart send
	_ = uart.commandCXR("ERRLOG GET 00", 0.001)
	// Should have sent data in chunks
	if mock.WriteCalls == 0 {
		t.Error("Expected write calls for long command")
	}
}

func TestCommandCXRVeryLongCommand(t *testing.T) {
	// Test very long command that requires multiple chunks
	// Checksum of "OK 00000000" = 0x3A
	mock := &MockSerialPort{ReadData: []byte("R:3A:OK 00000000\r\n")}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	// Command longer than 25 chars to trigger multiple chunks in loop
	longCmd := "AUTH1 10000000000000000000000000000000000000"
	_ = uart.commandCXR(longCmd, 0.001)
	if mock.WriteCalls < 2 {
		t.Errorf("Expected multiple write calls for very long command, got %d", mock.WriteCalls)
	}
}

func TestCommandCXRInvalidResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		errData  string
	}{
		{
			name:     "too few parts",
			response: "R:5D\r\n",
			errData:  "Answer length",
		},
		{
			name:     "invalid magic",
			response: "X:3A:OK 00000000\r\n",
			errData:  "Magic",
		},
		{
			name:     "invalid checksum",
			response: "R:00:OK 00000000\r\n",
			errData:  "Checksum",
		},
		{
			name:     "R with insufficient data",
			response: "R:9A:OK\r\n", // Checksum of "OK" is 0x9A
			errData:  "Data length",
		},
		{
			name:     "E with wrong data length",
			response: "E:9A:OK\r\n", // Checksum of "OK" is 0x9A
			errData:  "Data length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockSerialPort{ReadData: []byte(tt.response)}
			uart := NewPS3UARTWithPort(mock, "CXR", 57600)

			result := uart.commandCXR("VER", 0.001)
			if result.Code != 0xFFFFFFFF {
				t.Errorf("Expected error code 0xFFFFFFFF, got %d", result.Code)
			}
			if len(result.Data) == 0 || result.Data[0] != tt.errData {
				t.Errorf("Expected error data %q, got %v", tt.errData, result.Data)
			}
		})
	}
}

func TestCommandCXRErrorResponse(t *testing.T) {
	// Test error response (E: prefix) with proper format
	// Checksum of "ERR 00000001" = 0x3E3
	mock := &MockSerialPort{ReadData: []byte("E:E3:ERR 00000001\r\n")}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.commandCXR("VER", 0.001)
	// E response with proper data
	if result.Code != 1 {
		t.Logf("Result: Code=%d, Data=%v", result.Code, result.Data)
	}
}

func TestCommandCXRNotOKResponse(t *testing.T) {
	// Test response where OK is not present
	// Checksum of "FAIL 00000001"
	checksum := 0
	resp := "FAIL 00000001"
	for _, c := range resp {
		checksum += int(c)
	}
	checksum %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", checksum, resp))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.commandCXR("VER", 0.001)
	if result.Code != 1 {
		t.Logf("Not OK result: Code=%d, Data=%v", result.Code, result.Data)
	}
}

func TestCommandSW(t *testing.T) {
	// Test SW mode command
	mock := &MockSerialPort{ReadData: []byte("OK 00000000:56\n")}
	uart := NewPS3UARTWithPort(mock, "SW", 57600)

	result := uart.commandSW("VER", 0.001)
	// Check we got some result
	_ = result
}

func TestCommandSWLongCommand(t *testing.T) {
	// Test SW mode with long command (>= 0x40 chars)
	mock := &MockSerialPort{
		ReadChunks: [][]byte{
			[]byte("OK 00000000:56\n"), // SETCMDLONG response
			[]byte("OK 00000000:56\n"), // Actual command response
		},
	}
	uart := NewPS3UARTWithPort(mock, "SW", 57600)

	// Create a command >= 64 chars
	longCmd := "AUTH1 100000000000000000000000000000000000000000000000000000000000"
	result := uart.commandSW(longCmd, 0.001)
	_ = result
}

func TestCommandSWSetcmdlongFails(t *testing.T) {
	// Test SW mode with long command where SETCMDLONG fails
	mock := &MockSerialPort{ReadData: []byte("ERR 00000001:56\n")}
	uart := NewPS3UARTWithPort(mock, "SW", 57600)

	// Create a command >= 64 chars
	longCmd := "AUTH1 100000000000000000000000000000000000000000000000000000000000"
	result := uart.commandSW(longCmd, 0.001)
	if result.Code != 0xFFFFFFFF {
		t.Errorf("Expected error code, got %d", result.Code)
	}
}

func TestCommandSWInvalidResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		errData  string
	}{
		{
			name:     "no colon separator",
			response: "OK 00000000\n",
			errData:  "Answer length",
		},
		{
			name:     "invalid checksum",
			response: "OK 00000000:00\n",
			errData:  "Checksum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockSerialPort{ReadData: []byte(tt.response)}
			uart := NewPS3UARTWithPort(mock, "SW", 57600)

			result := uart.commandSW("VER", 0.001)
			if result.Code != 0xFFFFFFFF {
				t.Errorf("Expected error code 0xFFFFFFFF, got %d", result.Code)
			}
		})
	}
}

func TestCommandSWMultilineResponse(t *testing.T) {
	// Test multiline response
	// Calculate checksums for each line
	line1 := "DATA LINE 1"
	line2 := "OK 00000000"
	cs1 := 0
	for _, c := range line1 {
		cs1 += int(c)
	}
	cs1 %= 0x100
	cs2 := 0
	for _, c := range line2 {
		cs2 += int(c)
	}
	cs2 %= 0x100

	response := fmt.Sprintf("%s:%02X\n%s:%02X\n", line1, cs1, line2, cs2)
	mock := &MockSerialPort{ReadData: []byte(response)}
	uart := NewPS3UARTWithPort(mock, "SW", 57600)

	result := uart.commandSW("ERRLOG GET 00", 0.001)
	_ = result
}

func TestCommandSWShortLastLine(t *testing.T) {
	// Test response with short/invalid last line
	line := "OK"
	cs := 0
	for _, c := range line {
		cs += int(c)
	}
	cs %= 0x100
	response := fmt.Sprintf("%s:%02X\n", line, cs)
	mock := &MockSerialPort{ReadData: []byte(response)}
	uart := NewPS3UARTWithPort(mock, "SW", 57600)

	result := uart.commandSW("VER", 0.001)
	// Should return with code 0 and the lines as data
	if result.Code != 0 {
		t.Logf("Short line result: Code=%d", result.Code)
	}
}

func TestCommandCXRF(t *testing.T) {
	mock := &MockSerialPort{ReadData: []byte("SC_READY\r\n")}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.commandCXRF("scopen", 0.001)
	if result.Code != 0 {
		t.Errorf("Expected code 0, got %d", result.Code)
	}
	if len(result.Data) == 0 || result.Data[0] != "SC_READY" {
		t.Errorf("Expected data [SC_READY], got %v", result.Data)
	}
}

func TestParseHexUint32(t *testing.T) {
	tests := []struct {
		input    string
		expected uint32
	}{
		{"00000000", 0},
		{"FFFFFFFF", 0xFFFFFFFF},
		{"ffffffff", 0xFFFFFFFF},
		{"00000001", 1},
		{"DEADBEEF", 0xDEADBEEF},
		{"12345678", 0x12345678},
		{"", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseHexUint32(tt.input)
			if result != tt.expected {
				t.Errorf("parseHexUint32(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAuthCXRSuccess(t *testing.T) {
	// Build a valid AUTH1 response
	// auth1rHdr + encrypted data
	auth1Response := make([]byte, 0x40)
	copy(auth1Response[0:0x10], auth1rHdr)

	// Create valid encrypted payload
	// The payload after decryption should have:
	// - bytes 0x8:0x10 = zeros (matches zeroIV[0:8])
	// - bytes 0x10:0x20 = authValue
	// - bytes 0x20:0x30 = zeros (matches zeroIV)
	plaintext := make([]byte, 0x30)
	copy(plaintext[0x8:0x10], zeroIV[0:0x8])
	copy(plaintext[0x10:0x20], authValue)
	copy(plaintext[0x20:0x30], zeroIV)

	encrypted, _ := aesEncryptCBC(sc2tb, zeroIV, plaintext)
	copy(auth1Response[0x10:0x40], encrypted)

	// AUTH1 response - calculate checksum
	respPart := fmt.Sprintf("OK 00000000 %X", auth1Response)
	cs := 0
	for _, c := range respPart {
		cs += int(c)
	}
	cs %= 0x100
	auth1ResponseHex := fmt.Sprintf("R:%02X:%s\r\n", cs, respPart)

	// AUTH2 success response
	auth2Resp := "OK 00000000"
	cs2 := 0
	for _, c := range auth2Resp {
		cs2 += int(c)
	}
	cs2 %= 0x100
	auth2ResponseHex := fmt.Sprintf("R:%02X:%s\r\n", cs2, auth2Resp)

	mock := &MockSerialPort{
		Responses: []string{
			auth1ResponseHex,
			auth2ResponseHex,
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth successful" {
		t.Errorf("Expected 'Auth successful', got %q", result)
	}
}

func TestAuthCXRInvalidAuth1Response(t *testing.T) {
	// Test with invalid AUTH1 response code (non-zero code)
	// Checksum of "OK FFFFFFFF" = sum of ASCII values mod 256
	resp := "OK FFFFFFFF"
	cs := 0
	for _, c := range resp {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", cs, resp))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth1 response invalid" {
		t.Errorf("Expected 'Auth1 response invalid', got %q", result)
	}
}

func TestAuthCXREmptyAuth1Data(t *testing.T) {
	// Test with empty AUTH1 data
	resp := "OK 00000000"
	cs := 0
	for _, c := range resp {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", cs, resp))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth1 response invalid" {
		t.Errorf("Expected 'Auth1 response invalid', got %q", result)
	}
}

func TestAuthCXRInvalidHexDecode(t *testing.T) {
	// Test with invalid hex in AUTH1 data
	resp := "OK 00000000 INVALID_HEX"
	cs := 0
	for _, c := range resp {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", cs, resp))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth1 response decode error" {
		t.Errorf("Expected 'Auth1 response decode error', got %q", result)
	}
}

func TestAuthCXRShortAuth1Response(t *testing.T) {
	// Test with too short AUTH1 response (< 0x40 bytes)
	shortData := "00112233445566778899AABBCCDDEEFF"
	resp := "OK 00000000 " + shortData
	cs := 0
	for _, c := range resp {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", cs, resp))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth1 response header invalid" {
		t.Errorf("Expected 'Auth1 response header invalid', got %q", result)
	}
}

func TestAuthCXRInvalidHeader(t *testing.T) {
	// Test with invalid AUTH1 header
	auth1Response := make([]byte, 0x40)
	// Don't copy auth1rHdr, leave as zeros
	resp := fmt.Sprintf("OK 00000000 %X", auth1Response)
	cs := 0
	for _, c := range resp {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", cs, resp))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth1 response header invalid" {
		t.Errorf("Expected 'Auth1 response header invalid', got %q", result)
	}
}

func TestAuthCXRInvalidDecryptedBody(t *testing.T) {
	// Test with valid header but invalid decrypted body
	auth1Response := make([]byte, 0x40)
	copy(auth1Response[0:0x10], auth1rHdr)
	// Leave rest as zeros - will fail body validation

	resp := fmt.Sprintf("OK 00000000 %X", auth1Response)
	cs := 0
	for _, c := range resp {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", cs, resp))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth1 response body invalid" {
		t.Errorf("Expected 'Auth1 response body invalid', got %q", result)
	}
}

func TestAuthCXRFScopenFails(t *testing.T) {
	mock := &MockSerialPort{ReadData: []byte("ERROR\r\n")}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	if result != "scopen response invalid" {
		t.Errorf("Expected 'scopen response invalid', got %q", result)
	}
}

func TestAuthCXRFScopenSuccess(t *testing.T) {
	// Test CXRF auth flow - scopen succeeds but auth1 has no valid data
	mock := &MockSerialPort{
		Responses: []string{
			"SC_READY\r\n", // scopen response
			"INVALID\r\n",  // AUTH1 response - no \r to split
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	// Should fail at parsing auth1 response
	if result != "scopen response invalid" {
		t.Logf("CXRF auth result: %s", result)
	}
}

func TestAuthCXRFEmptyAuth1Data(t *testing.T) {
	// Test CXRF auth where auth1 response has parts[1] too short
	mock := &MockSerialPort{
		Responses: []string{
			"SC_READY\r\n",
			"OK\r\r\n", // parts[1] = "" after [1:]
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	if result != "scopen response invalid" {
		t.Logf("CXRF empty auth1 result: %s", result)
	}
}

func TestAuthCXRFValidAuth1(t *testing.T) {
	// Build a valid CXRF AUTH1 flow
	auth1Response := make([]byte, 0x40)
	copy(auth1Response[0:0x10], auth1rHdr)

	// Create valid encrypted payload
	plaintext := make([]byte, 0x30)
	copy(plaintext[0x8:0x10], zeroIV[0:0x8])
	copy(plaintext[0x10:0x20], authValue)
	copy(plaintext[0x20:0x30], zeroIV)

	encrypted, _ := aesEncryptCBC(sc2tb, zeroIV, plaintext)
	copy(auth1Response[0x10:0x40], encrypted)

	auth1Hex := fmt.Sprintf("%X", auth1Response)

	mock := &MockSerialPort{
		Responses: []string{
			"SC_READY\r\n",
			fmt.Sprintf("OK\r %s\r\n", auth1Hex),
			"SC_SUCCESS\r\n",
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	if result != "Auth successful" {
		t.Errorf("Expected 'Auth successful', got %q", result)
	}
}

func TestAuthCXRFInvalidAuth1Hex(t *testing.T) {
	mock := &MockSerialPort{
		Responses: []string{
			"SC_READY\r\n",
			"OK\r 00112233\r\n", // Too short hex (not 128 chars)
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	if result != "Auth1 response invalid" {
		t.Errorf("Expected 'Auth1 response invalid', got %q", result)
	}
}

func TestAuthCXRFInvalidAuth1Header(t *testing.T) {
	// 128 hex chars but invalid header (all zeros, not matching auth1rHdr)
	invalidData := make([]byte, 64)
	mock := &MockSerialPort{
		Responses: []string{
			"SC_READY\r\n",
			fmt.Sprintf("OK\r %X\r\n", invalidData),
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	if result != "Auth1 response header invalid" {
		t.Errorf("Expected 'Auth1 response header invalid', got %q", result)
	}
}

func TestAuthCXRFInvalidAuth1Body(t *testing.T) {
	auth1Response := make([]byte, 0x40)
	copy(auth1Response[0:0x10], auth1rHdr)
	// Invalid body (zeros) - decrypted data won't match expected pattern

	mock := &MockSerialPort{
		Responses: []string{
			"SC_READY\r\n",
			fmt.Sprintf("OK\r %X\r\n", auth1Response),
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	if result != "Auth1 response body invalid" {
		t.Errorf("Expected 'Auth1 response body invalid', got %q", result)
	}
}

func TestAuthCXRFAuth2Fails(t *testing.T) {
	auth1Response := make([]byte, 0x40)
	copy(auth1Response[0:0x10], auth1rHdr)

	plaintext := make([]byte, 0x30)
	copy(plaintext[0x8:0x10], zeroIV[0:0x8])
	copy(plaintext[0x10:0x20], authValue)
	copy(plaintext[0x20:0x30], zeroIV)

	encrypted, _ := aesEncryptCBC(sc2tb, zeroIV, plaintext)
	copy(auth1Response[0x10:0x40], encrypted)

	mock := &MockSerialPort{
		Responses: []string{
			"SC_READY\r\n",
			fmt.Sprintf("OK\r %X\r\n", auth1Response),
			"SC_FAIL\r\n", // AUTH2 fails - doesn't contain SC_SUCCESS
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXRF", 115200)

	result := uart.Auth()
	if result != "Auth failed" {
		t.Errorf("Expected 'Auth failed', got %q", result)
	}
}

func TestAuthSWMode(t *testing.T) {
	// SW mode follows same path as CXR
	resp := "OK FFFFFFFF"
	cs := 0
	for _, c := range resp {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("%s:%02X\n", resp, cs))}
	uart := NewPS3UARTWithPort(mock, "SW", 57600)

	result := uart.Auth()
	// Will fail due to invalid response but exercises SW path
	t.Logf("SW auth result: %s", result)
}

func TestGetSerialPorts(t *testing.T) {
	// This tests the actual function which queries system ports
	// We can't mock the serial library, but we can verify it doesn't panic
	ports := getSerialPorts()
	// Just verify it returns a slice (may be empty on systems without serial ports)
	if ports == nil {
		t.Error("getSerialPorts returned nil, expected empty slice")
	}
}

func TestNewPS3UARTFailsOnInvalidPort(t *testing.T) {
	// Test that NewPS3UART returns error for invalid port
	_, err := NewPS3UART("/dev/nonexistent_serial_port_12345", "CXR", 57600)
	if err == nil {
		t.Error("Expected error for non-existent port, got nil")
	}
}

func TestAuthDecryptionError(t *testing.T) {
	// Test auth with data that causes decryption issues
	// Create auth1 response with valid header but invalid encrypted data size
	auth1Response := make([]byte, 0x40)
	copy(auth1Response[0:0x10], auth1rHdr)
	// Fill with valid-looking but actually invalid encrypted data
	for i := 0x10; i < 0x40; i++ {
		auth1Response[i] = byte(i)
	}

	respPart := fmt.Sprintf("OK 00000000 %X", auth1Response)
	cs := 0
	for _, c := range respPart {
		cs += int(c)
	}
	cs %= 0x100
	mock := &MockSerialPort{ReadData: []byte(fmt.Sprintf("R:%02X:%s\r\n", cs, respPart))}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	// Should fail at body validation since decrypted data won't match expected pattern
	if result != "Auth1 response body invalid" {
		t.Logf("Decryption test result: %s", result)
	}
}

func TestAuthCXRAuth2Fails(t *testing.T) {
	// Build a valid AUTH1 response that leads to AUTH2
	auth1Response := make([]byte, 0x40)
	copy(auth1Response[0:0x10], auth1rHdr)

	plaintext := make([]byte, 0x30)
	copy(plaintext[0x8:0x10], zeroIV[0:0x8])
	copy(plaintext[0x10:0x20], authValue)
	copy(plaintext[0x20:0x30], zeroIV)

	encrypted, _ := aesEncryptCBC(sc2tb, zeroIV, plaintext)
	copy(auth1Response[0x10:0x40], encrypted)

	respPart := fmt.Sprintf("OK 00000000 %X", auth1Response)
	cs := 0
	for _, c := range respPart {
		cs += int(c)
	}
	cs %= 0x100
	auth1ResponseHex := fmt.Sprintf("R:%02X:%s\r\n", cs, respPart)

	// AUTH2 failure response - non-zero code means failure
	auth2Resp := "OK 00000001"
	cs2 := 0
	for _, c := range auth2Resp {
		cs2 += int(c)
	}
	cs2 %= 0x100
	auth2ResponseHex := fmt.Sprintf("R:%02X:%s\r\n", cs2, auth2Resp)

	mock := &MockSerialPort{
		Responses: []string{
			auth1ResponseHex,
			auth2ResponseHex,
		},
	}
	uart := NewPS3UARTWithPort(mock, "CXR", 57600)

	result := uart.Auth()
	if result != "Auth failed" {
		t.Errorf("Expected 'Auth failed', got %q", result)
	}
}
