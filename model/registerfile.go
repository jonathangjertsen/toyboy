package model

type RegisterFile struct {
	A  Data8
	F  Data8
	B  Data8
	C  Data8
	D  Data8
	E  Data8
	H  Data8
	L  Data8
	PC Addr
	SP Addr
	IR Opcode

	TempZ Data8
	TempW Data8
}

func (regs *RegisterFile) setFlag(mask Data8, v bool) {
	if v {
		regs.F |= mask
	} else {
		regs.F &= ^mask
	}
}

func (regs *RegisterFile) getFlag(mask Data8) bool {
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

func (regs *RegisterFile) SetWZ(v Data16) {
	regs.TempW = v.MSB()
	regs.TempZ = v.LSB()
}

func (regs *RegisterFile) GetWZ() Data16 {
	return join16(regs.TempW, regs.TempZ)
}
