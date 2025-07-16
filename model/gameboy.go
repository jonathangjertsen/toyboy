package model

import (
	"context"
	"log/slog"
)

type Gameboy struct {
	CLK           *ClockRT
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
	ppuClock := clk.Divide(2)
	cpuClock := clk.Divide(4)

	bootROMLock := NewBootROMLock(clk)
	bootROM := NewBootROM(clk, config.Model)
	vram := NewMemoryRegion(clk, 0x8000, 0x2000)
	hram := NewMemoryRegion(clk, 0xff80, 0x007f)
	apu := NewAPU(clk)
	oam := NewMemoryRegion(clk, 0xfe00, 0xa0)
	cartridgeSlot := NewMemoryRegion(clk, 0x0000, 0x4000)

	bus := &Bus{}
	ppu := NewPPU(clk, ppuClock, bus, sysif)

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
