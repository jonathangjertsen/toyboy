package model

type Bus struct {
	Data    Data8
	Address Addr
	GB      *Gameboy
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
		return b.GB.APU.Read(addr)
	}
	if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		return b.GB.PPU.Read(addr)
	}
	if addr == AddrP1 {
		return b.GB.Joypad.Read(mem[AddrP1], addr)
	}
	return mem[addr]
}

func (b *Bus) ProbeRange(mem []Data8, begin, end Addr) []Data8 {
	if begin > end {
		return nil
	}
	if begin >= AddrAPUBegin && end <= AddrAPUEnd {
		return readRange(&b.GB.APU, begin, end)
	}
	if begin >= AddrPPUBegin && end <= AddrPPUEnd {
		return readRange(&b.GB.PPU, begin, end)
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
		if b.GB.BootROMLock.BootOff {
			b.GB.Cartridge.Write(mem, addr, v)
		}
		return
	}
	if addr <= AddrCartridgeBankNEnd {
		b.GB.Cartridge.Write(mem, addr, v)
		return
	}
	mem[addr] = v

	if addr == AddrBootROMLock {
		b.GB.BootROMLock.Write(mem, &b.GB.Debug, &b.GB.Cartridge, v)
	} else if addr == AddrP1 {
		b.GB.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		b.GB.Interrupts.IRQCheck(mem)
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		b.GB.APU.Write(addr, v)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		b.GB.PPU.Write(mem, addr, v, &b.GB.Interrupts)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		b.GB.Timer.Write(addr, v)
	}
}
