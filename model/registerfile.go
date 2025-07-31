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
	BC RegisterPair
	DE RegisterPair
	HL RegisterPair
	PC Addr
	SP Addr
	IR Opcode

	Temp     RegisterPair
	TempPtr  [2]*Data8
	TempCond bool
}

type RegisterPair struct {
	MSR Data8
	LSR Data8
}

func (rp RegisterPair) Data16() Data16 {
	return join16(rp.MSR, rp.LSR)
}

func (rp RegisterPair) Addr() Addr {
	return Addr(rp.Data16())
}

func (rp *RegisterPair) Inc() {
	rp.LSR++
	if rp.LSR == 0 {
		rp.MSR++
	}
}

func (rp *RegisterPair) Dec() {
	rp.LSR--
	if rp.LSR == 0xff {
		rp.MSR--
	}
}

func (rp *RegisterPair) Set(d Data16) {
	rp.MSR = d.MSB()
	rp.LSR = d.LSB()
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
