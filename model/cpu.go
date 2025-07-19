package model

type Interrupts struct {
	IF              uint8
	IE              uint8
	IME             bool
	setIMENextCycle bool
}

type CPU struct {
	PHI          *Clock
	Bus          *Bus
	Debugger     *Debugger
	Disassembler *Disassembler

	Regs       RegisterFile
	Interrupts Interrupts

	CBOp CBOp

	machineCycle int

	clockCycle                 Cycle
	wroteToAddressBusThisCycle bool

	handlers [256]InstructionHandling

	rewindBuffer    []ExecLogEntry
	rewindBufferIdx int

	nopCount    int
	nopCountMax int

	lastBranchResult int
}

func (cpu *CPU) Reset() {
	cpu.Regs = RegisterFile{}
	cpu.Regs.SP = 0xfffe
	cpu.Interrupts = Interrupts{}
	cpu.CBOp = CBOp{}
	cpu.machineCycle = 0
	cpu.clockCycle = Cycle{}
	cpu.Bus.Address = 0
	cpu.Bus.Data = 0
	cpu.Bus.inCoreDump = false
	cpu.wroteToAddressBusThisCycle = false
	clear(cpu.rewindBuffer)
	cpu.rewindBufferIdx = 0
	cpu.nopCount = 0
}

type ExecLogEntry struct {
	PC           uint16
	Opcode       Opcode
	BranchResult int
}

func (cpu *CPU) Sync(f func()) {
}

func (cpu *CPU) SetHL(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetHL must be called on rising edge")
	}
	cpu.Regs.H = uint8(v >> 8)
	cpu.Regs.L = uint8(v)
}

func (cpu *CPU) SetBC(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetBC must be called on rising edge")
	}
	cpu.Regs.B = uint8(v >> 8)
	cpu.Regs.C = uint8(v)
}

func (cpu *CPU) SetDE(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetDE must be called on rising edge")
	}
	cpu.Regs.D = uint8(v >> 8)
	cpu.Regs.E = uint8(v)
}

func (cpu *CPU) SetSP(v uint16) {
	if cpu.clockCycle.Falling {
		panic("SetSP must be called on rising edge")
	}
	cpu.Regs.SP = v
}

func (cpu *CPU) GetA() uint8 {
	v := cpu.Regs.A
	return v
}

func (cpu *CPU) GetB() uint8 {
	v := cpu.Regs.B
	return v
}

func (cpu *CPU) GetC() uint8 {
	v := cpu.Regs.C
	return v
}

func (cpu *CPU) GetD() uint8 {
	v := cpu.Regs.D
	return v
}

func (cpu *CPU) GetE() uint8 {
	v := cpu.Regs.E
	return v
}

func (cpu *CPU) GetH() uint8 {
	v := cpu.Regs.H
	return v
}

func (cpu *CPU) GetL() uint8 {
	v := cpu.Regs.L
	return v
}

func (cpu *CPU) GetBC() uint16 {
	v := join16(cpu.Regs.B, cpu.Regs.C)
	return v
}

func (cpu *CPU) GetDE() uint16 {
	v := join16(cpu.Regs.D, cpu.Regs.E)
	return v
}

func (cpu *CPU) GetHL() uint16 {
	v := join16(cpu.Regs.H, cpu.Regs.L)
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
	cpu.Regs.PC = pc
	cpu.Debugger.SetPC(pc)
}

func (cpu *CPU) IncPC() {
	if cpu.clockCycle.Falling {
		panic("IncPC must be called on rising edge")
	}
	cpu.SetPC(cpu.Regs.PC + 1)
}

func NewCPU(
	phi *Clock,
	bus *Bus,
	debugger *Debugger,
	disassembler *Disassembler,
) *CPU {
	cpu := &CPU{
		PHI:          phi,
		Bus:          bus,
		Debugger:     debugger,
		Disassembler: disassembler,
		nopCountMax:  4,
		rewindBuffer: make([]ExecLogEntry, 16),
	}
	cpu.Reset()
	cpu.handlers = handlers(cpu)
	phi.AttachDevice(cpu.fsm)
	return cpu
}

func (cpu *CPU) fsm(c Cycle) {
	cpu.wroteToAddressBusThisCycle = false

	cpu.clockCycle = c

	cpu.applyPendingIME()

	var fetch bool
	if c.C > 0 {
		opcode := cpu.Regs.IR

		// Detect runaway code
		nopCheck := opcode == OpcodeNop
		if cpu.Regs.PC < 0x100 {
			// In unmapped bootrom, allow NOPs here because soft reset puts PC at 0
			nopCheck = false
		}
		if nopCheck {
			cpu.nopCount++
			if cpu.nopCount > cpu.nopCountMax {
				panic("max nop count exceeded")
			}
		} else {
			cpu.nopCount = 0
		}

		if handler := cpu.handlers[opcode]; handler != nil {
			e := edge{cpu.machineCycle, c.Falling}
			//cpu.Debug("Handler", "e=%v", e)
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
			//cpu.Debug("PreFetch", "PC=%04x", cpu.Regs.PC)
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		} else {
			//cpu.Debug("ExecDone", "")
			cpu.Regs.SetWZ(0)
			cpu.machineCycle = 1
			cpu.Regs.IR = Opcode(cpu.Bus.Data)
			//cpu.Debug("ExecBegin", "%s", cpu.Regs.IR)
			//cpu.Debug(fmt.Sprintf("ExecBegin%s", cpu.Regs.IR), "")
			prevIdx := cpu.rewindBufferIdx
			if prevIdx > 0 {
				prevIdx--
			} else {
				prevIdx = len(cpu.rewindBuffer) - 1
			}
			cpu.rewindBuffer[prevIdx].BranchResult = cpu.lastBranchResult
			cpu.rewindBuffer[cpu.rewindBufferIdx] = ExecLogEntry{
				PC:     cpu.Regs.PC - 1,
				Opcode: cpu.Regs.IR,
			}
			cpu.lastBranchResult = 0
			cpu.rewindBufferIdx++
			if cpu.rewindBufferIdx >= len(cpu.rewindBuffer) {
				cpu.rewindBufferIdx = 0
			}
			cpu.Disassembler.SetPC(cpu.Regs.PC - 1)
		}
	} else if c.Falling {
		cpu.machineCycle++
	}
}

func (cpu *CPU) writeAddressBus(addr uint16) {
	if !cpu.Bus.inCoreDump {
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
	if cpu.Interrupts.setIMENextCycle {
		cpu.Interrupts.setIMENextCycle = false
		cpu.Interrupts.IME = true
	}
}
