# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-12-15

### Added
- Initial release of PS3 Syscon UART Tool
- Cross-platform GUI built with Go and Fyne framework
- Support for CXR, CXRF, and SW (Sherwood) syscon types
- Serial port selection with auto-detection and refresh
- AES-CBC authentication support for CXRF mode
- Built-in serial monitor for diagnostics
- Searchable command list with autocomplete
- Subcommand support for complex commands
- Custom dark theme with PS3-inspired design
- Pre-built binaries for Windows and macOS
- Comprehensive error code reference in README

### Technical Details
- Written in Go 1.21+
- Uses Fyne v2 GUI framework
- Embedded PNG icon for cross-platform consistency
- Thread-safe UI updates with fyne.Do()
