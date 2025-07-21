package model

type Prohibited struct {
	FEA0toFEFF MemoryRegion
	FF71toFF7F MemoryRegion
}

func NewProhibited(clk *ClockRT) *Prohibited {
	return &Prohibited{
		FEA0toFEFF: NewMemoryRegion(clk, AddrProhibitedBegin, SizeProhibited),
		FF71toFF7F: NewMemoryRegion(clk, 0xff71, 0xf),
	}
}

func (p *Prohibited) Read(addr Addr) Data8 {
	if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		return p.FEA0toFEFF.Read(addr)
	} else if addr >= 0xff71 && addr <= 0xff7f {
		return p.FF71toFF7F.Read(addr)
	}
	panicf("unknown address %s", addr.Hex())
	return 0
}

func (p *Prohibited) Write(addr Addr, v Data8) {
	if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		p.FEA0toFEFF.Write(addr, v)
		return
	} else if addr >= 0xff71 && addr <= 0xff7f {
		p.FF71toFF7F.Write(addr, v)
		return
	}
	panicf("unknown address %s", addr.Hex())
}

func (p *Prohibited) GetCounters(addr Addr) (uint64, uint64) {
	if addr >= AddrProhibitedBegin && addr <= AddrProhibitedEnd {
		return p.FEA0toFEFF.GetCounters(addr)
	} else if addr >= 0xff71 && addr <= 0xff7f {
		return p.FF71toFF7F.GetCounters(addr)
	}
	panicf("unknown address %s", addr.Hex())
	return 0, 0
}
