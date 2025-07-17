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
	Config model.HWConfig
	GB     *model.Gameboy

	ClockMeasurement *plugin.ClockMeasurement
	Theme            material.Theme

	SpeedInput      widget.Editor
	StartButton     widget.Clickable
	PauseButton     widget.Clickable
	ResetButton     widget.Clickable
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

func New(config model.HWConfig) *GUI {
	gui := &GUI{
		Config: config,
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
	defer func() {
		if e := recover(); e != nil {
			gb.CPU.Dump()
			panic(e)
		}
	}()
	f, err := os.ReadFile("assets/cartridges/hello-world.gb")
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
	gui.SpeedInput.SetText("100")
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
		gui.GB.Pause()
	}
	if gui.ResetButton.Clicked(gtx) {
		gui.GB.Stop()
		gui.initGameboy()
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
						Rigid(gui.memHead("APU")),
						Rigid(gui.mem(cd.PrintAPU, &gui.APUScroll, unit.Dp(300))),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.memHead("PPU")),
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
					gc := DefaultGridConfig
					gc.Show = false
					return Column(
						gtx,
						Rigid(gui.memHead(fmt.Sprintf("Viewport (X=%d,Y=%d)", cd.PPU.BackgroundFetcher.X, (cd.PPU.Registers[model.AddrLY-model.AddrPPUBegin].Value)/8))),
						Rigid(func(gtx C) D {
							vp := gui.GB.GetViewport()
							pixels := vp.Flatten()
							return gui.GBGraphics(
								gtx,
								160,
								144,
								pixels[:],
								2,
								gc,
								nil,
							)
						}),
					)
				}),
				Rigid(func(gtx C) D {
					return Column(
						gtx,
						Rigid(gui.memHead("Tile data")),
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
						Rigid(gui.memHead("Tile map 1")),
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
						Rigid(gui.memHead("Tile map 2")),
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
			)
		}),
		Rigid(func(gtx C) D {
			return Row(
				gtx,
				Rigid(func(gtx C) D {
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
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.StartButton, "Run")
				}),
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.PauseButton, "Pause")
				}),
				Rigid(func(gtx C) D {
					return gui.Button(gtx, &gui.ResetButton, "Reset")
				}),
			)
		}),
		Spacer(25, 0),
	)
}

func tiledata(vram []uint8) []model.Color {
	tileData := vram[:0x1800]
	tiles := make([]model.Tile, len(tileData)/16)
	for i := range tiles {
		tiles[i] = model.DecodeTile(tileData[i*16 : (i+1)*16])
	}
	return placetiles(tiles, 24, 16)
}

func tilemap(vram []uint8, addr uint16, signedAddressing bool) []model.Color {
	tileMap := vram[addr-model.AddrVRAMBegin : addr-model.AddrVRAMBegin+0x400]
	tiles := make([]model.Tile, len(tileMap))
	for i := range tiles {
		tileID := tileMap[i]
		var offset uint16
		if signedAddressing {
			offset = uint16(int32(0x1000) + 16*int32(int8(tileID)))
		} else {
			offset = 16 * uint16(tileID)
		}
		tile := vram[offset : offset+16]
		tiles[i] = model.DecodeTile(tile)
	}
	return placetiles(tiles, 32, 32)
}

func placetiles(tiles []model.Tile, w, h int) []model.Color {
	fb := make([]model.Color, h*w*8*8)
	for tileRow := range h {
		for tileCol := range w {
			tile := tiles[tileRow*w+tileCol]
			for rowInTile := range 8 {
				for colInTile := range 8 {
					col := tile[rowInTile][colInTile].Color
					fb[(tileRow*8+rowInTile)*(8*w)+tileCol*8+colInTile] = col
				}
			}
		}
	}
	return fb
}

func (gui *GUI) Button(gtx C, clickable *widget.Clickable, text string) D {
	return material.Button(&gui.Theme, clickable, text).Layout(gtx)
}

func (gui *GUI) NumberInput(gtx C, editor *widget.Editor, placeholder string, f func(text string)) D {
	gui.SpeedInput.SingleLine = true
	gui.SpeedInput.Alignment = text.Middle
	f(editor.Text())
	return material.Editor(&gui.Theme, editor, placeholder).Layout(gtx)
}
