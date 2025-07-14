package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/jonathangjertsen/gameboy/model"
	"github.com/lmittmann/tint"
)

var boost = 40.0
var realFreq = 4194304.0

var hwConfig = model.HWConfig{
	SystemClock: model.ClockConfig{
		Frequency: realFreq * boost,
	},
}

func main() {
	ctx := context.Background()

	var logWriter io.Writer = os.Stdout

	var logHandler slog.Handler = tint.NewHandler(logWriter, &tint.Options{})

	logger := slog.New(logHandler)

	gb := model.NewGameboy(ctx, logger, hwConfig)
	i := atomic.Uint64{}
	gb.PHI.AddRiseCallback(func(c model.Cycle) {
		i.Add(4)
	})

	f, err := os.ReadFile("assets/cartridges/hello-world.gb")
	if err != nil {
		panic(fmt.Sprintf("failed to load cartridge: %v", err))
	} else if len(f) != 0x8000 {
		panic(fmt.Sprintf("len(bootrom)=%d", len(f)))
	}

	gb.CartridgeSlot.InsertCartridge(f)
	gb.PowerOn()
	nSeconds := 50.0
	<-time.After(time.Second * time.Duration(nSeconds))
	gb.PowerOff()
	gb.CPU.Dump()
	fmt.Printf("ran %v ticks in %f s (%.02f %% speed)\n", i.Load(), nSeconds, 100*float64(i.Load())/(nSeconds*realFreq))
}
