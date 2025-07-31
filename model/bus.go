package model

import "fmt"

type Bus struct {
	Mem []Data8

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
	Timer       *Timer
}

func NewBus(as []Data8) *Bus {
	return &Bus{
		Mem: as,
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
	if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		return b.APU.Read(addr)
	}
	if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		return b.PPU.Read(addr)
	}
	if addr == AddrP1 {
		return b.Joypad.Read(addr)
	}
	return b.Mem[addr]
}

func (b *Bus) ProbeRange(begin, end Addr) []Data8 {
	if begin > end {
		return nil
	}
	if begin >= AddrAPUBegin && end <= AddrAPUEnd {
		return readRange(b.APU, begin, end)
	}
	if begin >= AddrPPUBegin && end <= AddrPPUEnd {
		return readRange(b.PPU, begin, end)
	}
	return b.Mem[begin : end+1]
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
		if b.BootROMLock.BootOff {
			b.Mem[addr] = v
			b.Cartridge.Write(addr, v)
		}
	} else if addr <= AddrCartridgeBankNEnd {
		b.Cartridge.Write(addr, v)
	} else if addr == AddrBootROMLock {
		b.Mem[addr] = v
		b.BootROMLock.Write(addr, v)
	} else if addr == AddrP1 {
		b.Mem[addr] = v
		b.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		b.Mem[addr] = v
		b.Interrupts.IRQCheck()
	} else if addr >= AddrVRAMBegin && addr <= AddrVRAMEnd {
		b.Mem[addr] = v
	} else if addr >= AddrHRAMBegin && addr <= AddrHRAMEnd {
		b.Mem[addr] = v
	} else if addr >= AddrWRAMBegin && addr <= AddrWRAMEnd {
		b.Mem[addr] = v
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.Mem[addr] = v
		b.APU.Write(addr, v)
	} else if addr >= AddrWaveRAMBegin && addr <= AddrWaveRAMEnd {
		b.Mem[addr] = v
	} else if addr >= AddrOAMBegin && addr <= AddrOAMEnd {
		b.Mem[addr] = v
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.Mem[addr] = v
		b.PPU.Write(addr, v)
	} else if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		b.Mem[addr] = v
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		b.Mem[addr] = v
		b.Timer.Write(addr, v)
	} else if (addr >= 0xff4c && addr <= 0xff4f) || (addr >= 0xff51 && addr <= 0xff70) {
		// GBC stuff
	} else if addr >= 0xff71 && addr <= 0xff7f {
		b.Mem[addr] = v
	} else if addr == AddrSB || addr == AddrSC {
		b.Mem[addr] = v
	} else {
		if !b.inCoreDump {
			//panicf("write to unmapped address %s", addr.Hex())
		}
	}
}
