package model

type CPU struct {
	Config     *Config
	Bus        CPUBusIF
	Debug      *Debug
	Interrupts *Interrupts

	handlers [256]InstructionHandling

	Regs                       RegisterFile
	CBOp                       CBOp
	halted                     bool
	machineCycle               int
	clockCycle                 uint
	wroteToAddressBusThisCycle bool // deprecated
	rewind                     Rewind
	lastBranchResult           int
}

func (cpu *CPU) LoadSave(save *SaveState, mem []Data8) {
	cpu.Regs = save.CPURegisterFile
	cpu.machineCycle = save.CPUMachineCycle
	cpu.clockCycle = save.CPUClockCycle
	cpu.wroteToAddressBusThisCycle = save.CPUWroteToAddressBusThisCycle
	cpu.lastBranchResult = save.CPULastBranchResult
	cpu.halted = save.CPUHalted
	cpu.rewind = save.CPURewindBuffer
	cpu.CBOp = save.CPUCBOp
}

func (cpu *CPU) Save(save *SaveState) {
	save.CPURegisterFile = cpu.Regs
	save.CPUMachineCycle = cpu.machineCycle
	save.CPUClockCycle = cpu.clockCycle
	save.CPUWroteToAddressBusThisCycle = cpu.wroteToAddressBusThisCycle
	save.CPULastBranchResult = cpu.lastBranchResult
	save.CPUHalted = cpu.halted
	save.CPURewindBuffer = cpu.rewind
	save.CPUCBOp = cpu.CBOp
}

type CPUBusIF interface {
	Reset()
	WriteAddress([]Data8, Addr)
	WriteData([]Data8, Data8)
	GetAddress() Addr
	GetData() Data8
	ProbeAddress([]Data8, Addr) Data8
	ProbeRange([]Data8, Addr, Addr) []Data8
}

func (cpu *CPU) CurrInstruction() (DisInstruction, int) {
	return cpu.rewind.Curr().Instruction, cpu.machineCycle
}

func (cpu *CPU) Reset() {
	cpu.Regs = RegisterFile{}
	cpu.Regs.SP = 0xfffe
	if cpu.Interrupts != nil {
		cpu.Interrupts.mem[AddrIE] = 0
		cpu.Interrupts.mem[AddrIF] = 0
		cpu.Interrupts.IME = false
	}
	cpu.CBOp = CBOp{}
	cpu.machineCycle = 0
	cpu.clockCycle = 0
	cpu.Bus.Reset()
	cpu.wroteToAddressBusThisCycle = false
	cpu.rewind.Reset()

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
	cpu.Regs.H, cpu.Regs.L = v.Split()
}

func (cpu *CPU) SetBC(v Data16) {
	cpu.Regs.B, cpu.Regs.C = v.Split()
}

func (cpu *CPU) SetDE(v Data16) {
	cpu.Regs.D, cpu.Regs.E = v.Split()
}

func (cpu *CPU) SetSP(v Addr) {
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
	cpu.Regs.PC = pc
}

func (cpu *CPU) IncPC() {
	cpu.SetPC(cpu.Regs.PC + 1)
}

// Must call Reset before starting
func NewCPU(
	clk *ClockRT,
	interrupts *Interrupts,
	bus CPUBusIF,
	config *Config,
	debug *Debug,
) *CPU {
	cpu := &CPU{
		Config:     config,
		Bus:        bus,
		Debug:      debug,
		Interrupts: interrupts,
		rewind:     NewRewind(8192),
	}
	cpu.handlers = handlers(cpu)
	clk.cpu = cpu
	return cpu
}

func (cpu *CPU) fsm(clk *ClockRT, mem []Data8) {
	cpu.wroteToAddressBusThisCycle = false
	cpu.clockCycle = clk.Cycle
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
	if cpu.clockCycle > 0 {
		if cpu.Interrupts != nil && cpu.Interrupts.PendingInterrupt != 0 {
			fetch = cpu.execTransferToISR(clk, mem)
		} else {
			fetch = cpu.handlers[cpu.Regs.IR](mem, cpu.machineCycle)
		}
		if fetch {
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			if cpu.Interrupts == nil || cpu.Interrupts.PendingInterrupt == 0 {
				cpu.instructionFetch(clk, mem)
			}
			cpu.IncPC()
			cpu.machineCycle = 0
		}
	} else {
		// initial instruction
		fetch = true
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.instructionFetch(clk, mem)
		cpu.IncPC()
	}

	cpu.machineCycle++
}

func (cpu *CPU) instructionFetch(clk *ClockRT, mem []Data8) {
	// Reset inter-instruction state
	cpu.Regs.SetWZ(0)

	// Read next instruction opcode
	rawOp := mem[cpu.Bus.GetAddress()]
	cpu.Regs.IR = Opcode(rawOp)
	cpu.Debug.SetIR(cpu.Regs.IR, clk)

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
		di.Raw[i] = mem[cpu.Regs.PC+Addr(i)]
	}

	// Update rewind buffer
	cpu.rewind.Curr().BranchResult = cpu.lastBranchResult
	cpu.lastBranchResult = 0
	entry := cpu.rewind.Push()
	entry.Instruction = di

	// Set PC
	cpu.Debug.SetPC(cpu.Regs.PC, clk)
}

func (cpu *CPU) execTransferToISR(clk *ClockRT, mem []Data8) bool {
	switch cpu.machineCycle {
	// wait states
	case 1, 2:
	// push MSB of PC to stack
	case 3:
		cpu.SetSP(cpu.Regs.SP - 1)
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.Bus.WriteData(mem, cpu.Regs.PC.MSB())
		// push LSB of PC to stack
	case 4:
		cpu.SetSP(cpu.Regs.SP - 1)
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.Bus.WriteData(mem, cpu.Regs.PC.LSB())
	case 5:
		isr := cpu.Interrupts.PendingInterrupt.ISR()
		cpu.SetPC(isr)
		cpu.Debug.SetPC(isr, clk)
		cpu.Interrupts.PendingInterrupt = 0
		return true
	default:
		panicv(cpu.machineCycle)
	}
	return false
}

func (cpu *CPU) writeAddressBus(mem []Data8, addr Addr) {
	cpu.Bus.WriteAddress(mem, addr)
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
