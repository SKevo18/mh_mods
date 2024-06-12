package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/utils"
)

func gameButtons(parent fyne.Window, game *utils.Game) fyne.CanvasObject {
	split := container.NewHSplit(
		importModButton(parent, game), 
		launchGameButton(game),
	)

	return container.NewPadded(split)
}

func importModButton(parent fyne.Window, game *utils.Game) *widget.Button {
	return widget.NewButton("Import Mod", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
		}, parent)
	})
}

func launchGameButton(game *utils.Game) *widget.Button {
	return widget.NewButton("Launch Game", func() {
		utils.LaunchGame(game.Id)
	})
}
