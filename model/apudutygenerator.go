package model

type DutyGenerator struct {
	waveform Data8
	output   int8
}

func NewDutyGenerator() DutyGenerator {
	dg := DutyGenerator{}
	return dg
}

func (dg *DutyGenerator) SetDuty(v Data8) {
	switch (v >> 6) & 0x3 {
	case 0:
		dg.waveform = 0b1111_1110 // 12.5%
	case 1:
		dg.waveform = 0b0111_1110 // 25.0%
	case 2:
		dg.waveform = 0b0111_1000 // 50.0%
	case 3:
		dg.waveform = 0b1000_0001 // 75.0%
	}
}

func (dg *DutyGenerator) clock() {
	if dg.waveform&Bit0 != 0 {
		dg.output = 1
		dg.waveform = (dg.waveform >> 1) | 0x80
	} else {
		dg.output = -1
		dg.waveform >>= 1
	}
}
