package model

type CartridgeSlot struct {
	MemoryRegion
}

func NewCartridgeSlot() *CartridgeSlot {
	return &CartridgeSlot{
		MemoryRegion: NewMemoryRegion("CARTRIDGE", 0x0000, 0x4000),
	}
}

func (cs *CartridgeSlot) Read(addr uint16) uint8 {
	return cs.MemoryRegion.Read(addr)
}

func (cs *CartridgeSlot) Write(addr uint16, v uint8) {
	cs.MemoryRegion.Write(addr, v)
}

func (cs *CartridgeSlot) InsertCartridge(data []uint8) {
	copy(cs.MemoryRegion.data, data)
}
