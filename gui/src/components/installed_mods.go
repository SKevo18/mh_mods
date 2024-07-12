package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/utils"
)

func installedModsTab(parent fyne.Window, game *utils.Game) *container.TabItem {
	installedTab := container.NewGridWithColumns(1,
		installedModsItems(parent, game)...,
	)

	return container.NewTabItem("Installed Mods", installedTab)
}

func installedModsItems(parent fyne.Window, game *utils.Game) []fyne.CanvasObject {
	modFolder := game.ModsFolder()
	mods, err := utils.GetFolderNames(modFolder)
	if err != nil {
		return utils.TextLabel(err.Error())
	}

	items := make([]fyne.CanvasObject, 0, len(mods))
	for _, mod := range mods {
		items = append(items, modItem(parent, game, mod, false))
	}

	if len(items) == 0 {
		items = append(items, widget.NewLabel("You don't have any mods installed."))
	}

	return items
}
