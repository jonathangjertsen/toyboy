package model

//go:generate go-enum --marshal --flag --values --nocomments

import (
	"fmt"
	"slices"
	"sync"
)

var coreDebugEvents = []string{
	"Panic",
	"SetPC",
	//"IncPC",
	//"SetFlagZ",
	//"ExecDone",
	//"WriteAddressBus",
	"PeriphRead",
	//"GetFlagZ",
	"PeriphWrite",
	"ExecCBOp",
	"SetHL",
	"SetA",
	"SetC",
}

// ENUM(
// Nop      = 0x00,
// LDCn     = 0x0e,
// JRNZe    = 0x20,
// LDHLnn   = 0x21,
// LDSPnn   = 0x31,
// LDHLADec = 0x32,
// LDAn     = 0x3e,
// XORA     = 0xAF,
// CB       = 0xCB,
// DI       = 0xF3,
// EI       = 0xFB,
// )
type Opcode uint8

const (
	cbRLC uint8 = iota
	cbRRC
	cbRL
	cbRR
	cbSLA
	cbSRA
	cbSWAP
	cbSRL
	cbBit0
	cbBit1
	cbBit2
	cbBit3
	cbBit4
	cbBit5
	cbBit6
	cbBit7
	cbBit8
	cbRes0
	cbRes1
	cbRes2
	cbRes3
	cbRes4
	cbRes5
	cbRes6
	cbRes7
	cbRes8
	cbSet0
	cbSet1
	cbSet2
	cbSet3
	cbSet4
	cbSet5
	cbSet6
	cbSet7
	cbSet8
)

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

	RegisterFile RegisterFile
	Interrupts   Interrupts

	Z    uint8
	W    uint8
	CBOp CBOp

	machineCycle int

	clockCycle Cycle
}

func (core *CPUCore) SetHL(v uint16) {
	core.Debug("SetHL", "0x%04x\n", v)
	core.RegisterFile.H = uint8(v >> 8)
	core.RegisterFile.L = uint8(v)
}

func (core *CPUCore) GetHL() uint16 {
	return join16(core.RegisterFile.H, core.RegisterFile.L)
}

func join16(msb, lsb uint8) uint16 {
	return (uint16(msb) << 8) | uint16(lsb)
}

func (core *CPUCore) setFlag(mask uint8, v bool) {
	if v {
		core.RegisterFile.F |= mask
	} else {
		core.RegisterFile.F &= ^mask
	}
}

func (core *CPUCore) getFlag(mask uint8) bool {
	return (core.RegisterFile.F & mask) == mask
}

func (core *CPUCore) GetFlagZ() bool {
	v := core.getFlag(0x80)
	core.Debug("GetFlagZ", "F=%x Z=%v\n", core.RegisterFile.F, v)
	return v
}

func (core *CPUCore) SetFlagZ(v bool) {
	core.setFlag(0x80, v)
	core.Debug("SetFlagZ", "F=%x Z=%v\n", core.RegisterFile.F, v)
}

func (core *CPUCore) GetFlagN() bool {
	return core.getFlag(0x40)
}

func (core *CPUCore) SetFlagN(v bool) {
	core.Debug("SetFlagN", "%v\n", v)
	core.setFlag(0x40, v)
}

func (core *CPUCore) GetFlagH() bool {
	return core.getFlag(0x20)
}

func (core *CPUCore) SetFlagH(v bool) {
	core.Debug("SetFlagH", "%v\n", v)
	core.setFlag(0x20, v)
}

func (core *CPUCore) GetFlagC() bool {
	return core.getFlag(0x10)
}

func (core *CPUCore) SetFlagC(v bool) {
	core.Debug("SetFlagC", "%v\n", v)
	core.setFlag(0x10, v)
}

func (core *CPUCore) SetA(v uint8) {
	core.Debug("SetA", "%v\n", v)
	core.RegisterFile.A = v
}

func (core *CPUCore) SetB(v uint8) {
	core.Debug("SetB", "%v\n", v)
	core.RegisterFile.B = v
}

func (core *CPUCore) SetC(v uint8) {
	core.Debug("SetC", "%v\n", v)
	core.RegisterFile.C = v
}

func (core *CPUCore) SetPC(pc uint16) {
	core.Debug("SetPC", "%v\n", pc)
	core.RegisterFile.PC = pc
}

func (core *CPUCore) IncPC() {
	core.RegisterFile.PC++
	core.Debug("IncPC", "%v\n", core.RegisterFile.PC)
}

func (core *CPUCore) Debug(event string, f string, v ...any) {
	if !slices.Contains(coreDebugEvents, event) {
		return
	}
	dir := "v"
	if core.clockCycle.Rising {
		dir = "^"
	}
	fmt.Printf("%d %s PC=%d %v mcycle=%v | %s | ", core.clockCycle.C, dir, core.RegisterFile.PC, core.RegisterFile.IR, core.machineCycle, event)
	fmt.Printf(f, v...)
}

func (core *CPUCore) SetWZ(v uint16) {
	core.W = uint8(v >> 8)
	core.Z = uint8(v)
	wz := (uint16(core.W) << 8) | uint16(core.Z)
	core.Debug("SetWZ", "%v\n", wz)
}

func (core *CPUCore) GetWZ() uint16 {
	wz := (uint16(core.W) << 8) | uint16(core.Z)
	return wz
}

// ENUM(B, C, D, E, H, L, IndirectHL, A)
type CBTarget uint8

type CBOp struct {
	Op     uint8
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
	OpcodeNop: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeJRNZe: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			core.Z = core.DataBus
			core.IncPC()
		case edge{2, false}:
		case edge{2, true}:
			if !core.GetFlagZ() {
				newPC := uint16(int16(core.RegisterFile.PC) + int16(int8(core.Z)))
				core.SetWZ(newPC)
			} else {
				return true
			}
		case edge{3, false}:
			if !core.GetFlagZ() {
			} else {
				panicv(e)
			}
		case edge{3, true}:
			if !core.GetFlagZ() {
				core.SetPC(core.GetWZ())
				return true
			} else {
				panicv(e)
			}
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDHLnn: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			core.Z = core.DataBus
			core.IncPC()
		case edge{2, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{2, true}:
			core.W = core.DataBus
			core.SetHL(join16(core.W, core.Z))
			core.IncPC()
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDSPnn: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			core.Z = core.DataBus
			core.IncPC()
		case edge{2, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{2, true}:
			core.W = core.DataBus
			core.RegisterFile.SP = uint16(core.W)<<8 | uint16(core.Z)
			core.IncPC()
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDHLADec: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
		case edge{1, false}:
			hl := core.GetHL()
			core.writeAddressBus(hl)
			core.SetHL(hl - 1)
		case edge{1, true}:
			core.writeDataBus(core.RegisterFile.A)
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeLDCn: ldrn(func(core *CPUCore, z uint8) { core.SetC(z) }),
	OpcodeLDAn: ldrn(func(core *CPUCore, z uint8) { core.SetA(z) }),
	OpcodeCB: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			opcode := core.DataBus
			core.CBOp = CBOp{Op: (opcode & 0xf8) >> 3, Target: CBTarget(opcode & 0x7)}
			core.IncPC()
		case edge{2, false}:
		case edge{2, true}:
			var val uint8
			switch core.CBOp.Target {
			case CBTargetB:
				val = core.RegisterFile.B
			case CBTargetC:
				val = core.RegisterFile.C
			case CBTargetD:
				val = core.RegisterFile.D
			case CBTargetE:
				val = core.RegisterFile.E
			case CBTargetH:
				val = core.RegisterFile.H
			case CBTargetL:
				val = core.RegisterFile.L
			case CBTargetIndirectHL:
				panic("indirect thru HL not implemented")
			case CBTargetA:
				val = core.RegisterFile.A
			default:
				panic("unknown CBOp target")
			}
			core.Debug("ExecCBOp", "op=%v target=%v\n", core.CBOp.Op, core.CBOp.Target)
			switch core.CBOp.Op {
			case cbBit7:
				core.SetFlagZ(val&0x80 == 0)
				core.SetFlagN(false)
				core.SetFlagH(true)
			default:
				panicf("unknown op = %+v", core.CBOp)
			}
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeXORA: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
		case edge{1, false}:
		case edge{1, true}:
			core.RegisterFile.A = 0
			core.SetFlagZ(true)
			core.SetFlagN(false)
			core.SetFlagH(false)
			core.SetFlagC(false)
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeDI: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
			core.Interrupts.setIMENextCycle = false
			core.Interrupts.IME = false
			return true
		default:
			panicv(e)
		}
		return false
	},
	OpcodeEI: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
			core.Interrupts.setIMENextCycle = true
			return true
		default:
			panicv(e)
		}
		return false
	},
}

func ldrn(f func(core *CPUCore, z uint8)) func(core *CPUCore, e edge) bool {
	return func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
			core.IncPC()
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			core.Z = core.DataBus
		case edge{2, false}:
		case edge{2, true}:
			core.IncPC()
			f(core, core.Z)
			return true
		default:
			panicv(e)
		}
		return false
	}
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

func (core *CPUCore) fsm(c Cycle) {
	defer func() {
		if e := recover(); e != nil {
			core.Debug("Panic", "%v\n", e)
			panic(e)
		}
	}()

	core.clockCycle = c

	core.applyPendingIME()

	if core.machineCycle == 0 {
		if c.Rising {
			core.writeAddressBus(core.RegisterFile.PC)
		} else {
			core.RegisterFile.IR = Opcode(core.DataBus)
		}
	}

	opcode := core.RegisterFile.IR
	if handler, ok := handlers[opcode]; ok {
		done := handler(core, edge{core.machineCycle, !c.Rising})
		if c.Rising {
			if done {
				panic("can't be done on rising edge")
			}
		} else {
			if done {
				core.Debug("ExecDone", "\n")
				core.machineCycle = 0
			} else {
				core.machineCycle++
			}
		}
	} else {
		panicf("not implemented opcode 0x%x", int(opcode))
	}
}

func (core *CPUCore) writeAddressBus(addr uint16) {
	core.Debug("WriteAddressBus", "0x%04x\n", addr)
	core.AddressBus = addr
	for _, p := range core.Peripherals {
		start, size := p.Range()
		if addr >= start && addr <= start+(size-1) {
			v := p.Read(addr)
			core.Debug("PeriphRead", "0x%02x from %s @ 0x%04x\n", v, p.Name(), addr)
			core.DataBus = v
			return
		}
	}
	for _, p := range core.Peripherals {
		start, size := p.Range()
		core.Debug("Panic", "start=%0x size=%0x last=%0x\n", start, size, start+(size-1))
	}
	panicf("no peripheral mapped to 0x%x", addr)
}

func (core *CPUCore) writeDataBus(v uint8) {
	core.DataBus = v
	addr := core.AddressBus
	for _, p := range core.Peripherals {
		start, size := p.Range()
		if addr >= start && addr <= start+(size-1) {
			p.Write(addr, v)
			core.Debug("PeriphWrite", "0x%02x to %s @ 0x%04x\n", v, p.Name(), addr)
			core.DataBus = v
			return
		}
	}
}

func (core *CPUCore) applyPendingIME() {
	if core.Interrupts.setIMENextCycle {
		core.Interrupts.setIMENextCycle = false
		core.Interrupts.IME = true
	}
}
