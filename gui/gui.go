package gui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/jonathangjertsen/toyboy/model"
	"github.com/jonathangjertsen/toyboy/plugin"
)

type GUI struct {
	GB               *model.Gameboy
	ClockMeasurement *plugin.ClockMeasurement
	Theme            material.Theme

	SpeedInput  widget.Editor
	StartButton widget.Clickable
	PauseButton widget.Clickable
	TimingGrid  component.GridState

	TargetPercent uint64
	LastFrameCPS  uint64
}

func New(gb *model.Gameboy) *GUI {
	cm := plugin.NewClockMeasurement()
	gb.PHI.AddRiseCallback(func(c model.Cycle) {
		cm.Clocked()
	})
	return &GUI{
		GB:               gb,
		ClockMeasurement: cm,
		Theme: material.NewTheme().WithPalette(material.Palette{
			Bg:         color.NRGBA{R: 245, G: 245, B: 245, A: 255},
			ContrastBg: color.NRGBA{R: 220, G: 220, B: 220, A: 255},
			Fg:         color.NRGBA{R: 45, G: 156, B: 219, A: 255},
			ContrastFg: color.NRGBA{R: 35, G: 146, B: 209, A: 255},
		}),
		TargetPercent: 100,
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
	go func() {
		d := float64(time.Millisecond) * 16.667
		u := uint64(d)
		for range time.NewTicker(time.Duration(u)).C {
			window.Invalidate()
		}
	}()

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
	cycles, duration := gui.ClockMeasurement.Stop()
	defer gui.ClockMeasurement.Start()
	us := uint64(duration / time.Microsecond)

	if us > 0 {
		gui.LastFrameCPS = cycles * 1_000_000 / us
	} else {
		gui.LastFrameCPS = 0
	}

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
					minSize := gtx.Dp(unit.Dp(200))
					border := widget.Border{
						Color: color.NRGBA{A: 255},
						Width: unit.Dp(1),
					}

					inset := layout.UniformInset(unit.Dp(2))

					// Configure a label styled to be a heading.
					headingLabel := material.Body1(&gui.Theme, "")
					headingLabel.Font.Weight = font.Bold
					headingLabel.Alignment = text.Middle
					headingLabel.MaxLines = 1

					// Configure a label styled to be a data element.
					dataLabel := material.Body1(&gui.Theme, "")
					dataLabel.Font.Typeface = "monospace"
					dataLabel.MaxLines = 1
					dataLabel.Alignment = text.End

					var headingText = []string{"Data", "value"}

					orig := gtx.Constraints
					gtx.Constraints.Min = image.Point{}
					macro := op.Record(gtx.Ops)
					dims := inset.Layout(gtx, headingLabel.Layout)
					_ = macro.Stop()
					gtx.Constraints = orig
					return component.Table(
						&gui.Theme,
						&gui.TimingGrid,
					).Layout(
						gtx,
						4,
						2,
						func(axis layout.Axis, index, constraint int) int {
							widthUnit := max(int(float32(constraint)/3), minSize)
							switch axis {
							case layout.Horizontal:
								switch index {
								case 0, 1:
									return int(widthUnit)
								case 2, 3:
									return int(widthUnit / 2)
								default:
									return 0
								}
							default:
								return dims.Size.Y
							}
						},
						func(gtx C, col int) D {
							return border.Layout(gtx, func(gtx C) D {
								return inset.Layout(gtx, func(gtx C) D {
									headingLabel.Text = headingText[col]
									return headingLabel.Layout(gtx)
								})
							})
						},
						func(gtx C, row, col int) D {
							switch row {
							case 0:
								switch col {
								case 0:
									dataLabel.Text = "System clock"
								case 1:
									dataLabel.Text = fmt.Sprintf("%d", 4*gui.LastFrameCPS)
								}
							case 1:
								switch col {
								case 0:
									dataLabel.Text = "CPU clock"
								case 1:
									dataLabel.Text = fmt.Sprintf("%d", gui.LastFrameCPS)
								}
							case 2:
								switch col {
								case 0:
									dataLabel.Text = "Target emulation speed"
								case 1:
									dataLabel.Text = fmt.Sprintf("%d%%", gui.TargetPercent)
								}
							case 3:
								switch col {
								case 0:
									dataLabel.Text = "Actual emulation speed"
								case 1:
									dataLabel.Text = fmt.Sprintf("%d%%", (100*4*gui.LastFrameCPS)/4194304)
								}
							}
							return dataLabel.Layout(gtx)
						},
					)
				}),
			)
		}),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Rigid(func(gtx C) D {
					gui.SpeedInput.SingleLine = true
					gui.SpeedInput.Alignment = text.Middle
					text := gui.SpeedInput.Text()
					targetPercent, err := strconv.ParseUint(text, 10, 64)
					if err == nil && targetPercent > 0 && targetPercent < 1000 {
						if gui.TargetPercent != targetPercent {
							gui.TargetPercent = targetPercent
							gui.GB.CLK.SetFrequency(float64(targetPercent) * 4194304.0 / 100)
						}
					} else {
						// Ignore
					}
					return material.Editor(&gui.Theme, &gui.SpeedInput, "speed").Layout(gtx)
				}),
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

func fmtClockSpeed(hz uint64) string {
	if hz == 0 {
		return "0"
	}
	if hz < 1200 {
		return fmt.Sprintf("%d Hz", hz)
	}
	if hz < 1_200_000 {
		return fmt.Sprintf("%.1f KHz", float64(hz)/1_000)
	}
	return fmt.Sprintf("%.1f MHz", float64(hz)/1_000_000)
}
