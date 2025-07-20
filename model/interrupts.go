package model

type Interrupts struct {
	MemIF MemoryRegion
	MemIE MemoryRegion

	IME bool
	// TODO: move IME pend handling out of CPU code
	setIMENextCycle bool
}

func NewInterrupts(clk *ClockRT) *Interrupts {
	return &Interrupts{
		MemIF: NewMemoryRegion(clk, AddrIF, 1),
		MemIE: NewMemoryRegion(clk, AddrIE, 1),
	}
}

func (ints *Interrupts) Read(addr Addr) Data8 {
	switch Addr(addr) {
	case AddrIF:
		return ints.MemIF.Read(addr)
	case AddrIE:
		return ints.MemIE.Read(addr)
	}
	panicv(addr)
	return 0
}

func (ints *Interrupts) Write(addr Addr, v Data8) {
	switch Addr(addr) {
	case AddrIF:
		ints.MemIF.Write(addr, v)
		ints.IRQCheck()
	case AddrIE:
		ints.MemIE.Write(addr, v)
		ints.IRQCheck()
	default:
		panicv(addr)
	}
}

func (ints *Interrupts) SetIME(v bool) {
	ints.IME = v
	if v {
		ints.IRQCheck()
	}
}

func (ints *Interrupts) GetCounters(addr Addr) (uint64, uint64) {
	switch Addr(addr) {
	case AddrIF:
		return ints.MemIF.GetCounters(addr)
	case AddrIE:
		return ints.MemIE.GetCounters(addr)
	}
	panicv(addr)
	return 0, 0
}

func (ints *Interrupts) ExecInterrupt(in Data8) {
	ints.IME = false
	ints.MemIF.Data[0] &= ^in

	panic("ISR call not implemented")
}

func (ints *Interrupts) IRQSet(in Data8) {
	if ints.MemIF.Data[0]&in != 0 {
		return
	}
	ints.MemIF.Data[0] |= in
	ints.IRQCheck()
}

func (ints *Interrupts) IRQCheck() {
	if !ints.IME {
		return
	}
	regIF := ints.MemIF.Data[0]
	regIE := ints.MemIE.Data[0]
	for idx := Data8(0); idx < 5; idx++ {
		in := Data8(1 << idx)
		if (regIF & regIE & in) != 0 {
			ints.ExecInterrupt(in)
		}
	}
}
