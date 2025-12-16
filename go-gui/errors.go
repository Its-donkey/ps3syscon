// Package errors provides sentinel errors for the PS3 Syscon UART tool.
package main

import "errors"

// Sentinel errors for common error conditions.
var (
	// ErrPortNotSelected indicates no serial port was selected.
	ErrPortNotSelected = errors.New("serial port not selected")

	// ErrModeNotSelected indicates no SC mode was selected.
	ErrModeNotSelected = errors.New("mode not selected")

	// ErrCommandEmpty indicates an empty command was provided.
	ErrCommandEmpty = errors.New("command is empty")

	// ErrCommandFailed indicates a command execution failed.
	ErrCommandFailed = errors.New("command failed")

	// ErrAuthFailed indicates authentication failed.
	ErrAuthFailed = errors.New("authentication failed")

	// ErrInvalidResponse indicates an invalid response from the device.
	ErrInvalidResponse = errors.New("invalid response")

	// ErrChecksumMismatch indicates a checksum verification failed.
	ErrChecksumMismatch = errors.New("checksum mismatch")

	// ErrDecryptionFailed indicates AES decryption failed.
	ErrDecryptionFailed = errors.New("decryption failed")

	// ErrEncryptionFailed indicates AES encryption failed.
	ErrEncryptionFailed = errors.New("encryption failed")

	// ErrInvalidAuthResponse indicates an invalid authentication response.
	ErrInvalidAuthResponse = errors.New("invalid auth response")

	// ErrSerialOpenFailed indicates failure to open serial port.
	ErrSerialOpenFailed = errors.New("failed to open serial port")
)
