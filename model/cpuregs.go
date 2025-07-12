package model

import (
	"fmt"
	"sync"
)

const (
	ocNop      = 0x00
	ocLDHLnn   = 0x21
	ocLDSPnn   = 0x31
	ocLDHLADec = 0x32
	ocXORA     = 0xAF
	ocCBPrefix = 0xCB
	ocDI       = 0xF3
	ocEI       = 0xFB
)

const (
	cbRLC = 0b00000
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
	IR uint8
}

func (rf *RegisterFile) HL(v ...uint16) uint16 {
	if len(v) == 1 {
		rf.H = uint8(v[0] >> 8)
		rf.L = uint8(v[0])
	}
	return (uint16(rf.H) << 8) | uint16(rf.L)
}

func (rf *RegisterFile) flag(mask uint8, v ...bool) bool {
	if len(v) == 1 {
		if z := v[0]; z {
			rf.F |= mask
		} else {
			rf.F ^= mask
		}
	}
	return rf.F&mask == mask
}

func (rf *RegisterFile) FlagZ(v ...bool) bool {
	return rf.flag(0x80, v...)
}

func (rf *RegisterFile) FlagN(v ...bool) bool {
	return rf.flag(0x40, v...)
}

func (rf *RegisterFile) FlagH(v ...bool) bool {
	return rf.flag(0x20, v...)
}

func (rf *RegisterFile) FlagC(v ...bool) bool {
	return rf.flag(0x10, v...)
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
	DataBus     *Bus[uint8]
	AddressBus  *Bus[uint16]
	Peripherals []Peripheral

	RegisterFile RegisterFile
	Interrupts   Interrupts

	Z uint8
	W uint8

	machineCycle int
}

type Peripheral interface {
	Range() (start, size uint16)
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
}

type edge struct {
	Cycle   int
	Falling bool
}

type InstructionHandling func(c *CPUCore, e edge) bool

var handlers = map[uint8]InstructionHandling{
	ocNop: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{1, false}:
			return true
		default:
			panic(e)
		}
		return false
	},
	ocLDHLnn: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			core.Z = core.DataBus.Read()
			core.RegisterFile.PC++
		case edge{2, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{2, true}:
			core.W = core.DataBus.Read()
			core.RegisterFile.H = core.W
			core.RegisterFile.L = core.Z
			fmt.Printf("wrote HL = %x:%x\n", core.W, core.Z)
			core.RegisterFile.PC++
			return true
		default:
			panic(e)
		}
		return false
	},
	ocLDSPnn: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			core.Z = core.DataBus.Read()
			core.RegisterFile.PC++
		case edge{2, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{2, true}:
			core.W = core.DataBus.Read()
			core.RegisterFile.SP = uint16(core.W)<<8 | uint16(core.Z)
			fmt.Printf("wrote SP = %x:%x\n", core.W, core.Z)
			core.RegisterFile.PC++
			return true
		default:
			panic(e)
		}
		return false
	},
	ocLDHLADec: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{0, true}:
		case edge{1, false}:
			hl := core.RegisterFile.HL()
			core.writeAddressBus(hl)
			core.RegisterFile.HL(hl - 1)
		case edge{1, true}:
			core.RegisterFile.A = core.DataBus.Read()
			return true
		default:
			panic(e)
		}
		return false
	},
	ocCBPrefix: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{0, true}:
		case edge{1, false}:
			core.writeAddressBus(core.RegisterFile.PC)
		case edge{1, true}:
			op := core.DataBus.Read()
			core.startCBOp(op)
			core.RegisterFile.PC++
			return true
		default:
			panic(e)
		}
		return false
	},
	ocXORA: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{0, true}:
		case edge{1, false}:
		case edge{1, true}:
			core.RegisterFile.A = 0
			core.RegisterFile.FlagZ(true)
			core.RegisterFile.FlagN(false)
			core.RegisterFile.FlagH(false)
			core.RegisterFile.FlagC(false)
			return true
		default:
			panic(e)
		}
		return false
	},
	ocDI: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{0, true}:
			core.Interrupts.setIMENextCycle = false
			core.Interrupts.IME = false
			return true
		default:
			panic(e)
		}
		return false
	},
	ocEI: func(core *CPUCore, e edge) bool {
		switch e {
		case edge{0, false}:
		case edge{0, true}:
			core.Interrupts.setIMENextCycle = true
			return true
		default:
			panic(e)
		}
		return false
	},
}

func NewCPUCore(
	phi *Clock,
	dataBus *Bus[uint8],
	addressBus *Bus[uint16],
) *CPUCore {
	core := &CPUCore{
		m:          &sync.Mutex{},
		PHI:        phi,
		DataBus:    dataBus,
		AddressBus: addressBus,
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
	core.applyPendingIME()

	if core.machineCycle == 0 {
		if c.Rising {
			core.writeAddressBus(core.RegisterFile.PC)
			core.RegisterFile.PC++
		} else {
			core.RegisterFile.IR = core.DataBus.Read()
		}
	}

	opcode := core.RegisterFile.IR
	if handler, ok := handlers[opcode]; ok {
		done := handler(core, edge{core.machineCycle, !c.Rising})
		if !c.Rising {
			if done {
				core.machineCycle = 0
			} else {
				core.machineCycle++
			}
		}
	} else {
		panic(fmt.Sprintf("not implemented opcode 0x%x", opcode))
	}
}

func (core *CPUCore) writeAddressBus(addr uint16) {
	core.AddressBus.Write(addr)
	for _, p := range core.Peripherals {
		start, size := p.Range()
		fmt.Printf("start=%v size=%v\n", start, size)
		if addr >= start && addr <= start+(size-1) {
			v := p.Read(addr)
			core.DataBus.Write(v)
			return
		}
	}
	panic(fmt.Sprintf("no peripheral mapped to 0x%x", addr))
}

func (core *CPUCore) applyPendingIME() {
	if core.Interrupts.setIMENextCycle {
		core.Interrupts.setIMENextCycle = false
		core.Interrupts.IME = true
	}
}

func (core *CPUCore) startCBOp(opcode uint8) {
	op := (opcode & 0xf8) >> 3
	target := opcode & 0x7

	switch op {
	default:
		panic(fmt.Sprintf("unknown op = %v with target %v", op, target))
	}
}
