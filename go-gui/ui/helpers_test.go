package ui

import (
	"testing"
)

func TestFilterOptions(t *testing.T) {
	options := []string{"AUTH1", "AUTH2", "VER", "EEP", "ERRLOG", "FAN", "HALT"}

	tests := []struct {
		name     string
		search   string
		expected []string
	}{
		{
			name:     "empty search returns all",
			search:   "",
			expected: options,
		},
		{
			name:     "exact match",
			search:   "VER",
			expected: []string{"VER"},
		},
		{
			name:     "lowercase search",
			search:   "ver",
			expected: []string{"VER"},
		},
		{
			name:     "partial match",
			search:   "AUTH",
			expected: []string{"AUTH1", "AUTH2"},
		},
		{
			name:     "no match",
			search:   "NOTFOUND",
			expected: nil,
		},
		{
			name:     "single char match",
			search:   "E",
			expected: []string{"VER", "EEP", "ERRLOG"},
		},
		{
			name:     "mixed case search",
			search:   "ErR",
			expected: []string{"ERRLOG"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterOptions(options, tt.search)
			if len(result) != len(tt.expected) {
				t.Errorf("FilterOptions(%v, %q) returned %d items, want %d",
					options, tt.search, len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("FilterOptions result[%d] = %q, want %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestFilterOptionsEmptyInput(t *testing.T) {
	var emptyOptions []string
	result := FilterOptions(emptyOptions, "test")
	if len(result) != 0 {
		t.Errorf("FilterOptions with empty input returned %d items, want 0", len(result))
	}
}

func TestFilterOptionsNilInput(t *testing.T) {
	result := FilterOptions(nil, "test")
	if result != nil {
		t.Errorf("FilterOptions with nil input returned %v, want nil", result)
	}
}

func TestBuildCXRCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		subCmd   string
		args     string
		expected string
	}{
		{
			name:     "command only",
			cmd:      "VER",
			subCmd:   "",
			args:     "",
			expected: "VER",
		},
		{
			name:     "command with subcommand",
			cmd:      "EEP",
			subCmd:   "GET",
			args:     "",
			expected: "EEP GET",
		},
		{
			name:     "command with subcommand and args",
			cmd:      "EEP",
			subCmd:   "GET",
			args:     "00",
			expected: "EEP GET 00",
		},
		{
			name:     "command with args only",
			cmd:      "AUTH1",
			subCmd:   "",
			args:     "10000000",
			expected: "AUTH1 10000000",
		},
		{
			name:     "empty command",
			cmd:      "",
			subCmd:   "GET",
			args:     "00",
			expected: "",
		},
		{
			name:     "whitespace command",
			cmd:      "   ",
			subCmd:   "GET",
			args:     "00",
			expected: "",
		},
		{
			name:     "command with whitespace args",
			cmd:      "EEP",
			subCmd:   "SET",
			args:     "  00 FF  ",
			expected: "EEP SET 00 FF",
		},
		{
			name:     "command with leading whitespace",
			cmd:      "  VER",
			subCmd:   "",
			args:     "",
			expected: "VER",
		},
		{
			name:     "command with trailing whitespace",
			cmd:      "VER  ",
			subCmd:   "",
			args:     "",
			expected: "VER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildCXRCommand(tt.cmd, tt.subCmd, tt.args)
			if result != tt.expected {
				t.Errorf("BuildCXRCommand(%q, %q, %q) = %q, want %q",
					tt.cmd, tt.subCmd, tt.args, result, tt.expected)
			}
		})
	}
}

func TestGetSerialSpeed(t *testing.T) {
	tests := []struct {
		scType   string
		expected int
	}{
		{"CXR", 57600},
		{"SW", 57600},
		{"CXRF", 115200},
		{"UNKNOWN", 57600},
		{"", 57600},
		{"cxrf", 57600}, // case sensitive, lowercase should return default
	}

	for _, tt := range tests {
		t.Run(tt.scType, func(t *testing.T) {
			result := GetSerialSpeed(tt.scType)
			if result != tt.expected {
				t.Errorf("GetSerialSpeed(%q) = %d, want %d", tt.scType, result, tt.expected)
			}
		})
	}
}

func TestFormatCommandOutput(t *testing.T) {
	tests := []struct {
		name     string
		scType   string
		result   CommandResult
		expected string
	}{
		{
			name:     "CXR mode",
			scType:   "CXR",
			result:   CommandResult{Code: 0, Data: []string{"DATA1", "DATA2"}},
			expected: "00000000 DATA1 DATA2",
		},
		{
			name:     "CXR mode with code",
			scType:   "CXR",
			result:   CommandResult{Code: 0x12345678, Data: []string{"VALUE"}},
			expected: "12345678 VALUE",
		},
		{
			name:     "CXR mode empty data",
			scType:   "CXR",
			result:   CommandResult{Code: 0, Data: []string{}},
			expected: "00000000 ",
		},
		{
			name:     "SW mode single line",
			scType:   "SW",
			result:   CommandResult{Code: 0, Data: []string{"DATA"}},
			expected: "00000000 DATA",
		},
		{
			name:     "SW mode multiline",
			scType:   "SW",
			result:   CommandResult{Code: 0, Data: []string{"LINE1\nLINE2"}},
			expected: "00000000\nLINE1\nLINE2",
		},
		{
			name:     "SW mode empty data",
			scType:   "SW",
			result:   CommandResult{Code: 0, Data: []string{}},
			expected: "00000000\n",
		},
		{
			name:     "CXRF mode with data",
			scType:   "CXRF",
			result:   CommandResult{Code: 0, Data: []string{"SC_READY"}},
			expected: "SC_READY",
		},
		{
			name:     "CXRF mode empty data",
			scType:   "CXRF",
			result:   CommandResult{Code: 0, Data: []string{}},
			expected: "",
		},
		{
			name:     "Unknown mode with data",
			scType:   "UNKNOWN",
			result:   CommandResult{Code: 0, Data: []string{"RESPONSE"}},
			expected: "RESPONSE",
		},
		{
			name:     "Unknown mode empty data",
			scType:   "UNKNOWN",
			result:   CommandResult{Code: 0, Data: []string{}},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCommandOutput(tt.scType, tt.result)
			if result != tt.expected {
				t.Errorf("FormatCommandOutput(%q, %v) = %q, want %q",
					tt.scType, tt.result, result, tt.expected)
			}
		})
	}
}

func TestCommandResultStruct(t *testing.T) {
	// Test that CommandResult can be created and accessed
	result := CommandResult{
		Code: 0x12345678,
		Data: []string{"data1", "data2"},
	}

	if result.Code != 0x12345678 {
		t.Errorf("CommandResult.Code = %x, want 0x12345678", result.Code)
	}

	if len(result.Data) != 2 {
		t.Errorf("CommandResult.Data length = %d, want 2", len(result.Data))
	}

	if result.Data[0] != "data1" || result.Data[1] != "data2" {
		t.Errorf("CommandResult.Data = %v, want [data1 data2]", result.Data)
	}
}
