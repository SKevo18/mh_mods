package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var supportedGames = []string{
	"Moorhuhn Kart 1",
	"Moorhuhn Kart 2",
	"Moorhuhn Kart 3",
	"Moorhuhn Kart 4",
}

func main() {
	modApp := app.New()
	modWindow := modApp.NewWindow("Moorhuhn Kart Modding Tool")

	gameList := widget.NewList(
		func() int { return len(supportedGames) },
		func() fyne.CanvasObject { return widget.NewLabel("Loading...") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(supportedGames[i])
		},
	)
	gameList.OnSelected = func(id widget.ListItemID) {
		dialog.ShowInformation("Game Selected", fmt.Sprintf("Game %d selected", id+1), modWindow)
	}

	addGameButton := widget.NewButtonWithIcon("   Add Game", theme.ContentAddIcon(), func() {
		openAddGameModal(modWindow)
	})

	leftSide := container.NewBorder(nil, addGameButton, nil, nil, gameList)

	installedTab := container.NewVBox(
		createModListItem(modWindow, "Mod 1", "Description 1"),
		createModListItem(modWindow, "Mod 2", "Description 2"),
	)

	browseNewTab := container.NewVBox(
		createModListItem(modWindow, "Mod A", "Description A"),
		createModListItem(modWindow, "Mod B", "Description B"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Downloaded Mods", installedTab),
		container.NewTabItem("Discover Mods", browseNewTab),
	)

	launchButton := widget.NewButton("Launch", func() {
		dialog.ShowInformation("Launch", "Launching the game...", modWindow)
	})
	launchButton.Importance = widget.HighImportance

	mainLayout := container.NewBorder(nil, launchButton, leftSide, nil, tabs)
	modWindow.SetContent(mainLayout)
	modWindow.Resize(fyne.NewSize(1000, 600))
	modWindow.ShowAndRun()
}

func createModListItem(modWindow fyne.Window, title, description string) fyne.CanvasObject {
	checkbox := widget.NewCheck("", nil)
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle.Bold = true
	descriptionLabel := widget.NewLabel(description)

	downloadButton := widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {
		// ...
	})
	downloadButton.Importance = widget.HighImportance
	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		dialog.NewForm("Mod Settings", "Save", "Dismiss", []*widget.FormItem{
			{Text: "Setting 1", Widget: widget.NewEntry()},
			{Text: "Setting 2", Widget: widget.NewEntry()},
		}, func(b bool) {
			if b {
				dialog.ShowInformation("Settings", "Settings saved", modWindow)
			}
		}, modWindow).Show()
	})
	settingsButton.Importance = widget.WarningImportance
	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		confirmDeleteDialog()
	})
	deleteButton.Importance = widget.DangerImportance

	item := container.NewHBox(
		checkbox,
		container.NewVBox(
			titleLabel,
			descriptionLabel,
		),
		layout.NewSpacer(),
		downloadButton,
		settingsButton,
		deleteButton,
	)

	return item
}

func openAddGameModal(win fyne.Window) {
	gameSelect := widget.NewSelect([]string{"Game 1", "Game 2"}, func(value string) {})
	pathEntry := widget.NewEntry()
	pathButton := widget.NewButton("...", func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if file == nil {
				return
			}
			pathEntry.SetText(file.URI().Path())
		}, win)
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Select Game", Widget: gameSelect},
			{Text: "Data File Path", Widget: container.NewBorder(nil, nil, nil, pathButton, pathEntry)},
		},
	}

	dialog.ShowCustom("Add New Game", "Add", form, win)
}

func confirmDeleteDialog() {
	// ...
}
