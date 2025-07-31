package model

type Stat struct {
	Reg         Data8
	PrevStatInt bool
}

func (s *Stat) Write(v Data8) {
	s.Reg = maskedWrite(s.Reg, v, 0xf8)
}

func (s *Stat) SetMode(ints *Interrupts, mode PPUMode) {
	s.Reg = maskedWrite(s.Reg, Data8(mode), 0x7)
	s.CheckInterrupt(ints)
}

func (s *Stat) SetLYCEqLY(ints *Interrupts, equal bool) {
	if equal {
		s.Reg |= 1 << 2
	} else {
		s.Reg &= ^Data8(1 << 2)
	}
	s.CheckInterrupt(ints)
}

func (s *Stat) CheckInterrupt(ints *Interrupts) {
	statInt := false
	if s.Reg&0xb == 0x08 {
		// Mode 0 int selected and mode is 0
		statInt = true
	}
	if s.Reg&0x13 == 0x11 {
		// Mode 1 int selected and mode is 1
		statInt = true
	}
	if s.Reg&0x23 == 0x22 {
		// Mode 2 int selected and mode is 2
		statInt = true
	}
	if s.Reg&0x44 == 0x44 {
		// LYC==LY int selected and LYC==LY
		statInt = true
	}
	if !s.PrevStatInt && statInt {
		ints.IRQSet(IntSourceLCD)
	}
	s.PrevStatInt = statInt
}
