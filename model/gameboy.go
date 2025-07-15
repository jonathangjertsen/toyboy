package model

import (
	"context"
	"log/slog"
)

type Gameboy struct {
	CLK           *RealtimeClock
	PHI           *Clock
	CPU           *CPU
	CartridgeSlot *MemoryRegion
}

func (gb *Gameboy) PowerOn() {
	gb.CLK.Start()
}

func (gb *Gameboy) PowerOff() {
	gb.CLK.Stop()
}

func (gb *Gameboy) GetCoreDump() CoreDump {
	var cd CoreDump
	gb.CLK.Sync(func() {
		cd = gb.CPU.GetCoreDump()
	})
	return cd
}

type SysInterface interface {
	PPUHooks
}

func NewGameboy(
	ctx context.Context,
	logger *slog.Logger,
	config HWConfig,
	sysif SysInterface,
) *Gameboy {
	clk := NewRealtimeClock(config.SystemClock)
	ppuClock := clk.Divide(1)
	cpuClock := clk.Divide(2)

	bootROMLock := NewBootROMLock()
	bootROM := NewBootROM(config.Model)
	vram := NewMemoryRegion(0x8000, 0x2000)
	hram := NewMemoryRegion(0xff80, 0x007f)
	apu := NewAPU()
	oam := NewMemoryRegion(0xfe00, 0xa0)
	cartridgeSlot := NewMemoryRegion(0x0000, 0x4000)

	bus := &Bus{}
	ppu := NewPPU(ppuClock, bus, sysif)

	bus.BootROMLock = bootROMLock
	bus.BootROM = &bootROM
	bus.VRAM = &vram
	bus.HRAM = &hram
	bus.APU = apu
	bus.OAM = &oam
	bus.PPU = ppu
	bus.CartridgeSlot = &cartridgeSlot

	cpu := NewCPU(cpuClock, bus)

	soc := &Gameboy{}
	soc.CLK = clk
	soc.PHI = cpuClock
	soc.CPU = cpu
	soc.CartridgeSlot = &cartridgeSlot

	return soc
}
