package components

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"mhmods_gui/src/utils"
)

func installedModsTab(parent fyne.Window, gameId string, dataDir string) *container.TabItem {
	installedTab := container.NewGridWithColumns(1,
		installedModsItems(parent, gameId, dataDir)...,
	)

	return container.NewTabItem("Installed Mods", installedTab)
}

func installedModsItems(parent fyne.Window, gameId string, dataDir string) []fyne.CanvasObject {
	modFolder := filepath.Join(dataDir, "mods", gameId)
	mods, err := utils.GetInstalledMods(modFolder)
	if err != nil {
		return utils.TextLabel(err.Error())
	}

	items := make([]fyne.CanvasObject, 0, len(mods))
	for _, mod := range mods {
		items = append(items, modItem(parent, mod, "Description " + mod, false))
	}

	if len(items) == 0 {
		items = append(items, widget.NewLabel("You don't have any mods installed."))
	}

	return items
}
