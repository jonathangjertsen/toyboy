package model

type Gameboy struct {
	Config *Config

	CLK       *ClockRT
	Debug     *Debug
	CPU       *CPU
	PPU       *PPU
	Cartridge *Cartridge
	Joypad    *Joypad
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
		gb.CLK.Cycle.C = 0
		gb.CLK.Cycle.Falling = false
		gb.CPU.Reset()
		gb.PPU.Reset()
	})
}

func (gb *Gameboy) GetCoreDump() CoreDump {
	var cd CoreDump
	gb.CLK.Sync(func() {
		cd = gb.CPU.GetCoreDump()
		cd.Cycle = gb.CLK.Cycle
	})
	return cd
}

func NewGameboy(
	config *Config,
) *Gameboy {
	gameboy := &Gameboy{
		Config: config,
	}
	gameboy.Init()
	return gameboy
}

func (gb *Gameboy) Init() {
	clk := NewRealtimeClock(gb.Config.Clock)

	debug := NewDebug(clk, &gb.Config.Debug)

	interrupts := NewInterrupts(clk)

	bootROMLock := NewBootROMLock(clk)
	bootROM := NewBootROM(clk, gb.Config.BootROM)
	debug.SetProgram(ByteSlice(bootROM.Data))
	debug.SetPC(0)

	vram := NewMemoryRegion(clk, AddrVRAMBegin, SizeVRAM)
	hram := NewMemoryRegion(clk, AddrHRAMBegin, SizeHRAM)
	wram := NewMemoryRegion(clk, AddrWRAMBegin, SizeWRAM)
	apu := NewAPU(clk, gb.Config)
	oam := NewMemoryRegion(clk, AddrOAMBegin, SizeOAM)
	cartridge := NewCartridge(clk)
	joypad := NewJoypad(clk, interrupts)
	serial := NewSerial(clk)
	prohibited := NewProhibited(clk)
	timer := NewTimer(clk, apu, interrupts)

	bootROMLock.OnLock = func() {
		debug.SetProgram(ByteSlice(cartridge.CurrROMBank0.Data))
		debug.SetPC(0x100)
	}

	bus := &Bus{}

	cpu := NewCPU(clk, interrupts, bus, gb.Config, debug)

	ppu := NewPPU(clk, interrupts, bus, gb.Config, debug)

	bus.BootROMLock = bootROMLock
	bus.BootROM = &bootROM
	bus.VRAM = &vram
	bus.WRAM = &wram
	bus.HRAM = &hram
	bus.APU = apu
	bus.OAM = &oam
	bus.PPU = ppu
	bus.Cartridge = cartridge
	bus.Joypad = joypad
	bus.Interrupts = interrupts
	bus.Serial = serial
	bus.Prohibited = prohibited
	bus.Timer = timer
	bus.Config = gb.Config

	debug.HRAM.Source = hram.Data
	debug.WRAM.Source = wram.Data

	gb.CLK = clk
	gb.CPU = cpu
	gb.Cartridge = cartridge
	gb.PPU = ppu
	gb.Debug = debug
	gb.Joypad = joypad

	gb.CPU.Reset()

	clk.Onpanic = gb.CPU.Dump
}
