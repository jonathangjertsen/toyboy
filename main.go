package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

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

	// Create an instance of the app structure
	app := NewApp(&config)

	// Create application with options
	err = wails.Run(&options.App{
		Title:             "toyboy",
		Width:             1024,
		Height:            768,
		MinWidth:          1024,
		MinHeight:         768,
		DisableResize:     false,
		Fullscreen:        false,
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
