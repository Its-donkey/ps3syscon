# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-12-16

### Changed
- Simplified CXRF mode selection - now switches command UI immediately without automated wizard
- Users can manually authenticate using the "Authenticate" button when ready
- Removed automatic EEP 3961 checking/setting during mode switch
- Removed multi-step setup dialogs for CXRF mode

### Improved
- More responsive mode switching between CXR, CXRF, and SW
- Users have full control over the CXRF setup process

### Acknowledgements
- Thanks to Rambonz for helping with testing

## [1.0.0] - 2025-12-15

### Added
- Initial release of PS3 Syscon UART Tool
- Cross-platform GUI application built with Go and Fyne framework
- Pre-built binaries for Windows (64-bit) and macOS (64-bit)
- Serial port selection with refresh functionality
- Support for CXR, CXRF, and SW (Sherwood) syscon types
- Built-in serial monitor for diagnostics
- AES-CBC authentication support for syscon communication
- Searchable command list with autocomplete
- Subcommand support for complex commands
- Custom dark theme with PS3-inspired design
- Comprehensive documentation including:
  - PS3 UART connection guides
  - Syscon error code references
  - Motherboard test point diagrams for various PS3 models

### Technical Details
- Written in Go 1.21+
- Uses Fyne v2 GUI framework
- Embedded PNG icon for cross-platform consistency
- Thread-safe UI updates with fyne.Do()
