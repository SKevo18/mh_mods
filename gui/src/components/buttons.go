package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/utils"
)

// A container that holds "the two wide buttons at bottom" (Open Mods Folder and Launch Game).
func gameButtons(parent fyne.Window, game *utils.Game) fyne.CanvasObject {
	split := container.NewHSplit(
		openModsFolderButton(parent, game),
		launchGameButton(game),
	)

	return container.NewPadded(split)
}

// The "Open Mods Folder" button.
func openModsFolderButton(parent fyne.Window, game *utils.Game) *widget.Button {
	return widget.NewButton("Open Mods Folder", func() {
		if err := utils.OpenFolder(game.ModFolder()); err != nil {
			dialog.ShowError(err, parent)
		}

		ConfirmAppRestart(parent)
	})
}

// The "Launch Game" button.
func launchGameButton(game *utils.Game) *widget.Button {
	return widget.NewButton("Launch Game", func() {
		if err := game.Launch(); err != nil {
			dialog.ShowError(err, nil)
		}
	})
}
