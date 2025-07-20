package model

type Gameboy struct {
	Config HWConfig

	CLK           *ClockRT
	Debugger      *Debugger
	Disassembler  *Disassembler
	PHI           *Clock
	CPU           *CPU
	PPU           *PPU
	CartridgeSlot *MemoryRegion
	Joypad        *Joypad
}

func (gb *Gameboy) Start() {
	gb.CLK.Start()
}

func (gb *Gameboy) Pause() {
	gb.CLK.Pause()
}

func (gb *Gameboy) Step() {
	gb.CLK.pauseAfterCycle.Add(1)
	gb.CLK.Start()
}

func (gb *Gameboy) SoftReset() {
	gb.CLK.Sync(func() {
		gb.CLK.cycle.C = 0
		gb.CLK.cycle.Falling = false
		for i := range gb.CLK.divided {
			gb.CLK.divided[i].counter = 0
			gb.CLK.divided[i].cycle = 0
		}
		gb.CPU.Reset()
	})
}

func (gb *Gameboy) GetCoreDump() CoreDump {
	var cd CoreDump
	gb.CLK.Sync(func() {
		cd = gb.CPU.GetCoreDump()
		cd.Cycle = gb.CLK.cycle
	})
	return cd
}

func (gb *Gameboy) GetViewport() ViewPort {
	var vp ViewPort
	gb.CLK.Sync(func() {
		vp = gb.PPU.LastFrame
	})
	return vp
}

func NewGameboy(
	config HWConfig,
) *Gameboy {
	gameboy := &Gameboy{
		Config: config,
	}
	gameboy.init()
	return gameboy
}

func (gb *Gameboy) init() {
	clk := NewRealtimeClock(gb.Config.SystemClock)
	debugger := NewDebugger(clk)
	disassembler := NewDisassembler()

	interrupts := NewInterrupts(clk)

	bootROMLock := NewBootROMLock(clk)
	bootROM := NewBootROM(clk, gb.Config.Model)
	disassembler.SetProgram(ByteSlice(bootROM.Data))
	disassembler.SetPC(0)
	vram := NewMemoryRegion(clk, AddrVRAMBegin, SizeVRAM)
	hram := NewMemoryRegion(clk, AddrHRAMBegin, SizeHRAM)
	wram := NewMemoryRegion(clk, AddrWRAMBegin, SizeWRAM)
	apu := NewAPU(clk)
	oam := NewMemoryRegion(clk, AddrOAMBegin, SizeOAM)
	cartridgeSlot := NewMemoryRegion(clk, AddrZero, SizeCartridgeBank0)
	joypad := NewJoypad(clk, interrupts)

	bootROMLock.OnUnlock = func() {
		disassembler.SetProgram(ByteSlice(cartridgeSlot.Data))
		disassembler.SetPC(0x100)
	}

	bus := &Bus{}

	cpuClock := clk.Divide(4)
	cpu := NewCPU(cpuClock, interrupts, bus, debugger, disassembler)

	ppuClock := clk.Divide(2)
	ppu := NewPPU(clk, ppuClock, interrupts, bus, debugger)

	bus.BootROMLock = bootROMLock
	bus.BootROM = &bootROM
	bus.VRAM = &vram
	bus.WRAM = &wram
	bus.HRAM = &hram
	bus.APU = apu
	bus.OAM = &oam
	bus.PPU = ppu
	bus.CartridgeSlot = &cartridgeSlot
	bus.Joypad = joypad
	bus.Interrupts = interrupts

	gb.CLK = clk
	gb.Disassembler = disassembler
	gb.PHI = cpuClock
	gb.CPU = cpu
	gb.CartridgeSlot = &cartridgeSlot
	gb.PPU = ppu
	gb.Debugger = debugger
	gb.Joypad = joypad

	clk.Onpanic = gb.CPU.Dump
}
