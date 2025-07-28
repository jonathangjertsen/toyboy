package model

type PeriodCounter struct {
	periodDivider      Data16
	periodDividerReset Data16
}

func (pc *PeriodCounter) SetPeriodLow(v Data8) {
	// keep upper 3 bits, overwrite lower 8
	pc.periodDividerReset &= 0x0700
	pc.periodDividerReset |= join16(0x00, v)
}

func (pc *PeriodCounter) SetPeriodHigh(v Data8) {
	// keep lower 8 bits, overwrite upper 3
	pc.periodDividerReset &= 0x00ff
	pc.periodDividerReset |= join16(v&0x7, 0x00)
}

func (pc *PeriodCounter) clock() bool {
	pc.periodDivider++
	if pc.periodDivider == 0x800 {
		pc.periodDivider = pc.periodDividerReset
		return true
	}
	return false
}
