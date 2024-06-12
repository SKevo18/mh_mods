package components

import (
	"mhmods_gui/src"
	"mhmods_gui/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

func GameTabs(parent fyne.Window) *container.AppTabs {
	games := make([]*container.TabItem, 0, len(src.SupportedGames))
	for _, gameId := range src.SupportedGames {
		game := utils.Game{}
		game.FromConfig(gameId)

		games = append(games, gameTab(parent, &game))
	}

	return container.NewAppTabs(games...)
}

func gameTab(parent fyne.Window, game *utils.Game) *container.TabItem {
	gameTab := container.NewVBox(
		container.NewAppTabs(
			installedModsTab(parent, game),
			downloadModsTab(parent, game),
		),
		layout.NewSpacer(),
		gameButtons(parent, game),
	)

	return container.NewTabItemWithIcon(game.Name, theme.MenuIcon(), gameTab)
}
