// Package ui provides utility functions for UI operations.
package ui

import (
	"fmt"
	"strings"
)

// FilterOptions filters a list of options based on a search string.
func FilterOptions(options []string, search string) []string {
	if search == "" {
		return options
	}
	search = strings.ToUpper(search)
	var filtered []string
	for _, opt := range options {
		if strings.Contains(strings.ToUpper(opt), search) {
			filtered = append(filtered, opt)
		}
	}
	return filtered
}

// BuildCXRCommand builds a command string from parts for CXR mode.
func BuildCXRCommand(cmd, subCmd, args string) string {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}
	parts := []string{cmd}
	if subCmd != "" {
		parts = append(parts, subCmd)
	}
	args = strings.TrimSpace(args)
	if args != "" {
		parts = append(parts, args)
	}
	return strings.Join(parts, " ")
}

// GetSerialSpeed returns the appropriate baud rate for the SC type.
func GetSerialSpeed(scType string) int {
	if scType == "CXRF" {
		return 115200
	}
	return 57600
}

// CommandResult holds the result of a command execution.
// This is duplicated here to avoid circular imports with serial package.
type CommandResult struct {
	Code uint32
	Data []string
}

// FormatCommandOutput formats the command result for display.
func FormatCommandOutput(scType string, result CommandResult) string {
	switch scType {
	case "CXR":
		return fmt.Sprintf("%08X %s", result.Code, strings.Join(result.Data, " "))
	case "SW":
		if len(result.Data) > 0 && !strings.Contains(result.Data[0], "\n") {
			return fmt.Sprintf("%08X %s", result.Code, strings.Join(result.Data, " "))
		}
		return fmt.Sprintf("%08X\n%s", result.Code, strings.Join(result.Data, ""))
	default:
		if len(result.Data) > 0 {
			return result.Data[0]
		}
		return ""
	}
}
