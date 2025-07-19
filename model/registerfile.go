package model

type RegisterFile struct {
	A  uint8
	F  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8
	PC uint16
	SP uint16
	IR Opcode

	TempZ uint8
	TempW uint8
}

func (regs *RegisterFile) setFlag(mask uint8, v bool) {
	if v {
		regs.F |= mask
	} else {
		regs.F &= ^mask
	}
}

func (regs *RegisterFile) getFlag(mask uint8) bool {
	return (regs.F & mask) == mask
}

func (regs *RegisterFile) GetFlagZ() bool {
	v := regs.getFlag(0x80)
	return v
}

func (regs *RegisterFile) SetFlags(res ALUResult) {
	regs.SetFlagZ(res.Z())
	regs.SetFlagC(res.C)
	regs.SetFlagH(res.H)
	regs.SetFlagN(res.N)
}

func (regs *RegisterFile) SetFlagsAndA(res ALUResult) {
	regs.A = res.Value
	regs.SetFlags(res)
}

func (regs *RegisterFile) SetFlagZ(v bool) {
	regs.setFlag(0x80, v)
}

func (regs *RegisterFile) GetFlagN() bool {
	return regs.getFlag(0x40)
}

func (regs *RegisterFile) SetFlagN(v bool) {
	regs.setFlag(0x40, v)
}

func (regs *RegisterFile) GetFlagH() bool {
	return regs.getFlag(0x20)
}

func (regs *RegisterFile) SetFlagH(v bool) {
	regs.setFlag(0x20, v)
}

func (regs *RegisterFile) GetFlagC() bool {
	return regs.getFlag(0x10)
}

func (regs *RegisterFile) SetFlagC(v bool) {
	regs.setFlag(0x10, v)
}

func (regs *RegisterFile) SetWZ(v uint16) {
	regs.TempW = uint8(v >> 8)
	regs.TempZ = uint8(v)
}

func (regs *RegisterFile) GetWZ() uint16 {
	wz := (uint16(regs.TempW) << 8) | uint16(regs.TempZ)
	return wz
}
