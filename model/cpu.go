package model

import (
	"fmt"
	"os"
	"slices"
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

type Interrupts struct {
	IF              uint8
	IE              uint8
	IME             bool
	setIMENextCycle bool
}

type CPU struct {
	PHI *Clock

	Bus *Bus

	Regs       RegisterFile
	Interrupts Interrupts

	CBOp CBOp

	machineCycle int

	clockCycle                 Cycle
	inCoreDump                 bool
	wroteToAddressBusThisCycle bool

	handlers [256]InstructionHandling

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

func NewCPU(
	phi *Clock,
	bus *Bus,
) *CPU {
	cpu := &CPU{
		PHI: phi,
		Bus: bus,
	}
	cpu.handlers = handlers(cpu)
	phi.AttachDevice(cpu.fsm)
	return cpu
}

func (cpu *CPU) Dump() {

	cd := cpu.GetCoreDump()
	f := os.Stdout
	fmt.Fprintf(f, "\n--------\nCore dump:\n")
	cd.PrintRegs(f)
	fmt.Fprintf(f, "--------\n")
	fmt.Fprintf(f, "Code (PC highlighted)\n")
	cd.PrintProgram(f)
	fmt.Fprintf(f, "--------\n")
	fmt.Fprintf(f, "HRAM (SP highlighted):\n")
	cd.PrintHRAM(f)
	fmt.Fprintf(f, "--------\n")
	fmt.Fprintf(f, "OAM:\n")
	cd.PrintOAM(f)
	fmt.Fprintf(f, "--------\n")

	fmt.Printf("Last executed instructions:\n")
	for i := (cpu.rewindBufferIdx + 1) % len(cpu.rewindBuffer); i != cpu.rewindBufferIdx; i = (i + 1) % len(cpu.rewindBuffer) {
		fmt.Printf("[PC=%04x] %s\n", cpu.rewindBuffer[i].PC, cpu.rewindBuffer[i].Opcode)
	}
	fmt.Printf("--------\n")
}

type MemInfo struct {
	Value        uint8
	ReadCounter  uint64
	WriteCounter uint64
}

func (cpu *CPU) GetCoreDump() CoreDump {
	cpu.inCoreDump = true
	cpu.Bus.CountdownDisable()
	defer func() {
		cpu.inCoreDump = false
		cpu.Bus.CountdownEnable()
	}()

	var cd CoreDump
	cd.Regs = cpu.Regs
	cd.ProgramStart = uint16(0)
	if cpu.Regs.PC > 0x40 {
		cd.ProgramStart = cpu.Regs.PC - 0x40
	}
	cd.ProgramStart = (cd.ProgramStart / 0x10) * 0x10

	cd.ProgramEnd = uint16(0xffff)
	if cpu.Regs.PC < 0xffff-0x40 {
		cd.ProgramEnd = cpu.Regs.PC + 0x40
	}
	cd.ProgramEnd = (cd.ProgramEnd/0x10)*0x10 + 0x10 - 1
	cd.Program = cpu.getmem(cd.ProgramStart, cd.ProgramEnd)
	cd.HRAM = cpu.getmem(0xff80, 0xfffe)
	cd.OAM = cpu.getmem(0xfe00, 0xfe99)
	cd.VRAM = cpu.getmem(0x8000, 0x9fff)
	return cd
}

func (cpu *CPU) getmem(start, end uint16) []MemInfo {
	address := cpu.Bus.Address
	data := cpu.Bus.Data
	defer func() {
		cpu.Bus.Address = address
		cpu.Bus.Data = data
	}()

	out := make([]MemInfo, end-start+1)
	for addr := start; addr <= end; addr++ {
		memInfo := &out[addr-start]
		memInfo.ReadCounter, memInfo.WriteCounter = cpu.Bus.GetCounters(addr)
		cpu.Bus.WriteAddress(addr)
		memInfo.Value = cpu.Bus.Data
	}

	return out
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
		if handler := cpu.handlers[opcode]; handler != nil {
			e := edge{cpu.machineCycle, c.Falling}
			cpu.Debug("Handler", "e=%v", e)
			fetch = handler(e)
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
			cpu.Regs.SetWZ(0)
			cpu.machineCycle = 1
			cpu.Regs.IR = Opcode(cpu.Bus.Data)
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
	cpu.Bus.WriteAddress(addr)
}

func (cpu *CPU) writeDataBus(v uint8) {
	if !cpu.clockCycle.Falling {
		panic("writeDataBus must be called on falling edge")
	}
	cpu.Bus.WriteData(v)
}

func (cpu *CPU) readDataBus() uint8 {
	if !cpu.clockCycle.Falling {
		panic("readDataBus must be called on falling edge")
	}
	return cpu.Bus.Data
}

func (cpu *CPU) applyPendingIME() {
	if cpu.Interrupts.setIMENextCycle {
		cpu.Interrupts.setIMENextCycle = false
		cpu.Interrupts.IME = true
	}
}
