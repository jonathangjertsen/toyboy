package model

import (
	"fmt"
	"os"
)

type BootROMLock struct {
	MemoryRegion

	BootOff bool
	OnLock  func()
}

func NewBootROMLock(clock *ClockRT) *BootROMLock {
	return &BootROMLock{
		MemoryRegion: NewMemoryRegion(clock, AddrBootROMLock, 0x0001),
	}
}

func (brl *BootROMLock) Read(addr Addr) Data8 {
	return brl.MemoryRegion.Read(addr)
}

func (brl *BootROMLock) Write(addr Addr, v Data8) {
	if brl.BootOff {
		return
	}
	brl.MemoryRegion.Write(addr, v)
	if v&1 == 1 {
		brl.Lock()
	}
}

func (brl *BootROMLock) Lock() {
	brl.BootOff = true
	if brl.OnLock != nil {
		brl.OnLock()
	}
}

func NewBootROM(clk *ClockRT, config ConfigBootROM) MemoryRegion {
	bootrom := NewMemoryRegion(clk, AddrZero, SizeBootROM)
	// todo: static fs
	f, err := os.ReadFile(config.Location)
	if err != nil {
		panic(fmt.Sprintf("failed to load boot rom: %v", err))
	} else if len(f) != 256 {
		panic(fmt.Sprintf("len(bootrom)=%d", len(f)))
	}
	bootrom.Data = Data8Slice(f)
	return bootrom
}
