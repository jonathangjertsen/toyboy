package model

import (
	"context"
	"log/slog"
)

type Gameboy struct {
	CLK     *RealtimeClock
	PHI     *Clock
	CPUCore *CPUCore
	Data    *Bus[uint8]
	Address *Bus[uint16]
}

func (gb *Gameboy) Boot() {
	gb.CLK.Start()
}

func NewGameboy(ctx context.Context, logger *slog.Logger, config HWConfig) *Gameboy {
	clk := NewRealtimeClock(config.SystemClock)
	phi := clk.Divide(4)
	core := NewCPUCore(phi)

	bootROMLock := NewBootROMLock()
	core.AttachPeripheral(bootROMLock)

	bootROM := NewBootROM(bootROMLock, config.Model)
	core.AttachPeripheral(bootROM)

	vram := NewMemoryRegion("VRAM", 0x8000, 0x2000)
	core.AttachPeripheral(&vram)

	hram := NewMemoryRegion("HRAM", 0xff80, 0x007f)
	core.AttachPeripheral(&hram)

	audio := NewAudioCtl()
	core.AttachPeripheral(audio)

	lcd := NewPPU()
	core.AttachPeripheral(lcd)

	cartridgeSlot := NewCartridgeSlot(bootROMLock)
	core.AttachPeripheral(cartridgeSlot)

	soc := &Gameboy{}
	soc.CLK = clk
	soc.PHI = phi
	soc.CPUCore = core

	return soc
}
