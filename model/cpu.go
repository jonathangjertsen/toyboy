package model

type CPU struct {
	PHI          *Clock
	Bus          *Bus
	Debugger     *Debugger
	Disassembler *Disassembler

	Regs       RegisterFile
	Interrupts *Interrupts

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
	cpu.Interrupts.MemIE.Data[0] = 0
	cpu.Interrupts.MemIF.Data[0] = 0
	cpu.Interrupts.IME = false
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
	PC           Addr
	Opcode       Opcode
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
	if v == 0 {
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

func msb(w uint16) uint8 {
	return uint8((w >> 8) & 0xff)
}

func lsb(w uint16) uint8 {
	return uint8(w & 0xff)
}

func (cpu *CPU) SetPC(pc Addr) {
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
	interrupts *Interrupts,
	bus *Bus,
	debugger *Debugger,
	disassembler *Disassembler,
) *CPU {
	cpu := &CPU{
		PHI:          phi,
		Bus:          bus,
		Debugger:     debugger,
		Disassembler: disassembler,
		Interrupts:   interrupts,
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
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		} else {
			cpu.Regs.SetWZ(0)
			cpu.machineCycle = 1
			cpu.Regs.IR = Opcode(cpu.Bus.Data)
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

func (cpu *CPU) writeAddressBus(addr Addr) {
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
		cpu.Interrupts.SetIME(true)
	}
}
