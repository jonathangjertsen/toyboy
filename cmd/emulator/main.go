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
	gb.Boot()
	<-time.After(time.Second)
	fmt.Printf("ran %v ticks in 1s (%.02f %% speed)\n", i.Load(), 100*float64(i.Load())/realFreq)
}
