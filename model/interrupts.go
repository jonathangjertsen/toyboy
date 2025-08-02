package model

//go:generate go-enum --marshal --flag --values --nocomments

type Interrupts struct {
	IME bool

	InISR bool

	// TODO: move IME pend handling out of CPU code
	SetIMENextCycle bool

	PendingInterrupt IntSource
}

// ENUM(
// None   = 0
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

func (gb *Gameboy) SetIME(v bool) {
	gb.Interrupts.IME = v
	if v {
		gb.IRQCheck()
	}
}

func (gb *Gameboy) PendInterrupt(in IntSource) {
	gb.Interrupts.IME = false
	gb.Mem[AddrIF] &= ^in.Mask()
	gb.Interrupts.PendingInterrupt = in
	gb.CPU.UOpCycle = 1
}

func (gb *Gameboy) IRQSet(in IntSource) {
	if gb.Mem[AddrIF]&in.Mask() != 0 {
		return
	}
	gb.Mem[AddrIF] |= in.Mask()
	gb.IRQCheck()
}

func (gb *Gameboy) IRQCheck() {
	regIF := gb.Mem[AddrIF]
	regIE := gb.Mem[AddrIE]
	for is := IntSource(0); is < 5; is++ {
		if (regIF & regIE & is.Mask()) != 0 {
			if gb.Interrupts.IME {
				gb.PendInterrupt(is)
			}
			gb.CPU.Halted = false
			break
		}
	}
}

func (gb *Gameboy) applyPendingIME() {
	if gb.Interrupts.SetIMENextCycle {
		gb.Interrupts.SetIMENextCycle = false
		gb.SetIME(true)
	}
}
