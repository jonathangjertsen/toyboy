package model

type Gameboy struct {
	Config HWConfig

	CLK           *ClockRT
	Debugger      *Debugger
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

func (gb *Gameboy) ButtonRight(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Direction &= ^uint8(1 << 0)
		} else {
			gb.Joypad.Direction |= 1 << 0
		}
	})
}

func (gb *Gameboy) ButtonLeft(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Direction &= ^uint8(1 << 1)
		} else {
			gb.Joypad.Direction |= 1 << 1
		}
	})
}

func (gb *Gameboy) ButtonUp(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Direction &= ^uint8(1 << 2)
		} else {
			gb.Joypad.Direction |= 1 << 2
		}
	})
}

func (gb *Gameboy) ButtonDown(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Direction &= ^uint8(1 << 3)
		} else {
			gb.Joypad.Direction |= 1 << 3
		}
	})
}

func (gb *Gameboy) ButtonA(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Action &= ^uint8(1 << 0)
		} else {
			gb.Joypad.Action |= 1 << 0
		}
	})
}

func (gb *Gameboy) ButtonB(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Action &= ^uint8(1 << 1)
		} else {
			gb.Joypad.Action |= 1 << 1
		}
	})
}

func (gb *Gameboy) ButtonSelect(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Action &= ^uint8(1 << 2)
		} else {
			gb.Joypad.Action |= 1 << 2
		}
	})
}

func (gb *Gameboy) ButtonStart(pressed bool) {
	gb.CLK.Sync(func() {
		if pressed {
			gb.Joypad.Action &= ^uint8(1 << 3)
		} else {
			gb.Joypad.Action |= 1 << 3
		}
	})
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

	ppuClock := clk.Divide(2)
	cpuClock := clk.Divide(4)

	bootROMLock := NewBootROMLock(clk)
	bootROM := NewBootROM(clk, gb.Config.Model)
	vram := NewMemoryRegion(clk, AddrVRAMBegin, SizeVRAM)
	hram := NewMemoryRegion(clk, AddrHRAMBegin, SizeHRAM)
	wram := NewMemoryRegion(clk, AddrWRAMBegin, AddrWRAMEnd)
	apu := NewAPU(clk)
	oam := NewMemoryRegion(clk, AddrOAMBegin, SizeOAM)
	cartridgeSlot := NewMemoryRegion(clk, AddrCartridgeBank0Begin, AddrCartridgeBank0Size)
	joypad := NewJoypad(clk)

	bus := &Bus{}
	ppu := NewPPU(clk, ppuClock, bus, debugger)

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

	cpu := NewCPU(cpuClock, bus, debugger)

	gb.CLK = clk
	gb.PHI = cpuClock
	gb.CPU = cpu
	gb.CartridgeSlot = &cartridgeSlot
	gb.PPU = ppu
	gb.Debugger = debugger
	gb.Joypad = joypad

	clk.Onpanic = gb.CPU.Dump
}
