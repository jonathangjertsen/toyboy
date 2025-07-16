package model

type Gameboy struct {
	Config HWConfig

	CLK           *ClockRT
	PHI           *Clock
	CPU           *CPU
	PPU           *PPU
	CartridgeSlot *MemoryRegion
}

func (gb *Gameboy) PowerOn() {
	gb.CLK.Start()
}

func (gb *Gameboy) Pause() {
	gb.CLK.Pause()
}

func (gb *Gameboy) Stop() {
	gb.CLK.Stop()
}

func (gb *Gameboy) GetCoreDump() CoreDump {
	var cd CoreDump
	gb.CLK.Sync(func() {
		cd = gb.CPU.GetCoreDump()
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
	ppuClock := clk.Divide(2)
	cpuClock := clk.Divide(4)

	bootROMLock := NewBootROMLock(clk)
	bootROM := NewBootROM(clk, gb.Config.Model)
	vram := NewMemoryRegion(clk, AddrVRAMBegin, SizeVRAM)
	hram := NewMemoryRegion(clk, AddrHRAMBegin, SizeHRAM)
	apu := NewAPU(clk)
	oam := NewMemoryRegion(clk, AddrOAMBegin, SizeOAM)
	cartridgeSlot := NewMemoryRegion(clk, AddrCartridgeBank0Begin, AddrCartridgeBank0Size)

	bus := &Bus{}
	ppu := NewPPU(clk, ppuClock, bus)

	bus.BootROMLock = bootROMLock
	bus.BootROM = &bootROM
	bus.VRAM = &vram
	bus.HRAM = &hram
	bus.APU = apu
	bus.OAM = &oam
	bus.PPU = ppu
	bus.CartridgeSlot = &cartridgeSlot

	cpu := NewCPU(cpuClock, bus)

	gb.CLK = clk
	gb.PHI = cpuClock
	gb.CPU = cpu
	gb.CartridgeSlot = &cartridgeSlot
	gb.PPU = ppu
}
