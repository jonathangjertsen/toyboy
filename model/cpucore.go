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
	"ExecDone",
	"ExecBegin",
	//"SetBC",
	// "Handler",
	// "SetPC",
	//"ExecBeginCPn",
	//"GetHL",
	// "IncPC",
	//"SetFlagZ",
	// "WriteAddressBus",
	//"PeriphRead",
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
// INCB     = 0x04,
// DECB     = 0x05,
// LDBn     = 0x06,
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
// INCL     = 0x2C,
// DECL     = 0x2D,
// LDLn     = 0x2e,
// LDSPnn   = 0x31,
// LDHLADec = 0x32,
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
// XORA     = 0xAF,
// POPBC    = 0xC1,
// PUSHBC   = 0xC5,
// RET      = 0xC9,
// CB       = 0xCB,
// CALLnn   = 0xCD,
// LDHnA    = 0xE0,
// LDHCA    = 0xE2,
// LDnnA    = 0xEA,
// LDHAn    = 0xF0,
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

type CPUCore struct {
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
}

func (core *CPUCore) SetHL(v uint16) {
	if core.clockCycle.Falling {
		panic("SetHL must be called on rising edge")
	}
	core.Debug("SetHL", "0x%04x", v)
	core.Regs.H = uint8(v >> 8)
	core.Regs.L = uint8(v)
}

func (core *CPUCore) SetBC(v uint16) {
	if core.clockCycle.Falling {
		panic("SetBC must be called on rising edge")
	}
	core.Debug("SetBC", "0x%04x", v)
	core.Regs.B = uint8(v >> 8)
	core.Regs.C = uint8(v)
}

func (core *CPUCore) SetDE(v uint16) {
	if core.clockCycle.Falling {
		panic("SetDE must be called on rising edge")
	}
	core.Debug("SetDE", "0x%04x", v)
	core.Regs.D = uint8(v >> 8)
	core.Regs.E = uint8(v)
}

func (core *CPUCore) SetSP(v uint16) {
	if core.clockCycle.Falling {
		panic("SetSP must be called on rising edge")
	}
	core.Debug("SetSP", "0x%04x", v)
	core.Regs.SP = v
}

func (core *CPUCore) GetA() uint8 {
	v := core.Regs.A
	core.Debug("GetA", "0x%02x", v)
	return v
}

func (core *CPUCore) GetB() uint8 {
	v := core.Regs.B
	core.Debug("GetB", "0x%02x", v)
	return v
}

func (core *CPUCore) GetC() uint8 {
	v := core.Regs.C
	core.Debug("GetC", "0x%02x", v)
	return v
}

func (core *CPUCore) GetD() uint8 {
	v := core.Regs.D
	core.Debug("GetD", "0x%02x", v)
	return v
}

func (core *CPUCore) GetE() uint8 {
	v := core.Regs.E
	core.Debug("GetE", "0x%02x", v)
	return v
}

func (core *CPUCore) GetH() uint8 {
	v := core.Regs.H
	core.Debug("GetE", "0x%02x", v)
	return v
}

func (core *CPUCore) GetL() uint8 {
	v := core.Regs.L
	core.Debug("GetE", "0x%02x", v)
	return v
}

func (core *CPUCore) GetDE() uint16 {
	v := join16(core.Regs.D, core.Regs.E)
	core.Debug("GetDE", "0x%04x", v)
	return v
}

func (core *CPUCore) GetHL() uint16 {
	v := join16(core.Regs.H, core.Regs.L)
	core.Debug("GetHL", "0x%04x", v)
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

func (core *CPUCore) setFlag(mask uint8, v bool) {
	if v {
		core.Regs.F |= mask
	} else {
		core.Regs.F &= ^mask
	}
}

func (core *CPUCore) getFlag(mask uint8) bool {
	return (core.Regs.F & mask) == mask
}

func (core *CPUCore) GetFlagZ() bool {
	v := core.getFlag(0x80)
	core.Debug("GetFlagZ", "F=%x Z=%v", core.Regs.F, v)
	return v
}

func (core *CPUCore) SetFlagZ(v bool) {
	core.setFlag(0x80, v)
	core.Debug("SetFlagZ", "F=%x Z=%v", core.Regs.F, v)
}

func (core *CPUCore) GetFlagN() bool {
	return core.getFlag(0x40)
}

func (core *CPUCore) SetFlagN(v bool) {
	core.Debug("SetFlagN", "%v", v)
	core.setFlag(0x40, v)
}

func (core *CPUCore) GetFlagH() bool {
	return core.getFlag(0x20)
}

func (core *CPUCore) SetFlagH(v bool) {
	core.Debug("SetFlagH", "%v", v)
	core.setFlag(0x20, v)
}

func (core *CPUCore) TODOFlagN() {
	core.Debug("NotImplemented", "Set BCD subtraction flag")
}

func (core *CPUCore) TODOFlagH() {
	core.Debug("NotImplemented", "Set BCD half-carry flag")
}

func (core *CPUCore) GetFlagC() bool {
	return core.getFlag(0x10)
}

func (core *CPUCore) SetFlagC(v bool) {
	core.Debug("SetFlagC", "%v", v)
	core.setFlag(0x10, v)
}

func (core *CPUCore) SetPC(pc uint16) {
	if core.clockCycle.Falling {
		panic("SetPC must be called on rising edge")
	}
	core.Debug("SetPC", "0x%04x", pc)
	core.Regs.PC = pc
}

func (core *CPUCore) IncPC() {
	if core.clockCycle.Falling {
		panic("IncPC must be called on rising edge")
	}
	core.Regs.PC++
	core.Debug("IncPC", "0x%04x", core.Regs.PC)
}

func (core *CPUCore) Debug(event string, f string, v ...any) {
	if core.inCoreDump {
		return
	}
	if slices.Contains(coreDebugEvents, event) || slices.Contains(coreDumpEvents, event) {
		dir := "^"
		if core.clockCycle.Falling {
			dir = "v"
		}
		fmt.Printf("%d %s PC=0x%04x %v mcycle=%v | %s | ", core.clockCycle.C, dir, core.Regs.PC, core.Regs.IR, core.machineCycle, event)
		fmt.Printf(f, v...)
		fmt.Printf("\n")
	}
	if slices.Contains(coreDumpEvents, event) {
		core.Dump()
	}
}

func (core *CPUCore) SetWZ(v uint16) {
	core.W = uint8(v >> 8)
	core.Z = uint8(v)
	wz := (uint16(core.W) << 8) | uint16(core.Z)
	core.Debug("SetWZ", "%v", wz)
}

func (core *CPUCore) GetWZ() uint16 {
	wz := (uint16(core.W) << 8) | uint16(core.Z)
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

type InstructionHandling func(c *CPUCore, e edge) bool

var handlers = map[Opcode]InstructionHandling{
	OpcodeNop: singleCycle(func(core *CPUCore) {
		if core.clockCycle.C > 0 {
			panic("unexpected nop")
		}
	}),
	OpcodeLDAA: singleCycle(func(core *CPUCore) {}),
	OpcodeLDBB: singleCycle(func(core *CPUCore) {}),
	OpcodeLDCC: singleCycle(func(core *CPUCore) {}),
	OpcodeLDDD: singleCycle(func(core *CPUCore) {}),
	OpcodeLDEE: singleCycle(func(core *CPUCore) {}),
	OpcodeLDHH: singleCycle(func(core *CPUCore) {}),
	OpcodeLDLL: singleCycle(func(core *CPUCore) {}),
	OpcodeLDAB: singleCycle(func(core *CPUCore) { core.Regs.A = core.Regs.B }),
	OpcodeLDAC: singleCycle(func(core *CPUCore) { core.Regs.A = core.Regs.C }),
	OpcodeLDAD: singleCycle(func(core *CPUCore) { core.Regs.A = core.Regs.D }),
	OpcodeLDAE: singleCycle(func(core *CPUCore) { core.Regs.A = core.Regs.E }),
	OpcodeLDAH: singleCycle(func(core *CPUCore) { core.Regs.A = core.Regs.H }),
	OpcodeLDAL: singleCycle(func(core *CPUCore) { core.Regs.A = core.Regs.L }),
	OpcodeLDBA: singleCycle(func(core *CPUCore) { core.Regs.B = core.Regs.A }),
	OpcodeLDBC: singleCycle(func(core *CPUCore) { core.Regs.B = core.Regs.C }),
	OpcodeLDBD: singleCycle(func(core *CPUCore) { core.Regs.B = core.Regs.D }),
	OpcodeLDBE: singleCycle(func(core *CPUCore) { core.Regs.B = core.Regs.E }),
	OpcodeLDBH: singleCycle(func(core *CPUCore) { core.Regs.B = core.Regs.H }),
	OpcodeLDBL: singleCycle(func(core *CPUCore) { core.Regs.B = core.Regs.L }),
	OpcodeLDCA: singleCycle(func(core *CPUCore) { core.Regs.C = core.Regs.A }),
	OpcodeLDCB: singleCycle(func(core *CPUCore) { core.Regs.C = core.Regs.B }),
	OpcodeLDCD: singleCycle(func(core *CPUCore) { core.Regs.C = core.Regs.D }),
	OpcodeLDCE: singleCycle(func(core *CPUCore) { core.Regs.C = core.Regs.E }),
	OpcodeLDCH: singleCycle(func(core *CPUCore) { core.Regs.C = core.Regs.H }),
	OpcodeLDCL: singleCycle(func(core *CPUCore) { core.Regs.C = core.Regs.L }),
	OpcodeLDDA: singleCycle(func(core *CPUCore) { core.Regs.D = core.Regs.A }),
	OpcodeLDDB: singleCycle(func(core *CPUCore) { core.Regs.D = core.Regs.B }),
	OpcodeLDDC: singleCycle(func(core *CPUCore) { core.Regs.D = core.Regs.C }),
	OpcodeLDDE: singleCycle(func(core *CPUCore) { core.Regs.D = core.Regs.E }),
	OpcodeLDDH: singleCycle(func(core *CPUCore) { core.Regs.D = core.Regs.H }),
	OpcodeLDDL: singleCycle(func(core *CPUCore) { core.Regs.D = core.Regs.L }),
	OpcodeLDEA: singleCycle(func(core *CPUCore) { core.Regs.E = core.Regs.A }),
	OpcodeLDEB: singleCycle(func(core *CPUCore) { core.Regs.E = core.Regs.B }),
	OpcodeLDEC: singleCycle(func(core *CPUCore) { core.Regs.E = core.Regs.C }),
	OpcodeLDED: singleCycle(func(core *CPUCore) { core.Regs.E = core.Regs.D }),
	OpcodeLDEH: singleCycle(func(core *CPUCore) { core.Regs.E = core.Regs.H }),
	OpcodeLDEL: singleCycle(func(core *CPUCore) { core.Regs.E = core.Regs.L }),
	OpcodeLDHA: singleCycle(func(core *CPUCore) { core.Regs.H = core.Regs.A }),
	OpcodeLDHB: singleCycle(func(core *CPUCore) { core.Regs.H = core.Regs.B }),
	OpcodeLDHC: singleCycle(func(core *CPUCore) { core.Regs.H = core.Regs.C }),
	OpcodeLDHD: singleCycle(func(core *CPUCore) { core.Regs.H = core.Regs.D }),
	OpcodeLDHE: singleCycle(func(core *CPUCore) { core.Regs.H = core.Regs.E }),
	OpcodeLDHL: singleCycle(func(core *CPUCore) { core.Regs.H = core.Regs.L }),
	OpcodeLDLA: singleCycle(func(core *CPUCore) { core.Regs.L = core.Regs.A }),
	OpcodeLDLB: singleCycle(func(core *CPUCore) { core.Regs.L = core.Regs.B }),
	OpcodeLDLC: singleCycle(func(core *CPUCore) { core.Regs.L = core.Regs.C }),
	OpcodeLDLD: singleCycle(func(core *CPUCore) { core.Regs.L = core.Regs.D }),
	OpcodeLDLE: singleCycle(func(core *CPUCore) { core.Regs.L = core.Regs.E }),
	OpcodeLDLH: singleCycle(func(core *CPUCore) { core.Regs.L = core.Regs.H }),
	OpcodeXORA: singleCycle(func(core *CPUCore) {
		core.SetFlagZ(true)
		core.SetFlagN(false)
		core.SetFlagH(false)
		core.SetFlagC(false)
		core.Regs.A = 0
	}),
	OpcodeRLA: singleCycle(func(core *CPUCore) {
		a := core.Regs.A
		bit7 := a & 0x80
		a <<= 1
		if core.GetFlagC() {
			a |= 0x01
		}
		core.SetFlagZ(false)
		core.SetFlagN(false)
		core.SetFlagH(false)
		core.SetFlagC(bit7 != 0)
		core.Regs.A = a
	}),
	OpcodeDECA: decreg(func(core *CPUCore) *uint8 { return &core.Regs.A }),
	OpcodeDECB: decreg(func(core *CPUCore) *uint8 { return &core.Regs.B }),
	OpcodeDECC: decreg(func(core *CPUCore) *uint8 { return &core.Regs.C }),
	OpcodeDECD: decreg(func(core *CPUCore) *uint8 { return &core.Regs.D }),
	OpcodeDECE: decreg(func(core *CPUCore) *uint8 { return &core.Regs.E }),
	OpcodeDECH: decreg(func(core *CPUCore) *uint8 { return &core.Regs.H }),
	OpcodeDECL: decreg(func(core *CPUCore) *uint8 { return &core.Regs.L }),
	OpcodeINCA: increg(func(core *CPUCore) *uint8 { return &core.Regs.A }),
	OpcodeINCB: increg(func(core *CPUCore) *uint8 { return &core.Regs.B }),
	OpcodeINCC: increg(func(core *CPUCore) *uint8 { return &core.Regs.C }),
	OpcodeINCD: increg(func(core *CPUCore) *uint8 { return &core.Regs.D }),
	OpcodeINCE: increg(func(core *CPUCore) *uint8 { return &core.Regs.E }),
	OpcodeINCH: increg(func(core *CPUCore) *uint8 { return &core.Regs.H }),
	OpcodeINCL: increg(func(core *CPUCore) *uint8 { return &core.Regs.L }),
	OpcodeDI: singleCycle(func(core *CPUCore) {
		core.Interrupts.setIMENextCycle = false
		core.Interrupts.IME = false
	}),
	OpcodeEI: singleCycle(func(core *CPUCore) {
		core.Interrupts.setIMENextCycle = true
	}),
	OpcodeJRe: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			core.Z = core.readDataBus()
		// TODO: this impl is not exactly correct
		case edge{2, false}:
		case edge{2, true}:
		case edge{3, false}:
			core.SetPC(uint16(int16(core.Regs.PC) + int16(int8(core.Z))))
			return true
		case edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeJRZe:  jrcce(func(core *CPUCore) bool { return core.GetFlagZ() }),
	OpcodeJRNZe: jrcce(func(core *CPUCore) bool { return !core.GetFlagZ() }),
	OpcodeINCDE: iduOp(func(core *CPUCore) {
		de := core.GetDE()
		de++
		core.SetDE(de)
	}),
	OpcodeINCHL: iduOp(func(core *CPUCore) {
		hl := core.GetHL()
		hl++
		core.SetHL(hl)
	}),
	OpcodeCALLnn: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{2, true}:
			core.W = core.readDataBus()
		case edge{3, false}:
			core.SetSP(core.Regs.SP - 1)
		case edge{3, true}:
		case edge{4, false}:
			core.writeAddressBus(core.Regs.SP)
			core.SetSP(core.Regs.SP - 1)
		case edge{4, true}:
			core.writeDataBus(msb(core.Regs.PC))
		case edge{5, false}:
			core.writeAddressBus(core.Regs.SP)
		case edge{5, true}:
			core.writeDataBus(lsb(core.Regs.PC))
		case edge{6, false}:
			core.SetPC(core.GetWZ())
			return true
		case edge{6, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeRET: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.SP)
			core.SetSP(core.Regs.SP + 1)
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.writeAddressBus(core.Regs.SP)
			core.SetSP(core.Regs.SP + 1)
		case edge{2, true}:
			core.W = core.readDataBus()
		case edge{3, false}:
			core.SetPC(core.GetWZ())
		case edge{3, true}:
		case edge{4, false}, edge{4, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodePUSHBC: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.SetSP(core.Regs.SP - 1)
			core.writeAddressBus(core.Regs.SP)
		case edge{1, true}:
			core.writeDataBus(core.Regs.B)
		case edge{2, false}:
			core.SetSP(core.Regs.SP - 1)
			core.writeAddressBus(core.Regs.SP)
		case edge{2, true}:
			core.writeDataBus(core.Regs.C)
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
	OpcodePOPBC: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.SP)
			core.SetSP(core.Regs.SP + 1)
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.writeAddressBus(core.Regs.SP)
			core.SetSP(core.Regs.SP + 1)
		case edge{2, true}:
			core.W = core.readDataBus()
		case edge{3, false}:
			core.SetBC(core.GetWZ())
			return true
		case edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDBCnn:   ldxxnn(func(core *CPUCore, wz uint16) { core.SetBC(wz) }),
	OpcodeLDDEnn:   ldxxnn(func(core *CPUCore, wz uint16) { core.SetDE(wz) }),
	OpcodeLDHLnn:   ldxxnn(func(core *CPUCore, wz uint16) { core.SetHL(wz) }),
	OpcodeLDSPnn:   ldxxnn(func(core *CPUCore, wz uint16) { core.SetSP(wz) }),
	OpcodeLDHLAInc: ldhla(func(core *CPUCore) { core.SetHL(core.GetHL() + 1) }),
	OpcodeLDHLADec: ldhla(func(core *CPUCore) { core.SetHL(core.GetHL() - 1) }),
	OpcodeLDHLA:    ldhla(func(core *CPUCore) {}),
	OpcodeLDHCA: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(join16(0xff, core.Regs.C))
		case edge{1, true}:
			core.writeDataBus(core.Regs.A)
		case edge{2, false}, edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDnnA: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{2, true}:
			core.W = core.readDataBus()
		case edge{3, false}:
			core.writeAddressBus(core.GetWZ())
		case edge{3, true}:
			core.writeDataBus(core.Regs.A)
		case edge{4, false}, edge{4, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeCPn: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			return true
		case edge{2, true}:
			core.Debug("CPn", "A=%02x n=%02x", core.Regs.A, core.Z)
			carry := core.Regs.A < core.Z
			result := core.Regs.A - core.Z
			core.SetFlagZ(result == 0)
			core.SetFlagN(true)
			core.TODOFlagH()
			core.SetFlagC(carry)
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDHnA: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.writeAddressBus(join16(0xff, core.Z))
			core.IncPC()
		case edge{2, true}:
			core.writeDataBus(core.Regs.A)
		case edge{3, false}, edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDHAn: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.writeAddressBus(join16(0xff, core.Z))
		case edge{2, true}:
			core.Z = core.readDataBus()
		case edge{3, false}:
			return true
		case edge{3, true}:
			core.Regs.A = core.Z
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDADE: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(join16(core.Regs.D, core.Regs.E))
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			return true
		case edge{2, true}:
			core.Regs.A = core.Z
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDAn: ldrn(func(core *CPUCore, z uint8) { core.Regs.A = z }),
	OpcodeLDBn: ldrn(func(core *CPUCore, z uint8) { core.Regs.B = z }),
	OpcodeLDCn: ldrn(func(core *CPUCore, z uint8) { core.Regs.C = z }),
	OpcodeLDDn: ldrn(func(core *CPUCore, z uint8) { core.Regs.D = z }),
	OpcodeLDEn: ldrn(func(core *CPUCore, z uint8) { core.Regs.E = z }),
	OpcodeLDHn: ldrn(func(core *CPUCore, z uint8) { core.Regs.H = z }),
	OpcodeLDLn: ldrn(func(core *CPUCore, z uint8) { core.Regs.L = z }),
	OpcodeCB: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			opcode := core.readDataBus()
			core.CBOp = CBOp{Op: cb((opcode & 0xf8) >> 3), Target: CBTarget(opcode & 0x7)}
		case edge{2, false}:
			return true
		case edge{2, true}:
			var val uint8
			switch core.CBOp.Target {
			case CBTargetB:
				val = core.Regs.B
			case CBTargetC:
				val = core.Regs.C
			case CBTargetD:
				val = core.Regs.D
			case CBTargetE:
				val = core.Regs.E
			case CBTargetH:
				val = core.Regs.H
			case CBTargetL:
				val = core.Regs.L
			case CBTargetIndirectHL:
				panic("indirect thru HL not implemented")
			case CBTargetA:
				val = core.Regs.A
			default:
				panic("unknown CBOp target")
			}
			core.Debug("ExecCBOp", "op=%v target=%v", core.CBOp.Op, core.CBOp.Target)
			switch core.CBOp.Op {
			case CbRL:
				bit7 := val & 0x80
				val <<= 1
				if core.GetFlagC() {
					val |= 0x01
				}
				core.SetFlagZ(val == 0)
				core.SetFlagN(false)
				core.SetFlagH(false)
				core.SetFlagC(bit7 != 0)
			case CbBit0:
				cbbit(core, val, 0x01)
			case CbBit1:
				cbbit(core, val, 0x02)
			case CbBit2:
				cbbit(core, val, 0x04)
			case CbBit3:
				cbbit(core, val, 0x08)
			case CbBit4:
				cbbit(core, val, 0x10)
			case CbBit5:
				cbbit(core, val, 0x20)
			case CbBit6:
				cbbit(core, val, 0x40)
			case CbBit7:
				cbbit(core, val, 0x80)
			default:
				panicf("unknown op = %+v", core.CBOp)
			}
			switch core.CBOp.Target {
			case CBTargetB:
				core.Regs.B = val
			case CBTargetC:
				core.Regs.C = val
			case CBTargetD:
				core.Regs.D = val
			case CBTargetE:
				core.Regs.E = val
			case CBTargetH:
				core.Regs.H = val
			case CBTargetL:
				core.Regs.L = val
			case CBTargetIndirectHL:
				panic("indirect thru HL not implemented")
			case CBTargetA:
				core.Regs.A = val
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

func singleCycle(f func(core *CPUCore)) func(core *CPUCore, e edge) bool {
	return func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			return true
		case edge{1, true}:
			f(core)
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func jrcce(f func(core *CPUCore) bool) func(core *CPUCore, e edge) bool {
	return func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			if f(core) {
			} else {
				return true
			}
		case edge{2, true}:
			if f(core) {
				newPC := uint16(int16(core.Regs.PC) + int16(int8(core.Z)))
				core.SetWZ(newPC)
			} else {
				return true
			}
		case edge{3, false}:
			if f(core) {
				core.SetPC(core.GetWZ())
				return true
			} else {
				panicv(e)
			}
		case edge{3, true}:
			if f(core) {
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

func decreg(f func(core *CPUCore) *uint8) func(core *CPUCore, e edge) bool {
	return singleCycle(func(core *CPUCore) {
		reg := f(core)
		*reg -= 1
		core.SetFlagZ(*reg == 0)
		core.SetFlagN(true)
		core.TODOFlagH()
	})
}

func increg(f func(core *CPUCore) *uint8) func(core *CPUCore, e edge) bool {
	return singleCycle(func(core *CPUCore) {
		reg := f(core)
		*reg += 1
		core.SetFlagZ(*reg == 0)
		core.SetFlagN(false)
		core.TODOFlagH()
	})
}

func iduOp(f func(core *CPUCore)) func(core *CPUCore, e edge) bool {
	return func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			f(core)
		case edge{1, true}:
		case edge{2, false}, edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func ldrn(f func(core *CPUCore, z uint8)) func(core *CPUCore, e edge) bool {
	return func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.IncPC()
			return true
		case edge{2, true}:
			f(core, core.Z)
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func ldhla(f func(core *CPUCore)) func(core *CPUCore, e edge) bool {
	return func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.GetHL())
			f(core)
		case edge{1, true}:
			core.writeDataBus(core.Regs.A)
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

func ldxxnn(f func(core *CPUCore, wz uint16)) func(core *CPUCore, e edge) bool {
	return func(core *CPUCore, e edge) bool {
		switch e {
		case edge{1, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{1, true}:
			core.Z = core.readDataBus()
		case edge{2, false}:
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		case edge{2, true}:
			core.W = core.readDataBus()
		case edge{3, false}:
			f(core, join16(core.W, core.Z))
			return true
		case edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func cbbit(core *CPUCore, val, mask uint8) {
	core.SetFlagZ(val&mask == 0)
	core.SetFlagN(false)
	core.SetFlagH(true)
}

func NewCPUCore(
	phi *Clock,
) *CPUCore {
	core := &CPUCore{
		m:   &sync.Mutex{},
		PHI: phi,
	}
	phi.AddRiseCallback(core.fsm)
	phi.AddFallCallback(core.fsm)
	return core
}

func (core *CPUCore) AttachPeripheral(p Peripheral) {
	core.m.Lock()
	core.Peripherals = append(core.Peripherals, p)
	core.m.Unlock()
}

func (core *CPUCore) Dump() {
	core.inCoreDump = true
	defer func() { core.inCoreDump = false }()

	fmt.Printf("\n--------\nCore dump:\n")
	fmt.Printf("PC = 0x%04x\n", core.Regs.PC)
	fmt.Printf("SP = 0x%04x\n", core.Regs.SP)
	fmt.Printf("A  =   0x%02x\n", core.Regs.A)
	fmt.Printf("F  =   0x%02x (Z=%v, H=%v, N=%v C=%v)\n", core.Regs.F, core.GetFlagZ(), core.GetFlagH(), core.GetFlagN(), core.GetFlagC())
	fmt.Printf("B  =   0x%02x\n", core.Regs.B)
	fmt.Printf("C  =   0x%02x\n", core.Regs.C)
	fmt.Printf("D  =   0x%02x\n", core.Regs.D)
	fmt.Printf("E  =   0x%02x\n", core.Regs.E)
	fmt.Printf("H  =   0x%02x\n", core.Regs.H)
	fmt.Printf("L  =   0x%02x\n", core.Regs.L)
	fmt.Printf("W  =   0x%02x\n", core.W)
	fmt.Printf("Z  =   0x%02x\n", core.Z)
	fmt.Printf("IR =   0x%02x\n", uint8(core.Regs.IR))
	fmt.Printf("--------\n")
	fmt.Printf("Code (PC highlighted)\n")
	start := uint16(0)
	if core.Regs.PC > 0x40 {
		start = core.Regs.PC - 0x40
	}
	end := uint16(0xffff)
	if core.Regs.PC < 0xffff-0x40 {
		end = core.Regs.PC + 0x40
	}
	core.memdump(start, end, core.Regs.PC-1)
	fmt.Printf("--------\n")
	fmt.Printf("HRAM (SP highlighted):\n")
	core.memdump(0xff80, 0xfffe, core.Regs.SP)
	fmt.Printf("--------\n")
}

func (core *CPUCore) memdump(start, end, highlight uint16) {
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
		core.writeAddressBus(addr)
		if highlight == addr {
			fmt.Printf("[%02x]", core.DataBus)
		} else {
			fmt.Printf(" %02x ", core.DataBus)
		}
	}

	alignedEnd := (end/0x10)*0x10 + 0x10 - 1
	for addr := end; addr < alignedEnd; addr++ {
		fmt.Printf(" .. ")
	}
	fmt.Printf("\n")
}

func (core *CPUCore) fsm(c Cycle) {
	defer func() {
		if e := recover(); e != nil {
			core.Debug("Panic", "%v", e)
			core.Dump()
			panic(e)
		}
	}()
	core.wroteToAddressBusThisCycle = false

	core.clockCycle = c

	core.applyPendingIME()

	var fetch bool
	if c.C > 0 {
		opcode := core.Regs.IR
		if handler, ok := handlers[opcode]; ok {
			e := edge{core.machineCycle, c.Falling}
			core.Debug("Handler", "e=%v", e)
			fetch = handler(core, e)
		} else {
			panicf("not implemented opcode %v", opcode)
		}
	} else {
		// initial instruction
		fetch = true
	}

	if fetch {
		if !c.Falling {
			core.Debug("PreFetch", "PC=%04x", core.Regs.PC)
			core.writeAddressBus(core.Regs.PC)
			core.IncPC()
		} else {
			core.Debug("ExecDone", "")
			core.W = 0
			core.Z = 0
			core.machineCycle = 1
			core.Regs.IR = Opcode(core.DataBus)
			core.Debug("ExecBegin", "%s", core.Regs.IR)
			core.Debug(fmt.Sprintf("ExecBegin%s", core.Regs.IR), "")
		}
	} else if c.Falling {
		core.machineCycle++
	}
}

func (core *CPUCore) writeAddressBus(addr uint16) {
	if !core.inCoreDump {
		if core.clockCycle.Falling {
			panic("writeAddressBus must be called on rising edge")
		}
		if core.wroteToAddressBusThisCycle {
			panic("more than one call to writeAddressBus this cycle")
		}
	}
	core.wroteToAddressBusThisCycle = true
	core.Debug("WriteAddressBus", "0x%04x", addr)
	core.AddressBus = addr
	for _, p := range core.Peripherals {
		start, size := p.Range()
		if addr >= start && addr <= start+(size-1) {
			v := p.Read(addr)
			core.Debug("PeriphRead", "0x%02x from %s @ 0x%04x", v, p.Name(), addr)
			core.DataBus = v
			return
		}
	}
	for _, p := range core.Peripherals {
		start, size := p.Range()
		core.Debug("Panic", "start=%0x size=%0x last=%0x", start, size, start+(size-1))
	}
	panicf("no peripheral mapped to 0x%x", addr)
}

func (core *CPUCore) writeDataBus(v uint8) {
	if !core.clockCycle.Falling {
		panic("writeDataBus must be called on falling edge")
	}
	core.DataBus = v
	addr := core.AddressBus
	for _, p := range core.Peripherals {
		start, size := p.Range()
		if addr >= start && addr <= start+(size-1) {
			p.Write(addr, v)
			core.Debug("PeriphWrite", "0x%02x to %s @ 0x%04x", v, p.Name(), addr)
			core.Debug(fmt.Sprintf("Watch%04x", addr), "wrote %02x", v)
			core.DataBus = v
			return
		}
	}
}

func (core *CPUCore) readDataBus() uint8 {
	if !core.clockCycle.Falling {
		panic("readDataBus must be called on falling edge")
	}
	return core.DataBus
}

func (core *CPUCore) applyPendingIME() {
	if core.Interrupts.setIMENextCycle {
		core.Interrupts.setIMENextCycle = false
		core.Interrupts.IME = true
	}
}
