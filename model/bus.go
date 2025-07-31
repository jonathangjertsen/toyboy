package model

type Bus struct {
	Mem []Data8

	Data    Data8
	Address Addr

	Config      *Config
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

func (bus *Bus) LoadSave(save *SaveState, mem []Data8) {
	bus.Data = save.BusData
	bus.Address = save.BusAddress
	bus.Mem = mem
}

func (bus *Bus) Save(save *SaveState) {
	save.BusData = bus.Data
	save.BusAddress = bus.Address
}

func (b *Bus) Reset() {
	b.Address = 0
	b.Data = 0

	if b.Config.BootROM.Skip {
		b.BootROMLock.Lock()
	}
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
			b.Cartridge.Write(addr, v)
		}
		return
	}
	if addr <= AddrCartridgeBankNEnd {
		b.Cartridge.Write(addr, v)
		return
	}
	b.Mem[addr] = v

	if addr == AddrBootROMLock {
		b.BootROMLock.Write(addr, v)
	} else if addr == AddrP1 {
		b.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		b.Interrupts.IRQCheck()
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.APU.Write(addr, v)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.PPU.Write(addr, v)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		b.Timer.Write(addr, v)
	}
}
