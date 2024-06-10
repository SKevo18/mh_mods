package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func installedModsTab(parent fyne.Window, gameId string) *container.TabItem {
	installedTab := container.NewGridWithColumns(1,
		modItem(parent, "Mod 1", "Description 1"),
		modItem(parent, "Mod 2", "Description 2"),
	)

	return container.NewTabItem("Installed Mods", installedTab)
}
