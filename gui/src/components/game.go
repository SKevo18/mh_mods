package components

import (
	"mhmods_gui/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// Creates a tabbed view for all games.
func GameTabs(parent fyne.Window) (*container.AppTabs, error) {
	// get game folders
	mod_paths, err := utils.GetFolderNames(utils.DataDir)
	if err != nil {
		return nil, err
	}

	// load game from config and create tab
	games := make([]*container.TabItem, 0, len(mod_paths))
	for _, gameId := range mod_paths {
		game, err := utils.LoadGame(gameId)
		if err != nil {
			return nil, err
		}

		games = append(games, gameTab(parent, game))
	}

	return container.NewAppTabs(games...), nil
}

// Creates a tab for the given game.
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
