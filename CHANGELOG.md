# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Removed
- Removed Linux binary (ps3syscon-linux) - platform not currently supported

## [1.0.0] - 2025-12-15

### Added
- Initial release of PS3 Syscon UART Tool
- Cross-platform GUI application built with Go and Fyne framework
- Pre-built binaries for Windows (64-bit) and macOS (64-bit)
- Serial port selection with refresh functionality
- Support for CXR, CXRF, and SW (Sherwood) syscon types
- Built-in serial monitor for diagnostics
- AES-CBC authentication support for syscon communication
- Comprehensive documentation including:
  - PS3 UART connection guides (V1 and V2)
  - Syscon error code references
  - Motherboard test point diagrams for various PS3 models
  - CELL, RSX, and Southbridge pinout diagrams
  - Power flow charts and diagnostic guides
  - Fan table references for different motherboard revisions
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
