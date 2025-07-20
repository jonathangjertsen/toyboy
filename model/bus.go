package model

type Bus struct {
	Data    Data8
	Address Addr

	inCoreDump bool

	BootROMLock   *BootROMLock
	BootROM       *MemoryRegion
	VRAM          *MemoryRegion
	HRAM          *MemoryRegion
	WRAM          *MemoryRegion
	APU           *APU
	OAM           *MemoryRegion
	PPU           *PPU
	CartridgeSlot *MemoryRegion
	Joypad        *Joypad
	Interrupts    *Interrupts
}

func (b *Bus) WriteAddress(addr Addr) {
	b.Address = addr

	if addr <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			b.Data = b.CartridgeSlot.Read(addr)
		} else {
			b.Data = b.BootROM.Read(addr)
		}
	} else if addr <= AddrCartridgeBank0End {
		b.Data = b.CartridgeSlot.Read(addr)
	} else if addr >= AddrVRAMBegin && addr <= AddrVRAMEnd {
		b.Data = b.VRAM.Read(addr)
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		b.Data = b.HRAM.Read(addr)
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		b.Data = b.WRAM.Read(addr)
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.Data = b.APU.Read(addr)
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		b.Data = b.OAM.Read(addr)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.Data = b.PPU.Read(addr)
	} else if addr == AddrP1 {
		b.Data = b.Joypad.Read(addr)
	} else if addr == AddrIF || addr == AddrIE {
		b.Data = b.Interrupts.Read(addr)
	} else if addr == AddrBootROMLock {
		b.Data = b.BootROMLock.Read(addr)
	} else {
		if !b.inCoreDump {
			panicf("Read from unmapped address %s", addr.Hex())
		}
	}
}

func (b *Bus) WriteData(v Data8) {
	b.Data = v
	addr := b.Address
	if addr <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			panicf("Attempted write to cartridge (addr=%s v=%s)", addr.Hex(), v.Hex())
		} else {
			panicf("Attempted write to bootrom (addr=%s v=%s)", addr.Hex(), v.Hex())
		}
	} else if addr <= AddrCartridgeBank0End {
		b.Data = b.CartridgeSlot.Read(addr)
		panicf("Attempted write to cartridge (addr=%s v=%s)", addr.Hex(), v.Hex())
	} else if addr == AddrBootROMLock {
		b.BootROMLock.Write(addr, v)
	} else if addr == AddrP1 {
		b.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		b.Interrupts.Write(addr, v)
	} else if addr >= AddrVRAMBegin && addr <= AddrVRAMEnd {
		b.VRAM.Write(addr, v)
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		b.HRAM.Write(addr, v)
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		b.WRAM.Write(addr, v)
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.APU.Write(addr, v)
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		b.OAM.Write(addr, v)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.PPU.Write(addr, v)
	} else {
		if !b.inCoreDump {
			panicf("write to unmapped address %s", addr.Hex())
		}
	}
}

func (b *Bus) CoreDumpBegin() {
	b.inCoreDump = true
	b.BootROMLock.CountdownDisable = true
	b.BootROM.CountdownDisable = true
	b.VRAM.CountdownDisable = true
	b.HRAM.CountdownDisable = true
	b.WRAM.CountdownDisable = true
	b.APU.CountdownDisable = true
	b.OAM.CountdownDisable = true
	b.PPU.CountdownDisable = true
	b.CartridgeSlot.CountdownDisable = true
}

func (b *Bus) CoreDumpEnd() {
	b.inCoreDump = false
	b.BootROMLock.CountdownDisable = false
	b.BootROM.CountdownDisable = false
	b.VRAM.CountdownDisable = false
	b.HRAM.CountdownDisable = false
	b.WRAM.CountdownDisable = false
	b.APU.CountdownDisable = false
	b.OAM.CountdownDisable = false
	b.PPU.CountdownDisable = false
	b.CartridgeSlot.CountdownDisable = false
}

func (b *Bus) GetCounters(addr Addr) (uint64, uint64) {
	if addr <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			return b.CartridgeSlot.GetCounters(addr)
		} else {
			return b.BootROM.GetCounters(addr)
		}
	} else if addr <= AddrCartridgeBank0End {
		return b.CartridgeSlot.GetCounters(addr)
	} else if addr == AddrBootROMLock {
		return b.BootROMLock.GetCounters(addr)
	} else if addr >= AddrVRAMBegin && addr <= AddrVRAMEnd {
		return b.VRAM.GetCounters(addr)
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		return b.HRAM.GetCounters(addr)
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		return b.WRAM.GetCounters(addr)
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		return b.APU.GetCounters(addr)
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		return b.OAM.GetCounters(addr)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		return b.PPU.GetCounters(addr)
	} else if addr == AddrIF || addr == AddrIE {
		return b.Interrupts.GetCounters(addr)
	}
	if !b.inCoreDump {
		panicf("GetCounters for unmapped address %s", addr.Hex())
	}
	return 0, 0
}
