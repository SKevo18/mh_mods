package components

import (
	"mhmods_gui/src"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

func GameTabs(parent fyne.Window) *container.AppTabs {
	games := make([]*container.TabItem, 0, len(src.SupportedGames))
	for gameId, gameName := range src.SupportedGames {
		games = append(games, gameTab(parent, gameName, gameId))
	}

	return container.NewAppTabs(games...)
}

func gameTab(parent fyne.Window, gameName string, gameId string) *container.TabItem {
	gameTab := container.NewVBox(
		container.NewAppTabs(
			installedModsTab(parent, gameId),
			browseModsTab(parent, gameId),
		),
		layout.NewSpacer(),
		gameButtons(parent, gameId),
	)

	return container.NewTabItemWithIcon(gameName, theme.MenuIcon(), gameTab)
}
