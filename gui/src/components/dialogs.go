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
// If force is true, the app will restart without asking
func ConfirmAppRestart(parent fyne.Window, force bool) {
	if force {
		HandleError(utils.RestartApp(), parent)
		return
	}

	dialog.ShowConfirm("Restart", "Restart the app to apply the changes?", func(restart bool) {
		if restart {
			HandleError(utils.RestartApp(), parent)
		}
	}, parent)
}
