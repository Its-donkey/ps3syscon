package main

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrPortNotSelected", ErrPortNotSelected, "serial port not selected"},
		{"ErrModeNotSelected", ErrModeNotSelected, "mode not selected"},
		{"ErrCommandEmpty", ErrCommandEmpty, "command is empty"},
		{"ErrCommandFailed", ErrCommandFailed, "command failed"},
		{"ErrAuthFailed", ErrAuthFailed, "authentication failed"},
		{"ErrInvalidResponse", ErrInvalidResponse, "invalid response"},
		{"ErrChecksumMismatch", ErrChecksumMismatch, "checksum mismatch"},
		{"ErrDecryptionFailed", ErrDecryptionFailed, "decryption failed"},
		{"ErrEncryptionFailed", ErrEncryptionFailed, "encryption failed"},
		{"ErrInvalidAuthResponse", ErrInvalidAuthResponse, "invalid auth response"},
		{"ErrSerialOpenFailed", ErrSerialOpenFailed, "failed to open serial port"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s is nil", tt.name)
				return
			}
			if tt.err.Error() != tt.msg {
				t.Errorf("%s.Error() = %q, want %q", tt.name, tt.err.Error(), tt.msg)
			}
		})
	}
}

func TestErrorsAreDistinct(t *testing.T) {
	allErrors := []error{
		ErrPortNotSelected,
		ErrModeNotSelected,
		ErrCommandEmpty,
		ErrCommandFailed,
		ErrAuthFailed,
		ErrInvalidResponse,
		ErrChecksumMismatch,
		ErrDecryptionFailed,
		ErrEncryptionFailed,
		ErrInvalidAuthResponse,
		ErrSerialOpenFailed,
	}

	for i, err1 := range allErrors {
		for j, err2 := range allErrors {
			if i != j && errors.Is(err1, err2) {
				t.Errorf("Error %d and %d should be distinct but are equal", i, j)
			}
		}
	}
}

func TestErrorsCanBeWrapped(t *testing.T) {
	wrapped := errors.New("wrapper: " + ErrPortNotSelected.Error())
	if wrapped.Error() != "wrapper: serial port not selected" {
		t.Errorf("Wrapped error message incorrect: %s", wrapped.Error())
	}
}
