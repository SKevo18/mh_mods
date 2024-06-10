package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func browseModsTab(parent fyne.Window, gameId string) *container.TabItem {
	browseTab := container.NewGridWithColumns(1,
		modItem(parent, "Mod A", "Description A"),
		modItem(parent, "Mod B", "Description B"),
	)

	return container.NewTabItem("Browse Mods", browseTab)
}
