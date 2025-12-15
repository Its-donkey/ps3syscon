package main

import (
	"testing"
)

func TestGetCommandNames(t *testing.T) {
	names := GetCommandNames()

	// Check that we get the right number of commands
	if len(names) != len(MullionCommands) {
		t.Errorf("GetCommandNames() returned %d names, want %d", len(names), len(MullionCommands))
	}

	// Check that specific known commands are in the list
	expectedNames := []string{"AUTH1", "AUTH2", "VER", "EEP", "ERRLOG", "FAN", "HALT", "SHUTDOWN"}
	for _, expected := range expectedNames {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %q not found in GetCommandNames() result", expected)
		}
	}
}

func TestGetCommand(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedFound bool
		expectedName  string
	}{
		{
			name:          "exact match uppercase",
			input:         "AUTH1",
			expectedFound: true,
			expectedName:  "AUTH1",
		},
		{
			name:          "lowercase input",
			input:         "auth1",
			expectedFound: true,
			expectedName:  "AUTH1",
		},
		{
			name:          "mixed case input",
			input:         "Auth1",
			expectedFound: true,
			expectedName:  "AUTH1",
		},
		{
			name:          "with leading/trailing spaces",
			input:         "  VER  ",
			expectedFound: true,
			expectedName:  "VER",
		},
		{
			name:          "non-existent command",
			input:         "NOTACOMMAND",
			expectedFound: false,
			expectedName:  "",
		},
		{
			name:          "empty string",
			input:         "",
			expectedFound: false,
			expectedName:  "",
		},
		{
			name:          "just whitespace",
			input:         "   ",
			expectedFound: false,
			expectedName:  "",
		},
		{
			name:          "EEP command",
			input:         "EEP",
			expectedFound: true,
			expectedName:  "EEP",
		},
		{
			name:          "ERRLOG command",
			input:         "errlog",
			expectedFound: true,
			expectedName:  "ERRLOG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCommand(tt.input)
			if tt.expectedFound {
				if result == nil {
					t.Errorf("GetCommand(%q) = nil, want command %q", tt.input, tt.expectedName)
				} else if result.Name != tt.expectedName {
					t.Errorf("GetCommand(%q).Name = %q, want %q", tt.input, result.Name, tt.expectedName)
				}
			} else {
				if result != nil {
					t.Errorf("GetCommand(%q) = %v, want nil", tt.input, result)
				}
			}
		})
	}
}

func TestGetCommandReturnsCorrectSubcommands(t *testing.T) {
	tests := []struct {
		name                string
		cmdName             string
		expectedSubcommands []string
	}{
		{
			name:                "AUTH1 has no subcommands",
			cmdName:             "AUTH1",
			expectedSubcommands: nil,
		},
		{
			name:                "EEP has subcommands",
			cmdName:             "EEP",
			expectedSubcommands: []string{"GET", "SET", "INIT"},
		},
		{
			name:                "ERRLOG has subcommands",
			cmdName:             "ERRLOG",
			expectedSubcommands: []string{"GET", "CLEAR", "START", "STOP"},
		},
		{
			name:                "FAN has many subcommands",
			cmdName:             "FAN",
			expectedSubcommands: []string{"GETDUTY", "GETPOLICY", "SETDUTY", "SETPOLICY", "START", "STOP"},
		},
		{
			name:                "VER has no subcommands",
			cmdName:             "VER",
			expectedSubcommands: nil,
		},
		{
			name:                "AUTHVER has subcommands",
			cmdName:             "AUTHVER",
			expectedSubcommands: []string{"GET", "SET"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetCommand(tt.cmdName)
			if cmd == nil {
				t.Fatalf("GetCommand(%q) returned nil", tt.cmdName)
			}

			if tt.expectedSubcommands == nil {
				if cmd.Subcommands != nil {
					t.Errorf("GetCommand(%q).Subcommands = %v, want nil", tt.cmdName, cmd.Subcommands)
				}
			} else {
				if len(cmd.Subcommands) != len(tt.expectedSubcommands) {
					t.Errorf("GetCommand(%q).Subcommands length = %d, want %d",
						tt.cmdName, len(cmd.Subcommands), len(tt.expectedSubcommands))
				}
				for i, expected := range tt.expectedSubcommands {
					if i >= len(cmd.Subcommands) || cmd.Subcommands[i] != expected {
						t.Errorf("GetCommand(%q).Subcommands[%d] = %q, want %q",
							tt.cmdName, i, cmd.Subcommands[i], expected)
					}
				}
			}
		})
	}
}

func TestHasSubcommands(t *testing.T) {
	tests := []struct {
		name     string
		cmdName  string
		expected bool
	}{
		{
			name:     "AUTH1 has no subcommands",
			cmdName:  "AUTH1",
			expected: false,
		},
		{
			name:     "AUTH2 has no subcommands",
			cmdName:  "AUTH2",
			expected: false,
		},
		{
			name:     "EEP has subcommands",
			cmdName:  "EEP",
			expected: true,
		},
		{
			name:     "ERRLOG has subcommands",
			cmdName:  "ERRLOG",
			expected: true,
		},
		{
			name:     "FAN has subcommands",
			cmdName:  "FAN",
			expected: true,
		},
		{
			name:     "VER has no subcommands",
			cmdName:  "VER",
			expected: false,
		},
		{
			name:     "HALT has no subcommands",
			cmdName:  "HALT",
			expected: false,
		},
		{
			name:     "BOOT has subcommands",
			cmdName:  "BOOT",
			expected: true,
		},
		{
			name:     "CSAREA has subcommands",
			cmdName:  "CSAREA",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetCommand(tt.cmdName)
			if cmd == nil {
				t.Fatalf("GetCommand(%q) returned nil", tt.cmdName)
			}

			result := cmd.HasSubcommands()
			if result != tt.expected {
				t.Errorf("Command %q HasSubcommands() = %v, want %v", tt.cmdName, result, tt.expected)
			}
		})
	}
}

func TestMullionCommandsIntegrity(t *testing.T) {
	// Test that all commands have valid structure
	for _, cmd := range MullionCommands {
		if cmd.Name == "" {
			t.Error("Found command with empty name")
		}

		// Check that permission is not zero (all commands should have some permission)
		if cmd.Permission == 0 {
			t.Errorf("Command %q has zero permission", cmd.Name)
		}
	}
}

func TestMullionCommandsAllNames(t *testing.T) {
	// Verify all 33 commands are present
	expectedCommands := []string{
		"AUTH1", "AUTH2", "AUTHVER", "BOOT", "BOOTENABLE", "BUZ", "CID",
		"CSAREA", "ECID", "EEP", "ERRLOG", "FAN", "HALT", "KSV", "PDAREA",
		"PORTSTAT", "R8", "R16", "R32", "RBE", "REV", "SERVFAN", "SHUTDOWN",
		"SPU", "VER", "VID", "W8", "W16", "W32", "WBE",
	}

	for _, expected := range expectedCommands {
		cmd := GetCommand(expected)
		if cmd == nil {
			t.Errorf("Expected command %q not found in MullionCommands", expected)
		}
	}
}

func TestCommandPermissions(t *testing.T) {
	// Test specific known permissions
	tests := []struct {
		name               string
		cmdName            string
		expectedPermission uint32
	}{
		{
			name:               "AUTH1 permission",
			cmdName:            "AUTH1",
			expectedPermission: 0x0000C0EF,
		},
		{
			name:               "AUTH2 permission",
			cmdName:            "AUTH2",
			expectedPermission: 0x0000C0EF,
		},
		{
			name:               "VER permission",
			cmdName:            "VER",
			expectedPermission: 0x0000C0FF,
		},
		{
			name:               "HALT permission",
			cmdName:            "HALT",
			expectedPermission: 0x0000C0D5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetCommand(tt.cmdName)
			if cmd == nil {
				t.Fatalf("GetCommand(%q) returned nil", tt.cmdName)
			}

			if cmd.Permission != tt.expectedPermission {
				t.Errorf("Command %q Permission = 0x%08X, want 0x%08X",
					tt.cmdName, cmd.Permission, tt.expectedPermission)
			}
		})
	}
}

func TestGetCommandNamesLength(t *testing.T) {
	names := GetCommandNames()
	if len(names) < 30 {
		t.Errorf("Expected at least 30 commands, got %d", len(names))
	}
}

func TestGetCommandNamesNoDuplicates(t *testing.T) {
	names := GetCommandNames()
	seen := make(map[string]bool)
	for _, name := range names {
		if seen[name] {
			t.Errorf("Duplicate command name found: %q", name)
		}
		seen[name] = true
	}
}
