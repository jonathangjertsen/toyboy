package model

type Bus struct {
	data    uint8
	Address uint16

	BootROMLock   *BootROMLock
	BootROM       *MemoryRegion
	VRAM          *MemoryRegion
	HRAM          *MemoryRegion
	APU           *APU
	OAM           *MemoryRegion
	PPU           *PPU
	CartridgeSlot *CartridgeSlot
}

func (b *Bus) WriteAddress(addr uint16) {
	b.Address = addr

	if addr < 0x100 {
		if b.BootROMLock.BootOff {
			b.data = b.CartridgeSlot.Read(addr)
		} else {
			b.data = b.BootROM.Read(addr)
		}
	} else if addr <= 0x4000 {
		b.data = b.CartridgeSlot.Read(addr)
	} else if addr >= 0x8000 && addr < 0xa000 {
		b.data = b.VRAM.Read(addr)
	} else if addr >= 0xff80 && addr < 0xffff {
		b.data = b.HRAM.Read(addr)
	} else if addr >= 0xff10 && addr < 0xff28 {
		b.data = b.APU.Read(addr)
	} else if addr >= 0xfe00 && addr < 0xfea0 {
		b.data = b.OAM.Read(addr)
	} else if addr >= 0xff40 && addr < 0xff4c {
		b.data = b.PPU.Read(addr)
	} else if addr == 0xff50 {
		b.data = b.BootROMLock.Read(addr)
		return
	} else {
		panicf("read from unknown peripheral at 0x%x", addr)
	}
}

func (b *Bus) WriteData(v uint8) {
	b.data = v
	addr := b.Address
	if addr <= 0x4000 {
		if addr < 0x100 {
			panicf("Attempted write to bootrom (addr=0x%04x v=%02x)", addr, v)
		}
		panicf("Attempted write to cartridge (addr=0x%04x v=%02x)", addr, v)
	} else if addr == 0xff50 {
		b.BootROMLock.Write(addr, v)
	} else if addr >= 0x8000 && addr < 0xa000 {
		b.VRAM.Write(addr, v)
	} else if addr >= 0xff80 && addr < 0xffff {
		b.HRAM.Write(addr, v)
	} else if addr >= 0xff10 && addr < 0xff28 {
		b.APU.Write(addr, v)
	} else if addr >= 0xfe00 && addr < 0xfea0 {
		b.OAM.Write(addr, v)
	} else if addr >= 0xff40 && addr < 0xff4c {
		b.PPU.Write(addr, v)
	} else {
		panicf("write to unknown peripheral at 0x%x", addr)
	}
}

func (b *Bus) Data() uint8 {
	return b.data
}
