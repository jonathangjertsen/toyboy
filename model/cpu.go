package model

import "fmt"

type CPU struct {
	Config *Config
	PHI    *Clock
	Bus    CPUBusIF
	Debug  *Debug

	Regs       RegisterFile
	Interrupts *Interrupts

	CBOp CBOp

	halted bool

	machineCycle int

	clockCycle                 Cycle
	wroteToAddressBusThisCycle bool

	handlers [256]InstructionHandling

	rewindBuffer     []ExecLogEntry
	rewindBufferIdx  int
	rewindBufferFull bool

	nopCount int

	lastBranchResult int
}

type CPUBusIF interface {
	Reset()
	BeginCoreDump() func()
	InCoreDump() bool
	WriteAddress(Addr)
	WriteData(Data8)
	GetAddress() Addr
	GetData() Data8
	GetCounters(Addr) (uint64, uint64)
	GetPeripheral(any)
	PushState() func()
}

func (cpu *CPU) Reset() {
	cpu.Regs = RegisterFile{}
	cpu.Regs.SP = 0xfffe
	if cpu.Interrupts != nil {
		cpu.Interrupts.MemIE.Data[0] = 0
		cpu.Interrupts.MemIF.Data[0] = 0
		cpu.Interrupts.IME = false
	}
	cpu.CBOp = CBOp{}
	cpu.machineCycle = 0
	cpu.clockCycle = Cycle{}
	cpu.Bus.Reset()
	cpu.wroteToAddressBusThisCycle = false
	clear(cpu.rewindBuffer)
	cpu.rewindBufferIdx = 0
	cpu.nopCount = 0

	if cpu.Config.BootROM.Skip {
		cpu.Regs.A = 0x01
		cpu.Regs.F = 0xB0
		cpu.Regs.B = 0x00
		cpu.Regs.C = 0x13
		cpu.Regs.D = 0x00
		cpu.Regs.E = 0xD8
		cpu.Regs.H = 0x01
		cpu.Regs.L = 0x4D
		cpu.Regs.PC = 0x00ff
		cpu.Regs.SP = 0xFFFE
	}
}

type ExecLogEntry struct {
	Instruction  DisInstruction
	BranchResult int
}

func (cpu *CPU) SetHL(v Data16) {
	if cpu.clockCycle.Falling {
		panic("SetHL must be called on rising edge")
	}
	cpu.Regs.H, cpu.Regs.L = v.Split()
}

func (cpu *CPU) SetBC(v Data16) {
	if cpu.clockCycle.Falling {
		panic("SetBC must be called on rising edge")
	}
	cpu.Regs.B, cpu.Regs.C = v.Split()
}

func (cpu *CPU) SetDE(v Data16) {
	if cpu.clockCycle.Falling {
		panic("SetDE must be called on rising edge")
	}
	cpu.Regs.D, cpu.Regs.E = v.Split()
}

func (cpu *CPU) SetSP(v Addr) {
	if cpu.clockCycle.Falling {
		panic("SetSP must be called on rising edge")
	}
	if v == 0 && cpu.Config.Debug.PanicOnStackUnderflow {
		panic("stack underflow")
	}
	cpu.Regs.SP = v
}

func (cpu *CPU) GetBC() Data16 {
	v := join16(cpu.Regs.B, cpu.Regs.C)
	return v
}

func (cpu *CPU) GetDE() Data16 {
	v := join16(cpu.Regs.D, cpu.Regs.E)
	return v
}

func (cpu *CPU) GetHL() Data16 {
	v := join16(cpu.Regs.H, cpu.Regs.L)
	return v
}

func join16(msb, lsb Data8) Data16 {
	return (Data16(msb) << 8) | Data16(lsb)
}

func (cpu *CPU) SetPC(pc Addr) {
	if cpu.clockCycle.Falling {
		panic("SetPC must be called on rising edge")
	}
	cpu.Regs.PC = pc
	cpu.Debug.SetPC(pc)
}

func (cpu *CPU) IncPC() {
	if cpu.clockCycle.Falling {
		panic("IncPC must be called on rising edge")
	}
	cpu.SetPC(cpu.Regs.PC + 1)
}

// Must call Reset before starting
func NewCPU(
	phi *Clock,
	interrupts *Interrupts,
	bus CPUBusIF,
	config *Config,
	debug *Debug,
) *CPU {
	cpu := &CPU{
		Config:       config,
		PHI:          phi,
		Bus:          bus,
		Debug:        debug,
		Interrupts:   interrupts,
		rewindBuffer: make([]ExecLogEntry, 8192),
	}
	cpu.handlers = handlers(cpu)
	phi.AttachDevice(cpu.fsm)
	return cpu
}

func (cpu *CPU) fsm(c Cycle) {
	cpu.wroteToAddressBusThisCycle = false
	cpu.clockCycle = c
	cpu.applyPendingIME()

	if cpu.halted {
		if cpu.Interrupts == nil {
			panic("can't be HALTed without interrupts")
		}
		if cpu.Interrupts.PendingInterrupt == 0 {
			return
		}
		cpu.halted = false
	}

	var fetch bool
	if c.C > 0 {
		cpu.detectRunawayCode()
		if cpu.Interrupts != nil && cpu.Interrupts.PendingInterrupt != 0 {
			fetch = cpu.execTransferToISR()
		} else {
			fetch = cpu.execCurrentInstruction()
		}
	} else {
		// initial instruction
		fetch = true
		if c.Falling {
			cpu.machineCycle++
		}
	}

	if fetch {
		if !c.Falling {
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		} else if cpu.Interrupts == nil || cpu.Interrupts.PendingInterrupt == 0 {
			cpu.instructionFetch()
		}
	} else if c.Falling {
		cpu.machineCycle++
	}

	if cpu.Config.Debug.GBD.Enable && !c.Falling && fetch && cpu.machineCycle == 1 {
		cpu.doGBDLog()
	}
}

func (cpu *CPU) doGBDLog() {
	end := cpu.Bus.BeginCoreDump()
	defer end()

	pc := cpu.Regs.PC - 1

	origAddr := cpu.Bus.GetAddress()
	var pcmem [4]Data8
	for i := range Addr(4) {
		cpu.writeAddressBus(pc + i)
		pcmem[i] = cpu.Bus.GetData()
	}
	cpu.writeAddressBus(origAddr)

	cpu.Config.Debug.GBD.GBDLog(fmt.Sprintf(
		"A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X\n",
		uint(cpu.Regs.A),
		uint(cpu.Regs.F),
		uint(cpu.Regs.B),
		uint(cpu.Regs.C),
		uint(cpu.Regs.D),
		uint(cpu.Regs.E),
		uint(cpu.Regs.H),
		uint(cpu.Regs.L),
		uint(cpu.Regs.SP),
		uint(pc),
		uint(pcmem[0]),
		uint(pcmem[1]),
		uint(pcmem[2]),
		uint(pcmem[3]),
	))
}

func (cpu *CPU) detectRunawayCode() {
	// Detect runaway code
	nopCheck := cpu.Regs.IR == OpcodeNop
	if cpu.Regs.PC < 0x100 {
		// In unmapped bootrom, allow NOPs here because soft reset puts PC at 0
		nopCheck = false
	}
	if nopCheck && !cpu.clockCycle.Falling {
		cpu.nopCount++
		if cpu.nopCount > cpu.Config.Debug.MaxNOPCount {
			panic("max nop count exceeded")
		}
	} else {
		cpu.nopCount = 0
	}
}

func (cpu *CPU) instructionFetch() {
	// Reset inter-instruction state
	cpu.Regs.SetWZ(0)
	cpu.machineCycle = 1

	// Read next instruction opcode
	rawOp := cpu.Bus.GetData()
	cpu.Regs.IR = Opcode(rawOp)
	cpu.Debug.SetIR(cpu.Regs.IR)

	di := DisInstruction{
		Address: cpu.Regs.PC - 1,
		Opcode:  cpu.Regs.IR,
		Raw:     [3]Data8{rawOp, 0, 0},
	}
	size := instSize[cpu.Regs.IR]
	if size == 0 {
		panicf("no size set for %v", cpu.Regs.IR)
	}
	pop := cpu.Bus.PushState()
	for i := Size16(1); i < size; i++ {
		cpu.Bus.WriteAddress(cpu.Regs.PC - 1 + Addr(i))
		di.Raw[i] = cpu.Bus.GetData()
	}
	pop()

	// Update rewind buffer
	prevIdx := cpu.rewindBufferIdx
	if prevIdx > 0 {
		prevIdx--
	} else {
		prevIdx = len(cpu.rewindBuffer) - 1
	}
	cpu.rewindBuffer[prevIdx].BranchResult = cpu.lastBranchResult
	cpu.rewindBuffer[cpu.rewindBufferIdx] = ExecLogEntry{Instruction: di}

	cpu.lastBranchResult = 0
	cpu.rewindBufferIdx++
	if cpu.rewindBufferIdx >= len(cpu.rewindBuffer) {
		cpu.rewindBufferIdx = 0
		cpu.rewindBufferFull = true
	}

	// Set PC
	cpu.Debug.SetPC(cpu.Regs.PC - 1)
}

func (cpu *CPU) execCurrentInstruction() bool {
	opcode := cpu.Regs.IR
	if handler := cpu.handlers[opcode]; handler != nil {
		e := edge{cpu.machineCycle, cpu.clockCycle.Falling}
		return handler(e)
	}
	panicf("not implemented opcode %v", opcode)
	return false
}

func (cpu *CPU) execTransferToISR() bool {
	e := edge{cpu.machineCycle, cpu.clockCycle.Falling}
	switch e {
	// wait states
	case edge{1, false}, edge{1, true}, edge{2, false}, edge{2, true}:
	// push MSB of PC to stack
	case edge{3, false}:
		cpu.SetSP(cpu.Regs.SP - 1)
		cpu.writeAddressBus(cpu.Regs.SP)
	case edge{3, true}:
		cpu.Bus.WriteData(cpu.Regs.PC.MSB())
		// push LSB of PC to stack
	case edge{4, false}:
		cpu.SetSP(cpu.Regs.SP - 1)
		cpu.writeAddressBus(cpu.Regs.SP)
	case edge{4, true}:
		cpu.Bus.WriteData(cpu.Regs.PC.LSB())
	case edge{5, false}:
		cpu.SetPC(cpu.Interrupts.PendingInterrupt.ISR())
		return true
	case edge{5, true}:
		cpu.Interrupts.PendingInterrupt = 0
		return true
	default:
		panicv(e)
	}
	return false
}

func (cpu *CPU) writeAddressBus(addr Addr) {
	if !cpu.Bus.InCoreDump() {
		if cpu.clockCycle.Falling {
			panic("writeAddressBus must be called on rising edge")
		}
		if cpu.wroteToAddressBusThisCycle {
			panic("more than one call to writeAddressBus this cycle")
		}
	}
	cpu.wroteToAddressBusThisCycle = true
	cpu.Bus.WriteAddress(addr)
}

func (cpu *CPU) applyPendingIME() {
	if cpu.Interrupts == nil {
		return
	}
	if cpu.Interrupts.setIMENextCycle {
		cpu.Interrupts.setIMENextCycle = false
		cpu.Interrupts.SetIME(true)
	}
}
