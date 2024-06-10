package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func gameButtons(parent fyne.Window, gameId string) fyne.CanvasObject {
	split := container.NewHSplit(
		importModButton(parent, gameId), 
		launchGameButton(parent, gameId),
	)

	return container.NewPadded(split)
}

func importModButton(parent fyne.Window, gameId string) *widget.Button {
	return widget.NewButton("Import Mod", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			defer reader.Close()

			// ...
		}, parent)
	})
}

func launchGameButton(parent fyne.Window, gameId string) *widget.Button {
	return widget.NewButton("Launch Game", func() {
		// ...
	})
}
