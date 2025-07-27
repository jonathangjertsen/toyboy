package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonathangjertsen/toyboy/model"
	"github.com/jonathangjertsen/toyboy/plugin"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	config *Config

	ButtonMapping ButtonMapping
	speed         float64

	GB               *model.Gameboy
	ClockMeasurement *plugin.ClockMeasurement
	GBFPSMeasurement *plugin.ClockMeasurement
}

func NewApp(config *Config) *App {
	return &App{
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
		config: config,
	}
}

func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
}

func (app App) domReady(ctx context.Context) {

}

func (app *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (app *App) shutdown(ctx context.Context) {
}

func (app *App) StartGB() {
	app.GB = model.NewGameboy(&app.config.Model)

	runtime.LogPrintf(app.ctx, "Started gameboy")

	app.ClockMeasurement = plugin.NewClockMeasurement(&app.GB.CLK.Cycle.C)

	runtime.LogPrintf(app.ctx, "Started Clock Measurement")

	app.GBFPSMeasurement = plugin.NewClockMeasurement(&app.GB.PPU.FrameCount)
	runtime.LogPrintf(app.ctx, "Started Gameboy FPS Measurement")

	f, err := os.ReadFile("assets/cartridges/unbricked.gb")
	if err != nil {
		panic(err)
	}
	app.GB.Cartridge.LoadROM(f)

	app.GB.CLK.Start()
	app.StartWebSocketServer()
}

func (app *App) GetConfig() *Config {
	return app.config
}

type Frame struct {
	Buffer   [144 * 160 * 4]uint8
	SpeedPct string
	US       uint64
}

func (app *App) GetSpeedPct() string {
	return fmt.Sprintf("%.1f", app.speed)
}

func (app *App) SetKeyState(in map[string]bool) {
	jp := app.ButtonMapping.JoypadState(in)
	app.GB.Joypad.SetState(jp)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (app *App) StartWebSocketServer() {
	http.HandleFunc("/vp", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		i := 60
		for {
			app.GB.PPU.Sync(func(vp *model.ViewPort) {
				grayscale := vp.Grayscale()
				i--
				if i == 0 {
					cycles, fdur := app.ClockMeasurement.Stop()
					app.ClockMeasurement.Start()
					cps := float64(cycles) * 1_000_000 / float64(uint64(fdur/time.Microsecond))
					app.speed = (100 * cps) / 4194304
					i = 60
				}

				conn.WriteMessage(websocket.BinaryMessage, grayscale[:])
			})
		}
	})

	go http.ListenAndServe(":8081", nil)
}

func (app *App) GetCoreDump() model.CoreDump {
	return app.GB.GetCoreDump()
}
