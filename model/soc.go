package model

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type SOC struct {
	CLK     *RootClock
	PHI     *Clock
	CPUCore *CPUCore
	Data    *Bus[uint8]
	Address *Bus[uint16]
}

func NewSOC(ctx context.Context, logger *slog.Logger, config HWConfig) *SOC {
	clk := NewRootClock(config.SystemClock)
	phi := clk.Divide(4)
	core := NewCPUCore(phi)

	bootROM := NewMemoryRegion("BOOTROM", 0x0000, 0x0100)
	switch config.Model {
	case DMG:
		// todo: static fs
		f, err := os.ReadFile("assets/bootrom/dmg_boot.bin")
		if err != nil {
			panic(fmt.Sprintf("failed to load boot rom: %v", err))
		} else if len(f) != 256 {
			panic(fmt.Sprintf("len(bootrom)=%d", len(f)))
		}
		copy(bootROM.data, f)
	}
	core.AttachPeripheral(bootROM)

	vram := NewMemoryRegion("VRAM", 0x8000, 0x2000)
	core.AttachPeripheral(vram)

	audio := NewAudio()
	core.AttachPeripheral(audio)

	soc := &SOC{}
	soc.CLK = clk
	soc.PHI = phi
	soc.CPUCore = core

	return soc
}
