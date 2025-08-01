package model

type Bus struct {
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

func (b *Bus) Reset() {
	b.Address = 0
	b.Data = 0

	if b.Config.BootROM.Skip {
		b.BootROMLock.Lock()
	}
}

func (b *Bus) GetData() Data8 {
	return b.Data
}

func (b *Bus) WriteAddress(mem []Data8, addr Addr) {
	b.Address = addr
	b.Data = b.ProbeAddress(mem, addr)
}

func (b *Bus) ProbeAddress(mem []Data8, addr Addr) Data8 {
	if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		return b.APU.Read(addr)
	}
	if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		return b.PPU.Read(addr)
	}
	if addr == AddrP1 {
		return b.Joypad.Read(mem[AddrP1], addr)
	}
	return mem[addr]
}

func (b *Bus) ProbeRange(mem []Data8, begin, end Addr) []Data8 {
	if begin > end {
		return nil
	}
	if begin >= AddrAPUBegin && end <= AddrAPUEnd {
		return readRange(b.APU, begin, end)
	}
	if begin >= AddrPPUBegin && end <= AddrPPUEnd {
		return readRange(b.PPU, begin, end)
	}
	return mem[begin : end+1]
}

func readRange(device interface{ Read(Addr) Data8 }, begin, end Addr) []Data8 {
	out := make([]Data8, 0, end-begin+1)
	for addr := begin; addr <= end; addr++ {
		out = append(out, device.Read(addr))
	}
	return out
}

func (b *Bus) WriteData(mem []Data8, v Data8) {
	b.Data = v
	addr := b.Address

	if addr <= AddrBootROMEnd {
		if b.BootROMLock.BootOff {
			b.Cartridge.Write(mem, addr, v)
		}
		return
	}
	if addr <= AddrCartridgeBankNEnd {
		b.Cartridge.Write(mem, addr, v)
		return
	}
	mem[addr] = v

	if addr == AddrBootROMLock {
		b.BootROMLock.Write(addr, v)
	} else if addr == AddrP1 {
		b.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		b.Interrupts.IRQCheck()
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.APU.Write(addr, v)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.PPU.Write(addr, v, b.Interrupts)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		b.Timer.Write(addr, v)
	}
}
