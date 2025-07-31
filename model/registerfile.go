package model

const (
	A = iota
	F
	B
	C
	D
	E
	H
	L
	PC
	SP
	IR
	TempZ
	TempW
)

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

const (
	FlagBitZ = 7
	FlagBitN = 6
	FlagBitH = 5
	FlagBitC = 4
)

func (regs *RegisterFile) setFlag(bit int, v bool) {
	mask := Data8(1 << bit)
	if v {
		regs.F |= mask
	} else {
		regs.F &= ^mask
	}
}

func (regs *RegisterFile) getFlag(bit uint) bool {
	return regs.F.Bit(bit)
}

func (regs *RegisterFile) GetFlagZ() bool {
	v := regs.getFlag(FlagBitZ)
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
	regs.setFlag(FlagBitZ, v)
}

func (regs *RegisterFile) GetFlagN() bool {
	return regs.getFlag(FlagBitN)
}

func (regs *RegisterFile) SetFlagN(v bool) {
	regs.setFlag(FlagBitN, v)
}

func (regs *RegisterFile) GetFlagH() bool {
	return regs.getFlag(FlagBitH)
}

func (regs *RegisterFile) SetFlagH(v bool) {
	regs.setFlag(FlagBitH, v)
}

func (regs *RegisterFile) GetFlagC() bool {
	return regs.getFlag(FlagBitC)
}

func (regs *RegisterFile) SetFlagC(v bool) {
	regs.setFlag(FlagBitC, v)
}

func (regs *RegisterFile) SetWZ(v Data16) {
	regs.TempW = v.MSB()
	regs.TempZ = v.LSB()
}

func (regs *RegisterFile) GetWZ() Data16 {
	return join16(regs.TempW, regs.TempZ)
}
