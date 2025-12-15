# PS3 Syscon Diagnostic Tool
## UART Interface Setup & Command Reference (V3)

## Disclaimer and Limitation of Liability
> [!CAUTION]
> **PLEASE READ THIS DISCLAIMER CAREFULLY BEFORE PROCEEDING**
>
> This guide and the accompanying software are provided "AS IS" without warranty of any kind, either expressed or implied, including but not limited to the implied warranties of merchantability, fitness for a particular purpose, or non-infringement.

By using this guide, software, or any associated materials, you acknowledge and agree that:

1. **Assumption of Risk:** You assume all risks associated with modifying your PlayStation 3 hardware. This includes but is not limited to permanent damage to the motherboard, SYSCON, EEPROM, or other components.

2. **No Warranty:** The authors and contributors make no representations or warranties regarding the accuracy, completeness, or suitability of the information provided. Hardware modification carries inherent risks that cannot be fully anticipated or documented.

3. **Limitation of Liability:** In no event shall the authors, contributors, or distributors be liable for any direct, indirect, incidental, special, exemplary, or consequential damages (including but not limited to procurement of substitute goods or services, loss of use, data, or profits, or business interruption) however caused and on any theory of liability, whether in contract, strict liability, or tort (including negligence or otherwise) arising in any way out of the use of this guide or software.

4. **Expertise Required:** This procedure requires soldering skills and knowledge of electronics. Improper connections can result in damage to your PS3 or USB serial adapter. If you are not comfortable with soldering or electronics work, seek assistance from a qualified technician.

5. **Backup Responsibility:** You are solely responsible for backing up any data before proceeding. The authors are not responsible for any data loss.

> [!CAUTION]
> **BY PROCEEDING PAST THIS POINT, YOU ACKNOWLEDGE THAT YOU HAVE READ, UNDERSTOOD, AND AGREED TO THIS DISCLAIMER.**

---

## Why Do This?

- Diagnose YLOD (Yellow Light of Death) issues commonly found on older PS3 boards
- Adjust default fan speed policies to improve cooling
- Read error logs to identify hardware failures
- Understand the internal workings of the PS3 system
- Perform syscon remarry procedures
- Many other diagnostic and repair purposes

---

## PS3 Motherboard Types

### Mullion SYSCONs (A - K Models)
- **COK-001, COK-002** - Fat models (CECHA - CECHE)
- **SEM-001** - Fat models (CECHG)
- **DIA-001, DIA-002** - Fat models (CECHH - CECHK)

### Sherwood SYSCONs (L Models Onward)
- **VER-001** - Slim models (CECHL - CECHQ)
- **DYN-001** - Slim models (CECH-20xx)
- **SUR-001, JTP-001, JSD-001, KTE-001** - Slim models (CECH-21xx to CECH-40xx)
- **SW** - All Super Slim models

---

## Required Equipment

### Hardware
- 1x USB to TTL 3.3V Serial Converter Cable (FTDI or CH340 based)
- 4x AWG 30 single core wires (keep length short, under 15cm recommended)
- Soldering iron and solder
- Wire cutters/strippers
- Electrical tape or heat shrink tubing (for insulation)

### Software
- **PS3 Syscon Tool** - Download the pre-built GUI for your platform:
  - Windows: `ps3syscon-win.exe`
  - macOS: `ps3syscon-macos`
  - Linux: `ps3syscon-linux`

---

## Identify the Serial Connection Points

Refer to the images in the `PS3-Serial-Connection-Points` folder for your specific motherboard:

| Motherboard | Folder |
|-------------|--------|
| COK-001, COK-002 | `COK-001-002-(A-E-Models)/` |
| SEM-001 | `SEM-001-(G-Models)/` |
| DIA-001, DIA-002 | `DIA-001-002-(H-K-Models)/` |
| VER-001 | `VER-001-(L-Q-Models)/` |
| DYN-001 | `DYN-001-(20xx-Model)/` |
| SUR/JTP/JSD/KTE-001 | `SUR-001-JTP-001-JSD-001-KTE-001-(21xx-40xx-Models)/` |
| Super Slim (SW) | `All-Super-Slim-Models/` |

---

## Connecting the Serial Wires

You need 4 wires: RxD, TxD, DIAG, and GND

Cut and strip one end of each wire to solder onto the motherboard test points.

### Pin Connections

| Motherboard Pin | USB TTL Cable Pin | Description |
|-----------------|-------------------|-------------|
| RxD | TX | Receive (MB receives from cable TX) |
| TxD | RX | Transmit (MB transmits to cable RX) |
| GND | GND | Ground connection |
| DIAG | Connect to GND | Enables diagnostic/internal mode (not used on SW models) |
| 3.3V | **DO NOT CONNECT!** | Never connect the 3.3V pin |

**Important:** The RX/TX are crossed - motherboard RxD goes to cable TX, and motherboard TxD goes to cable RX.

---

## Software Installation

### Windows
1. Download `ps3syscon-win.exe`
2. Double-click to run (no installation needed)
3. If Windows SmartScreen appears, click "More info" then "Run anyway"

### macOS
1. Download `ps3syscon-macos`
2. Open Terminal and make it executable:
   ```bash
   chmod +x ps3syscon-macos
   ```
3. Run the application:
   ```bash
   ./ps3syscon-macos
   ```
4. If blocked by Gatekeeper, right-click the file and select "Open"

---

## SYSCON Types and Modes

The PS3 Syscon Tool supports three SYSCON types:

| Type | Baud Rate | Description |
|------|-----------|-------------|
| **CXR** | 57600 | External command mode (standard mode) |
| **CXRF** | 115200 | Internal command mode (diagnostic mode with DIAG grounded) |
| **SW** | 115200 | Sherwood syscons (Slim/Super Slim models, L onwards) |

### Which Mode to Use?

- **CXR**: Use for initial connection without DIAG grounded
- **CXRF**: Use after enabling internal mode (DIAG grounded, offset 0x3961 set to 00)
- **SW**: Use for all Slim and Super Slim models (VER-001 and later)

**Note:** SW (Sherwood) models do NOT have a DIAG mode - they connect directly.

---

## Using the PS3 Syscon Tool

### Basic Operation

1. Connect your USB TTL cable to the PS3 motherboard (RxD, TxD, GND)
2. Plug the USB end into your computer
3. Launch the PS3 Syscon Tool
4. Select your serial port from the dropdown (click Refresh if needed)
5. Select your SC Type (CXR, CXRF, or SW)
6. Click "Connect"

### Serial Monitor

The tool includes a built-in serial monitor window for viewing raw UART output. This is useful for:
- Verifying the connection is working
- Viewing boot messages
- Debugging communication issues

### Authentication

Many commands require authentication. Click the "Auth" button or type `AUTH` in the command field. You should see:
```
Auth successful
```

**Troubleshooting:** If you get "Auth1 response invalid":
1. Swap your RX and TX wires
2. Power cycle the PS3
3. Reconnect and try again

---

## External Command Mode (CXR)

### Setup
1. Connect RxD, TxD, and GND (leave DIAG disconnected)
2. Power on the PS3
3. Select "CXR" in the tool and connect
4. Run `AUTH` to authenticate

> [!TIP]
> Full details of all syscon commands are available at: https://www.psdevwiki.com/ps3/Syscon_Firmware#Command_list

### Available Commands

| Command | Subcommand | Description |
|---------|------------|-------------|
| BOOT | MODE | Boot mode control |
| BOOT | CONT | Continue boot |
| SHUTDOWN | - | Shutdown the system |
| HALT | - | Halt the system |
| AUTH1 | - | Authentication step 1 |
| AUTH2 | - | Authentication step 2 |
| EEP | GET/SET/INIT | EEPROM read/write operations |
| VID | GET | Get voltage ID |
| CID | GET | Get chip ID |
| ECID | GET | Get electronic chip ID |
| FAN | SETPOLICY/GETPOLICY | Fan policy control |
| FAN | SETDUTY/GETDUTY | Fan duty cycle control |
| ERRLOG | START/STOP/GET/CLEAR | Error log operations |
| VER | - | Get version info |
| BUZ | - | Activate buzzer |

### Enabling Internal Command Mode

To access the full command set, you need to enable internal mode:

1. Connect and authenticate in CXR mode
2. Check the current offset value:
   ```
   EEP GET 3961 01
   ```
   Should return: `00000000 FF`

3. Set the offset to enable internal mode:
   ```
   EEP SET 3961 01 00
   ```

4. Verify the change:
   ```
   EEP GET 3961 01
   ```
   Should return: `00000000 00`

5. Power off the PS3
6. Ground the DIAG pin to GND
7. Continue to Internal Command Mode

**Warning:** Setting this offset will temporarily prevent the PS3 from booting until the EEPROM checksum is corrected in internal mode!

---

## Internal Command Mode (CXRF)

### Setup
1. Ground the DIAG pin to GND
2. Power on the PS3 (LED will flash red - this is normal)
3. Select "CXRF" in the tool and connect
4. Run `auth` to authenticate

### Fixing the EEPROM Checksum

After enabling internal mode, you must fix the checksum:

1. Run the checksum command:
   ```
   eepcsum
   ```

   Example output:
   ```
   Addr:0x000032fe should be 0x528c
   Addr:0x000034fe should be 0x7115
   sum:0x0100
   Addr:0x000039fe should be 0x0038
   Addr:0x00003dfe should be 0x00ff
   Addr:0x00003ffe should be 0x00ff
   ```

2. The `sum:0x0100` line indicates an error. Fix the incorrect checksum using the write command.

   **Important:** Values are little-endian, so swap the bytes when writing!

   If `Addr:0x000039fe should be 0x0038`, write:
   ```
   w 39FE 38 00
   ```

3. Verify the fix:
   ```
   eepcsum
   ```
   The `sum:0x0100` line should disappear when the checksum is correct.

### Common Internal Commands

#### Diagnostic Commands

| Command | Description |
|---------|-------------|
| `version` | Show syscon firmware version |
| `errlog` | Display full error log |
| `lasterrlog` | Show last error since boot |
| `clearerrlog` | Clear the error log |
| `eepcsum` | Show EEPROM checksum status |
| `powerstate` | Show current power states |
| `bringup` | Power on the system |
| `shutdown` | Power off immediately |
| `syspowdown 0 0 0` | System power down |

#### Temperature Commands

| Command | Description |
|---------|-------------|
| `tmp 0` | CELL temperature |
| `tmp 1` | RSX temperature |
| `tsensor 0` | CELL thermal sensor raw data |
| `tsensor 1` | RSX thermal sensor raw data |
| `tsensor 3` | Southbridge temperature |
| `tzone` | List all thermal zones |
| `tshutdown get 0` | CELL thermal shutdown threshold |
| `tshutdown get 1` | RSX thermal shutdown threshold |

#### Fan Commands

| Command | Description |
|---------|-------------|
| `duty get 0` | Current fan duty for CELL cooling |
| `duty get 1` | Current fan duty for RSX cooling |
| `fanconpolicy get 0` | Fan policy for CELL |
| `fanconpolicy get 1` | Fan policy for RSX |
| `fantbl get` | Get current fan table |
| `fantbl set` | Set fan table values |

#### Memory Read/Write Commands

| Command | Description |
|---------|-------------|
| `r [offset] [length]` | Read bytes from EEPROM |
| `w [offset] [value]` | Write byte to EEPROM |
| `r16 [offset] [length]` | Read words (16-bit) |
| `w16 [offset] [value]` | Write word (16-bit) |
| `r32 [offset] [length]` | Read dwords (32-bit) |
| `w32 [offset] [value]` | Write dword (32-bit) |

---

## Sherwood Mode (SW)

For Slim and Super Slim models (VER-001 onward):

1. Connect RxD, TxD, and GND only (no DIAG pin on these models)
2. Power on the PS3
3. Select "SW" in the tool and connect
4. Run `auth` to authenticate

SW models have the same internal commands as CXRF mode but do not require the DIAG jumper or the offset modification.

---

## Common Diagnostic Procedures

### Reading Error Logs

```
auth
errlog
```

This shows all recorded errors with timestamps. Use the error code reference in the README.md file to interpret the codes.

### Checking System Health

```
auth
powerstate
tmp 0
tmp 1
tsensor 3
```

### Modifying Fan Settings

To increase minimum fan speed:
```
auth
duty get 0
duty setmin 0 [value]
```

**Note:** After modifying fan or thermal settings, run `eepcsum` and fix any checksum errors.

---

## EEPROM Memory Map

| Offset Range | Description | Notes |
|--------------|-------------|-------|
| 0x2600 - 0x27FF | System Info | Encrypted |
| 0x2800 - 0x2BFF | Patch Part 1 | Encrypted |
| 0x2F00 - 0x2FFF | Industry Area | |
| 0x3000 - 0x30FF | Customer Service Area | |
| 0x3100 - 0x31FF | Platform Config | |
| 0x3300 - 0x33FF | Fan/Thermal Config | |
| 0x3400 - 0x34FF | Fan/Thermal Config | |
| 0x3600 - 0x36FF | On/Off Count, On-Time | |
| 0x3800 - 0x38FF | Error Log | |
| 0x3900 - 0x39FF | Board Config | |
| 0x3A00 - 0x3AFF | HDMI/DVE Config | |

---

## Troubleshooting

### "Auth1 response invalid"
- Swap RX and TX wires
- Power cycle the PS3 and reconnect

### No response from syscon
- Verify wiring connections
- Check that the USB TTL cable is 3.3V (not 5V!)
- Try a different USB port
- Verify the correct serial port is selected

### PS3 won't boot after enabling internal mode
- This is expected! The EEPROM checksum is invalid
- Ground the DIAG pin and connect in CXRF mode
- Run `eepcsum` and fix the checksum using the `w` command

### Commands return errors
- Make sure you've run `auth` first
- Some commands require specific board revisions
- Check if the command is available for your syscon type

---

## Additional Resources

- Error code reference: See `README.md` in the project root
- Syscon error codes wiki: https://www.psdevwiki.com/ps3/Syscon_Error_Codes
- Syscon commands wiki: https://www.psdevwiki.com/ps3/Syscon_Firmware#Command_list

---

## Credits

- Original Python scripts and research by the PS3 homebrew community
- GUI tool development for cross-platform compatibility
- Documentation compiled from various sources including PSX-Place and PS3 Developer Wiki

---

*PS3 Syscon UART Guide V3 - Updated for the PS3 Syscon Tool GUI Application*
