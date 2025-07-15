package model

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

type Peripheral interface {
	Name() string
	Range() (start, size uint16)
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
}

func NewCPU(
	phi *Clock,
) *CPU {
	cpu := &CPU{
		m:   &sync.Mutex{},
		PHI: phi,
	}
	phi.AttachDevice(cpu.fsm)
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
