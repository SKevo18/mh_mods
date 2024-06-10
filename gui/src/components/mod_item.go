package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func modItem(parent fyne.Window, title, description string, downloadOnly bool) fyne.CanvasObject {
	// enable/disable checkbox
	checkbox := widget.NewCheck(title, func(enabled bool) {
		if enabled {
			// ...
		}
	})

	// description
	descriptionLabel := widget.NewLabel(description)

	// download
	downloadButton := widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {
		// ...
	})
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

	return container.NewHBox(
		checkbox,
		descriptionLabel,
		layout.NewSpacer(),
		buttons,
	)
}
