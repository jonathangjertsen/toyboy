package model

import (
	"sync/atomic"

	"github.com/jonathangjertsen/toyboy/assets"
)

type Gameboy struct {
	Config *Config

	Running atomic.Bool

	CLK       *ClockRT
	Bus       *Bus
	Debug     *Debug
	CPU       *CPU
	PPU       *PPU
	APU       *APU
	Cartridge *Cartridge
	Joypad    *Joypad
}

func (gb *Gameboy) Start() {
	gb.CLK.Start()
	gb.Running.Store(true)
}

func (gb *Gameboy) Pause() {
	gb.CLK.Pause()
	gb.Running.Store(false)
}

func (gb *Gameboy) Step() {
	gb.CLK.pauseAfterCycle.Add(1)
	gb.CLK.Start()
}

func (gb *Gameboy) SoftReset() {
	gb.CLK.Sync(func() {
		gb.CLK.Cycle = 0
		gb.CPU.Reset()
		gb.PPU.Reset()
		gb.APU.Reset()
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
	audio *Audio,
) *Gameboy {
	gameboy := &Gameboy{
		Config: config,
	}
	gameboy.Init(audio)
	return gameboy
}

func (gb *Gameboy) Init(audio *Audio) {
	clk := NewRealtimeClock(gb.Config.Clock, audio)

	debug := NewDebug(clk, &gb.Config.Debug)

	mem := NewAddressSpace()
	interrupts := NewInterrupts(clk, mem)

	if gb.Config.BootROM.Variant == "DMGBoot" {
		bootROM := Data8Slice(assets.DMGBoot)
		copy(mem[:SizeBootROM], bootROM)
		debug.SetProgram(bootROM)
		debug.SetPC(0)
	} else {
		panic("unknown boot ROM")
	}
	cartridge := NewCartridge(clk, mem)
	bootROMLock := NewBootROMLock(mem, cartridge, debug)
	apu := NewAPU(clk, gb.Config, mem)
	joypad := NewJoypad(clk, interrupts, mem)
	serial := NewSerial(clk)
	timer := NewTimer(clk, mem, apu, interrupts)

	bus := NewBus(mem)

	cpu := NewCPU(clk, interrupts, bus, gb.Config, debug, OpcodeNop)

	ppu := NewPPU(clk, interrupts, bus, gb.Config, debug)

	bus.BootROMLock = bootROMLock
	bus.APU = apu
	bus.PPU = ppu
	bus.Cartridge = cartridge
	bus.Joypad = joypad
	bus.Interrupts = interrupts
	bus.Serial = serial
	bus.Timer = timer
	bus.Config = gb.Config

	debug.HRAM.Source = bus.Mem[AddrHRAMBegin : AddrHRAMEnd+1]
	debug.WRAM.Source = bus.Mem[AddrWRAMBegin : AddrWRAMEnd+1]

	gb.Bus = bus
	gb.CLK = clk
	gb.CPU = cpu
	gb.APU = apu
	gb.Cartridge = cartridge
	gb.PPU = ppu
	gb.Debug = debug
	gb.Joypad = joypad

	gb.CPU.Reset()

	clk.Onpanic = gb.CPU.Dump
}
