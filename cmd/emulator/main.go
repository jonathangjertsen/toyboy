package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"
)

func main() {
	config, err := LoadConfig(DefaultConfig.Location)
	fmt.Printf("config=%+v\n", config)
	if err != nil {
		fmt.Println("failed to load config, loading default")
		config = DefaultConfig
	}
	config.Save()

	if config.PProfURL != "" {
		fmt.Println("Enabling pprof for profiling")
		go func() {
			fmt.Printf("Exited, err=%s\n", http.ListenAndServe(config.PProfURL, nil))
		}()
	}

	gui := NewGUI(&config)
	go gui.Run()
	Main()
}
