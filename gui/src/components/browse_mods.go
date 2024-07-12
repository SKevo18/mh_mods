package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/api"
	"mhmods_gui/src/utils"
)

// Creates a tab for downloading/updating mods
func downloadModsTab(parent fyne.Window, game *utils.Game) *container.TabItem {
	downloadTab := container.NewGridWithColumns(1,
		downloadModsItems(parent, game)...,
	)

	return container.NewTabItem("Download Mods", downloadTab)
}

// Creates a list of mod items that can be downloaded/updated
func downloadModsItems(parent fyne.Window, game *utils.Game) []fyne.CanvasObject {
	// fetch mods from API (https://mh-mods.svit.ac/mods/<game.Id>)
	mods, err := api.GetMods(game.Id)
	if err != nil {
		return utils.TextLabel(err.Error())
	}

	items := make([]fyne.CanvasObject, 0, len(mods))
	for mod, hash := range mods {
		items = append(items, modItem(parent, mod+" "+hash, true))
	}

	if len(items) == 0 {
		items = append(items, widget.NewLabel("No mods found."))
	}

	return items
}
