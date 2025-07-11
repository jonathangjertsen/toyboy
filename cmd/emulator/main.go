package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/jonathangjertsen/gameboy/model"
	"github.com/lmittmann/tint"
)

var boost = 20.0
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

	soc := model.NewSOC(ctx, logger, hwConfig)
	i := 0
	soc.CLK.AddCallback <- func(c model.Cycle) {
		i += 1
	}
	<-time.After(time.Second)
	fmt.Printf("ran %v ticks in 1s (%.02f %% speed)\n", i, 100*float64(i)/realFreq)
}
