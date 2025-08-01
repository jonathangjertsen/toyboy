package model

import (
	"sync/atomic"

	"github.com/jonathangjertsen/toyboy/assets"
)

type Gameboy struct {
	Config *Config

	Running atomic.Bool

	Mem        []Data8
	CLK        ClockRT
	Bus        Bus
	Debug      *Debug
	CPU        *CPU
	PPU        *PPU
	APU        *APU
	Cartridge  *Cartridge
	Joypad     *Joypad
	FrameSync  *FrameSync
	Interrupts *Interrupts
	Audio      *Audio
	Timer      *Timer
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
	gb.CLK = ClockRT{
		resume:  make(chan struct{}),
		pause:   make(chan struct{}),
		stop:    make(chan struct{}),
		jobs:    make(chan func()),
		Onpanic: func(mem []Data8) {},
	}

	if gb.Config.BootROM.Variant == "DMGBoot" {
		bootROM := Data8Slice(assets.DMGBoot)
		copy(mem[:SizeBootROM], bootROM)
		debug.SetProgram(bootROM)
		debug.SetPC(0, &gb.CLK)
	} else {
		panic("unknown boot ROM")
	}
	cartridge := NewCartridge(&gb.CLK, mem)
	bootROMLock := NewBootROMLock(mem, cartridge, debug)
	var apu APU
	joypad := NewJoypad(interrupts, mem)
	var timer Timer

	cpu := NewCPU(&gb.CLK, interrupts, &gb.Bus, gb.Config, debug)

	ppu := NewPPU(interrupts)

	gb.Bus.BootROMLock = bootROMLock
	gb.Bus.APU = &apu
	gb.Bus.PPU = ppu
	gb.Bus.Cartridge = cartridge
	gb.Bus.Joypad = joypad
	gb.Bus.Interrupts = interrupts
	gb.Bus.Timer = &timer
	gb.Bus.Config = gb.Config

	debug.HRAM.Source = mem[AddrHRAMBegin : AddrHRAMEnd+1]
	debug.WRAM.Source = mem[AddrWRAMBegin : AddrWRAMEnd+1]

	gb.Mem = mem
	gb.CPU = cpu
	gb.APU = &apu
	gb.Cartridge = cartridge
	gb.PPU = ppu
	gb.Debug = debug
	gb.Joypad = joypad
	gb.FrameSync = fs
	gb.Interrupts = interrupts
	gb.Audio = audio
	gb.Timer = &timer

	go gb.CLK.run(gb)

	gb.CPU.Reset()

	gb.CLK.Onpanic = gb.CPU.Dump
}
