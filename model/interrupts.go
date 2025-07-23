package model

//go:generate go-enum --marshal --flag --values --nocomments

type Interrupts struct {
	MemIF MemoryRegion
	MemIE MemoryRegion

	IME bool

	// TODO: move IME pend handling out of CPU code
	setIMENextCycle bool

	PendingInterrupt IntSource
}

// ENUM(
// VBlank = 1
// LCD    = 2
// Timer  = 3
// Serial = 4
// Joypad = 5
// )
type IntSource uint8

func (is IntSource) Mask() Data8 {
	return Data8(1 << (is - 1))
}

func (is IntSource) ISR() Addr {
	switch is {
	case IntSourceVBlank:
		return 0x0040
	case IntSourceLCD:
		return 0x0048
	case IntSourceTimer:
		return 0x0050
	case IntSourceSerial:
		return 0x0058
	case IntSourceJoypad:
		return 0x0060
	}
	panicv(is)
	return 0
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

func (ints *Interrupts) PendInterrupt(in IntSource) {
	ints.IME = false
	ints.MemIF.Data[0] &= ^in.Mask()
	ints.PendingInterrupt = in
}

func (ints *Interrupts) IRQSet(in IntSource) {
	if ints.MemIF.Data[0]&in.Mask() != 0 {
		return
	}
	ints.MemIF.Data[0] |= in.Mask()
	ints.IRQCheck()
}

func (ints *Interrupts) IRQCheck() {
	if !ints.IME {
		return
	}
	regIF := ints.MemIF.Data[0]
	regIE := ints.MemIE.Data[0]
	for is := IntSource(0); is < 5; is++ {
		if (regIF & regIE & is.Mask()) != 0 {
			ints.PendInterrupt(is)
			break
		}
	}
}
