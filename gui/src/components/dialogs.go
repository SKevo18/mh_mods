package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"mhmods_gui/src/utils"
)

// Shows an error dialog if the given error is not nil
func HandleError(err error, parent fyne.Window) {
	if err != nil {
		dialog.ShowError(err, parent)
	}
}

// Shows a confirmation dialog to restart the app
func ConfirmAppRestart(parent fyne.Window) {
	dialog.ShowConfirm("Restart", "Restart the app to apply the changes?", func(restart bool) {
		if restart {
			HandleError(utils.RestartApp(), parent)
		}
	}, parent)
}
