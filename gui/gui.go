package gui

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/jonathangjertsen/toyboy/model"
)

type GUI struct {
	GB          *model.Gameboy
	Theme       material.Theme
	StartButton widget.Clickable
	PauseButton widget.Clickable
}

func New(gb *model.Gameboy) *GUI {
	return &GUI{
		GB: gb,
		Theme: material.NewTheme().WithPalette(material.Palette{
			Bg:         color.NRGBA{R: 245, G: 245, B: 245, A: 255},
			ContrastBg: color.NRGBA{R: 220, G: 220, B: 220, A: 255},
			Fg:         color.NRGBA{R: 45, G: 156, B: 219, A: 255},
			ContrastFg: color.NRGBA{R: 35, G: 146, B: 209, A: 255},
		}),
	}
}

func (gui *GUI) Run() {
	window := new(app.Window)
	window.Option(app.Title("toyboy"))
	window.Option(app.Size(unit.Dp(1080), unit.Dp(720)))
	err := run(window, gui)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func Main() {
	app.Main()
}

func run(window *app.Window, gui *GUI) error {
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			gui.Interactions(gtx)
			gui.Render(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (gui *GUI) Interactions(gtx C) {
	if gui.StartButton.Clicked(gtx) {
		gui.GB.PowerOn()
	}
	if gui.PauseButton.Clicked(gtx) {
		gui.GB.PowerOff()
	}
}

func (gui *GUI) Render(gtx C) {
	Column(
		gtx,
		Spacer(25, 0),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.StartButton, "Run")
				}),
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.PauseButton, "Pause")
				}),
			)
		}),
		Spacer(25, 0),
	)
}

func (gui *GUI) Button(gtx C, clickable *widget.Clickable, text string) D {
	return material.Button(&gui.Theme, clickable, text).Layout(gtx)
}
