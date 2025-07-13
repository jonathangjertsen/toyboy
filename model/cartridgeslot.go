package model

type CartridgeSlot struct {
	MemoryRegion
	brl *BootROMLock
}

func NewCartridgeSlot(brl *BootROMLock) *CartridgeSlot {
	return &CartridgeSlot{
		MemoryRegion: NewMemoryRegion("CARTRIDGE", 0x0000, 0x4000),
		brl:          brl,
	}
}

func (cs *CartridgeSlot) Name() string {
	name := cs.MemoryRegion.name
	return name + " (empty)"
}

func (cs *CartridgeSlot) Range() (uint16, uint16) {
	start, size := cs.MemoryRegion.Range()
	if !cs.brl.BootOff {
		start += 0x100
		size -= 0x100
	}
	return start, size
}

func (cs *CartridgeSlot) Read(addr uint16) uint8 {
	if !cs.brl.BootOff && addr < 0x100 {
		panic("Read from BOOTROM range while bootrom not locked")
	}
	return cs.MemoryRegion.Read(addr)
}

func (cs *CartridgeSlot) Write(addr uint16, v uint8) {
	if !cs.brl.BootOff && addr < 0x100 {
		panic("Write to BOOTROM range while bootrom not locked")
	}
	cs.MemoryRegion.Write(addr, v)
}
