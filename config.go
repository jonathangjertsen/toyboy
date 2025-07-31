package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"github.com/jonathangjertsen/toyboy/model"
)

var DefaultConfig = Config{
	Location: "config.json",
	Model:    model.DefaultConfig,
	GUI: ConfigGUI{
		Graphics: ConfigGraphicsGlobal{
			Overlay:       false,
			DashLen:       4,
			GridThickness: 1,
			GridColor:     color.RGBA{136, 136, 136, 255},
			FillColor:     color.RGBA{240, 240, 240, 255},
			Font:          "Basic",
		},
		VRAMMem: ConfigVRAM{
			Box: ConfigBox{
				Show: true,
			},
		},
		WRAMMem: ConfigWRAM{
			Box: ConfigBox{
				Show: true,
			},
		},
		ProgMem: ConfigProgram{
			Box: ConfigBox{
				Show: true,
			},
		},
		HRAMMem: ConfigHRAM{
			Box: ConfigBox{
				Show: true,
			},
		},
		OAMMem: ConfigOAM{
			Box: ConfigBox{
				Show: true,
			},
		},
		Rewind: ConfigRewind{
			Box: ConfigBox{
				Show: true,
			},
		},
		Disassembly: ConfigDisassembly{
			Box: ConfigBox{
				Show: true,
			},
		},
		Timing: ConfigTiming{
			Box: ConfigBox{
				Show: true,
			},
		},
		Registers: ConfigRegisters{
			Box: ConfigBox{
				Show: true,
			},
		},
		APU: ConfigAPU{
			Box: ConfigBox{
				Show: true,
			},
		},
		PPU: ConfigPPU{
			Box: ConfigBox{
				Show: true,
			},
		},
		ViewPort: ConfigViewPort{
			Box: ConfigBox{
				Show: true,
			},
		},
		JoyPad: ConfigJoyPad{
			Box: ConfigBox{
				Show: true,
			},
		},
		Debugger: ConfigDebugger{
			Box: ConfigBox{
				Show: true,
			},
		},
		OAMBuffer: ConfigOAMBuffer{
			Box: ConfigBox{
				Show: true,
			},
		},
		OAMGraphics: ConfigOAMGraphics{
			Box: ConfigBox{
				Show: true,
			},
		},
		OAMList: ConfigOAMList{
			Box: ConfigBox{
				Show: true,
			},
		},
	},
}

type Config struct {
	Location string
	Model    model.Config
	PProfURL string
	GUI      ConfigGUI
}

type ConfigGUI struct {
	Graphics    ConfigGraphicsGlobal
	VRAMMem     ConfigVRAM
	WRAMMem     ConfigWRAM
	ProgMem     ConfigProgram
	HRAMMem     ConfigHRAM
	OAMMem      ConfigOAM
	Rewind      ConfigRewind
	Disassembly ConfigDisassembly
	Debugger    ConfigDebugger
	Timing      ConfigTiming
	Registers   ConfigRegisters
	APU         ConfigAPU
	PPU         ConfigPPU
	ViewPort    ConfigViewPort
	JoyPad      ConfigJoyPad
	TileData    ConfigTileData
	TileMap1    ConfigTileMap
	TileMap2    ConfigTileMap
	OAMBuffer   ConfigOAMBuffer
	OAMGraphics ConfigOAMGraphics
	OAMList     ConfigOAMList
}

type ConfigVRAM struct {
	Box ConfigBox
}

type ConfigWRAM struct {
	Box ConfigBox
}

type ConfigProgram struct {
	Box ConfigBox
}

type ConfigHRAM struct {
	Box ConfigBox
}

type ConfigOAM struct {
	Box ConfigBox
}

type ConfigRewind struct {
	Box ConfigBox
}

type ConfigDisassembly struct {
	Box ConfigBox
}

type ConfigDebugger struct {
	Box ConfigBox
}

type ConfigAPU struct {
	Box ConfigBox
}

type ConfigPPU struct {
	Box ConfigBox
}

type ConfigBox struct {
	Show   bool
	Height int
	Width  int
}

type ConfigTiming struct {
	Box ConfigBox
}

type ConfigRegisters struct {
	Box ConfigBox
}

type ConfigViewPort struct {
	Box      ConfigBox
	Graphics ConfigGraphics
}

type ConfigJoyPad struct {
	Box      ConfigBox
	Graphics ConfigGraphics
}

type ConfigTileData struct {
	Box      ConfigBox
	Graphics ConfigGraphics
}

type ConfigTileMap struct {
	Box      ConfigBox
	Graphics ConfigGraphics
}

type ConfigOAMBuffer struct {
	Box      ConfigBox
	Graphics ConfigGraphics
}

type ConfigOAMGraphics struct {
	Box      ConfigBox
	Graphics ConfigGraphics
}

type ConfigOAMList struct {
	Box ConfigBox
}

type ConfigGraphicsGlobal struct {
	Overlay       bool
	GridColor     color.RGBA // RGBA grid line color
	FillColor     color.RGBA // RGBA fill/background
	DashLen       int
	GridThickness int
	Font          string
}

type ConfigGraphics struct {
	ShowGrid  bool
	BlockSize int
	Scale     int

	ShowAddress    bool
	StartAddress   model.Addr
	BlockIncrement model.Addr
	LineIncrement  model.Addr
	DecimalAddress bool

	ShowOffsets bool
}

func LoadConfig(location string) (Config, error) {
	var config Config
	jsonData, err := os.ReadFile(location)
	if err != nil {
		return config, fmt.Errorf("failed to load config file in %s", location)
	}
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return config, fmt.Errorf("config in %s is corrupted", location)
	}
	config.Location = location
	return config, nil
}

func (conf *Config) Save() error {
	jsonData, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	if err := os.WriteFile(conf.Location, jsonData, 0o666); err != nil {
		return fmt.Errorf("failed to save config file: %v", err)
	}
	return nil
}
