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
	modApp := app.New()
	modWindow := modApp.NewWindow(src.AppName + " - " + src.AppVersion)

	if err := utils.EnsureDataDirs(); err != nil {
		dialog.ShowError(err, modWindow)
		return
	}

	gameTabs := components.GameTabs(modWindow)
	mainLayout := container.NewStack(gameTabs)
	modWindow.SetContent(mainLayout)

	modWindow.Resize(fyne.NewSize(800, 600))
	modWindow.ShowAndRun()
}
