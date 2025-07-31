package model

import (
	"sync/atomic"

	"github.com/jonathangjertsen/toyboy/assets"
)

type Gameboy struct {
	Config *Config

	Running atomic.Bool

	Mem        []Data8
	CLK        *ClockRT
	Bus        *Bus
	Debug      *Debug
	CPU        *CPU
	PPU        *PPU
	APU        *APU
	Cartridge  *Cartridge
	Joypad     *Joypad
	FrameSync  *FrameSync
	Interrupts *Interrupts
	Audio      *Audio
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
	mem := NewAddressSpace()
	interrupts := NewInterrupts(mem)
	debug := NewDebug(&gb.Config.Debug)
	fs := NewFrameSync()
	clk := NewRealtimeClock()

	if gb.Config.BootROM.Variant == "DMGBoot" {
		bootROM := Data8Slice(assets.DMGBoot)
		copy(mem[:SizeBootROM], bootROM)
		debug.SetProgram(bootROM)
		debug.SetPC(0, clk)
	} else {
		panic("unknown boot ROM")
	}
	cartridge := NewCartridge(clk, mem)
	bootROMLock := NewBootROMLock(mem, cartridge, debug)
	var apu APU
	joypad := NewJoypad(interrupts, mem)
	var timer Timer

	bus := NewBus()

	cpu := NewCPU(clk, interrupts, bus, gb.Config, debug)

	ppu := NewPPU(interrupts)

	bus.BootROMLock = bootROMLock
	bus.APU = &apu
	bus.PPU = ppu
	bus.Cartridge = cartridge
	bus.Joypad = joypad
	bus.Interrupts = interrupts
	bus.Timer = &timer
	bus.Config = gb.Config

	debug.HRAM.Source = mem[AddrHRAMBegin : AddrHRAMEnd+1]
	debug.WRAM.Source = mem[AddrWRAMBegin : AddrWRAMEnd+1]

	gb.Mem = mem
	gb.Bus = bus
	gb.CLK = clk
	gb.CPU = cpu
	gb.APU = &apu
	gb.Cartridge = cartridge
	gb.PPU = ppu
	gb.Debug = debug
	gb.Joypad = joypad
	gb.FrameSync = fs
	gb.Interrupts = interrupts
	gb.Audio = audio

	go clk.run(gb.Config.Clock.SpeedPercent, interrupts, debug, mem, fs, audio, &apu, ppu, cpu, &timer)

	gb.CPU.Reset()

	clk.Onpanic = gb.CPU.Dump
}
