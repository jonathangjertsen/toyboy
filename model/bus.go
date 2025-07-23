package model

import "fmt"

type Bus struct {
	Data    Data8
	Address Addr
	Config  *Config

	inCoreDump bool

	BootROMLock *BootROMLock
	BootROM     *MemoryRegion
	VRAM        *MemoryRegion
	HRAM        *MemoryRegion
	WRAM        *MemoryRegion
	APU         *APU
	OAM         *MemoryRegion
	PPU         *PPU
	Cartridge   *Cartridge
	Joypad      *Joypad
	Interrupts  *Interrupts
	Serial      *Serial
	Prohibited  *Prohibited
	Timer       *Timer
}

func (b *Bus) GetPeripheral(ptr any) {
	switch p := ptr.(type) {
	case **PPU:
		*p = b.PPU
	case **APU:
		*p = b.APU
	default:
		panic(fmt.Errorf("no peripheral of type %T", ptr))
	}
}

func (b *Bus) Reset() {
	b.Address = 0
	b.Data = 0
	b.inCoreDump = false

	if b.Config.BootROM.Skip {
		b.BootROMLock.Lock()
	}
}

func (b *Bus) PushState() func() {
	addr := b.Address
	data := b.Data
	pop := func() {
		b.Address = addr
		b.Data = data
	}
	return pop
}

func (b *Bus) BeginCoreDump() func() {
	b.inCoreDump = true
	b.BootROMLock.CountdownDisable = true
	b.BootROM.CountdownDisable = true
	b.VRAM.CountdownDisable = true
	b.HRAM.CountdownDisable = true
	b.WRAM.CountdownDisable = true
	b.APU.CountdownDisable = true
	b.OAM.CountdownDisable = true
	b.PPU.MemoryRegion.CountdownDisable = true
	b.Prohibited.FEA0toFEFF.CountdownDisable = true
	b.Cartridge.Bank0.CountdownDisable = true
	b.Timer.Mem.CountdownDisable = true
	return func() {
		b.inCoreDump = false
		b.BootROMLock.CountdownDisable = false
		b.BootROM.CountdownDisable = false
		b.VRAM.CountdownDisable = false
		b.HRAM.CountdownDisable = false
		b.WRAM.CountdownDisable = false
		b.APU.CountdownDisable = false
		b.OAM.CountdownDisable = false
		b.PPU.MemoryRegion.CountdownDisable = false
		b.Prohibited.FEA0toFEFF.CountdownDisable = false
		b.Cartridge.Bank0.CountdownDisable = false
		b.Timer.Mem.CountdownDisable = false
	}
}

func (b *Bus) InCoreDump() bool {
	return b.inCoreDump
}

func (b *Bus) GetAddress() Addr {
	return b.Address
}

func (b *Bus) GetData() Data8 {
	return b.Data
}

func (b *Bus) WriteAddress(addr Addr) {
	b.Address = addr

	if addr <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			b.Data = b.Cartridge.Read(addr)
		} else {
			b.Data = b.BootROM.Read(addr)
		}
	} else if addr <= AddrCartridgeBankNEnd {
		b.Data = b.Cartridge.Read(addr)
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
	} else if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		b.Data = b.Prohibited.Read(addr)
	} else if addr >= 0xff71 && addr <= 0xff7f {
		b.Data = b.Prohibited.Read(addr)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		b.Data = b.Timer.Read(addr)
	} else if addr == AddrP1 {
		b.Data = b.Joypad.Read(addr)
	} else if addr == AddrIF || addr == AddrIE {
		b.Data = b.Interrupts.Read(addr)
	} else if addr == AddrBootROMLock {
		b.Data = b.BootROMLock.Read(addr)
	} else if addr == AddrSB || addr == AddrSC {
		b.Data = b.Serial.Read(addr)
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
		if !b.inCoreDump {
			if b.BootROMLock.BootOff {
				panicf("Attempted write to cartridge (addr=%s v=%s)", addr.Hex(), v.Hex())
			} else {
				panicf("Attempted write to bootrom (addr=%s v=%s)", addr.Hex(), v.Hex())
			}
		}
	} else if addr <= AddrCartridgeBankNEnd {
		b.Cartridge.Write(addr, v)
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
	} else if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		b.Prohibited.Write(addr, v)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		b.Timer.Write(addr, v)
	} else if addr >= 0xff71 && addr <= 0xff7f {
		b.Prohibited.Write(addr, v)
	} else if addr == AddrSB || addr == AddrSC {
		b.Serial.Write(addr, v)
	} else {
		if !b.inCoreDump {
			panicf("write to unmapped address %s", addr.Hex())
		}
	}
}

func (b *Bus) GetCounters(addr Addr) (uint64, uint64) {
	if addr <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			return b.Cartridge.GetCounters(addr)
		} else {
			return b.BootROM.GetCounters(addr)
		}
	} else if addr <= AddrCartridgeBankNEnd {
		return b.Cartridge.GetCounters(addr)
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
		return b.PPU.MemoryRegion.GetCounters(addr)
	} else if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		return b.Prohibited.GetCounters(addr)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		return b.Timer.GetCounters(addr)
	} else if addr >= 0xff71 && addr <= 0xff7f {
		return b.Prohibited.GetCounters(addr)
	} else if addr == AddrIF || addr == AddrIE {
		return b.Interrupts.GetCounters(addr)
	} else if addr == AddrSB || addr == AddrSC {
		return b.Serial.GetCounters(addr)
	}
	if !b.inCoreDump {
		panicf("GetCounters for unmapped address %s", addr.Hex())
	}
	return 0, 0
}
