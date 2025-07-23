package model

import (
	"io"
	"os"
)

type Config struct {
	Clock   ConfigClock
	ROM     ConfigROM
	BootROM ConfigBootROM
	Debug   ConfigDebug
}

type ConfigBootROM struct {
	Location string
	Skip     bool
}

type ConfigClock struct {
	SpeedPercent float64
}

type ConfigROM struct {
	Location string
}

type ConfigDebug struct {
	MaxNOPCount  int
	GBD          ConfigGBD
	Disassembler ConfigDisassembler
}

type ConfigDisassembler struct {
	Trace bool
}

type ConfigGBD struct {
	Enable bool
	File   string
	gbdf   io.WriteCloser
}

func (cgbd *ConfigGBD) OpenGBDLogFile() error {
	f, err := os.Create(cgbd.File)
	cgbd.gbdf = f
	return err
}

func (cgbd *ConfigGBD) CloseGBDLogFile() error {
	return cgbd.gbdf.Close()
}

func (cgbd *ConfigGBD) GBDLog(s string) {
	_, _ = cgbd.gbdf.Write([]byte(s))
}

var DefaultConfig = Config{
	Clock: ConfigClock{
		SpeedPercent: 100.0,
	},
	ROM: ConfigROM{
		"assets/cartridges/01-special.gb",
	},
	BootROM: ConfigBootROM{
		Skip:     false,
		Location: "assets/bootrom/dmg_boot.bin",
	},
	Debug: ConfigDebug{
		MaxNOPCount: 10,
		GBD: ConfigGBD{
			Enable: false,
			File:   "bin/gbdlog.txt",
		},
		Disassembler: ConfigDisassembler{
			Trace: true,
		},
	},
}
