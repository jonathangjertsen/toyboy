package model

type CPU struct {
	Config *Config
	Bus    CPUBusIF
	Debug  *Debug

	Regs       RegisterFile
	Interrupts *Interrupts

	halted bool

	machineCycle int

	clockCycle                 uint
	wroteToAddressBusThisCycle bool

	rewind *Rewind

	lastBranchResult int

	handler CycleHandler
}

type CPUBusIF interface {
	Reset()
	BeginCoreDump() func()
	InCoreDump() bool
	WriteAddress(Addr)
	WriteData(Data8)
	GetAddress() Addr
	GetData() Data8
	ProbeAddress(Addr) Data8
	ProbeRange(Addr, Addr) []Data8
	GetPeripheral(any)
}

func (cpu *CPU) CurrInstruction() (DisInstruction, int) {
	return cpu.rewind.Curr().Instruction, cpu.machineCycle
}

func (cpu *CPU) Reset() {
	cpu.Regs = RegisterFile{}
	cpu.Regs.SP = 0xfffe
	if cpu.Interrupts != nil {
		cpu.Interrupts.MemIE.Data[0] = 0
		cpu.Interrupts.MemIF.Data[0] = 0
		cpu.Interrupts.IME = false
	}
	cpu.machineCycle = 0
	cpu.clockCycle = 0
	cpu.Bus.Reset()
	cpu.wroteToAddressBusThisCycle = false
	cpu.rewind.Reset()

	if cpu.Config.BootROM.Skip {
		cpu.Regs.A = 0x01
		cpu.Regs.F = 0xB0
		cpu.Regs.BC = RegisterPair{0x00, 0x13}
		cpu.Regs.DE = RegisterPair{0x00, 0xD8}
		cpu.Regs.HL = RegisterPair{0x01, 0x4D}
		cpu.Regs.PC = 0x00ff
		cpu.Regs.SP = 0xFFFE
	}
}

type ExecLogEntry struct {
	Instruction  DisInstruction
	BranchResult int
}

func (cpu *CPU) SetSP(v Addr) {
	if v == 0 && cpu.Config.Debug.PanicOnStackUnderflow {
		panic("stack underflow")
	}
	cpu.Regs.SP = v
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
	initOpcode Opcode,
) *CPU {
	cpu := &CPU{
		Config:     config,
		Bus:        bus,
		Debug:      debug,
		Interrupts: interrupts,
		rewind:     NewRewind(1024 * 16),
		handler:    Handlers[initOpcode],
	}

	clk.cpu = cpu
	return cpu
}

func (cpu *CPU) fsm(c uint) {
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
	if cpu.Interrupts != nil && cpu.Interrupts.PendingInterrupt != 0 {
		fetch = cpu.execTransferToISR()
	} else {
		if cpu.handler == nil {
			// Start new instruction
			// Can use either switch statement or jump table. Profile to see which is faster.

			// Uncomment to use switch statement version:
			cpu.handler = Handler(cpu)

			// Uncomment to use jump table version:
			//cpu.handler = Handlers[cpu.Regs.IR](cpu)
		} else {
			// Continue existing instructino
			cpu.handler = cpu.handler(cpu)
		}
		fetch = cpu.handler == nil
	}
	if fetch {
		cpu.Bus.WriteAddress(cpu.Regs.PC)
		if cpu.Interrupts == nil || cpu.Interrupts.PendingInterrupt == 0 {
			cpu.instructionFetch()
		}
		cpu.IncPC()
		cpu.machineCycle = 0
	}
	cpu.machineCycle++
}

func (cpu *CPU) instructionFetch() {
	// Reset inter-instruction state
	cpu.Regs.Temp = RegisterPair{}
	clear(cpu.Regs.TempPtr[:])
	cpu.Regs.TempCond = false

	// Read next instruction opcode
	rawOp := cpu.Bus.ProbeAddress(cpu.Bus.GetAddress())
	cpu.Regs.IR = Opcode(rawOp)
	cpu.Debug.SetIR(cpu.Regs.IR)

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
		di.Raw[i] = cpu.Bus.ProbeAddress(cpu.Regs.PC + Addr(i))
	}

	// Update rewind buffer
	cpu.rewind.Curr().BranchResult = cpu.lastBranchResult
	cpu.lastBranchResult = 0
	entry := cpu.rewind.Push()
	entry.Instruction = di

	// Set PC
	cpu.Debug.SetPC(cpu.Regs.PC)
}

func (cpu *CPU) execTransferToISR() bool {
	switch cpu.machineCycle {
	// wait states
	case 1, 2:
	// push MSB of PC to stack
	case 3:
		cpu.SetSP(cpu.Regs.SP - 1)
		cpu.Bus.WriteAddress(cpu.Regs.SP)
		cpu.Bus.WriteData(cpu.Regs.PC.MSB())
		// push LSB of PC to stack
	case 4:
		cpu.SetSP(cpu.Regs.SP - 1)
		cpu.Bus.WriteAddress(cpu.Regs.SP)
		cpu.Bus.WriteData(cpu.Regs.PC.LSB())
	case 5:
		isr := cpu.Interrupts.PendingInterrupt.ISR()
		cpu.SetPC(isr)
		cpu.Debug.SetPC(isr)
		cpu.Interrupts.PendingInterrupt = 0
		return true
	default:
		panicv(cpu.machineCycle)
	}
	return false
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
