package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/api"
	"mhmods_gui/src/utils"
)

func downloadModsTab(parent fyne.Window, gameId string) *container.TabItem {
	downloadTab := container.NewGridWithColumns(1,
		downloadModsItems(parent, gameId)...,
	)

	return container.NewTabItem("Download Mods", downloadTab)
}

func downloadModsItems(parent fyne.Window, gameId string) []fyne.CanvasObject {
	mods, err := api.GetMods(gameId)
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
