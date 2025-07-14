package model

//go:generate go-enum --marshal --flag --values --nocomments

import (
	"fmt"
	"slices"
	"sync"
)

var coreDebugEvents = []string{
	// "NotImplemented",
	"Panic",
	// "PreFetch",
	// "ExecDone",
	// "ExecBegin",
	//"SetBC",
	// "Handler",
	// "SetPC",
	//"ExecBeginCPn",
	//"GetHL",
	// "IncPC",
	//"SetFlagZ",
	// "WriteAddressBus",
	// "PeriphRead",
	//"GetFlagZ",
	//"PeriphWrite",
	//"ExecCBOp",
	// "SetHL",
	//"CPn",
	//"SetA",
	//"SetC",
	//"Watchfffc",
}

var coreDumpEvents = []string{
	// "Watchfffc",
}

// ENUM(
// Nop      = 0x00,
// LDBCnn   = 0x01,
// INCBC    = 0x03,
// INCB     = 0x04,
// DECB     = 0x05,
// LDBn     = 0x06,
// DECBC    = 0x0b,
// INCC     = 0x0c,
// DECC     = 0x0d,
// LDCn     = 0x0e,
// LDDEnn   = 0x11,
// INCDE    = 0x13,
// INCD     = 0x14,
// DECD     = 0x15,
// LDDn     = 0x16,
// RLA      = 0x17,
// JRe      = 0x18,
// LDADE    = 0x1a,
// DECDE    = 0x1b,
// INCE     = 0x1c,
// DECE     = 0x1d,
// LDEn     = 0x1e,
// JRNZe    = 0x20,
// LDHLnn   = 0x21,
// LDHLAInc = 0x22,
// INCHL    = 0x23,
// INCH     = 0x24,
// DECH     = 0x25,
// LDHn     = 0x26,
// JRZe     = 0x28,
// DECHL    = 0x2b,
// INCL     = 0x2C,
// DECL     = 0x2D,
// LDLn     = 0x2e,
// JRNCe    = 0x30,
// LDSPnn   = 0x31,
// LDHLADec = 0x32,
// INCSP    = 0x33,
// JRCe     = 0x38,
// DECSP    = 0x3b,
// INCA     = 0x3c,
// DECA     = 0x3d,
// LDAn     = 0x3e,
// LDBB     = 0x40,
// LDBC     = 0x41,
// LDBD     = 0x42,
// LDBE     = 0x43,
// LDBH     = 0x44,
// LDBL     = 0x45,
// LDBHL    = 0x46,
// LDBA     = 0x47,
// LDCB     = 0x48,
// LDCC     = 0x49,
// LDCD     = 0x4a,
// LDCE     = 0x4b,
// LDCH     = 0x4c,
// LDCL     = 0x4d,
// LDCHL    = 0x4e,
// LDCA     = 0x4f,
// LDDB     = 0x50,
// LDDC     = 0x51,
// LDDD     = 0x52,
// LDDE     = 0x53,
// LDDH     = 0x54,
// LDDL     = 0x55,
// LDDHL    = 0x56,
// LDDA     = 0x57,
// LDEB     = 0x58,
// LDEC     = 0x59,
// LDED     = 0x5a,
// LDEE     = 0x5b,
// LDEH     = 0x5c,
// LDEL     = 0x5d,
// LDEHL    = 0x5e,
// LDEA     = 0x5f,
// LDHB     = 0x60,
// LDHC     = 0x61,
// LDHD     = 0x62,
// LDHE     = 0x63,
// LDHH     = 0x64,
// LDHL     = 0x65,
// LDHHL    = 0x66,
// LDHA     = 0x67,
// LDLB     = 0x68,
// LDLC     = 0x69,
// LDLD     = 0x6a,
// LDLE     = 0x6b,
// LDLH     = 0x6c,
// LDLL     = 0x6d,
// LDLHL    = 0x6e,
// LDLA     = 0x6f,
// LDHLB    = 0x70,
// LDHLC    = 0x71,
// LDHLD    = 0x72,
// LDHLE    = 0x73,
// LDHLH    = 0x74,
// LDHLL    = 0x75,
// HALT     = 0x76,
// LDHLA    = 0x77,
// LDAB     = 0x78,
// LDAC     = 0x79,
// LDAD     = 0x7a,
// LDAE     = 0x7b,
// LDAH     = 0x7c,
// LDAL     = 0x7d,
// LDAHL    = 0x7e,
// LDAA     = 0x7f,
// ADDB     = 0x80,
// ADDC     = 0x81,
// ADDD     = 0x82,
// ADDE     = 0x83,
// ADDH     = 0x84,
// ADDL     = 0x85,
// ADDHL    = 0x86,
// ADDA     = 0x87,
// ADCB     = 0x88,
// ADCC     = 0x89,
// ADCD     = 0x8A,
// ADCE     = 0x8B,
// ADCH     = 0x8c,
// ADCL     = 0x8d,
// ADCHL    = 0x8e,
// ADCA     = 0x8f,
// SUBB     = 0x90,
// SUBC     = 0x91,
// SUBD     = 0x92,
// SUBE     = 0x93,
// SUBH     = 0x94,
// SUBL     = 0x95,
// SUBHL    = 0x96,
// SUBA     = 0x97,
// SBCB     = 0x98,
// SBCC     = 0x99,
// SBCD     = 0x9A,
// SBCE     = 0x9B,
// SBCH     = 0x9c,
// SBCL     = 0x9d,
// SBCHL    = 0x9e,
// SBCA     = 0x9f,
// ANDB     = 0xA0,
// ANDC     = 0xA1,
// ANDD     = 0xA2,
// ANDE     = 0xA3,
// ANDH     = 0xA4,
// ANDL     = 0xA5,
// ANDHL    = 0xA6,
// ANDA     = 0xA7,
// XORB     = 0xA8,
// XORC     = 0xA9,
// XORD     = 0xAA,
// XORE     = 0xAB,
// XORH     = 0xAc,
// XORL     = 0xAd,
// XORHL    = 0xAe,
// XORA     = 0xAf,
// ORB      = 0xB0,
// ORC      = 0xB1,
// ORD      = 0xB2,
// ORE      = 0xB3,
// ORH      = 0xB4,
// ORL      = 0xB5,
// ORHL     = 0xB6,
// ORA      = 0xB7,
// CPB      = 0xB8,
// CPC      = 0xB9,
// CPD      = 0xBA,
// CPE      = 0xBB,
// CPH      = 0xBc,
// CPL      = 0xBd,
// CPHL     = 0xBe,
// CPA      = 0xBf,
// POPBC    = 0xC1,
// JPNZnn   = 0xC2,
// JPnn     = 0xC3,
// PUSHBC   = 0xC5,
// RET      = 0xC9,
// JPZnn    = 0xCA,
// CB       = 0xCB,
// CALLnn   = 0xCD,
// JPCnn    = 0xDA,
// JPNCnn   = 0xD2,
// LDHnA    = 0xE0,
// LDHCA    = 0xE2,
// LDnnA    = 0xEA,
// LDHAn    = 0xF0,
// LDAnn    = 0xFA,
// DI       = 0xF3,
// EI       = 0xFB,
// CPn      = 0xFE,
// )
type Opcode uint8

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
}

type Interrupts struct {
	IF              uint8
	IE              uint8
	IME             bool
	setIMENextCycle bool
}

type CPU struct {
	m *sync.Mutex

	PHI         *Clock
	DataBus     uint8
	AddressBus  uint16
	Peripherals []Peripheral

	Regs       RegisterFile
	Interrupts Interrupts

	Z    uint8
	W    uint8
	CBOp CBOp

	machineCycle int

	clockCycle                 Cycle
	inCoreDump                 bool
	wroteToAddressBusThisCycle bool

	rewindBuffer    [16]ExecLogEntry
	rewindBufferIdx int
}

type ExecLogEntry struct {
	PC     uint16
	Opcode Opcode
}

func (cpu *CPU) SetHL(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetHL must be called on rising edge")
	}
	cpu.Debug("SetHL", "0x%04x", v)
	cpu.Regs.H = uint8(v >> 8)
	cpu.Regs.L = uint8(v)
}

func (cpu *CPU) SetBC(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetBC must be called on rising edge")
	}
	cpu.Debug("SetBC", "0x%04x", v)
	cpu.Regs.B = uint8(v >> 8)
	cpu.Regs.C = uint8(v)
}

func (cpu *CPU) SetDE(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetDE must be called on rising edge")
	}
	cpu.Debug("SetDE", "0x%04x", v)
	cpu.Regs.D = uint8(v >> 8)
	cpu.Regs.E = uint8(v)
}

func (cpu *CPU) SetSP(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetSP must be called on rising edge")
	}
	cpu.Debug("SetSP", "0x%04x", v)
	cpu.Regs.SP = v
}

func (cpu *CPU) GetA() uint8 {
	v := cpu.Regs.A
	cpu.Debug("GetA", "0x%02x", v)
	return v
}

func (cpu *CPU) GetB() uint8 {
	v := cpu.Regs.B
	cpu.Debug("GetB", "0x%02x", v)
	return v
}

func (cpu *CPU) GetC() uint8 {
	v := cpu.Regs.C
	cpu.Debug("GetC", "0x%02x", v)
	return v
}

func (cpu *CPU) GetD() uint8 {
	v := cpu.Regs.D
	cpu.Debug("GetD", "0x%02x", v)
	return v
}

func (cpu *CPU) GetE() uint8 {
	v := cpu.Regs.E
	cpu.Debug("GetE", "0x%02x", v)
	return v
}

func (cpu *CPU) GetH() uint8 {
	v := cpu.Regs.H
	cpu.Debug("GetE", "0x%02x", v)
	return v
}

func (cpu *CPU) GetL() uint8 {
	v := cpu.Regs.L
	cpu.Debug("GetE", "0x%02x", v)
	return v
}

func (cpu *CPU) GetBC() uint16 {
	v := join16(cpu.Regs.B, cpu.Regs.C)
	cpu.Debug("GetBC", "0x%04x", v)
	return v
}

func (cpu *CPU) GetDE() uint16 {
	v := join16(cpu.Regs.D, cpu.Regs.E)
	cpu.Debug("GetDE", "0x%04x", v)
	return v
}

func (cpu *CPU) GetHL() uint16 {
	v := join16(cpu.Regs.H, cpu.Regs.L)
	cpu.Debug("GetHL", "0x%04x", v)
	return v
}

func join16(msb, lsb uint8) uint16 {
	return (uint16(msb) << 8) | uint16(lsb)
}

func msb(w uint16) uint8 {
	return uint8((w >> 8) & 0xff)
}

func lsb(w uint16) uint8 {
	return uint8(w & 0xff)
}

func (cpu *CPU) setFlag(mask uint8, v bool) {
	if v {
		cpu.Regs.F |= mask
	} else {
		cpu.Regs.F &= ^mask
	}
}

func (cpu *CPU) getFlag(mask uint8) bool {
	return (cpu.Regs.F & mask) == mask
}

func (cpu *CPU) GetFlagZ() bool {
	v := cpu.getFlag(0x80)
	cpu.Debug("GetFlagZ", "F=%x Z=%v", cpu.Regs.F, v)
	return v
}

func (cpu *CPU) SetFlagZ(v bool) {
	cpu.setFlag(0x80, v)
	cpu.Debug("SetFlagZ", "F=%x Z=%v", cpu.Regs.F, v)
}

func (cpu *CPU) GetFlagN() bool {
	return cpu.getFlag(0x40)
}

func (cpu *CPU) SetFlagN(v bool) {
	cpu.Debug("SetFlagN", "%v", v)
	cpu.setFlag(0x40, v)
}

func (cpu *CPU) GetFlagH() bool {
	return cpu.getFlag(0x20)
}

func (cpu *CPU) SetFlagH(v bool) {
	cpu.Debug("SetFlagH", "%v", v)
	cpu.setFlag(0x20, v)
}

func (cpu *CPU) TODOFlagN() {
	cpu.Debug("NotImplemented", "Set BCD subtraction flag")
}

func (cpu *CPU) TODOFlagH() {
	cpu.Debug("NotImplemented", "Set BCD half-carry flag")
}

func (cpu *CPU) GetFlagC() bool {
	return cpu.getFlag(0x10)
}

func (cpu *CPU) SetFlagC(v bool) {
	cpu.Debug("SetFlagC", "%v", v)
	cpu.setFlag(0x10, v)
}

func (cpu *CPU) SetPC(pc uint16) {
	if cpu.clockCycle.Falling {
		panic("SetPC must be called on rising edge")
	}
	cpu.Debug("SetPC", "0x%04x", pc)
	cpu.Regs.PC = pc
}

func (cpu *CPU) IncPC() {
	if cpu.clockCycle.Falling {
		panic("IncPC must be called on rising edge")
	}
	cpu.Regs.PC++
	cpu.Debug("IncPC", "0x%04x", cpu.Regs.PC)
}

func (cpu *CPU) Debug(event string, f string, v ...any) {
	if cpu.inCoreDump {
		return
	}
	if slices.Contains(coreDebugEvents, event) || slices.Contains(coreDumpEvents, event) {
		dir := "^"
		if cpu.clockCycle.Falling {
			dir = "v"
		}
		fmt.Printf("%d %s PC=0x%04x %v mcycle=%v | %s | ", cpu.clockCycle.C, dir, cpu.Regs.PC, cpu.Regs.IR, cpu.machineCycle, event)
		fmt.Printf(f, v...)
		fmt.Printf("\n")
	}
	if slices.Contains(coreDumpEvents, event) {
		cpu.Dump()
	}
}

func (cpu *CPU) SetWZ(v uint16) {
	cpu.W = uint8(v >> 8)
	cpu.Z = uint8(v)
	wz := (uint16(cpu.W) << 8) | uint16(cpu.Z)
	cpu.Debug("SetWZ", "%v", wz)
}

func (cpu *CPU) GetWZ() uint16 {
	wz := (uint16(cpu.W) << 8) | uint16(cpu.Z)
	return wz
}

// ENUM(
// RLC,
// RRC,
// RL,
// RR,
// SLA,
// SRA,
// SWAP,
// SRL,
// Bit0,
// Bit1,
// Bit2,
// Bit3,
// Bit4,
// Bit5,
// Bit6,
// Bit7,
// Res0,
// Res1,
// Res2,
// Res3,
// Res4,
// Res5,
// Res6,
// Res7,
// Set0,
// Set1,
// Set2,
// Set3,
// Set4,
// Set5,
// Set6,
// Set7,
// )
type cb uint8

// ENUM(B, C, D, E, H, L, IndirectHL, A)
type CBTarget uint8

type CBOp struct {
	Op     cb
	Target CBTarget
}

type Peripheral interface {
	Name() string
	Range() (start, size uint16)
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
}

type edge struct {
	Cycle   int
	Falling bool
}

type InstructionHandling func(c *CPU, e edge) bool

var handlers = map[Opcode]InstructionHandling{
	OpcodeNop: singleCycle(func(cpu *CPU) {
		if cpu.clockCycle.C > 0 {
			panic("unexpected nop")
		}
	}),
	OpcodeLDAA: singleCycle(func(cpu *CPU) {}),
	OpcodeLDBB: singleCycle(func(cpu *CPU) {}),
	OpcodeLDCC: singleCycle(func(cpu *CPU) {}),
	OpcodeLDDD: singleCycle(func(cpu *CPU) {}),
	OpcodeLDEE: singleCycle(func(cpu *CPU) {}),
	OpcodeLDHH: singleCycle(func(cpu *CPU) {}),
	OpcodeLDLL: singleCycle(func(cpu *CPU) {}),
	OpcodeLDAB: singleCycle(func(cpu *CPU) { cpu.Regs.A = cpu.Regs.B }),
	OpcodeLDAC: singleCycle(func(cpu *CPU) { cpu.Regs.A = cpu.Regs.C }),
	OpcodeLDAD: singleCycle(func(cpu *CPU) { cpu.Regs.A = cpu.Regs.D }),
	OpcodeLDAE: singleCycle(func(cpu *CPU) { cpu.Regs.A = cpu.Regs.E }),
	OpcodeLDAH: singleCycle(func(cpu *CPU) { cpu.Regs.A = cpu.Regs.H }),
	OpcodeLDAL: singleCycle(func(cpu *CPU) { cpu.Regs.A = cpu.Regs.L }),
	OpcodeLDBA: singleCycle(func(cpu *CPU) { cpu.Regs.B = cpu.Regs.A }),
	OpcodeLDBC: singleCycle(func(cpu *CPU) { cpu.Regs.B = cpu.Regs.C }),
	OpcodeLDBD: singleCycle(func(cpu *CPU) { cpu.Regs.B = cpu.Regs.D }),
	OpcodeLDBE: singleCycle(func(cpu *CPU) { cpu.Regs.B = cpu.Regs.E }),
	OpcodeLDBH: singleCycle(func(cpu *CPU) { cpu.Regs.B = cpu.Regs.H }),
	OpcodeLDBL: singleCycle(func(cpu *CPU) { cpu.Regs.B = cpu.Regs.L }),
	OpcodeLDCA: singleCycle(func(cpu *CPU) { cpu.Regs.C = cpu.Regs.A }),
	OpcodeLDCB: singleCycle(func(cpu *CPU) { cpu.Regs.C = cpu.Regs.B }),
	OpcodeLDCD: singleCycle(func(cpu *CPU) { cpu.Regs.C = cpu.Regs.D }),
	OpcodeLDCE: singleCycle(func(cpu *CPU) { cpu.Regs.C = cpu.Regs.E }),
	OpcodeLDCH: singleCycle(func(cpu *CPU) { cpu.Regs.C = cpu.Regs.H }),
	OpcodeLDCL: singleCycle(func(cpu *CPU) { cpu.Regs.C = cpu.Regs.L }),
	OpcodeLDDA: singleCycle(func(cpu *CPU) { cpu.Regs.D = cpu.Regs.A }),
	OpcodeLDDB: singleCycle(func(cpu *CPU) { cpu.Regs.D = cpu.Regs.B }),
	OpcodeLDDC: singleCycle(func(cpu *CPU) { cpu.Regs.D = cpu.Regs.C }),
	OpcodeLDDE: singleCycle(func(cpu *CPU) { cpu.Regs.D = cpu.Regs.E }),
	OpcodeLDDH: singleCycle(func(cpu *CPU) { cpu.Regs.D = cpu.Regs.H }),
	OpcodeLDDL: singleCycle(func(cpu *CPU) { cpu.Regs.D = cpu.Regs.L }),
	OpcodeLDEA: singleCycle(func(cpu *CPU) { cpu.Regs.E = cpu.Regs.A }),
	OpcodeLDEB: singleCycle(func(cpu *CPU) { cpu.Regs.E = cpu.Regs.B }),
	OpcodeLDEC: singleCycle(func(cpu *CPU) { cpu.Regs.E = cpu.Regs.C }),
	OpcodeLDED: singleCycle(func(cpu *CPU) { cpu.Regs.E = cpu.Regs.D }),
	OpcodeLDEH: singleCycle(func(cpu *CPU) { cpu.Regs.E = cpu.Regs.H }),
	OpcodeLDEL: singleCycle(func(cpu *CPU) { cpu.Regs.E = cpu.Regs.L }),
	OpcodeLDHA: singleCycle(func(cpu *CPU) { cpu.Regs.H = cpu.Regs.A }),
	OpcodeLDHB: singleCycle(func(cpu *CPU) { cpu.Regs.H = cpu.Regs.B }),
	OpcodeLDHC: singleCycle(func(cpu *CPU) { cpu.Regs.H = cpu.Regs.C }),
	OpcodeLDHD: singleCycle(func(cpu *CPU) { cpu.Regs.H = cpu.Regs.D }),
	OpcodeLDHE: singleCycle(func(cpu *CPU) { cpu.Regs.H = cpu.Regs.E }),
	OpcodeLDHL: singleCycle(func(cpu *CPU) { cpu.Regs.H = cpu.Regs.L }),
	OpcodeLDLA: singleCycle(func(cpu *CPU) { cpu.Regs.L = cpu.Regs.A }),
	OpcodeLDLB: singleCycle(func(cpu *CPU) { cpu.Regs.L = cpu.Regs.B }),
	OpcodeLDLC: singleCycle(func(cpu *CPU) { cpu.Regs.L = cpu.Regs.C }),
	OpcodeLDLD: singleCycle(func(cpu *CPU) { cpu.Regs.L = cpu.Regs.D }),
	OpcodeLDLE: singleCycle(func(cpu *CPU) { cpu.Regs.L = cpu.Regs.E }),
	OpcodeLDLH: singleCycle(func(cpu *CPU) { cpu.Regs.L = cpu.Regs.H }),
	OpcodeRLA: singleCycle(func(cpu *CPU) {
		a := cpu.Regs.A
		bit7 := a & 0x80
		a <<= 1
		if cpu.GetFlagC() {
			a |= 0x01
		}
		cpu.SetFlagZ(false)
		cpu.SetFlagN(false)
		cpu.SetFlagH(false)
		cpu.SetFlagC(bit7 != 0)
		cpu.Regs.A = a
	}),
	OpcodeORA:  orreg(func(cpu *CPU) uint8 { return cpu.Regs.A }),
	OpcodeORB:  orreg(func(cpu *CPU) uint8 { return cpu.Regs.B }),
	OpcodeORC:  orreg(func(cpu *CPU) uint8 { return cpu.Regs.C }),
	OpcodeORD:  orreg(func(cpu *CPU) uint8 { return cpu.Regs.D }),
	OpcodeORE:  orreg(func(cpu *CPU) uint8 { return cpu.Regs.E }),
	OpcodeORH:  orreg(func(cpu *CPU) uint8 { return cpu.Regs.H }),
	OpcodeORL:  orreg(func(cpu *CPU) uint8 { return cpu.Regs.L }),
	OpcodeANDA: andreg(func(cpu *CPU) uint8 { return cpu.Regs.A }),
	OpcodeANDB: andreg(func(cpu *CPU) uint8 { return cpu.Regs.B }),
	OpcodeANDC: andreg(func(cpu *CPU) uint8 { return cpu.Regs.C }),
	OpcodeANDD: andreg(func(cpu *CPU) uint8 { return cpu.Regs.D }),
	OpcodeANDE: andreg(func(cpu *CPU) uint8 { return cpu.Regs.E }),
	OpcodeANDH: andreg(func(cpu *CPU) uint8 { return cpu.Regs.H }),
	OpcodeANDL: andreg(func(cpu *CPU) uint8 { return cpu.Regs.L }),
	OpcodeXORA: xorreg(func(cpu *CPU) uint8 { return cpu.Regs.A }),
	OpcodeXORB: xorreg(func(cpu *CPU) uint8 { return cpu.Regs.B }),
	OpcodeXORC: xorreg(func(cpu *CPU) uint8 { return cpu.Regs.C }),
	OpcodeXORD: xorreg(func(cpu *CPU) uint8 { return cpu.Regs.D }),
	OpcodeXORE: xorreg(func(cpu *CPU) uint8 { return cpu.Regs.E }),
	OpcodeXORH: xorreg(func(cpu *CPU) uint8 { return cpu.Regs.H }),
	OpcodeXORL: xorreg(func(cpu *CPU) uint8 { return cpu.Regs.L }),
	OpcodeSUBA: subreg(func(cpu *CPU) uint8 { return cpu.Regs.A }),
	OpcodeSUBB: subreg(func(cpu *CPU) uint8 { return cpu.Regs.B }),
	OpcodeSUBC: subreg(func(cpu *CPU) uint8 { return cpu.Regs.C }),
	OpcodeSUBD: subreg(func(cpu *CPU) uint8 { return cpu.Regs.D }),
	OpcodeSUBE: subreg(func(cpu *CPU) uint8 { return cpu.Regs.E }),
	OpcodeSUBH: subreg(func(cpu *CPU) uint8 { return cpu.Regs.H }),
	OpcodeSUBL: subreg(func(cpu *CPU) uint8 { return cpu.Regs.L }),
	OpcodeADDA: addreg(func(cpu *CPU) uint8 { return cpu.Regs.A }),
	OpcodeADDB: addreg(func(cpu *CPU) uint8 { return cpu.Regs.B }),
	OpcodeADDC: addreg(func(cpu *CPU) uint8 { return cpu.Regs.C }),
	OpcodeADDD: addreg(func(cpu *CPU) uint8 { return cpu.Regs.D }),
	OpcodeADDE: addreg(func(cpu *CPU) uint8 { return cpu.Regs.E }),
	OpcodeADDH: addreg(func(cpu *CPU) uint8 { return cpu.Regs.H }),
	OpcodeADDL: addreg(func(cpu *CPU) uint8 { return cpu.Regs.L }),
	OpcodeDECA: decreg(func(cpu *CPU) *uint8 { return &cpu.Regs.A }),
	OpcodeDECB: decreg(func(cpu *CPU) *uint8 { return &cpu.Regs.B }),
	OpcodeDECC: decreg(func(cpu *CPU) *uint8 { return &cpu.Regs.C }),
	OpcodeDECD: decreg(func(cpu *CPU) *uint8 { return &cpu.Regs.D }),
	OpcodeDECE: decreg(func(cpu *CPU) *uint8 { return &cpu.Regs.E }),
	OpcodeDECH: decreg(func(cpu *CPU) *uint8 { return &cpu.Regs.H }),
	OpcodeDECL: decreg(func(cpu *CPU) *uint8 { return &cpu.Regs.L }),
	OpcodeINCA: increg(func(cpu *CPU) *uint8 { return &cpu.Regs.A }),
	OpcodeINCB: increg(func(cpu *CPU) *uint8 { return &cpu.Regs.B }),
	OpcodeINCC: increg(func(cpu *CPU) *uint8 { return &cpu.Regs.C }),
	OpcodeINCD: increg(func(cpu *CPU) *uint8 { return &cpu.Regs.D }),
	OpcodeINCE: increg(func(cpu *CPU) *uint8 { return &cpu.Regs.E }),
	OpcodeINCH: increg(func(cpu *CPU) *uint8 { return &cpu.Regs.H }),
	OpcodeINCL: increg(func(cpu *CPU) *uint8 { return &cpu.Regs.L }),
	OpcodeDI: singleCycle(func(cpu *CPU) {
		cpu.Interrupts.setIMENextCycle = false
		cpu.Interrupts.IME = false
	}),
	OpcodeEI: singleCycle(func(cpu *CPU) {
		cpu.Interrupts.setIMENextCycle = true
	}),
	OpcodeJRe: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		// TODO: this impl is not exactly correct
		case edge{2, false}:
		case edge{2, true}:
		case edge{3, false}:
			cpu.SetPC(uint16(int16(cpu.Regs.PC) + int16(int8(cpu.Z))))
			return true
		case edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeJPnn: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
		case edge{3, true}:
		case edge{4, false}:
			cpu.SetPC(cpu.GetWZ())
			return true
		case edge{4, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeJRZe:   jrcce(func(cpu *CPU) bool { return cpu.GetFlagZ() }),
	OpcodeJRCe:   jrcce(func(cpu *CPU) bool { return cpu.GetFlagC() }),
	OpcodeJRNZe:  jrcce(func(cpu *CPU) bool { return !cpu.GetFlagZ() }),
	OpcodeJRNCe:  jrcce(func(cpu *CPU) bool { return !cpu.GetFlagC() }),
	OpcodeJPCnn:  jpccnn(func(cpu *CPU) bool { return cpu.GetFlagC() }),
	OpcodeJPNCnn: jpccnn(func(cpu *CPU) bool { return !cpu.GetFlagC() }),
	OpcodeJPZnn:  jpccnn(func(cpu *CPU) bool { return cpu.GetFlagZ() }),
	OpcodeJPNZnn: jpccnn(func(cpu *CPU) bool { return !cpu.GetFlagZ() }),
	OpcodeINCBC:  iduOp(func(cpu *CPU) { cpu.SetBC(cpu.GetBC() + 1) }),
	OpcodeINCDE:  iduOp(func(cpu *CPU) { cpu.SetDE(cpu.GetDE() + 1) }),
	OpcodeINCHL:  iduOp(func(cpu *CPU) { cpu.SetHL(cpu.GetHL() + 1) }),
	OpcodeINCSP:  iduOp(func(cpu *CPU) { cpu.Regs.SP++ }),
	OpcodeDECBC:  iduOp(func(cpu *CPU) { cpu.SetBC(cpu.GetBC() - 1) }),
	OpcodeDECDE:  iduOp(func(cpu *CPU) { cpu.SetDE(cpu.GetDE() - 1) }),
	OpcodeDECHL:  iduOp(func(cpu *CPU) { cpu.SetHL(cpu.GetHL() - 1) }),
	OpcodeDECSP:  iduOp(func(cpu *CPU) { cpu.Regs.SP-- }),
	OpcodeCALLnn: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
			cpu.SetSP(cpu.Regs.SP - 1)
		case edge{3, true}:
		case edge{4, false}:
			cpu.writeAddressBus(cpu.Regs.SP)
			cpu.SetSP(cpu.Regs.SP - 1)
		case edge{4, true}:
			cpu.writeDataBus(msb(cpu.Regs.PC))
		case edge{5, false}:
			cpu.writeAddressBus(cpu.Regs.SP)
		case edge{5, true}:
			cpu.writeDataBus(lsb(cpu.Regs.PC))
		case edge{6, false}:
			cpu.SetPC(cpu.GetWZ())
			return true
		case edge{6, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeRET: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.SP)
			cpu.SetSP(cpu.Regs.SP + 1)
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.SP)
			cpu.SetSP(cpu.Regs.SP + 1)
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
			cpu.SetPC(cpu.GetWZ())
		case edge{3, true}:
		case edge{4, false}, edge{4, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodePUSHBC: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.SetSP(cpu.Regs.SP - 1)
			cpu.writeAddressBus(cpu.Regs.SP)
		case edge{1, true}:
			cpu.writeDataBus(cpu.Regs.B)
		case edge{2, false}:
			cpu.SetSP(cpu.Regs.SP - 1)
			cpu.writeAddressBus(cpu.Regs.SP)
		case edge{2, true}:
			cpu.writeDataBus(cpu.Regs.C)
		case edge{3, false}:
		case edge{3, true}:
		case edge{4, false}:
			return true
		case edge{4, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodePOPBC: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.SP)
			cpu.SetSP(cpu.Regs.SP + 1)
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.SP)
			cpu.SetSP(cpu.Regs.SP + 1)
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
			cpu.SetBC(cpu.GetWZ())
			return true
		case edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDBCnn:   ldxxnn(func(cpu *CPU, wz uint16) { cpu.SetBC(wz) }),
	OpcodeLDDEnn:   ldxxnn(func(cpu *CPU, wz uint16) { cpu.SetDE(wz) }),
	OpcodeLDHLnn:   ldxxnn(func(cpu *CPU, wz uint16) { cpu.SetHL(wz) }),
	OpcodeLDSPnn:   ldxxnn(func(cpu *CPU, wz uint16) { cpu.SetSP(wz) }),
	OpcodeLDHLAInc: ldhla(func(cpu *CPU) { cpu.SetHL(cpu.GetHL() + 1) }),
	OpcodeLDHLADec: ldhla(func(cpu *CPU) { cpu.SetHL(cpu.GetHL() - 1) }),
	OpcodeLDHLA:    ldhla(func(cpu *CPU) {}),
	OpcodeLDHCA: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(join16(0xff, cpu.Regs.C))
		case edge{1, true}:
			cpu.writeDataBus(cpu.Regs.A)
		case edge{2, false}, edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeADDHL: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.GetHL())
		case edge{1, true}:
			data := cpu.readDataBus()
			carry := uint16(cpu.Regs.A)+uint16(data) > 256
			result := cpu.Regs.A + data
			cpu.Regs.A = result
			cpu.SetFlagZ(result == 0)
			cpu.SetFlagN(false)
			cpu.TODOFlagH()
			cpu.SetFlagC(carry)
		case edge{2, false}, edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeCPHL: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.GetHL())
		case edge{1, true}:
			data := cpu.readDataBus()
			carry := data > cpu.Regs.A
			result := cpu.Regs.A - data
			cpu.SetFlagZ(result == 0)
			cpu.SetFlagN(true)
			cpu.TODOFlagH()
			cpu.SetFlagC(carry)
		case edge{2, false}, edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDnnA: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
			cpu.writeAddressBus(cpu.GetWZ())
		case edge{3, true}:
			cpu.writeDataBus(cpu.Regs.A)
		case edge{4, false}, edge{4, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDAnn: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
			cpu.writeAddressBus(cpu.GetWZ())
		case edge{3, true}:
			cpu.Regs.A = cpu.readDataBus()
		case edge{4, false}, edge{4, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeCPn: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			return true
		case edge{2, true}:
			cpu.Debug("CPn", "A=%02x n=%02x", cpu.Regs.A, cpu.Z)
			carry := cpu.Regs.A < cpu.Z
			result := cpu.Regs.A - cpu.Z
			cpu.SetFlagZ(result == 0)
			cpu.SetFlagN(true)
			cpu.TODOFlagH()
			cpu.SetFlagC(carry)
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDHnA: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(join16(0xff, cpu.Z))
			cpu.IncPC()
		case edge{2, true}:
			cpu.writeDataBus(cpu.Regs.A)
		case edge{3, false}, edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDHAn: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(join16(0xff, cpu.Z))
		case edge{2, true}:
			cpu.Z = cpu.readDataBus()
		case edge{3, false}:
			return true
		case edge{3, true}:
			cpu.Regs.A = cpu.Z
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDADE: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(join16(cpu.Regs.D, cpu.Regs.E))
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			return true
		case edge{2, true}:
			cpu.Regs.A = cpu.Z
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDAn: ldrn(func(cpu *CPU, z uint8) { cpu.Regs.A = z }),
	OpcodeLDBn: ldrn(func(cpu *CPU, z uint8) { cpu.Regs.B = z }),
	OpcodeLDCn: ldrn(func(cpu *CPU, z uint8) { cpu.Regs.C = z }),
	OpcodeLDDn: ldrn(func(cpu *CPU, z uint8) { cpu.Regs.D = z }),
	OpcodeLDEn: ldrn(func(cpu *CPU, z uint8) { cpu.Regs.E = z }),
	OpcodeLDHn: ldrn(func(cpu *CPU, z uint8) { cpu.Regs.H = z }),
	OpcodeLDLn: ldrn(func(cpu *CPU, z uint8) { cpu.Regs.L = z }),
	OpcodeCB: func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			opcode := cpu.readDataBus()
			cpu.CBOp = CBOp{Op: cb((opcode & 0xf8) >> 3), Target: CBTarget(opcode & 0x7)}
		case edge{2, false}:
			return true
		case edge{2, true}:
			var val uint8
			switch cpu.CBOp.Target {
			case CBTargetB:
				val = cpu.Regs.B
			case CBTargetC:
				val = cpu.Regs.C
			case CBTargetD:
				val = cpu.Regs.D
			case CBTargetE:
				val = cpu.Regs.E
			case CBTargetH:
				val = cpu.Regs.H
			case CBTargetL:
				val = cpu.Regs.L
			case CBTargetIndirectHL:
				panic("indirect thru HL not implemented")
			case CBTargetA:
				val = cpu.Regs.A
			default:
				panic("unknown CBOp target")
			}
			cpu.Debug("ExecCBOp", "op=%v target=%v", cpu.CBOp.Op, cpu.CBOp.Target)
			switch cpu.CBOp.Op {
			case CbRL:
				bit7 := val & 0x80
				val <<= 1
				if cpu.GetFlagC() {
					val |= 0x01
				}
				cpu.SetFlagZ(val == 0)
				cpu.SetFlagN(false)
				cpu.SetFlagH(false)
				cpu.SetFlagC(bit7 != 0)
			case CbBit0:
				cbbit(cpu, val, 0x01)
			case CbBit1:
				cbbit(cpu, val, 0x02)
			case CbBit2:
				cbbit(cpu, val, 0x04)
			case CbBit3:
				cbbit(cpu, val, 0x08)
			case CbBit4:
				cbbit(cpu, val, 0x10)
			case CbBit5:
				cbbit(cpu, val, 0x20)
			case CbBit6:
				cbbit(cpu, val, 0x40)
			case CbBit7:
				cbbit(cpu, val, 0x80)
			default:
				panicf("unknown op = %+v", cpu.CBOp)
			}
			switch cpu.CBOp.Target {
			case CBTargetB:
				cpu.Regs.B = val
			case CBTargetC:
				cpu.Regs.C = val
			case CBTargetD:
				cpu.Regs.D = val
			case CBTargetE:
				cpu.Regs.E = val
			case CBTargetH:
				cpu.Regs.H = val
			case CBTargetL:
				cpu.Regs.L = val
			case CBTargetIndirectHL:
				panic("indirect thru HL not implemented")
			case CBTargetA:
				cpu.Regs.A = val
			default:
				panic("unknown CBOp target")
			}
			return true
		default:
			panicv(e)
		}
		return false
	},
}

func singleCycle(f func(cpu *CPU)) func(cpu *CPU, e edge) bool {
	return func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			return true
		case edge{1, true}:
			f(cpu)
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func jrcce(f func(cpu *CPU) bool) func(cpu *CPU, e edge) bool {
	return func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			if f(cpu) {
			} else {
				return true
			}
		case edge{2, true}:
			if f(cpu) {
				newPC := uint16(int16(cpu.Regs.PC) + int16(int8(cpu.Z)))
				cpu.SetWZ(newPC)
			} else {
				return true
			}
		case edge{3, false}:
			if f(cpu) {
				cpu.SetPC(cpu.GetWZ())
				return true
			} else {
				panicv(e)
			}
		case edge{3, true}:
			if f(cpu) {
				return true
			} else {
				panicv(e)
			}
		default:
			panicv(e)
		}
		return false
	}
}

func jpccnn(f func(cpu *CPU) bool) func(cpu *CPU, e edge) bool {
	return func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
			if f(cpu) {
				cpu.SetPC(cpu.GetWZ())
			} else {
				return true
			}
		case edge{3, true}:
			if f(cpu) {
			} else {
				return true
			}
		case edge{4, false}:
			if f(cpu) {
				return true
			} else {
				panicv(e)
			}
		case edge{4, true}:
			if f(cpu) {
				return true
			} else {
				panicv(e)
			}
		default:
			panicv(e)
		}
		return false
	}
}

func andreg(f func(cpu *CPU) uint8) func(cpu *CPU, e edge) bool {
	return singleCycle(func(cpu *CPU) {
		reg := f(cpu)
		cpu.Regs.A &= reg
		cpu.SetFlagZ(cpu.Regs.A == 0)
		cpu.SetFlagN(false)
		cpu.SetFlagH(true)
		cpu.SetFlagC(false)
	})
}

func xorreg(f func(cpu *CPU) uint8) func(cpu *CPU, e edge) bool {
	return singleCycle(func(cpu *CPU) {
		reg := f(cpu)
		cpu.Regs.A ^= reg
		cpu.SetFlagZ(cpu.Regs.A == 0)
		cpu.SetFlagN(false)
		cpu.SetFlagH(false)
		cpu.SetFlagC(false)
	})
}

func orreg(f func(cpu *CPU) uint8) func(cpu *CPU, e edge) bool {
	return singleCycle(func(cpu *CPU) {
		reg := f(cpu)
		cpu.Regs.A |= reg
		cpu.SetFlagZ(cpu.Regs.A == 0)
		cpu.SetFlagN(false)
		cpu.SetFlagH(false)
		cpu.SetFlagC(false)
	})
}

func addreg(f func(cpu *CPU) uint8) func(cpu *CPU, e edge) bool {
	return singleCycle(func(cpu *CPU) {
		reg := f(cpu)
		carry := uint16(reg)+uint16(cpu.Regs.A) > 256
		cpu.Regs.A += reg
		cpu.SetFlagZ(cpu.Regs.A == 0)
		cpu.SetFlagN(false)
		cpu.TODOFlagH()
		cpu.SetFlagC(carry)
	})
}

func subreg(f func(cpu *CPU) uint8) func(cpu *CPU, e edge) bool {
	return singleCycle(func(cpu *CPU) {
		reg := f(cpu)
		carry := reg > cpu.Regs.A
		cpu.Regs.A -= reg
		cpu.SetFlagZ(cpu.Regs.A == 0)
		cpu.SetFlagN(true)
		cpu.TODOFlagH()
		cpu.SetFlagC(carry)
	})
}

func decreg(f func(cpu *CPU) *uint8) func(cpu *CPU, e edge) bool {
	return singleCycle(func(cpu *CPU) {
		reg := f(cpu)
		*reg -= 1
		cpu.SetFlagZ(*reg == 0)
		cpu.SetFlagN(true)
		cpu.TODOFlagH()
	})
}

func increg(f func(cpu *CPU) *uint8) func(cpu *CPU, e edge) bool {
	return singleCycle(func(cpu *CPU) {
		reg := f(cpu)
		*reg += 1
		cpu.SetFlagZ(*reg == 0)
		cpu.SetFlagN(false)
		cpu.TODOFlagH()
	})
}

func iduOp(f func(cpu *CPU)) func(cpu *CPU, e edge) bool {
	return func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			f(cpu)
		case edge{1, true}:
		case edge{2, false}, edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func ldrn(f func(cpu *CPU, z uint8)) func(cpu *CPU, e edge) bool {
	return func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.IncPC()
			return true
		case edge{2, true}:
			f(cpu, cpu.Z)
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func ldhla(f func(cpu *CPU)) func(cpu *CPU, e edge) bool {
	return func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.GetHL())
			f(cpu)
		case edge{1, true}:
			cpu.writeDataBus(cpu.Regs.A)
		case edge{2, false}:
			return true
		case edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func ldxxnn(f func(cpu *CPU, wz uint16)) func(cpu *CPU, e edge) bool {
	return func(cpu *CPU, e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Z = cpu.readDataBus()
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.W = cpu.readDataBus()
		case edge{3, false}:
			f(cpu, join16(cpu.W, cpu.Z))
			return true
		case edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func cbbit(cpu *CPU, val, mask uint8) {
	cpu.SetFlagZ(val&mask == 0)
	cpu.SetFlagN(false)
	cpu.SetFlagH(true)
}

func NewCPU(
	phi *Clock,
) *CPU {
	cpu := &CPU{
		m:   &sync.Mutex{},
		PHI: phi,
	}
	phi.AddRiseCallback(cpu.fsm)
	phi.AddFallCallback(cpu.fsm)
	return cpu
}

func (cpu *CPU) AttachPeripheral(p Peripheral) {
	cpu.m.Lock()
	cpu.Peripherals = append(cpu.Peripherals, p)
	cpu.m.Unlock()
}

func (cpu *CPU) Dump() {
	cpu.inCoreDump = true
	defer func() { cpu.inCoreDump = false }()

	fmt.Printf("\n--------\nCore dump:\n")
	fmt.Printf("PC = 0x%04x\n", cpu.Regs.PC)
	fmt.Printf("SP = 0x%04x\n", cpu.Regs.SP)
	fmt.Printf("A  =   0x%02x\n", cpu.Regs.A)
	fmt.Printf("F  =   0x%02x (Z=%v, H=%v, N=%v C=%v)\n", cpu.Regs.F, cpu.GetFlagZ(), cpu.GetFlagH(), cpu.GetFlagN(), cpu.GetFlagC())
	fmt.Printf("B  =   0x%02x\n", cpu.Regs.B)
	fmt.Printf("C  =   0x%02x\n", cpu.Regs.C)
	fmt.Printf("D  =   0x%02x\n", cpu.Regs.D)
	fmt.Printf("E  =   0x%02x\n", cpu.Regs.E)
	fmt.Printf("H  =   0x%02x\n", cpu.Regs.H)
	fmt.Printf("L  =   0x%02x\n", cpu.Regs.L)
	fmt.Printf("W  =   0x%02x\n", cpu.W)
	fmt.Printf("Z  =   0x%02x\n", cpu.Z)
	fmt.Printf("IR =   0x%02x\n", uint8(cpu.Regs.IR))
	fmt.Printf("--------\n")
	fmt.Printf("Code (PC highlighted)\n")
	start := uint16(0)
	if cpu.Regs.PC > 0x40 {
		start = cpu.Regs.PC - 0x40
	}
	end := uint16(0xffff)
	if cpu.Regs.PC < 0xffff-0x40 {
		end = cpu.Regs.PC + 0x40
	}
	cpu.memdump(start, end, cpu.Regs.PC-1)
	fmt.Printf("--------\n")
	fmt.Printf("HRAM (SP highlighted):\n")
	cpu.memdump(0xff80, 0xfffe, cpu.Regs.SP)
	fmt.Printf("--------\n")
	fmt.Printf("Last executed instructions:\n")
	for i := (cpu.rewindBufferIdx + 1) % len(cpu.rewindBuffer); i != cpu.rewindBufferIdx; i = (i + 1) % len(cpu.rewindBuffer) {
		fmt.Printf("[PC=%04x] %s\n", cpu.rewindBuffer[i].PC, cpu.rewindBuffer[i].Opcode)
	}
	fmt.Printf("--------\n")
}

func (cpu *CPU) memdump(start, end, highlight uint16) {
	alignedStart := (start / 0x10) * 0x10
	for addr := alignedStart; addr < start; addr++ {
		if addr%0x10 == 0 {
			fmt.Printf("\n %04x |", addr)
		}

		fmt.Printf(" .. ")
	}

	for addr := start; addr <= end; addr++ {
		if addr%0x10 == 0 {
			fmt.Printf("\n %04x |", addr)
		}
		cpu.writeAddressBus(addr)
		if highlight == addr {
			fmt.Printf("[%02x]", cpu.DataBus)
		} else {
			fmt.Printf(" %02x ", cpu.DataBus)
		}
	}

	alignedEnd := (end/0x10)*0x10 + 0x10 - 1
	for addr := end; addr < alignedEnd; addr++ {
		fmt.Printf(" .. ")
	}
	fmt.Printf("\n")
}

func (cpu *CPU) fsm(c Cycle) {
	defer func() {
		if e := recover(); e != nil {
			cpu.Debug("Panic", "%v", e)
			cpu.Dump()
			panic(e)
		}
	}()
	cpu.wroteToAddressBusThisCycle = false

	cpu.clockCycle = c

	cpu.applyPendingIME()

	var fetch bool
	if c.C > 0 {
		opcode := cpu.Regs.IR
		if handler, ok := handlers[opcode]; ok {
			e := edge{cpu.machineCycle, c.Falling}
			cpu.Debug("Handler", "e=%v", e)
			fetch = handler(cpu, e)
		} else {
			panicf("not implemented opcode %v", opcode)
		}
	} else {
		// initial instruction
		fetch = true
	}

	if fetch {
		if !c.Falling {
			cpu.Debug("PreFetch", "PC=%04x", cpu.Regs.PC)
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		} else {
			cpu.Debug("ExecDone", "")
			cpu.W = 0
			cpu.Z = 0
			cpu.machineCycle = 1
			cpu.Regs.IR = Opcode(cpu.DataBus)
			cpu.Debug("ExecBegin", "%s", cpu.Regs.IR)
			cpu.Debug(fmt.Sprintf("ExecBegin%s", cpu.Regs.IR), "")

			cpu.rewindBuffer[cpu.rewindBufferIdx] = ExecLogEntry{
				PC:     cpu.Regs.PC - 1,
				Opcode: cpu.Regs.IR,
			}
			cpu.rewindBufferIdx++
			if cpu.rewindBufferIdx >= len(cpu.rewindBuffer) {
				cpu.rewindBufferIdx = 0
			}
		}
	} else if c.Falling {
		cpu.machineCycle++
	}
}

func (cpu *CPU) writeAddressBus(addr uint16) {
	if !cpu.inCoreDump {
		if cpu.clockCycle.Falling {
			panic("writeAddressBus must be called on rising edge")
		}
		if cpu.wroteToAddressBusThisCycle {
			panic("more than one call to writeAddressBus this cycle")
		}
	}
	cpu.wroteToAddressBusThisCycle = true
	cpu.Debug("WriteAddressBus", "0x%04x", addr)
	cpu.AddressBus = addr
	for _, p := range cpu.Peripherals {
		start, size := p.Range()
		if uint32(start)+uint32(size-1) > 0xffff {
			panic("")
		}
		if size > 0 && addr >= start && addr <= start+(size-1) {
			cpu.Debug("PeriphRead", "%s @ 0x%04x", p.Name(), addr)
			v := p.Read(addr)
			cpu.DataBus = v
			return
		}
	}
	for _, p := range cpu.Peripherals {
		start, size := p.Range()
		cpu.Debug("Panic", "start=%0x size=%0x last=%0x", start, size, start+(size-1))
	}
	panicf("no peripheral mapped to 0x%x", addr)
}

func (cpu *CPU) writeDataBus(v uint8) {
	if !cpu.clockCycle.Falling {
		panic("writeDataBus must be called on falling edge")
	}
	cpu.DataBus = v
	addr := cpu.AddressBus
	for _, p := range cpu.Peripherals {
		start, size := p.Range()
		if size > 0 && addr >= start && addr <= start+(size-1) {
			p.Write(addr, v)
			cpu.Debug("PeriphWrite", "0x%02x to %s @ 0x%04x,", v, p.Name(), addr)
			cpu.Debug(fmt.Sprintf("Watch%04x", addr), "wrote %02x", v)
			cpu.DataBus = v
			return
		}
	}
}

func (cpu *CPU) readDataBus() uint8 {
	if !cpu.clockCycle.Falling {
		panic("readDataBus must be called on falling edge")
	}
	return cpu.DataBus
}

func (cpu *CPU) applyPendingIME() {
	if cpu.Interrupts.setIMENextCycle {
		cpu.Interrupts.setIMENextCycle = false
		cpu.Interrupts.IME = true
	}
}
