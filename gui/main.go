package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	"mhmods_gui/src"
	"mhmods_gui/src/components"
	"mhmods_gui/src/utils"
)

func main() {
	// create
	modApp := app.New()
	modWindow := modApp.NewWindow(src.AppName + " - " + src.AppVersion)
	if err := utils.EnsureDataDirs(); err != nil {
		dialog.ShowError(err, modWindow)
		return
	}

	// build tabs
	gameTabs, err := components.GameTabs(modWindow)
	if err != nil {
		dialog.ShowError(err, modWindow)
		return
	}
	mainLayout := container.NewStack(gameTabs)
	modWindow.SetContent(mainLayout)

	// show
	modWindow.Resize(fyne.NewSize(800, 600))
	modWindow.ShowAndRun()
}
