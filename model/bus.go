package model

import "fmt"

type Bus struct {
	AddressSpace *AddressSpace

	Data    Data8
	Address Addr
	Config  *Config

	inCoreDump bool

	BootROMLock *BootROMLock
	APU         *APU
	PPU         *PPU
	Cartridge   *Cartridge
	Joypad      *Joypad
	Interrupts  *Interrupts
	Serial      *Serial
	Prohibited  *Prohibited
	Timer       *Timer
}

func NewBus(as *AddressSpace) *Bus {
	return &Bus{
		AddressSpace: as,
	}
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
	b.Data = b.ProbeAddress(addr)
}

func (b *Bus) ProbeAddress(addr Addr) Data8 {
	if addr <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			return b.Cartridge.Read(addr)
		} else {
			return b.AddressSpace[addr]
		}
	} else if addr <= AddrCartridgeBankNEnd {
		return b.Cartridge.Read(addr)
	} else if addr >= AddrVRAMBegin && addr <= AddrVRAMEnd {
		return b.AddressSpace[addr]
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		return b.AddressSpace[addr]
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		return b.AddressSpace[addr]
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		return b.APU.Read(addr)
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		return b.AddressSpace[addr]
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		return b.PPU.Read(addr)
	} else if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		return b.Prohibited.Read(addr)
	} else if (addr >= 0xff4c && addr <= 0xff4f) || (addr >= 0xff51 && addr <= 0xff70) {
		// GBC stuff
		return 0
	} else if addr >= 0xff71 && addr <= 0xff7f {
		return b.Prohibited.Read(addr)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		return b.Timer.Read(addr)
	} else if addr == AddrP1 {
		return b.Joypad.Read(addr)
	} else if addr == AddrIF || addr == AddrIE {
		return b.Interrupts.Read(addr)
	} else if addr == AddrBootROMLock {
		return b.AddressSpace[addr]
	} else if addr == AddrSB || addr == AddrSC {
		return b.Serial.Read(addr)
	} else {
		if !b.inCoreDump {
			//panicf("Read from unmapped address %s", addr.Hex())
		}
	}
	return b.Data
}

func (b *Bus) ProbeRange(begin, end Addr) []Data8 {
	if begin > end {
		return nil
	}
	len := end - begin + 1
	if end <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			return b.Cartridge.ReadRange(begin, end)
		} else {
			return b.AddressSpace[begin : end+1]
		}
	} else if end <= AddrCartridgeBankNEnd {
		return b.Cartridge.ReadRange(begin, end)
	} else if begin >= AddrVRAMBegin && end <= AddrVRAMEnd {
		return b.AddressSpace[begin : end+1]
	} else if begin >= AddrHRAMBegin && end <= AddrHRAMEnd {
		return b.AddressSpace[begin : end+1]
	} else if begin >= AddrWRAMBegin && end <= AddrWRAMEnd {
		return b.AddressSpace[begin : end+1]
	} else if begin >= AddrAPUBegin && end <= AddrAPUEnd {
		return readRange(b.APU, begin, end)
	} else if begin >= AddrOAMBegin && end <= AddrOAMEnd {
		return b.AddressSpace[begin : end+1]
	} else if begin >= AddrPPUBegin && end <= AddrPPUEnd {
		return readRange(b.PPU, begin, end)
	} else if begin >= AddrProhibitedBegin && end <= AddrProhibitedEnd {
		return b.Prohibited.FEA0toFEFF.Data[begin-AddrProhibitedBegin : end-AddrProhibitedBegin+1]
	} else if (begin >= 0xff4c && end <= 0xff4f) || (begin >= 0xff51 && end <= 0xff70) {
		// GBC stuff
		return nil
	} else if begin >= 0xff71 && end <= 0xff7f {
		return b.Prohibited.FF71toFF7F.Data[begin-0xff71 : end-0xff7f+1]
	} else if begin >= AddrTimerBegin && end <= AddrTimerEnd {
		return readRange(b.Timer, begin, end)
	} else if len == 1 && begin == AddrP1 {
		return readRange(b.Joypad, begin, end)
	} else if len == 1 && (begin == AddrIF || begin == AddrIE) {
		return readRange(b.Interrupts, begin, end)
	} else if len == 1 && begin == AddrBootROMLock {
		return b.AddressSpace[begin : end+1]
	} else if len == 1 && (begin == AddrSB || begin == AddrSC) {
		return readRange(b.Serial, begin, end)
	} else {
		if !b.inCoreDump {
			//panicf("Read from unmapped address %s", addr.Hex())
		}
	}
	return nil
}

func readRange(device interface{ Read(Addr) Data8 }, begin, end Addr) []Data8 {
	out := make([]Data8, 0, end-begin+1)
	for addr := begin; addr <= end; addr++ {
		out = append(out, device.Read(addr))
	}
	return out
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
		b.AddressSpace[b.Address] = v
		b.BootROMLock.Write(addr, v)
	} else if addr == AddrP1 {
		b.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		b.Interrupts.Write(addr, v)
	} else if addr >= AddrVRAMBegin && addr <= AddrVRAMEnd {
		b.AddressSpace[addr] = v
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		b.AddressSpace[addr] = v
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		b.AddressSpace[addr] = v
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.APU.Write(addr, v)
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		b.AddressSpace[addr] = v
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
