package main

import (
	"image/color"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow()
		w.Option(app.Title("mhmods"))
		err := run(w)
		if err != nil {
			if strings.Contains(err.Error(), "wl_display_connect failed") {
				log.Fatal("wl_display_connect failed - does your machine support GUIs?")
			}
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}


func run(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			title := material.Body1(th, "todo")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon
			title.Alignment = text.Middle

			title.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}