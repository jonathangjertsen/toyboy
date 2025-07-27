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

func (b *Bus) PushState() (Addr, Data8) {
	return b.Address, b.Data
}

func (b *Bus) PopState(addr Addr, data Data8) {
	b.Address, b.Data = addr, data
}

func (b *Bus) BeginCoreDump() func() {
	b.inCoreDump = true
	return func() {
		b.inCoreDump = false
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
			b.Data = b.BootROM.Data[addr]
		}
	} else if addr <= AddrCartridgeBankNEnd {
		b.Data = b.Cartridge.Read(addr)
	} else if addr >= AddrVRAMBegin && addr <= AddrVRAMEnd {
		b.Data = b.VRAM.Data[addr-AddrVRAMBegin]
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		b.Data = b.HRAM.Data[addr-AddrHRAMBegin]
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		b.Data = b.WRAM.Data[addr-AddrWRAMBegin]
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.Data = b.APU.Read(addr)
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		b.Data = b.OAM.Data[addr-AddrOAMBegin]
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.Data = b.PPU.Read(addr)
	} else if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		b.Data = b.Prohibited.Read(addr)
	} else if (addr >= 0xff4c && addr <= 0xff4f) || (addr >= 0xff51 && addr <= 0xff70) {
		// GBC stuff
		b.Data = 0
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
			//panicf("Read from unmapped address %s", addr.Hex())
		}
	}
}

func (b *Bus) WriteData(v Data8) {
	b.Data = v
	addr := b.Address

	if addr <= AddrBootROMEnd {
		if !b.inCoreDump {
			if b.BootROMLock.BootOff {
				//panicf("Attempted write to cartridge (addr=%s v=%s)", addr.Hex(), v.Hex())
			} else {
				//panicf("Attempted write to bootrom (addr=%s v=%s)", addr.Hex(), v.Hex())
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
		b.VRAM.Data[addr-AddrVRAMBegin] = v
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		b.HRAM.Data[addr-AddrHRAMBegin] = v
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		b.WRAM.Data[addr-AddrWRAMBegin] = v
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.APU.Write(addr, v)
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		b.OAM.Data[addr-AddrOAMBegin] = v
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.PPU.Write(addr, v)
	} else if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		b.Prohibited.Write(addr, v)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		b.Timer.Write(addr, v)
	} else if (addr >= 0xff4c && addr <= 0xff4f) || (addr >= 0xff51 && addr <= 0xff70) {
		// GBC stuff
	} else if addr >= 0xff71 && addr <= 0xff7f {
		b.Prohibited.Write(addr, v)
	} else if addr == AddrSB || addr == AddrSC {
		b.Serial.Write(addr, v)
	} else {
		if !b.inCoreDump {
			//panicf("write to unmapped address %s", addr.Hex())
		}
	}
}
