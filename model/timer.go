package model

type Timer struct {
	Mem              MemoryRegion
	APU              *APU
	Interrupts       *Interrupts
	DIV              Data16
	prevAndResult    bool
	preReloadCounter int
}

const (
	offsDIV  = 0
	offsTIMA = 1
	offsTMA  = 2
	offsTAC  = 3
)

var timerBitMask = [4]Data16{1 << 9, 1 << 3, 1 << 5, 1 << 7}

func NewTimer(clock *ClockRT, apu *APU, interrupts *Interrupts) *Timer {
	t := &Timer{
		Mem:        NewMemoryRegion(clock, AddrTimerBegin, SizeTimer),
		APU:        apu,
		Interrupts: interrupts,
	}

	clock.timer = t
	return t
}

// Tick the DIV timer
// Runs every cycle, so this code path is extremely hot
func (t *Timer) tickDIVTimer() {
	// https://gbdev.io/pandocs/Audio_details.html#div-apu
	// A “DIV-APU” counter is increased every time DIV’s bit 4 (5 in double-speed mode) goes from 1 to 0
	div := t.DIV
	if div&(Bit4|Bit3) == Bit4 { // bit 4 set, bit 3 clear => next time bit 4 will go low
		t.APU.incDIVAPU()
	}
	div++
	t.DIV = div

	t.Mem.Data[offsDIV] = Data8(div >> 8)
	tac := t.Mem.Data[offsTAC]
	var andResult bool
	if tac&Bit2 != 0 {
		switch tac & (Bit0 | Bit1) {
		case 0:
			andResult = (div&Bit9 != 0)
		case 1:
			andResult = (div&Bit3 != 0)
		case 2:
			andResult = (div&Bit5 != 0)
		case 3:
			andResult = (div&Bit7 != 0)
		}
	}
	if t.prevAndResult && !andResult {
		tima := &t.Mem.Data[offsTIMA]
		*tima++
		if *tima == 0 {
			t.preReloadCounter = 4
		}
	}
	t.prevAndResult = andResult
	if t.preReloadCounter > 0 {
		t.preReloadCounter--
		if t.preReloadCounter == 0 {
			t.Mem.Data[offsTIMA] = t.Mem.Data[offsTMA]
			t.Interrupts.IRQSet(IntSourceTimer)
		}
	}
}

func (t *Timer) Write(addr Addr, v Data8) {
	switch addr {
	case AddrDIV:
		t.DIV = 0
	case AddrTIMA:
		// If TIMA is written to on the same T-cycle on which the reload from TMA occurs
		// the write is ignored and the value in TMA will be loaded into TIMA.
		if t.preReloadCounter == 1 {
			return
		}
		// The reload of the TMA value as well as the interrupt request can be aborted by writing any value to TIMA during the four T-cycles until it is supposed to be reloaded.
		// The TIMA register contains whatever value was written to it
		// even after the 4 T-cycles have elapsed and no interrupt will be requested.
		t.preReloadCounter = 0
	case AddrTAC:
	}
	t.Mem.Write(addr, v)
}

func (t *Timer) Read(addr Addr) Data8 {
	return t.Mem.Read(addr)
}
