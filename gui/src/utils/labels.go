package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func TextLabel(text string) []fyne.CanvasObject {
	return []fyne.CanvasObject{widget.NewLabel(text)}
}
