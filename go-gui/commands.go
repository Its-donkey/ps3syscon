package main

import "strings"

// Command represents a Mullion external command.
type Command struct {
	Name        string
	Subcommands []string
	Description string
	Permission  uint32
}

// MullionCommands contains all known Mullion (CXR) external commands.
var MullionCommands = []Command{
	{Name: "AUTH1", Subcommands: nil, Permission: 0x0000C0EF},
	{Name: "AUTH2", Subcommands: nil, Permission: 0x0000C0EF},
	{Name: "AUTHVER", Subcommands: []string{"GET", "SET"}, Permission: 0x0000C0DF},
	{Name: "BOOT", Subcommands: []string{"MODE", "CONT"}, Permission: 0x000080D5},
	{Name: "BOOTENABLE", Subcommands: nil, Permission: 0x0000809A},
	{Name: "BUZ", Subcommands: nil, Permission: 0x00008096},
	{Name: "CID", Subcommands: []string{"GET"}, Permission: 0x0000C0D5},
	{Name: "CSAREA", Subcommands: []string{"GET", "SET"}, Permission: 0x0000C0DF},
	{Name: "ECID", Subcommands: []string{"GET"}, Permission: 0x0000C0D5},
	{Name: "EEP", Subcommands: []string{"GET", "SET", "INIT"}, Permission: 0x0000C0DF},
	{Name: "ERRLOG", Subcommands: []string{"GET", "CLEAR", "START", "STOP"}, Permission: 0x0000C0DF},
	{Name: "FAN", Subcommands: []string{"GETDUTY", "GETPOLICY", "SETDUTY", "SETPOLICY", "START", "STOP"}, Permission: 0x0000C0D7},
	{Name: "HALT", Subcommands: nil, Permission: 0x0000C0D5},
	{Name: "KSV", Subcommands: nil, Permission: 0x0000C0D5},
	{Name: "PDAREA", Subcommands: []string{"GET", "SET"}, Permission: 0x0000C0DF},
	{Name: "PORTSTAT", Subcommands: nil, Permission: 0x0000C0DF},
	{Name: "R8", Subcommands: nil, Permission: 0x0000C0DF},
	{Name: "R16", Subcommands: nil, Permission: 0x0000C0DF},
	{Name: "R32", Subcommands: nil, Permission: 0x0000C0DF},
	{Name: "RBE", Subcommands: nil, Permission: 0x0000C0D5},
	{Name: "REV", Subcommands: []string{"SB"}, Permission: 0x0000C0D5},
	{Name: "SERVFAN", Subcommands: nil, Permission: 0x0000C0D7},
	{Name: "SHUTDOWN", Subcommands: nil, Permission: 0x0000C0D5},
	{Name: "SPU", Subcommands: []string{"INFO"}, Permission: 0x0000C0D5},
	{Name: "VER", Subcommands: nil, Permission: 0x0000C0FF},
	{Name: "VID", Subcommands: []string{"GET"}, Permission: 0x0000C0D5},
	{Name: "W8", Subcommands: nil, Permission: 0x0000C0DF},
	{Name: "W16", Subcommands: nil, Permission: 0x0000C0DF},
	{Name: "W32", Subcommands: nil, Permission: 0x0000C0DF},
	{Name: "WBE", Subcommands: nil, Permission: 0x0000C0D5},
}

// CXRFCommands contains all known CXRF internal commands.
var CXRFCommands = []Command{
	{Name: "becount", Subcommands: nil, Description: "Display bringup/shutdown count + Power-on time"},
	{Name: "bepgoff", Subcommands: nil, Description: "BE power grid off"},
	{Name: "bepkt", Subcommands: []string{"show", "set", "unset", "mode", "debug", "help"}, Description: "Packet permissions"},
	{Name: "bestat", Subcommands: nil, Description: "Get status of BE"},
	{Name: "boardconfig", Subcommands: nil, Description: "Displays board configuration"},
	{Name: "bootbeep", Subcommands: []string{"stat", "on", "off"}, Description: "Boot beep"},
	{Name: "bringup", Subcommands: nil, Description: "Turn PS3 on"},
	{Name: "bsn", Subcommands: nil, Description: "Get board serial number"},
	{Name: "bstatus", Subcommands: nil, Description: "HDMI related status"},
	{Name: "buzz", Subcommands: nil, Description: "Activate buzzer [freq]"},
	{Name: "buzzpattern", Subcommands: nil, Description: "Buzzer pattern [freq] [pattern] [count]"},
	{Name: "clear_err", Subcommands: []string{"last", "eeprom", "all"}, Description: "Clear errors"},
	{Name: "clearerrlog", Subcommands: nil, Description: "Clears error log"},
	{Name: "comm", Subcommands: nil, Description: "Communication mode"},
	{Name: "commt", Subcommands: []string{"help", "start", "stop", "send"}, Description: "Manual BE communication"},
	{Name: "cp", Subcommands: []string{"ready", "busy", "reset", "beepremote", "beep2kn1n3", "beep2kn2n3"}, Description: "CP control commands"},
	{Name: "csum", Subcommands: nil, Description: "Firmware checksum"},
	{Name: "devpm", Subcommands: []string{"ata", "pci", "pciex", "rsx"}, Description: "Device power management"},
	{Name: "diag", Subcommands: nil, Description: "Diag (execute without param to show help)"},
	{Name: "disp_err", Subcommands: nil, Description: "Displays errors"},
	{Name: "duty", Subcommands: []string{"get", "set", "getmin", "setmin", "getmax", "setmax", "getinmin", "setinmin", "getinmax", "setinmax"}, Description: "Fan policy"},
	{Name: "dve", Subcommands: []string{"help", "set", "save", "show"}, Description: "DVE chip parameters"},
	{Name: "eepcsum", Subcommands: nil, Description: "Shows eeprom checksum"},
	{Name: "eepromcheck", Subcommands: nil, Description: "Check eeprom [id]"},
	{Name: "eeprominit", Subcommands: nil, Description: "Init eeprom [id]"},
	{Name: "ejectsw", Subcommands: nil, Description: "Eject switch"},
	{Name: "errlog", Subcommands: nil, Description: "Gets the error log"},
	{Name: "fancon", Subcommands: nil, Description: "Does nothing"},
	{Name: "fanconautotype", Subcommands: nil, Description: "Does nothing"},
	{Name: "fanconmode", Subcommands: []string{"get"}, Description: "Fan control mode"},
	{Name: "fanconpolicy", Subcommands: []string{"get", "set", "getini", "setini"}, Description: "Fan control policy"},
	{Name: "fandiag", Subcommands: nil, Description: "Fan test"},
	{Name: "faninictrl", Subcommands: nil, Description: "Does nothing"},
	{Name: "fanpol", Subcommands: nil, Description: "Does nothing"},
	{Name: "fanservo", Subcommands: nil, Description: "Does nothing"},
	{Name: "fantbl", Subcommands: []string{"get", "set", "getini", "setini", "gettable", "settable"}, Description: "Fan table"},
	{Name: "firmud", Subcommands: nil, Description: "Firmware update"},
	{Name: "geterrlog", Subcommands: nil, Description: "Gets error log [id]"},
	{Name: "getrtc", Subcommands: nil, Description: "Gets rtc"},
	{Name: "halt", Subcommands: nil, Description: "Halts syscon"},
	{Name: "hdmi", Subcommands: nil, Description: "HDMI (various commands, use help)"},
	{Name: "hdmiid", Subcommands: nil, Description: "Get HDMI id's"},
	{Name: "hdmiid2", Subcommands: nil, Description: "Get HDMI id's"},
	{Name: "hversion", Subcommands: nil, Description: "Platform ID"},
	{Name: "hyst", Subcommands: []string{"get", "set", "getini", "setini"}, Description: "Temperature zones"},
	{Name: "lasterrlog", Subcommands: nil, Description: "Last error from log"},
	{Name: "ledmode", Subcommands: nil, Description: "Get led mode [id] [id]"},
	{Name: "LS", Subcommands: nil, Description: "LabStation Mode"},
	{Name: "ltstest", Subcommands: []string{"get", "set be", "rsx"}, Description: "Temp related values"},
	{Name: "osbo", Subcommands: nil, Description: "Sets 0x2000F60"},
	{Name: "patchcsum", Subcommands: nil, Description: "Patch checksum"},
	{Name: "patchvereep", Subcommands: nil, Description: "Patch version eeprom"},
	{Name: "patchverram", Subcommands: nil, Description: "Patch version ram"},
	{Name: "poll", Subcommands: nil, Description: "Poll log"},
	{Name: "portscan", Subcommands: nil, Description: "Scan port [port]"},
	{Name: "powbtnmode", Subcommands: nil, Description: "Power button mode [mode (0/1)]"},
	{Name: "powerstate", Subcommands: nil, Description: "Get power state"},
	{Name: "powersw", Subcommands: nil, Description: "Power switch"},
	{Name: "powupcause", Subcommands: nil, Description: "Power up cause"},
	{Name: "printmode", Subcommands: nil, Description: "Set printmode [mode (0/1/2/3)]"},
	{Name: "printpatch", Subcommands: nil, Description: "Prints patch"},
	{Name: "r", Subcommands: nil, Description: "Read byte from SC [offset] [length]"},
	{Name: "r16", Subcommands: nil, Description: "Read word from SC [offset] [length]"},
	{Name: "r32", Subcommands: nil, Description: "Read dword from SC [offset] [length]"},
	{Name: "r64", Subcommands: nil, Description: "Read qword from SC [offset] [length]"},
	{Name: "r64d", Subcommands: nil, Description: "Read qword data from SC [offset] [length]"},
	{Name: "rbe", Subcommands: nil, Description: "Read from BE [offset]"},
	{Name: "recv", Subcommands: nil, Description: "Receive something"},
	{Name: "resetsw", Subcommands: nil, Description: "Reset switch"},
	{Name: "restartlogerrtoeep", Subcommands: nil, Description: "Reenable error logging to eeprom"},
	{Name: "revision", Subcommands: nil, Description: "Get softid"},
	{Name: "rrsxc", Subcommands: nil, Description: "Read from RSX [offset] [length]"},
	{Name: "rtcreset", Subcommands: nil, Description: "Reset RTC"},
	{Name: "scagv2", Subcommands: nil, Description: "Auth related"},
	{Name: "scasv2", Subcommands: nil, Description: "Auth related"},
	{Name: "scclose", Subcommands: nil, Description: "Auth related"},
	{Name: "scopen", Subcommands: nil, Description: "Auth related"},
	{Name: "send", Subcommands: nil, Description: "Send something [variable]"},
	{Name: "shutdown", Subcommands: nil, Description: "PS3 shutdown"},
	{Name: "startlogerrtsk", Subcommands: nil, Description: "Start error log task"},
	{Name: "stoplogerrtoeep", Subcommands: nil, Description: "Stop error logging to eeprom"},
	{Name: "stoplogerrtsk", Subcommands: nil, Description: "Stop error log task"},
	{Name: "syspowdown", Subcommands: nil, Description: "System power down (3 params 0 0 0)"},
	{Name: "task", Subcommands: nil, Description: "Print tasks"},
	{Name: "thalttest", Subcommands: nil, Description: "Does nothing"},
	{Name: "thermfatalmode", Subcommands: []string{"canboot", "cannotboot"}, Description: "Set thermal boot mode"},
	{Name: "therrclr", Subcommands: nil, Description: "Thermal register clear"},
	{Name: "thrm", Subcommands: nil, Description: "Does nothing"},
	{Name: "tmp", Subcommands: nil, Description: "Get temperature [zone]"},
	{Name: "trace", Subcommands: nil, Description: "Trace tasks (use help)"},
	{Name: "trp", Subcommands: []string{"get", "set", "getini", "setini"}, Description: "Temperature zones"},
	{Name: "tsensor", Subcommands: nil, Description: "Get raw temperature [sensor]"},
	{Name: "tshutdown", Subcommands: []string{"get", "set", "getini", "setini"}, Description: "Thermal shutdown"},
	{Name: "tshutdowntime", Subcommands: nil, Description: "Thermal shutdown time [time]"},
	{Name: "tzone", Subcommands: nil, Description: "Show thermal zones"},
	{Name: "version", Subcommands: nil, Description: "SC firmware version"},
	{Name: "w", Subcommands: nil, Description: "Write byte to SC [offset] [value]"},
	{Name: "w16", Subcommands: nil, Description: "Write word to SC [offset] [value]"},
	{Name: "w32", Subcommands: nil, Description: "Write dword to SC [offset] [value]"},
	{Name: "w64", Subcommands: nil, Description: "Write qword to SC [offset] [value]"},
	{Name: "wbe", Subcommands: nil, Description: "Write to BE [offset] [value]"},
	{Name: "wmmto", Subcommands: []string{"get"}, Description: "Get watch dog timeout"},
	{Name: "wrsxc", Subcommands: nil, Description: "Write to RSX [offset] [value]"},
	{Name: "xdrdiag", Subcommands: []string{"start", "info", "result"}, Description: "XDR diag"},
	{Name: "xiodiag", Subcommands: nil, Description: "XIO diag"},
	{Name: "xrcv", Subcommands: nil, Description: "Xmodem receive"},
}

// GetCommandNames returns a list of all CXR command names.
func GetCommandNames() []string {
	names := make([]string, len(MullionCommands))
	for i, cmd := range MullionCommands {
		names[i] = cmd.Name
	}
	return names
}

// GetCXRFCommandNames returns a list of all CXRF command names.
func GetCXRFCommandNames() []string {
	names := make([]string, len(CXRFCommands))
	for i, cmd := range CXRFCommands {
		names[i] = cmd.Name
	}
	return names
}

// GetCommand returns the CXR command with the given name (case-insensitive).
func GetCommand(name string) *Command {
	name = strings.ToUpper(strings.TrimSpace(name))
	for i := range MullionCommands {
		if MullionCommands[i].Name == name {
			return &MullionCommands[i]
		}
	}
	return nil
}

// GetCXRFCommand returns the CXRF command with the given name (case-sensitive for internal commands).
func GetCXRFCommand(name string) *Command {
	name = strings.TrimSpace(name)
	for i := range CXRFCommands {
		if CXRFCommands[i].Name == name {
			return &CXRFCommands[i]
		}
	}
	return nil
}

// HasSubcommands returns true if the command has subcommands.
func (c *Command) HasSubcommands() bool {
	return len(c.Subcommands) > 0
}
