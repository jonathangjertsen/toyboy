package main

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
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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

	ToggleButtons map[string]*widget.Clickable

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
		KeyboardControl: NewKeyboardControl(),
		Theme: material.NewTheme().WithPalette(material.Palette{
			Bg:         color.NRGBA{R: 245, G: 245, B: 245, A: 255},
			ContrastBg: color.NRGBA{R: 220, G: 220, B: 220, A: 255},
			Fg:         color.NRGBA{R: 45, G: 156, B: 219, A: 255},
			ContrastFg: color.NRGBA{R: 35, G: 146, B: 209, A: 255},
		}),
		ToggleButtons: make(map[string]*widget.Clickable),
	}
	return gui
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
	us := 
	if us > 0 {
		gui.LastFrameCPS = float64(cycles) * 1_000_000 / float64(uint64(duration / time.Microsecond))
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
		lbl := material.Label(&gui.Theme, unit.Sp(20), name)
		lbl.Font.Typeface = "monospace"
		lbl.Font.Weight = font.Black
		lbl.Alignment = text.Start
		return lbl.Layout(gtx)
	})
}

func (gui *GUI) mem(f func(io.Writer), list *widget.List) W {
	return func(gtx C) D {
		buf := bytes.Buffer{}
		f(&buf)
		txt := buf.String()
		lines := strings.Split(txt, "\n")
		return material.List(&gui.Theme, list).Layout(gtx, len(lines), func(gtx layout.Context, index int) layout.Dimensions {
			lbl := material.Label(&gui.Theme, unit.Sp(14), lines[index])
			lbl.Font.Typeface = "monospace"
			lbl.Alignment = text.Start
			return lbl.Layout(gtx)
		})
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
			memviews = append(memviews, R(func(gtx C) D {
				return Column(
					gtx,
					gui.Box("VRAM",
						&gui.Config.GUI.VRAMMem.Box,
						gui.mem(
							gui.CoreDump.PrintVRAM,
							&gui.VRAMScroll,
						)),
				)
			}))
			memviews = append(memviews, R(func(gtx C) D {
				return Column(
					gtx,
					gui.Box("WRAM",
						&gui.Config.GUI.WRAMMem.Box,
						gui.mem(
							gui.CoreDump.PrintWRAM,
							&gui.WRAMScroll,
						)),
				)
			}))
			memviews = append(memviews, R(func(gtx C) D {
				var smallViews []FC
				smallViews = append(smallViews,
					gui.Box("Program",
						&gui.Config.GUI.ProgMem.Box,
						gui.mem(
							gui.CoreDump.PrintProgram,
							&gui.ProgramScroll,
						)),
				)
				smallViews = append(smallViews,
					gui.Box("HRAM",
						&gui.Config.GUI.HRAMMem.Box,
						gui.mem(
							gui.CoreDump.PrintHRAM,
							&gui.HRAMScroll,
						)),
				)
				smallViews = append(smallViews,
					gui.Box("OAM",
						&gui.Config.GUI.OAMMem.Box,
						gui.mem(
							gui.CoreDump.PrintOAM,
							&gui.OAMScroll,
						)),
				)
				return Column(gtx, smallViews...)
			}))
			memviews = append(memviews,
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.Box("Execution log",
							&gui.Config.GUI.Rewind.Box,
							gui.mem(
								gui.CoreDump.PrintRewindBuffer,
								&gui.RewindBufferScroll,
							)),
					)
				}),
			)
			memviews = append(memviews,
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.Box("Disassembly",
							&gui.Config.GUI.Disassembly.Box,
							gui.mem(
								gui.CoreDump.PrintDisassembly,
								&gui.DisassemblyScroll,
							)),
					)
				}),
			)
			memviews = append(memviews,
				R(func(gtx C) D {
					var smallViews []FC
					smallViews = append(smallViews,
						gui.Box("Speed", &gui.Config.GUI.Timing.Box,
							gui.mem(func(w io.Writer) {
								fmt.Fprintf(w, "System clock:           %.0f\n", 4*gui.LastFrameCPS)
								fmt.Fprintf(w, "CPU clock:              %.0f\n", gui.LastFrameCPS)
								fmt.Fprintf(w, "Target emulation speed: %.2f%%\n", gui.Config.Model.Clock.SpeedPercent)
								fmt.Fprintf(w, "Actual emulation speed: %.2f%%\n", (100*4*gui.LastFrameCPS)/4194304)
								fmt.Fprintf(w, "Gameboy FPS:            %.0f\n", gui.LastFrameGBFPS)
								fmt.Fprintf(w, "GUI FPS:                %.0f\n", gui.LastFrameGUIFPS)
							}, &gui.TimingScroll),
						))
					smallViews = append(smallViews, R(func(gtx C) D {
						return Row(
							gtx,
							gui.Box("Registers",
								&gui.Config.GUI.Registers.Box,
								gui.mem(
									gui.CoreDump.PrintRegs,
									&gui.RegistersScroll,
								),
							),
							gui.Box("APU",
								&gui.Config.GUI.APU.Box,
								gui.mem(
									gui.CoreDump.PrintAPU,
									&gui.APUScroll,
								),
							),
						)
					}))
					return Column(gtx, smallViews...)
				}))
			memviews = append(memviews,
				R(func(gtx C) D {
					return Column(
						gtx,
						gui.Box(
							"PPU",
							&gui.Config.GUI.PPU.Box,
							gui.mem(gui.CoreDump.PrintPPU, &gui.PPUScroll),
						),
					)
				}))
			return Row(gtx, memviews...)
		}),
		gui.PlayArea(),
		gui.Graphics(),
		Spacer(25, 0),
	)
}
func (gui *GUI) Box(title string, config *ConfigBox, content ...W) FC {
	// Get or create toggle button for this box
	toggleBtn, ok := gui.ToggleButtons[title]
	if !ok {
		toggleBtn = new(widget.Clickable)
		gui.ToggleButtons[title] = toggleBtn
	}

	return R(func(gtx C) D {
		// Check if toggle button was clicked
		if toggleBtn.Clicked(gtx) {
			config.Show = !config.Show
			gui.Config.Save()
		}

		gtx.Constraints.Min.X = int(unit.Dp(config.Width))
		gtx.Constraints.Max.X = int(unit.Dp(config.Width))

		if config.Show {
			gtx.Constraints.Min.Y = int(unit.Dp(config.Height))
			gtx.Constraints.Max.Y = int(unit.Dp(config.Height))
		} else {
			gtx.Constraints.Min.Y = int(unit.Dp(70))
			gtx.Constraints.Max.Y = int(unit.Dp(70))
		}

		borderOuter := layout.UniformInset(8)
		return borderOuter.Layout(gtx, func(gtx C) D {
			borderInner := layout.UniformInset(8)
			inner := borderInner.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Title row with toggle button
					layout.Rigid(func(gtx C) D {
						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							gui.label(title),
							layout.Rigid(func(gtx C) D {
								buttonText := "hide"
								if !config.Show {
									buttonText = "show"
								}

								btn := material.Button(&gui.Theme, toggleBtn, buttonText)
								btn.Background = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
								btn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}

								gtx.Constraints.Min = image.Point{}
								gtx.Constraints.Max = image.Pt(80, 40)

								return btn.Layout(gtx)
							}),
						)

					}),
					// Content (only show if not collapsed)
					layout.Rigid(func(gtx C) D {
						if !config.Show {
							return D{}
						}

						// Add spacer between title and content
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							Spacer(10, 1),
							layout.Rigid(func(gtx C) D {
								// Layout all content widgets vertically
								widgets := make([]layout.FlexChild, len(content))
								for i, w := range content {
									w := w // capture loop variable
									widgets[i] = layout.Rigid(func(gtx C) D {
										return w(gtx)
									})
								}
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx, widgets...)
							}),
						)
					}),
				)
			})

			// Draw border around the whole box
			paint.FillShape(gtx.Ops, color.NRGBA{A: 128}, clip.Stroke{
				Path:  clip.Rect{Max: inner.Size}.Path(),
				Width: 2,
			}.Op())
			return inner
		})
	})
}

func (gui *GUI) Button(clickable *widget.Clickable, text string) FC {
	return R(func(gtx C) D {
		return material.Button(&gui.Theme, clickable, text).Layout(gtx)
	})
}
func (gui *GUI) NumberInput(editor *widget.Editor, what string, placeholder string, f func(text string)) FC {
	editor.SingleLine = true
	editor.Alignment = text.Middle
	f(editor.Text())
	return R(func(gtx C) D {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceBetween,
			Alignment: layout.Middle,
		}.Layout(gtx,
			R(func(gtx C) D {
				inset := layout.Inset{Right: unit.Dp(8)}
				return inset.Layout(gtx, func(gtx C) D {
					lbl := material.Label(&gui.Theme, unit.Sp(14), what)
					lbl.Font.Typeface = "monospace"
					lbl.Font.Weight = font.Black
					lbl.Alignment = text.Start
					return lbl.Layout(gtx)
				})
			}),
			R(func(gtx C) D {
				inset := layout.Inset{Left: unit.Dp(8)}
				return inset.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Dp(unit.Dp(200))
					gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))
					gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(20))
					gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(20))
					return widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}.Layout(gtx, material.Editor(&gui.Theme, editor, placeholder).Layout)
				})
			}),
		)
	})
}

func (gui *GUI) PlayArea() FC {
	return R(func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx C) D {
			var views []FC

			views = append(views,
				gui.Box(
					"Viewport",
					&gui.Config.GUI.ViewPort.Box,
					func(gtx C) D {
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
							gui.Config.GUI.ViewPort.Graphics,
							highlights,
						)
					},
				),
			)
			views = append(views, gui.Box("Joypad",
				&gui.Config.GUI.JoyPad.Box,
				func(gtx C) D {
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
						gui.Config.GUI.JoyPad.Graphics,
						highlights[:],
					)
				}))

			return Column(gtx, views...)
		})
	})
}

func (gui *GUI) Graphics() FC {
	return R(func(gtx C) D {
		var views []FC
		views = append(views, gui.TileData())
		views = append(views, gui.TileMap1())
		views = append(views, gui.TileMap2())
		views = append(views, gui.OAMColumn())
		views = append(views, gui.Debugger())
		return Row(gtx, views...)
	})
}

func (gui *GUI) TileData() FC {
	return R(func(gtx C) D {
		return Column(
			gtx,
			gui.Box("Tile data", &gui.Config.GUI.TileData.Box,
				func(gtx C) D {
					return gui.GBGraphics(
						gtx,
						192,
						128,
						tiledata(gui.VRAM),
						gui.Config.GUI.TileData.Graphics,
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
			gui.Box("Tile map 1", &gui.Config.GUI.TileMap1.Box,
				func(gtx C) D {
					return gui.GBGraphics(
						gtx,
						256,
						256,
						tilemap(gui.VRAM, 0x9800, gui.CoreDump.PPU.Registers[0].Value&model.Data8(1<<4) == 0),
						gui.Config.GUI.TileMap1.Graphics,
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
			gui.Box("Tile map 2", &gui.Config.GUI.TileMap2.Box,
				func(gtx C) D {
					return gui.GBGraphics(
						gtx,
						256,
						256,
						tilemap(gui.VRAM, 0x9c00, gui.CoreDump.PPU.Registers[0].Value&model.Data8(1<<4) == 0),
						gui.Config.GUI.TileMap2.Graphics,
						nil,
					)
				}),
		)
	})
}

func (gui *GUI) OAMColumn() FC {
	return R(func(gtx C) D {
		var views []FC
		views = append(views,
			gui.Box("OAM buffer", &gui.Config.GUI.OAMBuffer.Box,
				func(gtx C) D {
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
						gui.Config.GUI.OAMBuffer.Graphics,
						highlights,
					)
				}),
		)
		views = append(views,
			gui.Box("OAM graphics", &gui.Config.GUI.OAMGraphics.Box,
				func(gtx C) D {
					return gui.GBGraphics(
						gtx,
						80,
						32,
						oam(gui.VRAM, gui.CoreDump.OAM.Bytes()),
						gui.Config.GUI.OAMGraphics.Graphics,
						nil,
					)
				}),
		)
		views = append(views,
			gui.Box("OAM attributes", &gui.Config.GUI.OAMList.Box,
				gui.mem(
					gui.CoreDump.PrintOAMAttrs,
					&gui.OAMAttrScroll,
				)),
		)

		return Column(gtx, views...)
	})
}

func (gui *GUI) Debugger() FC {
	return gui.Box("Debugger", &gui.Config.GUI.Debugger.Box,
		func(gtx C) D {
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
				gui.NumberInput(&gui.StepCyclesInput, "Step Cycles", "cycles", func(text string) {
					cycles, err := strconv.ParseUint(text, 10, 64)
					if err == nil && cycles > 0 {
						gui.StepCycles = cycles
					}
				}),
				gui.NumberInput(&gui.SpeedInput, "Speed%", "speed", func(text string) {
					targetPercent, err := strconv.ParseFloat(text, 64)
					if err == nil && targetPercent > 0 && targetPercent < 10000 {
						if gui.Config.Model.Clock.SpeedPercent != targetPercent {
							gui.Config.Model.Clock.SpeedPercent = targetPercent
							gui.GB.CLK.SetSpeedPercent(gui.Config.Model.Clock.SpeedPercent)
							gui.Config.Save()
						}
					}
				}),
				gui.NumberInput(&gui.BreakPCInput, "PC breakpoint", "breakPC", func(text string) {
					breakPC, err := strconv.ParseInt(text, 16, 64)
					if err != nil {
						breakPC = -1
					}
					if gui.GB.Debug.BreakPC.Load() != breakPC {
						gui.GB.Debug.BreakPC.Store(breakPC)
					}
				}),
				gui.NumberInput(&gui.BreakIRInput, "IR breakpoint", "breakIR", func(text string) {
					breakIR, err := strconv.ParseInt(text, 16, 64)
					if err != nil {
						breakIR = -1
					}
					if gui.GB.Debug.BreakIR.Load() != breakIR {
						gui.GB.Debug.BreakIR.Store(breakIR)
					}
				}), gui.NumberInput(&gui.BreakXInput, "PPU X breakpoint", "breakX", func(text string) {
					breakX, err := strconv.ParseInt(text, 10, 64)
					if err != nil {
						breakX = -1
					}
					if gui.GB.Debug.BreakX.Load() != breakX {
						gui.GB.Debug.BreakX.Store(breakX)
					}
				}),
				gui.NumberInput(&gui.BreakYInput, "PPU Y breakpoint", "breakY", func(text string) {
					breakY, err := strconv.ParseInt(text, 10, 64)
					if err != nil {
						breakY = -1
					}
					if gui.GB.Debug.BreakY.Load() != breakY {
						gui.GB.Debug.BreakY.Store(breakY)
					}
				}),
			)
		})
}
