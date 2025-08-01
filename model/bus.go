package model

type Bus struct {
	Data    Data8
	Address Addr
}

func (b *Bus) GetData() Data8 {
	return b.Data
}

func (b *Bus) WriteAddress(gb *Gameboy, addr Addr) {
	b.Address = addr
	b.Data = b.ProbeAddress(gb, addr)
}

func (b *Bus) ProbeAddress(gb *Gameboy, addr Addr) Data8 {
	if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		return gb.APU.Read(addr)
	}
	if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		return gb.PPU.Read(addr)
	}
	if addr == AddrP1 {
		return gb.Joypad.Read(gb.Mem[AddrP1], addr)
	}
	return gb.Mem[addr]
}

func (b *Bus) ProbeRange(gb *Gameboy, begin, end Addr) []Data8 {
	if begin > end {
		return nil
	}
	if begin >= AddrAPUBegin && end <= AddrAPUEnd {
		return readRange(&gb.APU, begin, end)
	}
	if begin >= AddrPPUBegin && end <= AddrPPUEnd {
		return readRange(&gb.PPU, begin, end)
	}
	return gb.Mem[begin : end+1]
}

func readRange(device interface{ Read(Addr) Data8 }, begin, end Addr) []Data8 {
	out := make([]Data8, 0, end-begin+1)
	for addr := begin; addr <= end; addr++ {
		out = append(out, device.Read(addr))
	}
	return out
}

func (b *Bus) WriteData(gb *Gameboy, v Data8) {
	b.Data = v
	addr := b.Address

	if addr <= AddrBootROMEnd {
		if gb.BootROMLock.BootOff {
			gb.Cartridge.Write(gb.Mem, addr, v)
		}
		return
	}
	if addr <= AddrCartridgeBankNEnd {
		gb.Cartridge.Write(gb.Mem, addr, v)
		return
	}
	gb.Mem[addr] = v

	if addr == AddrBootROMLock {
		gb.BootROMLock.Write(gb.Mem, &gb.Debug, &gb.Cartridge, v)
	} else if addr == AddrP1 {
		gb.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		gb.Interrupts.IRQCheck(gb.Mem)
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		gb.APU.Write(addr, v)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		gb.PPU.Write(gb.Mem, addr, v, &gb.Interrupts)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		gb.Timer.Write(addr, v)
	}
}
