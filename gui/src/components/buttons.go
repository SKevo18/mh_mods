package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/api"
	"mhmods_gui/src/utils"
)

// A container that holds "the two wide buttons at bottom" (Open Mods Folder and Launch Game).
func GameButtons(parent fyne.Window, game *utils.Game) fyne.CanvasObject {
	split := container.NewHSplit(
		openModsFolderButton(parent, game),
		launchGameButton(game),
	)

	return container.NewPadded(split)
}

// The "Open Mods Folder" button.
func openModsFolderButton(parent fyne.Window, game *utils.Game) *widget.Button {
	return widget.NewButton("Open Mods Folder", func() {
		HandleError(utils.OpenFolder(game.ModsFolder()), parent)
		ConfirmAppRestart(parent, false)
	})
}

// The "Launch Game" button.
func launchGameButton(game *utils.Game) *widget.Button {
	return widget.NewButton("Launch Game", func() {
		HandleError(game.Launch(), nil)
	})
}

// The download/update mod button
func DownloadModButton(parent fyne.Window, game *utils.Game, modId string) *widget.Button {
	canUpdate, err := hasUpdates(game.Id, modId)
	HandleError(err, parent)

	btn := widget.NewButton("Download Mod", nil)
	btn.OnTapped = func() {
		HandleError(game.DownloadMod(modId), parent)
		dialog.ShowInformation("Mod Downloaded", "The mod has been downloaded and installed successfully.", parent)
		btn.Disable()
	}

	// disable button if no updates are available
	if !canUpdate {
		btn.Disable()
	}

	return btn
}

// Helper function that returns true if the API mod hash is different from the local hash (has updates)
func hasUpdates(gameId string, modId string) (bool, error) {
	currentHash, err := api.GetCurrentModHash(gameId, modId)
	if err != nil {
		return false, err
	}

	localHash, err := utils.ReadHash(gameId, modId)
	if err != nil {
		return false, err
	}

	return currentHash != localHash, nil
}
