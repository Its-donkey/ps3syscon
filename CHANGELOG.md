# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- Added go-gui source code to version control (previously ignored)

### Changed
- Refactored go-gui source code into separate files for better maintainability:
  - `crypto.go` - AES encryption/decryption and authentication keys
  - `serial.go` - PS3UART struct and serial communication logic
  - `ui.go` - Fyne GUI components and window management
  - `main.go` - Minimal entry point that wires dependencies and starts the app

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
