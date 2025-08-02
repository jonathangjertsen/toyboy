package model

type Timer struct {
	DIV              Data16
	PrevAndResult    bool
	PreReloadCounter int
}

// Tick the DIV timer
func (t *Timer) tickDIV(gb *Gameboy) {
	// https://gbdev.io/pandocs/Audio_details.html#div-apu
	// A “DIV-APU” counter is increased every time DIV’s bit 4 (5 in double-speed mode) goes from 1 to 0
	div := t.DIV
	if (div&Bit12 == Bit12) && ((div+1)&Bit12 == 0) { // bit 4 set, bit 3 clear => next time bit 4 will go low
		gb.APU.incDIVAPU()
	}
	div++
	t.DIV = div

	gb.Mem[AddrDIV] = Data8(div >> 8)
	tac := gb.Mem[AddrTAC]
	var andResult bool
	var clockSelect Data16
	if tac&Bit2 != 0 {
		clockSelect = clockSelectBits[tac&(Bit0|Bit1)]
	}
	andResult = (div&clockSelect != 0)
	if t.PrevAndResult && !andResult {
		tima := &gb.Mem[AddrTIMA]
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
			gb.Mem[AddrTIMA] = gb.Mem[AddrTMA]
			gb.IRQSet(IntSourceTimer)
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
