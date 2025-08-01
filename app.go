package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonathangjertsen/toyboy/model"
	"github.com/jonathangjertsen/toyboy/plugin"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(
// None = 0
// Viewport = 1
// CPURegisters = 2
// PPURegisters = 3
// APURegisters = 4
// Disassembly = 5
// HRAM = 6
// WRAM = 7
// OAM = 8
// CPUState = 9
// Clock = 10
// )
type DataID uint8

type App struct {
	ctx     context.Context
	config  *Config
	reqChan chan MachineStateRequest
	Audio   *AudioInterface

	needStateUpdate chan struct{}

	ButtonMapping ButtonMapping

	GB               *model.Gameboy
	ClockMeasurement *plugin.ClockMeasurement
	GBFPSMeasurement *plugin.ClockMeasurement
}

func NewApp(config *Config) *App {
	app := &App{
		ButtonMapping: ButtonMapping{
			A:              "l",
			B:              "k",
			Up:             "w",
			Left:           "a",
			Down:           "s",
			Right:          "d",
			Start:          "m",
			Select:         "n",
			SOCDResolution: SOCDResolutionOppositeNeutral,
		},
		reqChan:         make(chan MachineStateRequest, 10),
		needStateUpdate: make(chan struct{}, 1),
		config:          config,
		Audio:           NewAudio(),
	}
	return app
}

func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
	app.startGB()
	app.startWebSocketServer(app.GB.FrameSync)
}

func (app App) domReady(ctx context.Context) {

}

func (app *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (app *App) shutdown(ctx context.Context) {
}

func (app *App) startGB() {
	audio := &model.Audio{
		SampleInterval: time.Second / 44100,
		SampleBuffers:  model.NewSampleBuffers(512),
		SubSampling:    1024,
		Out:            app.Audio.In,
	}
	app.GB = model.NewGameboy(&app.config.Model, audio)

	runtime.LogPrintf(app.ctx, "Started gameboy")

	app.ClockMeasurement = plugin.NewClockMeasurement(&app.GB.CLK.Cycle)

	runtime.LogPrintf(app.ctx, "Started Clock Measurement")

	app.GBFPSMeasurement = plugin.NewClockMeasurement(&app.GB.PPU.FrameCount)
	runtime.LogPrintf(app.ctx, "Started Gameboy FPS Measurement")

	if err := model.LoadROM(
		"assets/cartridges/tetris.gb",
		app.GB.Mem,
		app.GB.Cartridge,
		app.GB.Bus.BootROMLock,
	); err != nil {
		panic(err)
	}

	app.GB.Start()
}

func (app *App) GetConfig() *Config {
	return app.config
}

type Frame struct {
	Buffer   [144 * 160 * 4]uint8
	SpeedPct string
	US       uint64
}

func (app *App) SetKeyState(in map[string]bool) {
	jp := app.ButtonMapping.JoypadState(in)
	app.GB.Joypad.SetState(&app.GB.CLK, app.GB.Interrupts, jp, app.GB.Mem)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (app *App) MachineStateRequest(req MachineStateRequest) {
	if req.ClickedNumber == "TargetSpeed" {
		app.GB.CLK.SetSpeedPercent(req.Numbers["TargetSpeed"], app.GB.Audio)
		fmt.Printf("Updated speed to %f\n", req.Numbers["TargetSpeed"])
	}
	app.reqChan <- req
}

type TimeoutState struct {
	Interval time.Duration
	Next     time.Time
}

func (ts *TimeoutState) Update() {
	if ts.Interval == 0 {
		ts.Interval = time.Second
	}
	ts.Next = time.Now().Add(ts.Interval)
}

func (app *App) startWebSocketServer(fs *model.FrameSync) {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		nGoroutines := 2
		exit := make(chan struct{}, nGoroutines)

		prevDisRange := Range{0, 0}
		timeouts := map[DataID]TimeoutState{
			DataIDCPURegisters: {Interval: time.Millisecond * 100},
			DataIDPPURegisters: {Interval: time.Millisecond * 100},
			DataIDAPURegisters: {Interval: time.Millisecond * 100},
			DataIDDisassembly:  {Interval: time.Millisecond * 1000},
			DataIDHRAM:         {Interval: time.Millisecond * 100},
			DataIDWRAM:         {Interval: time.Millisecond * 500},
			DataIDOAM:          {Interval: time.Millisecond * 100},
			DataIDCPUState:     {Interval: time.Millisecond * 1000},
		}
		var req MachineStateRequest
		mu := &sync.Mutex{}

		// Hammer the websocket with frames from the PPU
		go func() {
			for {
				app.GB.PPU.Sync(fs, func(vp *model.ViewPort) {
					grayscale := vp.Grayscale()
					if !sendData(conn, mu, DataIDViewport, grayscale[:]) {
						exit <- struct{}{}
						return
					}
				})
			}
		}()

		// Send other stuff at a leisurely pace
		go func() {
			ticker := time.NewTicker(time.Millisecond * 10)
			for {
				force := false
				select {
				case <-ticker.C:
				case <-app.needStateUpdate:
					force = true
				}

				select {
				case req = <-app.reqChan:
				default:
				}

				buffers := map[DataID]*bytes.Buffer{}
				for _, id := range DataIDValues() {
					if force || rateLimit(id, &req, timeouts) {
						buffers[id] = &bytes.Buffer{}
					}
				}
				app.GB.CLK.Sync(func() {
					if buf := buffers[DataIDCPUState]; buf != nil {
						if app.GB.Running.Load() {
							buf.WriteByte(1)
						} else {
							buf.WriteByte(0)
						}
					}
					if buf := buffers[DataIDCPURegisters]; buf != nil {
						model.PrintRegs(buf, app.GB.CPU.Regs)
						di, cycle := app.GB.CPU.CurrInstruction()
						fmt.Fprintf(buf, "\n%s\n   cycle=%d", di.Asm(), cycle)
						fmt.Fprintf(buf, "                                     ")
					}
					if buf := buffers[DataIDPPURegisters]; buf != nil {
						model.PrintPPU(buf, app.GB.PPU.GetDump(), app.GB.Mem)
					}
					if buf := buffers[DataIDAPURegisters]; buf != nil {
						model.PrintAPU(buf, app.GB.Mem, &app.GB.APU)
					}
					if buf := buffers[DataIDHRAM]; buf != nil {
						model.MemDump(
							buf,
							app.GB.Mem,
							model.AddrHRAMBegin,
							model.AddrHRAMEnd,
							app.GB.CPU.Regs.SP,
						)
					}
					if buf := buffers[DataIDOAM]; buf != nil {
						model.MemDump(
							buf,
							app.GB.Mem,
							model.AddrOAMBegin,
							model.AddrOAMEnd,
							0,
						)
					}
					if buf := buffers[DataIDDisassembly]; buf != nil {
						rng := req.Ranges[DataIDDisassembly.String()]
						rng = rng.Constrain(0x0000, 0xffff)
						if rng != prevDisRange {
							dis := app.GB.Debug.Disassembler.Disassembly(model.Addr(rng.Begin), model.Addr(rng.End))
							dis.Print(buf)
						}
						prevDisRange = rng
					}
					if buf := buffers[DataIDWRAM]; buf != nil {
						rang := req.Ranges[DataIDWRAM.String()].Constrain(
							uint(model.AddrWRAMBegin),
							uint(model.AddrWRAMEnd),
						)
						model.MemDump(
							buf,
							app.GB.Mem,
							model.Addr(rang.Begin),
							model.Addr(rang.End),
							app.GB.CPU.Regs.SP,
						)
					}
					if buf := buffers[DataIDClock]; buf != nil {
						cycles, fdur := app.ClockMeasurement.Stop()
						app.ClockMeasurement.Start()

						cps := float64(cycles) * 1_000_000 / float64(uint64(fdur/time.Microsecond))

						if app.GB.CLK.Cycle < 1_000_000 {
							fmt.Fprintf(buf, "Cycle: %d\n", app.GB.CLK.Cycle)
						} else if app.GB.CLK.Cycle < 1_000_000_000 {
							fmt.Fprintf(buf, "Cycle: %d M\n", app.GB.CLK.Cycle/1_000_000)
						} else {
							fmt.Fprintf(buf, "Cycle: %d G\n", app.GB.CLK.Cycle/1_000_000_000)
						}
						fmt.Fprintf(buf, "Speed: %.0f %%\n", (100*cps)/4194304)
					}
				})
				for id, buf := range buffers {
					if buf.Len() > 0 {
						if !sendData(conn, mu, id, buf.Bytes()) {
							exit <- struct{}{}
							return
						}
					}
				}
			}
		}()

		<-exit
		conn.Close()
	})

	go func() {
		err := http.ListenAndServe(":8081", nil)
		panic(err)
	}()
}

func rateLimit(id DataID, req *MachineStateRequest, timeouts map[DataID]TimeoutState) bool {
	if time.Now().Before(timeouts[id].Next) {
		return false
	}
	if !req.OpenBoxes[id.String()] {
		return false
	}
	ts := timeouts[id]
	ts.Update()
	timeouts[id] = ts
	return true
}

func sendData(conn *websocket.Conn, mu *sync.Mutex, id DataID, content []uint8) bool {
	ok := true
	mu.Lock()
	if conn.WriteMessage(websocket.TextMessage, []uint8(id.String())) != nil {
		ok = false
	}
	if conn.WriteMessage(websocket.BinaryMessage, content) != nil {
		ok = false
	}
	mu.Unlock()
	return ok
}

type MachineStateRequest struct {
	OpenBoxes     map[string]bool
	Numbers       map[string]float64
	Ranges        map[string]Range
	ClickedNumber string
	ClickedRange  string
}

type Range struct {
	Begin float64
	End   float64
}

func (r Range) Valid() bool {
	return r.Begin != 0 || r.End != 0
}

func (r Range) Constrain(begin, end uint) Range {
	if r.Begin == 0 && r.End == 0 {
		return Range{}
	}
	fbegin, fend := float64(begin), float64(end)
	if r.Begin < fbegin {
		r.Begin = fbegin
	}
	if r.End > fend {
		r.End = fend
	}
	if r.Begin > r.End {
		return Range{}
	}
	return r
}

func (app *App) Pause() {
	app.GB.Pause()
	select {
	case <-app.needStateUpdate:
	default:
	}
}

func (app *App) Start() {
	app.GB.Start()
	select {
	case <-app.needStateUpdate:
	default:
	}
}

func (app *App) Step() {
	app.GB.Step()
	select {
	case <-app.needStateUpdate:
	default:
	}
}

func (app *App) SoftReset() {
	app.GB.SoftReset()
}
