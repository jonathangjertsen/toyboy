package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/jonathangjertsen/toyboy/gui"
	"github.com/jonathangjertsen/toyboy/model"
	"github.com/lmittmann/tint"
)

var realFreq = 4194304.0

var hwConfig = model.HWConfig{
	SystemClock: model.ClockConfig{
		Frequency: realFreq,
	},
}

type sysInterface struct {
}

func (si *sysInterface) FrameCompleted(vp model.ViewPort) {
	if vp == (model.ViewPort{}) {
		return
	}

	fmt.Printf("FRAME:\n")
	for _, row := range vp {
		for _, col := range row {
			fmt.Printf("%d", int(col))
		}
		fmt.Printf("\n")
	}
}

func main() {
	ctx := context.Background()
	var logWriter io.Writer = os.Stdout
	var logHandler slog.Handler = tint.NewHandler(logWriter, &tint.Options{})
	logger := slog.New(logHandler)
	if os.Getenv("APP_ENV") == "development" {
		logger.Info("Enabling pprof for profiling")
		go func() {
			logger.Info("Exited", "err", http.ListenAndServe("localhost:6060", nil))
		}()
	}
	gb := model.NewGameboy(ctx, logger, hwConfig, &sysInterface{})
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
	g := gui.New(gb)
	go g.Run()
	gui.Main()
}
