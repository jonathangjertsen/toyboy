package model

type PeriodCounter struct {
	Counter Data16
	Reset   Data16
}

func (pc *PeriodCounter) SetPeriodLow(v Data8) {
	// keep upper 3 bits, overwrite lower 8
	pc.Reset &= 0x0700
	pc.Reset |= join16(0x00, v)
}

func (pc *PeriodCounter) SetPeriodHigh(v Data8) {
	// keep lower 8 bits, overwrite upper 3
	pc.Reset &= 0x00ff
	pc.Reset |= join16(v&0x7, 0x00)
}

func (pc *PeriodCounter) clock() bool {
	pc.Counter++
	if pc.Counter == 0x800 {
		pc.Counter = pc.Reset
		return true
	}
	return false
}
