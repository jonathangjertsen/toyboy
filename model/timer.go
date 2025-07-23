package model

type Timer struct {
	Mem              MemoryRegion
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

var timerBitPos = [4]int{9, 3, 5, 7}

func NewTimer(clock *ClockRT, interrupts *Interrupts) *Timer {
	t := &Timer{
		Mem: NewMemoryRegion(clock, AddrTimerBegin, SizeTimer),
	}

	// https://hacktix.github.io/GBEDG/timers/#timer-operation
	// We only act on the rising edge of the clock, so divide by 2
	// and act on both rising and falling edges of that divided clock
	clock.Divide(2).AttachDevice(func(c Cycle) {
		t.DIV++
		t.Mem.Data[offsDIV] = t.DIV.MSB()
		tac := t.Mem.Data[offsTAC]
		bit := t.DIV.Bit(timerBitPos[tac&0x3])
		enable := tac.Bit(2)
		andResult := bit && enable
		if t.prevAndResult && !andResult {
			t.Mem.Data[offsTIMA]++
			if t.Mem.Data[offsTIMA] == 0 {
				t.preReloadCounter = 4
			}
		}
		t.prevAndResult = andResult
		if t.preReloadCounter > 0 {
			t.preReloadCounter--
			if t.preReloadCounter == 0 {
				t.Mem.Data[offsTIMA] = t.Mem.Data[offsTMA]
				interrupts.IRQSet(IntSourceTimer)
			}
		}
	})
	return t
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

func (t *Timer) GetCounters(addr Addr) (uint64, uint64) {
	return t.Mem.GetCounters(addr)
}
