package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/jonathangjertsen/toyboy/model"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/src
var assets embed.FS

////go:embed build/appicon.png
//var icon []byte

func main() {
	config, err := LoadConfig(DefaultConfig.Location)
	fmt.Printf("config=%+v\n", config)
	if err != nil {
		fmt.Println("failed to load config, loading default")
		config = DefaultConfig
	}
	config.Save()

	if config.PProfURL != "" {
		go func() {
			http.ListenAndServe(config.PProfURL, nil)
		}()
	}

	if len(os.Args) == 2 && os.Args[1] == "debug" {
		gb := model.NewGameboy(&config.Model)

		if err := model.LoadROM(
			"assets/cartridges/tetris.gb",
			gb.Mem,
			&gb.Cartridge,
			gb.Bus.BootROMLock,
		); err != nil {
			panic(err)
		}

		gb.Start()
		select {}
	}

	// Create an instance of the app structure
	app := NewApp(&config)

	// Create application with options
	err = wails.Run(&options.App{
		Title:             "toyboy",
		MinWidth:          2000,
		MinHeight:         2000,
		DisableResize:     false,
		Fullscreen:        true,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnBeforeClose:     app.beforeClose,
		OnShutdown:        app.shutdown,
		WindowStartState:  options.Normal,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
