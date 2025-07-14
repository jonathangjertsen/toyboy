package model

import (
	"context"
	"log/slog"
)

type Gameboy struct {
	CLK           *RealtimeClock
	PHI           *Clock
	CPU           *CPU
	CartridgeSlot *CartridgeSlot
}

func (gb *Gameboy) PowerOn() {
	gb.CLK.Start()
}

func (gb *Gameboy) PowerOff() {
	gb.CLK.Stop()
}

func NewGameboy(ctx context.Context, logger *slog.Logger, config HWConfig) *Gameboy {
	clk := NewRealtimeClock(config.SystemClock)
	ppuClock := clk.Divide(2)
	cpuClock := clk.Divide(4)
	cpu := NewCPU(cpuClock)

	bootROMLock := NewBootROMLock()
	cpu.AttachPeripheral(bootROMLock)

	bootROM := NewBootROM(bootROMLock, config.Model)
	cpu.AttachPeripheral(bootROM)

	vram := NewMemoryRegion("VRAM", 0x8000, 0x2000)
	cpu.AttachPeripheral(&vram)

	hram := NewMemoryRegion("HRAM", 0xff80, 0x007f)
	cpu.AttachPeripheral(&hram)

	apu := NewAPU()
	cpu.AttachPeripheral(apu)

	oam := NewMemoryRegion("OAM", 0xfe00, 0xa0)
	cpu.AttachPeripheral(&oam)

	ppu := NewPPU(ppuClock, &oam, &vram)
	cpu.AttachPeripheral(ppu)

	cartridgeSlot := NewCartridgeSlot(bootROMLock)
	cpu.AttachPeripheral(cartridgeSlot)

	soc := &Gameboy{}
	soc.CLK = clk
	soc.PHI = cpuClock
	soc.CPU = cpu
	soc.CartridgeSlot = cartridgeSlot

	return soc
}
