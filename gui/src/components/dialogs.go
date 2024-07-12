package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"mhmods_gui/src/utils"
)

func ConfirmAppRestart(parent fyne.Window) {
	dialog.ShowConfirm("Restart", "Restart the app to apply the changes?", func(restart bool) {
		if restart {
			if err := utils.RestartApp(); err != nil {
				dialog.ShowError(err, parent)
			}
		}
	}, parent)
}
