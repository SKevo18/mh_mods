package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/utils"
)

// TODO: split into two functions - one for installed mods, another for downloadable mods
// installed mods - settings (simply opens config file for now), delete mod (asks for confirmation, removes dir)
// browse mods - download: if already downloaded, check updates instead
func modItem(parent fyne.Window, game *utils.Game, modId string, downloadOnly bool) *fyne.Container {
	modItemContainer := container.NewHBox()

	// enable/disable checkbox
	if !downloadOnly {
		checkbox := widget.NewCheck("", func(enabled bool) {
			if enabled {
				// ...
			}
		})
		modItemContainer.Add(checkbox)
	}

	label := widget.NewLabel(modId)

	// download
	downloadButton := DownloadModButton(parent, game, modId)
	downloadButton.Importance = widget.HighImportance

	buttons := container.NewHBox(downloadButton)

	if !downloadOnly {
		// settings
		settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
			dialog.NewForm("Mod Settings", "Save", "Dismiss", []*widget.FormItem{
				widget.NewFormItem("Name", widget.NewEntry()),
				widget.NewFormItem("Description", widget.NewEntry()),
				widget.NewFormItem("Enabled", widget.NewCheck("", func(enabled bool) {
					// ...
				})),
			},
				func(bool) {}, parent).Show()
		})
		settingsButton.Importance = widget.WarningImportance

		// delete
		deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			dialog.ShowConfirm("Delete Mod", "Are you sure you want to delete this mod?", func(b bool) {
				if b {
					// ...
				}
			}, parent)
		})
		deleteButton.Importance = widget.DangerImportance

		buttons.Add(settingsButton)
		buttons.Add(deleteButton)
	}

	modItemContainer.Add(label)
	modItemContainer.Add(layout.NewSpacer())
	modItemContainer.Add(buttons)

	return modItemContainer
}
