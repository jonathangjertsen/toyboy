package model

//go:generate go-enum --marshal --flag --values --nocomments

type Interrupts struct {
	mem []Data8

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

func NewInterrupts(mem []Data8) *Interrupts {
	return &Interrupts{
		mem: mem,
	}
}

func (ints *Interrupts) SetIME(v bool) {
	ints.IME = v
	if v {
		ints.IRQCheck()
	}
}

func (ints *Interrupts) PendInterrupt(in IntSource) {
	ints.IME = false
	ints.mem[AddrIF] &= ^in.Mask()
	ints.PendingInterrupt = in
}

func (ints *Interrupts) IRQSet(in IntSource) {
	if ints.mem[AddrIF]&in.Mask() != 0 {
		return
	}
	ints.mem[AddrIF] |= in.Mask()
	ints.IRQCheck()
}

func (ints *Interrupts) IRQCheck() {
	if !ints.IME {
		return
	}
	regIF := ints.mem[AddrIF]
	regIE := ints.mem[AddrIE]
	for is := IntSource(0); is < 5; is++ {
		if (regIF & regIE & is.Mask()) != 0 {
			ints.PendInterrupt(is)
			break
		}
	}
}
