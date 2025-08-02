package model

type CPU struct {
	Regs             RegisterFile
	CBOp             CBOp
	Halted           bool
	UOpCycle         int
	Rewind           Rewind
	LastBranchResult int
}

func (cpu *CPU) CurrInstruction() (DisInstruction, int) {
	return cpu.Rewind.Curr().Instruction, cpu.UOpCycle
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

func (cpu *CPU) fsm(clk *ClockRT, gb *Gameboy) {
	if cpu.Halted {
		return
	}
	var handler UOpHandler
	if gb.Interrupts.PendingInterrupt != 0 {
		handler = IRQHandler[cpu.UOpCycle-1]
	} else {
		handler = Handlers[cpu.Regs.IR][cpu.UOpCycle-1]
	}
	if done := handler(gb); done {
		gb.WriteAddress(cpu.Regs.PC)
		gb.instructionFetch(clk)
		cpu.UOpCycle = 0
	}
	cpu.UOpCycle++
}

func (gb *Gameboy) instructionFetch(clk *ClockRT) {
	cpu := &gb.CPU

	// Reset inter-instruction state
	cpu.Regs.SetWZ(0)

	// Read next instruction opcode
	rawOp := gb.Mem[cpu.Regs.PC]
	cpu.Regs.IR = Opcode(rawOp)
	gb.Debug.SetIR(gb, cpu.Regs.IR, clk)

	di := DisInstruction{
		Address: cpu.Regs.PC,
		Opcode:  cpu.Regs.IR,
		Raw:     [3]Data8{rawOp, 0, 0},
		InISR:   gb.Interrupts.InISR,
	}
	size := instSize[cpu.Regs.IR]
	if size == 0 {
		panicf("no size set for %v", cpu.Regs.IR)
	}
	for i := Size16(1); i < size; i++ {
		di.Raw[i] = gb.Mem[cpu.Regs.PC+Addr(i)]
	}

	// Update rewind buffer
	curr := cpu.Rewind.Curr()
	curr.BranchResult = cpu.LastBranchResult
	if curr.Instruction.Opcode == OpcodeNop && di.Opcode == OpcodeNop {
		curr.Instruction.NopCount++
	} else {
		cpu.LastBranchResult = 0
		entry := cpu.Rewind.Push()
		entry.Instruction = di
	}

	// Set PC
	gb.Debug.SetPC(gb, cpu.Regs.PC, clk)

	cpu.IncPC()
}
