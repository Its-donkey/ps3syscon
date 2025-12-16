// Package ui provides the startup disclaimer dialog.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const disclaimerText = `DISCLAIMER AND LIMITATION OF LIABILITY

This software is provided "AS IS" without warranty of any kind, either expressed or implied, including but not limited to the implied warranties of merchantability, fitness for a particular purpose, or non-infringement.

By using this software, you acknowledge and agree that:

1. ASSUMPTION OF RISK: You assume all risks associated with modifying your PlayStation 3 hardware. This includes but is not limited to permanent damage to the motherboard, SYSCON, EEPROM, or other components.

2. NO WARRANTY: The author makes no representations or warranties regarding the accuracy, completeness, or suitability of this software. Hardware modification carries inherent risks that cannot be fully anticipated.

3. LIMITATION OF LIABILITY: In no event shall the author be liable for any direct, indirect, incidental, special, exemplary, or consequential damages however caused.

4. EXPERTISE REQUIRED: This software requires soldering skills and knowledge of electronics. Improper connections can result in damage to your PS3 or USB serial adapter.

5. BACKUP RESPONSIBILITY: You are solely responsible for backing up any data before proceeding.

BY CLICKING "ACCEPT", YOU ACKNOWLEDGE THAT YOU HAVE READ, UNDERSTOOD, AND AGREED TO THIS DISCLAIMER.`

// ShowDisclaimer displays a disclaimer dialog that must be accepted to proceed.
// onAccept is called when the user accepts, onDecline when they decline.
func ShowDisclaimer(myWindow fyne.Window, onAccept, onDecline func()) {
	textWidget := widget.NewLabel(disclaimerText)
	textWidget.Wrapping = fyne.TextWrapWord

	scroll := container.NewVScroll(textWidget)
	scroll.SetMinSize(fyne.NewSize(500, 300))

	dialog.ShowCustomConfirm("Terms of Use", "Accept", "Decline", scroll, func(accepted bool) {
		if accepted {
			onAccept()
		} else {
			onDecline()
		}
	}, myWindow)
}
