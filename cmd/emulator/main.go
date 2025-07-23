package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/jonathangjertsen/toyboy/model"
	"github.com/lmittmann/tint"
)

var DefaultConfig = Config{
	Location: "config.json",
	Model:    model.DefaultConfig,
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

	config, err := LoadConfig(DefaultConfig.Location)
	fmt.Printf("config=%+v\n", config)
	if err != nil {
		fmt.Printf("failed to load config, loading default")
		config = DefaultConfig
	}
	config.Save()

	if config.Model.Debug.GBD.Enable {
		err := config.Model.Debug.GBD.OpenGBDLogFile()
		if err != nil {
			fmt.Printf("couldn't open logfile %s: %v", config.Model.Debug.GBD.File, err)
			return
		}
		defer func() {
			_ = config.Model.Debug.GBD.CloseGBDLogFile()
		}()
	}

	gui := NewGUI(&config)
	go gui.Run()
	Main()
}
