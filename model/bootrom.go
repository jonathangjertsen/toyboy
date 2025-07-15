package model

import (
	"fmt"
	"os"
)

type BootROMLock struct {
	RegBootROMLock uint8

	BootOff bool
}

func NewBootROMLock() *BootROMLock {
	return &BootROMLock{}
}

func (brl *BootROMLock) Read(addr uint16) uint8 {
	if addr == 0xff50 {
		return brl.RegBootROMLock
	}
	panicv(addr)
	return 0
}

func (brl *BootROMLock) Write(addr uint16, v uint8) {
	if brl.BootOff {
		return
	}
	if v&1 == 1 {
		brl.RegBootROMLock = 0x01
		brl.BootOff = true
	}
}

func NewBootROM(model Model) MemoryRegion {
	bootrom := NewMemoryRegion("BOOTROM", 0x0000, 0x0100)
	switch model {
	case DMG:
		// todo: static fs
		f, err := os.ReadFile("assets/bootrom/dmg_boot.bin")
		if err != nil {
			panic(fmt.Sprintf("failed to load boot rom: %v", err))
		} else if len(f) != 256 {
			panic(fmt.Sprintf("len(bootrom)=%d", len(f)))
		}
		copy(bootrom.data, f)
	}
	return bootrom
}
