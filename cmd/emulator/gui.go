package main

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
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/jonathangjertsen/toyboy/model"
	"github.com/jonathangjertsen/toyboy/plugin"
	"golang.org/x/image/font/basicfont"
)

type GUI struct {
	Config      *Config
	GB          *model.Gameboy
	JoypadState model.JoypadState

	ClockMeasurement  *plugin.ClockMeasurement
	GBFPSMeasurement  *plugin.ClockMeasurement
	GUIFPSMeasurement *plugin.ClockMeasurement
	Theme             material.Theme

	KeyboardControl *KeyboardControl

	SpeedInput      widget.Editor
	BreakPCInput    widget.Editor
	BreakIRInput    widget.Editor
	BreakXInput     widget.Editor
	BreakYInput     widget.Editor
	StepCyclesInput widget.Editor
	StepButton      widget.Clickable
	StartButton     widget.Clickable
	PauseButton     widget.Clickable
	ResetButton     widget.Clickable
	Registers       widget.Label

	VRAMScroll         widget.List
	WRAMScroll         widget.List
	HRAMScroll         widget.List
	OAMScroll          widget.List
	OAMAttrScroll      widget.List
	ProgramScroll      widget.List
	RewindBufferScroll widget.List
	DisassemblyScroll  widget.List
	RegistersScroll    widget.List
	PPUScroll          widget.List
	APUScroll          widget.List
	TimingScroll       widget.List

	LastFrameCPS    float64
	LastFrameGBFPS  float64
	LastFrameGUIFPS float64
	LastFrameT      time.Time
	StepCycles      uint64

	CoreDump model.CoreDump
	VRAM     []model.Data8
}

func NewGUI(config *Config) *GUI {
	gui := &GUI{
		Config: config,
		KeyboardControl: NewKeyboardControl(ButtonMapping{
			A:              "L",
			B:              "K",
			Up:             "W",
			Left:           "A",
			Down:           "S",
			Right:          "D",
			Start:          "M",
			Select:         "N",
			SOCDResolution: SOCDResolutionOppositeNeutral,
		}),
		Theme: material.NewTheme().WithPalette(material.Palette{
			Bg:         color.NRGBA{R: 245, G: 245, B: 245, A: 255},
			ContrastBg: color.NRGBA{R: 220, G: 220, B: 220, A: 255},
			Fg:         color.NRGBA{R: 45, G: 156, B: 219, A: 255},
			ContrastFg: color.NRGBA{R: 35, G: 146, B: 209, A: 255},
		}),
	}
	gui.initGameboy()
	return gui
}

func (gui *GUI) initGameboy() {
	gb := model.NewGameboy(&gui.Config.Model)

	gui.GB = gb

	// Measure clock freq
	gui.ClockMeasurement = plugin.NewClockMeasurement()
	gb.PHI.AttachDevice(func(c model.Cycle) {
		if c.Falling {
			gui.ClockMeasurement.Clocked()
		}
	})

	// Measure gameboy FPS
	gui.GBFPSMeasurement = plugin.NewClockMeasurement()
	gb.PPU.FrameClock.AttachDevice(func(c model.Cycle) {
		if c.Falling {
			gui.GBFPSMeasurement.Clocked()
		}
	})
}

func (gui *GUI) Run() {
	window := new(app.Window)
	window.Option(app.Title("toyboy"))
	window.Option(app.Size(unit.Dp(3000), unit.Dp(1920)))
	window.Option(app.Fullscreen.Option())
	gui.SpeedInput.Filter = "0123456789."
	gui.BreakXInput.Filter = "0123456789"
	gui.BreakYInput.Filter = "0123456789"
	gui.BreakPCInput.Filter = "0123456789abcdefABCDEF"
	gui.BreakIRInput.Filter = "0123456789abcdefABCDEF"
	gui.VRAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.WRAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.HRAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.ProgramScroll.List = layout.List{Axis: layout.Vertical}
	gui.RewindBufferScroll.List = layout.List{Axis: layout.Vertical}
	gui.DisassemblyScroll.List = layout.List{Axis: layout.Vertical}
	gui.OAMScroll.List = layout.List{Axis: layout.Vertical}
	gui.OAMAttrScroll.List = layout.List{Axis: layout.Vertical}
	gui.RegistersScroll.List = layout.List{Axis: layout.Vertical}
	gui.PPUScroll.List = layout.List{Axis: layout.Vertical}
	gui.APUScroll.List = layout.List{Axis: layout.Vertical}
	gui.TimingScroll.List = layout.List{Axis: layout.Vertical}
	gui.Reset()
	err := run(window, gui)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func (gui *GUI) Reset() {
	gui.SpeedInput.SetText(fmt.Sprintf("%.0f", gui.Config.Model.Clock.SpeedPercent))
	gui.BreakXInput.SetText("")
	gui.BreakYInput.SetText("")
	gui.initGameboy()
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
	// TODO: avg over more frames
	t := time.Now()
	frameElapsed := t.Sub(gui.LastFrameT)
	gui.LastFrameT = t
	gui.LastFrameGUIFPS = float64(time.Second / frameElapsed)

	// TODO: avg over more frames
	frames, duration := gui.GBFPSMeasurement.Stop()
	defer gui.GBFPSMeasurement.Start()
	if duration > 0 {
		gui.LastFrameGBFPS = float64((time.Duration(frames) * time.Second) / duration)
	} else {
		gui.LastFrameGBFPS = 0
	}

	// TODO: avg over more frames
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
	if gui.ResetButton.Clicked(gtx) {
		gui.Reset()
	}

	jps := gui.KeyboardControl.Frame(gtx)
	gui.JoypadState = jps
	gui.GB.Joypad.SetState(jps)
}

func (gui *GUI) label(name string) FC {
	return R(func(gtx C) D {
		lbl := material.Label(&gui.Theme, unit.Sp(14), name)
		lbl.Font.Typeface = "monospace"
		lbl.Font.Weight = font.Black
		lbl.Alignment = text.Start
		return lbl.Layout(gtx)
	})
}

func (gui *GUI) mem(f func(io.Writer), list *widget.List, config ConfigBox) W {
	return func(gtx C) D {
		buf := bytes.Buffer{}
		f(&buf)
		txt := buf.String()
		lines := strings.Split(txt, "\n")
		return layout.Stack{}.Layout(
			gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = int(unit.Dp(config.Width))
				gtx.Constraints.Max.X = int(unit.Dp(config.Width))
				gtx.Constraints.Min.Y = int(unit.Dp(config.Height))
				gtx.Constraints.Max.Y = int(unit.Dp(config.Height))
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
	gui.CoreDump = gui.GB.GetCoreDump()
	gui.VRAM = gui.CoreDump.VRAM.Bytes()

	Column(
		gtx,
		Spacer(25, 0),
		R(func(gtx C) D {
			var memviews []FC
			if gui.Config.GUI.VRAMMem.Box.Show {
				memviews = append(memviews, R(func(gtx C) D {
					return Column(
						gtx,
						gui.Box("VRAM",
							gui.mem(
								gui.CoreDump.PrintVRAM,
								&gui.VRAMScroll,
								gui.Config.GUI.VRAMMem.Box,
							)),
					)
				}))
			}
			if gui.Config.GUI.WRAMMem.Box.Show {
				memviews = append(memviews, R(func(gtx C) D {
					return Column(
						gtx,
						gui.Box("WRAM", gui.mem(
							gui.CoreDump.PrintWRAM,
							&gui.WRAMScroll,
							gui.Config.GUI.WRAMMem.Box,
						)),
					)
				}))
			}
			prog := gui.Config.GUI.ProgMem.Box.Show
			hram := gui.Config.GUI.HRAMMem.Box.Show
			oam := gui.Config.GUI.OAMMem.Box.Show
			if prog || hram || oam {
				memviews = append(memviews, R(func(gtx C) D {
					var smallViews []FC
					if prog {
						smallViews = append(smallViews,
							gui.Box("Program", gui.mem(
								gui.CoreDump.PrintProgram,
								&gui.ProgramScroll,
								gui.Config.GUI.ProgMem.Box,
							)),
						)
					}
					if hram {
						smallViews = append(smallViews,
							gui.Box("HRAM", gui.mem(
								gui.CoreDump.PrintHRAM,
								&gui.HRAMScroll,
								gui.Config.GUI.HRAMMem.Box,
							)),
						)
					}
					if oam {
						smallViews = append(smallViews,
							gui.Box("OAM", gui.mem(
								gui.CoreDump.PrintOAM,
								&gui.OAMScroll,
								gui.Config.GUI.OAMMem.Box,
							)),
						)
					}
					return Column(gtx, smallViews...)
				}))
			}
			if gui.Config.GUI.Rewind.Box.Show {
				memviews = append(memviews,
					R(func(gtx C) D {
						return Column(
							gtx,
							gui.Box("Last executed instructions", gui.mem(
								gui.CoreDump.PrintRewindBuffer,
								&gui.RewindBufferScroll,
								gui.Config.GUI.Rewind.Box,
							)),
						)
					}),
				)
			}
			if gui.Config.GUI.Disassembly.Box.Show {
				memviews = append(memviews,
					R(func(gtx C) D {
						return Column(
							gtx,
							gui.Box("Disassembly", gui.mem(
								gui.CoreDump.PrintDisassembly,
								&gui.DisassemblyScroll,
								gui.Config.GUI.Disassembly.Box,
							)),
						)
					}),
				)
			}
			timing := gui.Config.GUI.Timing.Box.Show
			registers := gui.Config.GUI.Registers.Box.Show
			apu := gui.Config.GUI.APU.Box.Show
			if timing || registers || apu {
				memviews = append(memviews,
					R(func(gtx C) D {
						var smallViews []FC
						if timing {
							smallViews = append(smallViews,
								gui.Box("Speed", gui.mem(func(w io.Writer) {
									fmt.Fprintf(w, "System clock:           %.0f\n", 4*gui.LastFrameCPS)
									fmt.Fprintf(w, "CPU clock:              %.0f\n", gui.LastFrameCPS)
									fmt.Fprintf(w, "Target emulation speed: %.2f%%\n", gui.Config.Model.Clock.SpeedPercent)
									fmt.Fprintf(w, "Actual emulation speed: %.2f%%\n", (100*4*gui.LastFrameCPS)/4194304)
									fmt.Fprintf(w, "Gameboy FPS:            %.0f\n", gui.LastFrameGBFPS)
									fmt.Fprintf(w, "GUI FPS:                %.0f\n", gui.LastFrameGUIFPS)
								}, &gui.TimingScroll, gui.Config.GUI.Timing.Box),
								))
						}
						if registers {
							smallViews = append(smallViews,
								gui.Box("Registers", gui.mem(
									gui.CoreDump.PrintRegs,
									&gui.RegistersScroll,
									gui.Config.GUI.Registers.Box,
								)),
							)
						}
						if apu {
							smallViews = append(smallViews,
								gui.Box("APU", gui.mem(
									gui.CoreDump.PrintAPU,
									&gui.APUScroll,
									gui.Config.GUI.APU.Box,
								)),
							)
						}
						return Column(gtx, smallViews...)
					}))
			}
			if gui.Config.GUI.PPU.Box.Show {
				memviews = append(memviews,
					R(func(gtx C) D {
						return Column(
							gtx,
							gui.Box(
								"PPU",
								gui.mem(gui.CoreDump.PrintPPU, &gui.PPUScroll, gui.Config.GUI.PPU.Box),
							),
						)
					}))
			}
			return Row(gtx, memviews...)
		}),
		gui.Graphics(),
		Spacer(25, 0),
	)
}

func (gui *GUI) Box(title string, content ...W) FC {
	return R(func(gtx C) D {
		out := make([]FC, 1+len(content))
		out[0] = gui.label(title)
		for i, w := range content {
			out[i+1] = R(w)
		}
		return Column(gtx, out...)
	})
}

func (gui *GUI) Button(clickable *widget.Clickable, text string) FC {
	return R(func(gtx C) D {
		return material.Button(&gui.Theme, clickable, text).Layout(gtx)
	})
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

func (gui *GUI) PlayArea() FC {
	return R(func(gtx C) D {
		confViewport := DefaultGridConfig
		confViewport.ShowGrid = false
		confViewport.ShowAddress = true
		confViewport.StartAddress = 0
		confViewport.BlockIncrement = 8
		confViewport.LineIncrement = 1
		confViewport.ShowOffsets = true
		confViewport.DecimalAddress = true

		confJoypad := DefaultGridConfig
		confJoypad.ShowAddress = false
		confJoypad.ShowGrid = false

		var views []FC

		if gui.Config.GUI.ViewPort.Box.Show {
			views = append(views,
				gui.label(fmt.Sprintf("Viewport (X=%d,Y=%d)", gui.CoreDump.PPU.BackgroundFetcher.X, (gui.CoreDump.PPU.Registers[model.AddrLY-model.AddrPPUBegin].Value)/8)),
				R(func(gtx C) D {
					vp := gui.GB.GetViewport()
					pixels := vp.Flatten()
					var highlights []Highlight
					if !gui.GB.CLK.Running.Load() {
						highlights = append(highlights, Highlight{
							BlockX: int(gui.CoreDump.PPU.PixelShifter.X) / 8,
							BlockY: int(gui.CoreDump.PPU.Registers[model.AddrLY-model.AddrPPUBegin].Value) / 8,
							Color:  color.RGBA{R: 255, A: 128},
						})
					}
					return gui.GBGraphics(
						gtx,
						160,
						144,
						pixels[:],
						4,
						confViewport,
						highlights,
					)
				}),
			)
		}
		if gui.Config.GUI.JoyPad.Box.Show {
			views = append(views, gui.Box("Joypad", func(gtx C) D {
				joypadViz := make([]model.Color, 8*6*8*4)
				highlights := [8]Highlight{}

				setButton := func(h *Highlight, x, y int, pressed bool, text string) {
					h.BlockX = x
					h.BlockY = y
					h.Text = text
					h.Font = basicfont.Face7x13
					h.TextColor = color.RGBA{R: 0, B: 0, G: 0, A: 255}
					if pressed {
						h.Color = color.RGBA{R: 255, B: 0, G: 0, A: 80}
					} else {
						h.Color = color.RGBA{R: 0, B: 0, G: 0, A: 40}
					}
				}
				setButton(&highlights[0], 1, 0, gui.JoypadState.Up, "Up")
				setButton(&highlights[1], 0, 1, gui.JoypadState.Left, "Lf")
				setButton(&highlights[2], 2, 1, gui.JoypadState.Right, "Ri")
				setButton(&highlights[3], 1, 2, gui.JoypadState.Down, "Dn")
				setButton(&highlights[4], 4, 1, gui.JoypadState.B, "B")
				setButton(&highlights[5], 5, 0, gui.JoypadState.A, "A")
				setButton(&highlights[6], 5, 3, gui.JoypadState.Start, "St")
				setButton(&highlights[7], 3, 3, gui.JoypadState.Select, "Se")

				return gui.GBGraphics(
					gtx,
					8*6,
					8*4,
					joypadViz,
					8,
					confJoypad,
					highlights[:],
				)
			}))
		}

		return Column(gtx, views...)
	})
}

func (gui *GUI) Graphics() FC {
	return R(func(gtx C) D {
		var views []FC
		viewPort := gui.Config.GUI.ViewPort.Box.Show
		joyPad := gui.Config.GUI.JoyPad.Box.Show
		if viewPort || joyPad {
			views = append(views, gui.PlayArea())
		}
		if gui.Config.GUI.TileData.Box.Show {
			views = append(views, gui.TileData())
		}
		if gui.Config.GUI.TileMap1.Box.Show {
			views = append(views, gui.TileMap1())
		}
		if gui.Config.GUI.TileMap2.Box.Show {
			views = append(views, gui.TileMap2())
		}
		oamBuf := gui.Config.GUI.OAMBuffer.Box.Show
		oamGraphics := gui.Config.GUI.OAMGraphics.Box.Show
		oamList := gui.Config.GUI.OAMList.Box.Show
		if oamBuf || oamGraphics || oamList {
			views = append(views, gui.OAMColumn())
		}
		if gui.Config.GUI.Debugger.Box.Show {
			views = append(views,
				gui.Debugger(),
			)
		}
		return Row(gtx, views...)
	})
}

func (gui *GUI) TileData() FC {
	return R(func(gtx C) D {
		return Column(
			gtx,
			gui.Box("Tile data", func(gtx C) D {
				return gui.GBGraphics(
					gtx,
					192,
					128,
					tiledata(gui.VRAM),
					2,
					DefaultGridConfig.WithMem(0x8000, 16),
					nil,
				)
			}),
		)
	})
}

func (gui *GUI) TileMap1() FC {
	return R(func(gtx C) D {
		return Column(
			gtx,
			gui.Box("Tile map 1", func(gtx C) D {
				return gui.GBGraphics(
					gtx,
					256,
					256,
					tilemap(gui.VRAM, 0x9800, gui.CoreDump.PPU.Registers[0].Value&model.Data8(1<<4) == 0),
					2,
					DefaultGridConfig.WithMem(0x9800, 1),
					nil,
				)
			}),
		)
	})
}

func (gui *GUI) TileMap2() FC {
	return R(func(gtx C) D {
		return Column(
			gtx,
			gui.Box("Tile map 2", func(gtx C) D {
				return gui.GBGraphics(
					gtx,
					256,
					256,
					tilemap(gui.VRAM, 0x9c00, gui.CoreDump.PPU.Registers[0].Value&model.Data8(1<<4) == 0),
					2,
					DefaultGridConfig.WithMem(0x9c00, 1),
					nil,
				)
			}),
		)
	})
}

func (gui *GUI) OAMColumn() FC {
	return R(func(gtx C) D {
		var views []FC
		if gui.Config.GUI.OAMBuffer.Box.Show {
			views = append(views,
				gui.Box("OAM buffer", func(gtx C) D {
					highlights := make([]Highlight, 10-gui.CoreDump.PPU.OAMBuffer.Level)
					for i := range highlights {
						highlights[i].BlockX = gui.CoreDump.PPU.OAMBuffer.Level + i
						highlights[i].Color = color.RGBA{R: 255, A: 10}
					}
					return gui.GBGraphics(
						gtx,
						80,
						8,
						oambuffer(gui.VRAM, gui.CoreDump.PPU.OAMBuffer),
						4,
						DefaultGridConfig,
						highlights,
					)
				}),
			)
		}
		if gui.Config.GUI.OAMGraphics.Box.Show {
			views = append(views,
				gui.Box("OAM", func(gtx C) D {
					return gui.GBGraphics(
						gtx,
						80,
						32,
						oam(gui.VRAM, gui.CoreDump.OAM.Bytes()),
						4,
						DefaultGridConfig.WithMem(model.AddrOAMBegin, 4),
						nil,
					)
				}),
			)
		}
		if gui.Config.GUI.OAMList.Box.Show {
			views = append(views,
				gui.Box("OAM attributes", gui.mem(
					gui.CoreDump.PrintOAMAttrs,
					&gui.OAMAttrScroll,
					gui.Config.GUI.OAMList.Box,
				)),
			)
		}

		return Column(gtx, views...)
	})
}

func (gui *GUI) Debugger() FC {
	return R(func(gtx C) D {
		return Column(
			gtx,
			gui.label(fmt.Sprintf("Running=%v", gui.GB.CLK.Running.Load())),
			gui.label(fmt.Sprintf("Clock: %d, Falling=%v", gui.CoreDump.Cycle.C, gui.CoreDump.Cycle.Falling)),
			R(func(gtx C) D {
				return Row(
					gtx,
					gui.Button(&gui.StepButton, "Step"),
					gui.Button(&gui.StartButton, "Run"),
					gui.Button(&gui.PauseButton, "Pause"),
					gui.Button(&gui.ResetButton, "Reset"),
				)
			}),
			R(func(gtx C) D {
				return Row(
					gtx,
					gui.Box("Step cycles", func(gtx C) D {
						return gui.NumberInput(gtx, &gui.StepCyclesInput, "cycles", func(text string) {
							cycles, err := strconv.ParseUint(text, 10, 64)
							if err == nil && cycles > 0 {
								gui.StepCycles = cycles
							}
						})
					}),
				)
			}),
			R(func(gtx C) D {
				return Row(
					gtx,
					gui.Box("Speed%", func(gtx C) D {
						return gui.NumberInput(gtx, &gui.SpeedInput, "speed", func(text string) {
							targetPercent, err := strconv.ParseFloat(text, 64)
							if err == nil && targetPercent > 0 && targetPercent < 10000 {
								if gui.Config.Model.Clock.SpeedPercent != targetPercent {
									gui.Config.Model.Clock.SpeedPercent = targetPercent
									gui.GB.CLK.SetSpeedPercent(gui.Config.Model.Clock.SpeedPercent)
									gui.Config.Save()
								}
							}
						})
					}),
				)
			}),
			R(func(gtx C) D {
				return Row(
					gtx,
					gui.Box("PC breakpoint", func(gtx C) D {
						return gui.NumberInput(gtx, &gui.BreakPCInput, "breakPC", func(text string) {
							breakPC, err := strconv.ParseInt(text, 16, 64)
							if err != nil {
								breakPC = -1
							}
							if gui.GB.Debug.BreakPC.Load() != breakPC {
								gui.GB.Debug.BreakPC.Store(breakPC)
							}
						})
					}),
				)
			}),
			R(func(gtx C) D {
				return Row(
					gtx,
					gui.Box("Opcode breakpoint", func(gtx C) D {
						return gui.NumberInput(gtx, &gui.BreakIRInput, "breakIR", func(text string) {
							breakIR, err := strconv.ParseInt(text, 16, 64)
							if err != nil {
								breakIR = -1
							}
							if gui.GB.Debug.BreakIR.Load() != breakIR {
								gui.GB.Debug.BreakIR.Store(breakIR)
							}
						})
					}),
				)
			}),
			R(func(gtx C) D {
				return Row(
					gtx,
					gui.Box("PPU breakpoint", func(gtx C) D {
						return gui.NumberInput(gtx, &gui.BreakXInput, "breakX", func(text string) {
							breakX, err := strconv.ParseInt(text, 10, 64)
							if err != nil {
								breakX = -1
							}
							if gui.GB.Debug.BreakX.Load() != breakX {
								gui.GB.Debug.BreakX.Store(breakX)
							}
						})
					}, func(gtx C) D {
						return gui.NumberInput(gtx, &gui.BreakYInput, "breakY", func(text string) {
							breakY, err := strconv.ParseInt(text, 10, 64)
							if err != nil {
								breakY = -1
							}
							if gui.GB.Debug.BreakY.Load() != breakY {
								gui.GB.Debug.BreakY.Store(breakY)
							}
						})
					}),
				)
			}),
		)
	})
}
