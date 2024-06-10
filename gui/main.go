package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"mhmods_gui/src"
	"mhmods_gui/src/components"
)

func main() {
	modApp := app.New()
	modWindow := modApp.NewWindow(src.AppName + " - " + src.AppVersion)

	gameTabs := components.GameTabs(modWindow)
	mainLayout := container.NewBorder(nil, nil, nil, nil, gameTabs)
	modWindow.SetContent(mainLayout)

	modWindow.Resize(fyne.NewSize(800, 600))
	modWindow.ShowAndRun()
}
