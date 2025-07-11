package main

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/jonathangjertsen/gameboy/model"
	"github.com/lmittmann/tint"
)

var hwConfig = model.HWConfig{
	SystemClockFrequency: 4194304.0,
}

func main() {
	ctx := context.Background()

	var logWriter io.Writer = os.Stdout

	var logHandler slog.Handler = tint.NewHandler(logWriter, &tint.Options{})

	logger := slog.New(logHandler)

	soc := model.NewSOC(ctx, logger, hwConfig)
	_ = soc
	select {}
}
