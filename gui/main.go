package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"mhmods_gui/src"
	"mhmods_gui/src/components"
	"mhmods_gui/src/utils"
)

func main() {
	// create
	modApp := app.New()
	modWindow := modApp.NewWindow(src.AppName + " - " + src.AppVersion)
	components.HandleError(utils.EnsureDataDirs(), modWindow)

	// build tabs
	gameTabs, err := components.GameTabs(modWindow)
	components.HandleError(err, modWindow)

	mainLayout := container.NewStack(gameTabs)
	modWindow.SetContent(mainLayout)

	// show
	modWindow.Resize(fyne.NewSize(800, 600))
	modWindow.ShowAndRun()
}
