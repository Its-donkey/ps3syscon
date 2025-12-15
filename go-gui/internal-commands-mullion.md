| Command | Address | Perms | SubCommands | Description |
|---|---:|---:|---|---|
| becount | 0xCA7D | 0xDD0C0000 | - | Display bringup/shutdown count + Power-on time |
| bepgoff | 0xA4E7 | 0xD00C0000 | - | BE power grid off |
| bepkt | 0x2435D | 0xDC0C0000 | show/set/unset/mode/debug/help | Packet permissions |
| bestat | 0xD413 | 0xFD0F0000 | - | Get status of BE |
| boardconfig | 0x99C7 | 0xDC0C0000 | - | Displays board configuration |
| bootbeep | 0x1EA67 | 0xF0000000 | stat/on/off | Boot beep |
| bringup | 0xD597 | 0xFD0F0000 | - | Turn PS3 on |
| bsn | 0xD805 | 0xF00F0000 | - | Get board serial number |
| bstatus | 0x24269 | 0xDD0C0000 | - | HDMI related status |
| buzz | 0xA4FF | 0xDC0C0000 | [freq] | Activate buzzer |
| buzzpattern | 0xA8B7 | 0xDC0C0000 | [freq] [pattern] [count] | Buzzer pattern |
| clear_err | 0x2595B | 0xDD0C0000 | last/eeprom/all | Clear errors |
| clearerrlog | 0xB8CB | 0xDD0C0000 | - | Clears error log |
| comm | 0x9919 | 0xDC0C0000 | - | Communication mode |
| commt | 0x24907 | 0xDC0C0000 | help/start/stop/send | Manual BE communication |
| cp | 0x1E077 | 0xF0000000 | ready/busy/reset/beepremote/beep2kn1n3/beep2kn2n3 | CP control commands |
| csum | 0xD687 | 0xFF0F0000 | - | Firmware checksum |
| devpm | 0xD053 | 0xDD0C0000 | ata/pci/pciex/rsx | Device power management |
| diag | 0x9AAD | 0xD00C0000 | ... | Diag (execute without param to show help) |
| disp_err | 0x25911 | 0xDD0C0000 | - | Displays errors |
| duty | 0x9B23 | 0xDD0C0000 | get/set/getmin/setmin/getmax/setmax/getinmin/setinmin/getinmax/setinmax | Fan policy |
| dve | 0x2995D | 0xDC0C0000 | help/set/save/show | DVE chip parameters |
| eepcsum | 0xAA65 | 0xDD0C0000 | - | Shows eeprom checksum |
| eepromcheck | 0x9A1D | 0x000C0000 | [id] | Check eeprom |
| eeprominit | 0x9A65 | 0x000C0000 | [id] | Init eeprom |
| ejectsw | 0xD611 | 0xFD0F0000 | - | Eject switch |
| errlog | 0xB7ED | 0xFF0C0000 | - | Gets the error log |
| fancon | 0xD26D | 0x0D000000 | - | Does nothing |
| fanconautotype | 0xC075 | 0xDD0C0000 | - | Does nothing |
| fanconmode | 0xBF35 | 0xDD0C0000 | get | Fan control mode |
| fanconpolicy | 0xBBC9 | 0xDD0C0000 | get/set/getini/setini | Fan control policy |
| fandiag | 0x1E91B | 0xF0000000 | - | Fan test |
| faninictrl | 0xD3D9 | 0x0D000000 | - | Does nothing |
| fanpol | 0xCA31 | 0xDD0C0000 | - | Does nothing |
| fanservo | 0xBF29 | 0xDD0C0000 | - | Does nothing |
| fantbl | 0xC087 | 0xDD0C0000 | get/set/getini/setini/gettable/settable | Fan table |
| firmud | 0xD61D | 0xFDFF0000 | - | Firmware update |
| geterrlog | 0xB84F | 0xDD0C0000 | [id] | Gets error log |
| getrtc | 0xA6F3 | 0xDD0C0000 | - | Gets rtc |
| halt | 0x1E107 | 0xF0000000 | - | Halts syscon |
| hdmi | 0x29F39 | 0xDD0C0000 | ... | HDMI (various commands, use help) |
| hdmiid | 0x29D1D | 0xDC0F0000 | - | Get HDMI id's |
| hdmiid2 | 0x29D81 | 0xDC0F0000 | - | Get HDMI id's |
| hversion | 0x2422F | 0xDD0C0000 | - | Platform ID |
| hyst | 0xAEF5 | 0xDD0C0000 | get/set/getini/setini | Temperature zones |
| lasterrlog | 0xB7FF | 0xDD0C0000 | - | Last error from log |
| ledmode | 0xA80B | 0xDC0C0000 | [id] [id] | Get led mode |
| LS | 0x2421B | 0xDD0C0000 | - | LabStation Mode |
| ltstest | 0xCB97 | 0xDD0C0000 | get/set be/rsx | ?Temp related? values |
| osbo | 0x1EA3F | 0xF0000000 | - | Sets 0x2000F60 |
| patchcsum | 0xD9F7 | 0xDD0C0000 | - | Patch checksum |
| patchvereep | 0xD9B1 | 0xDD0C0000 | - | Patch version eeprom |
| patchverram | 0xD965 | 0xDD0C0000 | - | Patch version ram |
| poll | 0x240E3 | 0xDD0C0000 | - | Poll log |
| portscan | 0xDA0D | 0xDD0C0000 | [port] | Scan port |
| powbtnmode | 0xB911 | 0xDC0C0000 | [mode (0/1)] | Power button mode |
| powerstate | 0xCE6F | 0xDD0C0000 | - | Get power state |
| powersw | 0xD5F9 | 0xFD0F0000 | - | Power switch |
| powupcause | 0xB621 | 0xDD0C0000 | - | Power up cause |
| printmode | 0x99D9 | 0xDC0C0000 | [mode (0/1/2/3)] | Set printmode |
| printpatch | 0xD94F | 0xDD0C0000 | - | Prints patch |
| r | 0x8CA5 | 0xDD0C0000 | [offset] [length] | Read byte from SC |
| r16 | 0x8ED5 | 0xDD0C0000 | [offset] [length] | Read word from SC |
| r32 | 0x9191 | 0xDD0C0000 | [offset] [length] | Read dword from SC |
| r64 | 0x935D | 0xDD0C0000 | [offset] [length] | Read qword from SC |
| r64d | 0x948F | 0xDD0C0000 | [offset] [length] | Read ?qword data? from SC |
| rbe | 0x96F9 | 0xDD0C0000 | [offset] | Read from BE |
| recv | 0x24135 | 0xDD0C0000 | - | Receive something |
| resetsw | 0xD605 | 0xFC0F0000 | - | Reset switch |
| restartlogerrtoeep | 0xB903 | 0xDD0C0000 | - | Reenable error logging to eeprom |
| revision | 0xD7E1 | 0xFFFF0000 | - | Get softid |
| rrsxc | 0xD313 | 0xDD0C0000 | [offset] [length] | Read from RSX |
| rtcreset | 0xA7BB | 0x000C0000 | - | Reset RTC |
| scagv2 | 0xE24F | 0xFF000000 | - | Auth related? |
| scasv2 | 0xE207 | 0xDD000000 | - | Auth related? |
| scclose | 0xE1EF | 0xFF000000 | - | Auth related? |
| scopen | 0xE121 | 0xFF000000 | - | Auth related? |
| send | 0x2416F | 0xDD0C0000 | [variable] | Send something |
| shutdown | 0xD5C5 | 0xFD0F0000 | - | PS3 shutdown |
| startlogerrtsk | 0xB8E7 | 0xDD0C0000 | - | Start error log task |
| stoplogerrtoeep | 0xB8F5 | 0xDD0C0000 | - | Stop error logging to eeprom |
| stoplogerrtsk | 0xB8D9 | 0xDD0C0000 | - | Stop error log task |
| syspowdown | 0xB6E9 | 0xDD0C0000 | 3 params 0 0 0 | System power down |
| task | 0x15005 | 0xDD0C0000 | - | Print tasks |
| thalttest | 0xD813 | 0x000F0000 | - | Does nothing |
| thermfatalmode | 0xCA3B | 0xDD0C0000 | canboot/cannotboot | Set thermal boot mode |
| therrclr | 0xD3E5 | 0xDD0C0000 | - | Thermal register clear |
| thrm | 0xBF1D | 0xDD0C0000 | - | Does nothing |
| tmp | 0xAA69 | 0xDD0C0000 | [zone] | Get temperature |
| trace | 0xB951 | 0xDD0C0000 | ... | Trace tasks (use help) |
| trp | 0xAB2F | 0xDD0C0000 | get/set/getini/setini | Temperature zones |
| tsensor | 0xA279 | 0xDD0C0000 | [sensor] | Get raw temperature |
| tshutdown | 0xB2A1 | 0xDD0C0000 | get/set/getini/setini | Thermal shutdown |
| tshutdowntime | 0xC95D | 0xDD0C0000 | [time] | Thermal shutdown time |
| tzone | 0xB5E1 | 0xDD0C0000 | - | Show thermal zones |
| version | 0xD65F | 0xFFFF0000 | - | SC firmware version |
| w | 0x8BF9 | 0xDD0C0000 | [offset] [value] | Write byte to SC |
| w16 | 0x8E2D | 0xDD0C0000 | [offset] [value] | Write word to SC |
| w32 | 0x8FED | 0xDD0C0000 | [offset] [value] | Write dword to SC |
| w64 | 0x92A9 | 0xDD0C0000 | [offset] [value] | Write qword to SC |
| wbe | 0x9665 | 0xDD0C0000 | [offset] [value] | Write to BE |
| wmmto | 0xCB3B | 0xDC0C0000 | get | Get watch dog timeout |
| wrsxc | 0xD279 | 0xDD0C0000 | [offset] [value] | Write to RSX |
| xdrdiag | 0x1E711 | 0xF0000000 | start/info/result | XDR diag |
| xiodiag | 0x1E875 | 0xF0000000 | - | XIO diag |
| xrcv | 0x25313 | 0xDC0C0000 | - | Xmodem receive |


