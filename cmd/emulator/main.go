package main

import (
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

func main() {
	var logWriter io.Writer = os.Stdout
	var logHandler slog.Handler = tint.NewHandler(logWriter, &tint.Options{})
	logger := slog.New(logHandler)
	if os.Getenv("APP_ENV") == "development" {
		logger.Info("Enabling pprof for profiling")
		go func() {
			logger.Info("Exited", "err", http.ListenAndServe("localhost:6060", nil))
		}()
	}

	g := gui.New(hwConfig)
	go g.Run()
	gui.Main()
}
