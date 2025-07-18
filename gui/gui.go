package gui

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/jonathangjertsen/toyboy/model"
	"github.com/jonathangjertsen/toyboy/plugin"
)

type GUI struct {
	Config model.HWConfig
	GB     *model.Gameboy

	ClockMeasurement *plugin.ClockMeasurement
	Theme            material.Theme

	KeyboardControl *KeyboardControl

	SpeedInput      widget.Editor
	BreakXInput     widget.Editor
	BreakYInput     widget.Editor
	StepCyclesInput widget.Editor
	StepButton      widget.Clickable
	StartButton     widget.Clickable
	PauseButton     widget.Clickable
	SoftResetButton widget.Clickable
	Registers       widget.Label

	VRAMScroll      widget.List
	HRAMScroll      widget.List
	OAMScroll       widget.List
	OAMAttrScroll   widget.List
	ProgramScroll   widget.List
	RegistersScroll widget.List
	PPUScroll       widget.List
	APUScroll       widget.List
	TimingScroll    widget.List

	TargetPercent float64
	LastFrameCPS  float64
	StepCycles    uint64
}

func New(config model.HWConfig) *GUI {
	gui := &GUI{
		Config: config,
		KeyboardControl: NewKeyboardControl(ButtonMapping{
			A:              "L",
			B:              "K",
			Up:             "W",
			Left:           "A",
			Down:           "S",
			Right:          "D",
			Start:          key.NameEnter,
			Select:         key.NameSpace,
			SOCDResolution: SOCDResolutionOppositeNeutral,
		}),
		Theme: material.NewTheme().WithPalette(material.Palette{
			Bg:         color.NRGBA{R: 245, G: 245, B: 245, A: 255},
			ContrastBg: color.NRGBA{R: 220, G: 220, B: 220, A: 255},
			Fg:         color.NRGBA{R: 45, G: 156, B: 219, A: 255},
			ContrastFg: color.NRGBA{R: 35, G: 146, B: 209, A: 255},
		}),
		TargetPercent: 100,
	}
	gui.initGameboy()
	return gui
}

func (gui *GUI) initGameboy() {
	gb := model.NewGameboy(gui.Config)
	f, err := os.ReadFile("assets/cartridges/unbricked.gb")
	if err != nil {
		panic(fmt.Sprintf("failed to load cartridge: %v", err))
	} else if len(f) != 0x8000 {
		panic(fmt.Sprintf("len(bootrom)=%d", len(f)))
	}
	copy(gb.CartridgeSlot.Data, f)

	cm := plugin.NewClockMeasurement()
	gb.PHI.AttachDevice(func(c model.Cycle) {
		if c.Falling {
			cm.Clocked()
		}
	})

	gui.GB = gb
	gui.ClockMeasurement = cm
}

func (gui *GUI) Run() {
	window := new(app.Window)
	window.Option(app.Title("toyboy"))
	window.Option(app.Size(unit.Dp(3000), unit.Dp(1920)))
	window.Option(app.Fullscreen.Option())
	gui.SpeedInput.SetText("100")
	gui.BreakXInput.SetText("")
	gui.BreakYInput.SetText("")
	gui.SpeedInput.Filter = "0123456789."
	gui.BreakXInput.Filter = "0123456789"
	gui.BreakYInput.Filter = "0123456789"
	gui.VRAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.HRAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.ProgramScroll.List = layout.List{Axis: layout.Vertical}
	gui.OAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.OAMAttrScroll.List = layout.List{Axis: layout.Vertical}
	gui.RegistersScroll.List = layout.List{Axis: layout.Vertical}
	gui.PPUScroll.List = layout.List{Axis: layout.Vertical}
	gui.APUScroll.List = layout.List{Axis: layout.Vertical}
	gui.TimingScroll.List = layout.List{Axis: layout.Vertical}
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
		gui.GB.Start()
	}
	if gui.PauseButton.Clicked(gtx) {
		gui.GB.Pause()
	}
	if gui.StepButton.Clicked(gtx) {
		for range gui.StepCycles {
			gui.GB.Step()
		}
	}
	if gui.SoftResetButton.Clicked(gtx) {
		gui.GB.SoftReset()
	}

	jps := gui.KeyboardControl.Frame(gtx)
	gui.GB.Joypad.SetState(jps)
}

func (gui *GUI) label(name string) func(gtx C) D {
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
	cd := gui.GB.GetCoreDump()
	vramBytes := cd.VRAM.Bytes()

	Column(
		gtx,
		Spacer(25, 0),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("VRAM")),
						Rigid(gui.mem(cd.PrintVRAM, &gui.VRAMScroll, unit.Dp(660))),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("Program")),
						Rigid(gui.mem(cd.PrintProgram, &gui.ProgramScroll, unit.Dp(200))),
						Rigid(gui.label("HRAM")),
						Rigid(gui.mem(cd.PrintHRAM, &gui.HRAMScroll, unit.Dp(200))),
						Rigid(gui.label("OAM")),
						Rigid(gui.mem(cd.PrintOAM, &gui.OAMScroll, unit.Dp(220))),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("Registers")),
						Rigid(gui.mem(cd.PrintRegs, &gui.RegistersScroll, unit.Dp(300))),
						Rigid(gui.label("APU")),
						Rigid(gui.mem(cd.PrintAPU, &gui.APUScroll, unit.Dp(300))),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("PPU")),
						Rigid(gui.mem(cd.PrintPPU, &gui.PPUScroll, unit.Dp(800))),
					)
				}),
			)
		}),
		Spacer(25, 0),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Spacer(25, 700),
				Rigid(
					gui.mem(func(w io.Writer) {
						fmt.Fprintf(w, "System clock:           %.0f\n", 4*gui.LastFrameCPS)
						fmt.Fprintf(w, "CPU clock:              %.0f\n", gui.LastFrameCPS)
						fmt.Fprintf(w, "Target emulation speed: %.2f%%\n", gui.TargetPercent)
						fmt.Fprintf(w, "Actual emulation speed: %.2f%%\n", (100*4*gui.LastFrameCPS)/4194304)
					}, &gui.TimingScroll, unit.Dp(400)),
				),
			)
		}),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Rigid(func(gtx C) D {
					gc := DefaultGridConfig
					gc.Show = false
					gc.ShowAddress = true
					gc.StartAddress = 0
					gc.BlockIncrement = 8
					gc.LineIncrement = 1
					gc.ShowOffsets = true
					gc.DecimalAddress = true
					return Column(
						gtx,
						Rigid(gui.label(fmt.Sprintf("Viewport (X=%d,Y=%d)", cd.PPU.BackgroundFetcher.X, (cd.PPU.Registers[model.AddrLY-model.AddrPPUBegin].Value)/8))),
						Rigid(func(gtx C) D {
							vp := gui.GB.GetViewport()
							pixels := vp.Flatten()
							return gui.GBGraphics(
								gtx,
								160,
								144,
								pixels[:],
								4,
								gc,
								nil,
								/*
									[]Highlight{
										{
											BlockX: int(cd.PPU.PixelShifter.X) / 8,
											BlockY: int(cd.PPU.Registers[model.AddrLY-model.AddrPPUBegin].Value) / 8,
											Color:  color.RGBA{R: 255, A: 128},
										},
									},
								*/
							)
						}),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("Tile data")),
						Rigid(func(gtx C) D {
							return gui.GBGraphics(
								gtx,
								192,
								128,
								tiledata(vramBytes),
								2,
								DefaultGridConfig.WithMem(0x8000, 16),
								nil,
							)
						}),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("Tile map 1")),
						Rigid(func(gtx C) D {
							return gui.GBGraphics(
								gtx,
								256,
								256,
								tilemap(vramBytes, 0x9800, cd.PPU.Registers[0].Value&uint8(1<<4) == 0),
								2,
								DefaultGridConfig.WithMem(0x9800, 1),
								nil,
							)
						}),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("Tile map 2")),
						Rigid(func(gtx C) D {
							return gui.GBGraphics(
								gtx,
								256,
								256,
								tilemap(vramBytes, 0x9c00, cd.PPU.Registers[0].Value&uint8(1<<4) == 0),
								2,
								DefaultGridConfig.WithMem(0x9c00, 1),
								nil,
							)
						}),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label("OAM buffer")),
						Rigid(func(gtx C) D {
							highlights := make([]Highlight, 10-cd.PPU.OAMBuffer.Level)
							for i := range highlights {
								highlights[i].BlockX = cd.PPU.OAMBuffer.Level + i
								highlights[i].Color = color.RGBA{R: 255, A: 10}
							}
							return gui.GBGraphics(
								gtx,
								80,
								8,
								oambuffer(vramBytes, cd.PPU.OAMBuffer),
								4,
								DefaultGridConfig,
								highlights,
							)
						}),
						Rigid(gui.label("OAM")),
						Rigid(func(gtx C) D {
							return gui.GBGraphics(
								gtx,
								80,
								32,
								oam(vramBytes, cd.OAM.Bytes()),
								4,
								DefaultGridConfig.WithMem(model.AddrOAMBegin, 4),
								nil,
							)
						}),
						Rigid(gui.mem(cd.PrintOAMAttrs, &gui.OAMAttrScroll, unit.Dp(100))),
					)
				}),
			)
		}),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Rigid(gui.label(fmt.Sprintf("Running=%v", gui.GB.Running))),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(func(gtx C) D {
							return Row(
								gtx,
								Rigid(gui.label("Speed%")),
								Rigid(func(gtx layout.Context) layout.Dimensions {
									return gui.NumberInput(gtx, &gui.SpeedInput, "speed", func(text string) {
										targetPercent, err := strconv.ParseFloat(text, 64)
										if err == nil && targetPercent > 0 && targetPercent < 10000 {
											if gui.TargetPercent != targetPercent {
												gui.TargetPercent = targetPercent
												gui.Config.SystemClock.Frequency = targetPercent * 4194304.0 / 100
												gui.GB.CLK.SetFrequency(gui.Config.SystemClock.Frequency)
											}
										}
									})
								}),
							)
						}),
						Rigid(func(gtx C) D {
							return Row(
								gtx,
								Rigid(gui.label("PPU breakpoint")),
								Rigid(func(gtx C) D {
									return gui.NumberInput(gtx, &gui.BreakXInput, "breakX", func(text string) {
										breakX, err := strconv.ParseInt(text, 10, 64)
										if err != nil {
											breakX = -1
										}
										if gui.GB.Debugger.BreakX.Load() != breakX {
											gui.GB.Debugger.BreakX.Store(breakX)
										}
									})
								}),
								Rigid(func(gtx C) D {
									return gui.NumberInput(gtx, &gui.BreakYInput, "breakY", func(text string) {
										breakY, err := strconv.ParseInt(text, 10, 64)
										if err != nil {
											breakY = -1
										}
										if gui.GB.Debugger.BreakY.Load() != breakY {
											gui.GB.Debugger.BreakY.Store(breakY)
										}
									})
								}),
							)
						}),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.label(fmt.Sprintf("Clock: %d, Falling=%v", cd.Cycle.C, cd.Cycle.Falling))),
						Rigid(func(gtx C) D {
							return gui.Button(gtx, &gui.StepButton, "Step")
						}),
						Rigid(func(gtx C) D {
							return Row(
								gtx,
								Rigid(gui.label("Step cycles")),
								Rigid(func(gtx C) D {
									return gui.NumberInput(gtx, &gui.StepCyclesInput, "cycles", func(text string) {
										cycles, err := strconv.ParseUint(text, 10, 64)
										if err == nil && cycles > 0 {
											gui.StepCycles = cycles
										}
									})
								}),
							)
						}),
					)
				}),
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.StartButton, "Run")
				}),
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.PauseButton, "Pause")
				}),
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.SoftResetButton, "SoftReset")
				}),
			)
		}),
		Spacer(25, 0),
	)
}

func (gui *GUI) Button(gtx C, clickable *widget.Clickable, text string) D {
	return material.Button(&gui.Theme, clickable, text).Layout(gtx)
}

func (gui *GUI) NumberInput(gtx C, editor *widget.Editor, placeholder string, f func(text string)) D {
	editor.SingleLine = true
	editor.Alignment = text.Middle
	f(editor.Text())
	return layout.Inset{
		Top:    unit.Dp(0),
		Right:  unit.Dp(170),
		Bottom: unit.Dp(40),
		Left:   unit.Dp(170),
	}.Layout(gtx, func(gtx C) D {
		return widget.Border{
			Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
			CornerRadius: unit.Dp(3),
			Width:        unit.Dp(2),
		}.Layout(gtx, material.Editor(&gui.Theme, editor, placeholder).Layout)
	})
}
