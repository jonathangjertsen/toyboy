package gui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

	SpeedInput      widget.Editor
	StartButton     widget.Clickable
	PauseButton     widget.Clickable
	TimingGrid      component.GridState
	Registers       widget.Label
	VRAMScroll      widget.List
	HRAMScroll      widget.List
	OAMScroll       widget.List
	ProgramScroll   widget.List
	RegistersScroll widget.List
	PPUScroll       widget.List
	APUScroll       widget.List

	TargetPercent float64
	LastFrameCPS  float64
}

func New(gb *model.Gameboy) *GUI {
	cm := plugin.NewClockMeasurement()
	gb.PHI.AttachDevice(func(c model.Cycle) {
		if c.Falling {
			cm.Clocked()
		}
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
	window.Option(app.Size(unit.Dp(1440), unit.Dp(1080)))
	gui.SpeedInput.SetText("9999")
	gui.VRAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.HRAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.ProgramScroll.List = layout.List{Axis: layout.Vertical}
	gui.OAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.RegistersScroll.List = layout.List{Axis: layout.Vertical}
	gui.PPUScroll.List = layout.List{Axis: layout.Vertical}
	gui.APUScroll.List = layout.List{Axis: layout.Vertical}
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
		gui.LastFrameCPS = float64(cycles) * 1_000_000 / float64(us)
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

func (gui *GUI) memHead(name string) func(gtx C) D {
	return func(gtx C) D {
		lbl := material.Label(&gui.Theme, unit.Sp(14), name)
		lbl.Font.Typeface = "monospace"
		lbl.Font.Weight = font.Black
		lbl.Alignment = text.Start
		return lbl.Layout(gtx)
	}
}

func (gui *GUI) mem(f func(io.Writer), list *widget.List, height unit.Dp) func(gtx C) D {
	return func(gtx C) D {
		buf := bytes.Buffer{}
		f(&buf)
		txt := buf.String()
		lines := strings.Split(txt, "\n")
		return layout.Stack{}.Layout(
			gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = int(height)
				gtx.Constraints.Max.Y = int(height)
				return material.List(&gui.Theme, list).Layout(gtx, len(lines), func(gtx layout.Context, index int) layout.Dimensions {
					lbl := material.Label(&gui.Theme, unit.Sp(14), lines[index])
					lbl.Font.Typeface = "monospace"
					lbl.Alignment = text.Start
					return lbl.Layout(gtx)
				})
			}),
		)
	}
}

func (gui *GUI) Render(gtx C) {
	Column(
		gtx,
		Spacer(25, 0),
		Rigid(func(gtx C) D {
			cd := gui.GB.GetCoreDump()
			return Row(
				gtx,
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.memHead("VRAM")),
						Rigid(gui.mem(cd.PrintVRAM, &gui.VRAMScroll, unit.Dp(660))),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.memHead("Program")),
						Rigid(gui.mem(cd.PrintProgram, &gui.ProgramScroll, unit.Dp(200))),
						Rigid(gui.memHead("HRAM")),
						Rigid(gui.mem(cd.PrintHRAM, &gui.HRAMScroll, unit.Dp(200))),
						Rigid(gui.memHead("OAM")),
						Rigid(gui.mem(cd.PrintOAM, &gui.OAMScroll, unit.Dp(220))),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.memHead("Registers")),
						Rigid(gui.mem(cd.PrintRegs, &gui.RegistersScroll, unit.Dp(300))),
						Rigid(gui.memHead("PPU")),
						Rigid(gui.mem(cd.PrintPPU, &gui.PPUScroll, unit.Dp(300))),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.memHead("APU")),
						Rigid(gui.mem(cd.PrintAPU, &gui.APUScroll, unit.Dp(400))),
					)
				}),
			)
		}),
		Spacer(25, 0),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Spacer(25, 700),
				Rigid(func(gtx C) D {
					minSize := gtx.Dp(unit.Dp(100))

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
							return inset.Layout(gtx, func(gtx C) D {
								return headingLabel.Layout(gtx)
							})
						},
						func(gtx C, row, col int) D {
							switch row {
							case 0:
								switch col {
								case 0:
									dataLabel.Text = "System clock"
								case 1:
									dataLabel.Text = fmt.Sprintf("%.0f", 4*gui.LastFrameCPS)
								}
							case 1:
								switch col {
								case 0:
									dataLabel.Text = "CPU clock"
								case 1:
									dataLabel.Text = fmt.Sprintf("%.0f", gui.LastFrameCPS)
								}
							case 2:
								switch col {
								case 0:
									dataLabel.Text = "Target emulation speed"
								case 1:
									dataLabel.Text = fmt.Sprintf("%.2f%%", gui.TargetPercent)
								}
							case 3:
								switch col {
								case 0:
									dataLabel.Text = "Actual emulation speed"
								case 1:
									dataLabel.Text = fmt.Sprintf("%.2f%%", (100*4*gui.LastFrameCPS)/4194304)
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
					targetPercent, err := strconv.ParseFloat(text, 64)
					if err == nil && targetPercent > 0 && targetPercent < 10000 {
						if gui.TargetPercent != targetPercent {
							gui.TargetPercent = targetPercent
							gui.GB.CLK.SetFrequency(targetPercent * 4194304.0 / 100)
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
