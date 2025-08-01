package model

type CPU struct {
	Regs                       RegisterFile
	CBOp                       CBOp
	Halted                     bool
	MachineCycle               int
	ClockCycle                 uint
	WroteToAddressBusThisCycle bool // deprecated
	Rewind                     Rewind
	LastBranchResult           int
}

func (cpu *CPU) CurrInstruction() (DisInstruction, int) {
	return cpu.Rewind.Curr().Instruction, cpu.MachineCycle
}

type ExecLogEntry struct {
	Instruction  DisInstruction
	BranchResult int
}

func (cpu *CPU) SetHL(v Data16) {
	cpu.Regs.H, cpu.Regs.L = v.Split()
}

func (cpu *CPU) SetBC(v Data16) {
	cpu.Regs.B, cpu.Regs.C = v.Split()
}

func (cpu *CPU) SetDE(v Data16) {
	cpu.Regs.D, cpu.Regs.E = v.Split()
}

func (cpu *CPU) SetSP(v Addr) {
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

func (cpu *CPU) GetWZ() Data16 {
	v := join16(cpu.Regs.TempW, cpu.Regs.TempZ)
	return v
}

func join16(msb, lsb Data8) Data16 {
	return (Data16(msb) << 8) | Data16(lsb)
}

func (cpu *CPU) SetPC(pc Addr) {
	cpu.Regs.PC = pc
}

func (cpu *CPU) IncPC() {
	cpu.SetPC(cpu.Regs.PC + 1)
}

func (cpu *CPU) fsm(clk *ClockRT, gb *Gameboy, handlers *HandlerArray) {
	cpu.WroteToAddressBusThisCycle = false
	cpu.ClockCycle = clk.Cycle
	cpu.applyPendingIME(gb)

	if cpu.Halted {
		if gb.Interrupts.PendingInterrupt == 0 {
			return
		}
		cpu.Halted = false
	}

	var fetch bool
	if cpu.ClockCycle > 0 {
		if gb.Interrupts.PendingInterrupt != 0 {
			fetch = cpu.execTransferToISR(clk, gb)
		} else {
			handler := handlers[cpu.Regs.IR][cpu.MachineCycle-1]
			if h1, ok := handler.(func(gb *Gameboy, e int) bool); ok {
				// patch
				for i := 1; i < 6; i++ {
					handlers[cpu.Regs.IR][i] = handler
				}
				fetch = h1(gb, cpu.MachineCycle)
			} else if h2, ok := handler.(func(gb *Gameboy) bool); ok {
				fetch = h2(gb)
			} else {
				panicf("%T @ %s %d", handler, cpu.Regs.IR, cpu.MachineCycle-1)
			}
		}
		if fetch {
			gb.WriteAddress(cpu.Regs.PC)
			if gb.Interrupts.PendingInterrupt == 0 {
				cpu.instructionFetch(clk, gb)
			}
			cpu.IncPC()
			cpu.MachineCycle = 0
		}
	} else {
		// initial instruction
		fetch = true
		gb.WriteAddress(cpu.Regs.PC)
		cpu.instructionFetch(clk, gb)
		cpu.IncPC()
	}

	cpu.MachineCycle++
}

func (cpu *CPU) instructionFetch(clk *ClockRT, gb *Gameboy) {
	// Reset inter-instruction state
	cpu.Regs.SetWZ(0)

	// Read next instruction opcode
	rawOp := gb.Mem[cpu.Regs.PC]
	cpu.Regs.IR = Opcode(rawOp)
	gb.Debug.SetIR(cpu.Regs.IR, clk)

	di := DisInstruction{
		Address: cpu.Regs.PC,
		Opcode:  cpu.Regs.IR,
		Raw:     [3]Data8{rawOp, 0, 0},
	}
	size := instSize[cpu.Regs.IR]
	if size == 0 {
		panicf("no size set for %v", cpu.Regs.IR)
	}
	for i := Size16(1); i < size; i++ {
		di.Raw[i] = gb.Mem[cpu.Regs.PC+Addr(i)]
	}

	// Update rewind buffer
	cpu.Rewind.Curr().BranchResult = cpu.LastBranchResult
	cpu.LastBranchResult = 0
	entry := cpu.Rewind.Push()
	entry.Instruction = di

	// Set PC
	gb.Debug.SetPC(cpu.Regs.PC, clk)
}

func (cpu *CPU) execTransferToISR(clk *ClockRT, gb *Gameboy) bool {
	switch cpu.MachineCycle {
	// wait states
	case 1, 2:
	// push MSB of PC to stack
	case 3:
		cpu.SetSP(cpu.Regs.SP - 1)
		gb.WriteAddress(cpu.Regs.SP)
		gb.WriteData(cpu.Regs.PC.MSB())
		// push LSB of PC to stack
	case 4:
		cpu.SetSP(cpu.Regs.SP - 1)
		gb.WriteAddress(cpu.Regs.SP)
		gb.WriteData(cpu.Regs.PC.LSB())
	case 5:
		isr := gb.Interrupts.PendingInterrupt.ISR()
		cpu.SetPC(isr)
		gb.Debug.SetPC(isr, clk)
		gb.Interrupts.PendingInterrupt = 0
		return true
	default:
		panicv(cpu.MachineCycle)
	}
	return false
}

func (cpu *CPU) applyPendingIME(gb *Gameboy) {
	if gb.Interrupts.SetIMENextCycle {
		gb.Interrupts.SetIMENextCycle = false
		gb.Interrupts.SetIME(gb.Mem, true)
	}
}
