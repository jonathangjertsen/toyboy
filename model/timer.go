package model

type Timer struct {
	DIV              Data16
	PrevAndResult    bool
	PreReloadCounter int
}

// Tick the DIV timer
func (t *Timer) tickDIV(mem []Data8, ints *Interrupts, apu *APU) {
	// https://gbdev.io/pandocs/Audio_details.html#div-apu
	// A “DIV-APU” counter is increased every time DIV’s bit 4 (5 in double-speed mode) goes from 1 to 0
	div := t.DIV
	if (div&Bit4 == Bit4) && ((div+1)&Bit4 == 0) { // bit 4 set, bit 3 clear => next time bit 4 will go low
		apu.incDIVAPU()
	}
	div++
	t.DIV = div

	mem[AddrDIV] = Data8(div >> 8)
	tac := mem[AddrTAC]
	var andResult bool
	var clockSelect Data16
	if tac&Bit2 != 0 {
		clockSelect = clockSelectBits[tac&(Bit0|Bit1)]
	}
	andResult = (div&clockSelect != 0)
	if t.PrevAndResult && !andResult {
		tima := &mem[AddrTIMA]
		if *tima == 0xFF {
			t.PreReloadCounter = 4
			*tima = 0
		} else {
			*tima++
		}
	}
	t.PrevAndResult = andResult
	if t.PreReloadCounter > 0 {
		t.PreReloadCounter--
		if t.PreReloadCounter == 0 {
			mem[AddrTIMA] = mem[AddrTMA]
			ints.IRQSet(mem, IntSourceTimer)
		}
	}
}

var clockSelectBits = [4]Data16{
	Bit9,
	Bit3,
	Bit5,
	Bit7,
}

func (t *Timer) Write(addr Addr, v Data8) {
	switch addr {
	case AddrDIV:
		t.DIV = 0
	case AddrTIMA:
		// If TIMA is written to on the same T-cycle on which the reload from TMA occurs
		// the write is ignored and the value in TMA will be loaded into TIMA.
		if t.PreReloadCounter == 1 {
			return
		}
		// The reload of the TMA value as well as the interrupt request can be aborted by writing any value to TIMA during the four T-cycles until it is supposed to be reloaded.
		// The TIMA register contains whatever value was written to it
		// even after the 4 T-cycles have elapsed and no interrupt will be requested.
		t.PreReloadCounter = 0
	case AddrTAC:
	}
}
