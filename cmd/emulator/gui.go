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

func (gui *GUI) mem(f func(io.Writer), list *widget.List, height, width unit.Dp) func(gtx C) D {
	return func(gtx C) D {
		buf := bytes.Buffer{}
		f(&buf)
		txt := buf.String()
		lines := strings.Split(txt, "\n")
		return layout.Stack{}.Layout(
			gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = int(width)
				gtx.Constraints.Max.X = int(width)
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
	gui.CoreDump = gui.GB.GetCoreDump()
	gui.VRAM = gui.CoreDump.VRAM.Bytes()

	Column(
		gtx,
		Spacer(25, 0),
		R(func(gtx C) D {
			return Row(
				gtx,
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.label("VRAM"),
						R(gui.mem(gui.CoreDump.PrintVRAM, &gui.VRAMScroll, unit.Dp(660), unit.Dp(600))),
					)
				}),
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.label("WRAM"),
						R(gui.mem(gui.CoreDump.PrintWRAM, &gui.WRAMScroll, unit.Dp(660), unit.Dp(600))),
					)
				}),
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.label("Program"),
						R(gui.mem(gui.CoreDump.PrintProgram, &gui.ProgramScroll, unit.Dp(200), unit.Dp(600))),
						gui.label("HRAM"),
						R(gui.mem(gui.CoreDump.PrintHRAM, &gui.HRAMScroll, unit.Dp(200), unit.Dp(600))),
						gui.label("OAM"),
						R(gui.mem(gui.CoreDump.PrintOAM, &gui.OAMScroll, unit.Dp(220), unit.Dp(600))),
					)
				}),
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.label("Last executed instructions"),
						R(gui.mem(gui.CoreDump.PrintRewindBuffer, &gui.RewindBufferScroll, unit.Dp(300), unit.Dp(400))),
					)
				}),
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.label("Disassembly"),
						R(gui.mem(gui.CoreDump.PrintDisassembly, &gui.DisassemblyScroll, unit.Dp(800), unit.Dp(500))),
					)
				}),
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.label("Speed"),
						R(
							gui.mem(func(w io.Writer) {
								fmt.Fprintf(w, "System clock:           %.0f\n", 4*gui.LastFrameCPS)
								fmt.Fprintf(w, "CPU clock:              %.0f\n", gui.LastFrameCPS)
								fmt.Fprintf(w, "Target emulation speed: %.2f%%\n", gui.Config.Model.Clock.SpeedPercent)
								fmt.Fprintf(w, "Actual emulation speed: %.2f%%\n", (100*4*gui.LastFrameCPS)/4194304)
								fmt.Fprintf(w, "Gameboy FPS:            %.0f\n", gui.LastFrameGBFPS)
								fmt.Fprintf(w, "GUI FPS:                %.0f\n", gui.LastFrameGUIFPS)
							}, &gui.TimingScroll, unit.Dp(100), unit.Dp(400)),
						),
						gui.label("Registers"),
						R(gui.mem(gui.CoreDump.PrintRegs, &gui.RegistersScroll, unit.Dp(300), unit.Dp(200))),
						gui.label("APU"),
						R(gui.mem(gui.CoreDump.PrintAPU, &gui.APUScroll, unit.Dp(300), unit.Dp(200))),
					)
				}),
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.label("PPU"),
						R(gui.mem(gui.CoreDump.PrintPPU, &gui.PPUScroll, unit.Dp(800), unit.Dp(200))),
					)
				}),
			)
		}),
		gui.Graphics(),
		R(func(gtx C) D {
			return Row(
				gtx,
				gui.label(fmt.Sprintf("Running=%v", gui.GB.CLK.Running.Load())),
				R(func(gtx C) D {
					return Column(
						gtx,
					)
				}),
				gui.Debugger(),
			)
		}),
		Spacer(25, 0),
	)
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
		confViewport.Show = false
		confViewport.ShowAddress = true
		confViewport.StartAddress = 0
		confViewport.BlockIncrement = 8
		confViewport.LineIncrement = 1
		confViewport.ShowOffsets = true
		confViewport.DecimalAddress = true

		confJoypad := DefaultGridConfig
		confJoypad.ShowAddress = false
		confJoypad.Show = false

		return Column(
			gtx,
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
			gui.label("Joypad"),
			R(func(gtx C) D {
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
			}),
		)
	})
}

func (gui *GUI) Graphics() FC {
	return R(func(gtx C) D {
		return Row(
			gtx,
			gui.PlayArea(),
			gui.TileData(),
			gui.TileMap1(),
			gui.TileMap2(),
			gui.OAMColumn(),
		)
	})
}

func (gui *GUI) TileData() FC {
	return R(func(gtx C) D {
		return Column(
			gtx,
			gui.label("Tile data"),
			R(func(gtx C) D {
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
			gui.label("Tile map 1"),
			R(func(gtx C) D {
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
			gui.label("Tile map 2"),
			R(func(gtx C) D {
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
		return Column(
			gtx,
			gui.label("OAM buffer"),
			R(func(gtx C) D {
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
			gui.label("OAM"),
			R(func(gtx C) D {
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
			R(gui.mem(gui.CoreDump.PrintOAMAttrs, &gui.OAMAttrScroll, unit.Dp(100), unit.Dp(200))),
		)
	})
}

func (gui *GUI) Debugger() FC {
	return R(func(gtx C) D {
		return Column(
			gtx,
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
					gui.label("Step cycles"),
					R(func(gtx C) D {
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
					gui.label("Speed%"),
					R(func(gtx layout.Context) layout.Dimensions {
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
					gui.label("PC breakpoint"),
					R(func(gtx C) D {
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
					gui.label("Opcode breakpoint"),
					R(func(gtx C) D {
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
					gui.label("PPU breakpoint"),
					R(func(gtx C) D {
						return gui.NumberInput(gtx, &gui.BreakXInput, "breakX", func(text string) {
							breakX, err := strconv.ParseInt(text, 10, 64)
							if err != nil {
								breakX = -1
							}
							if gui.GB.Debug.BreakX.Load() != breakX {
								gui.GB.Debug.BreakX.Store(breakX)
							}
						})
					}),
					R(func(gtx C) D {
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
