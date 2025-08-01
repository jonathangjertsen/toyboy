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
	Debug      Debug
	CPU        CPU
	PPU        PPU
	APU        APU
	Cartridge  Cartridge
	Joypad     Joypad
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
	gb.Debug = Debug{
		Debugger:     NewDebugger(),
		Disassembler: NewDisassembler(&gb.Config.Debug.Disassembler),
		Warnings:     map[string]UserMessage{},
	}
	gb.Debug.Init()
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
		gb.Debug.SetProgram(bootROM)
		gb.Debug.SetPC(0, &gb.CLK)
	} else {
		panic("unknown boot ROM")
	}
	gb.Debug.HRAM.Source = mem[AddrHRAMBegin : AddrHRAMEnd+1]
	gb.Debug.WRAM.Source = mem[AddrWRAMBegin : AddrWRAMEnd+1]

	gb.Cartridge = Cartridge{
		mem:             mem,
		BankNo1:         1,
		SelectedROMBank: 1,
	}
	bootROMLock := NewBootROMLock(mem, &gb.Cartridge, &gb.Debug)
	gb.Joypad.Action = 0xf
	gb.Joypad.Direction = 0xf
	mem[AddrP1] = 0x1f
	var timer Timer

	gb.CPU = CPU{
		Config:     gb.Config,
		Bus:        &gb.Bus,
		Debug:      &gb.Debug,
		Interrupts: interrupts,
		rewind:     NewRewind(8192),
	}
	gb.CPU.handlers = handlers(&gb.CPU)

	gb.PPU.SpriteFetcher.Suspended = true
	gb.PPU.SpriteFetcher.DoneX = 0xff
	gb.PPU.beginFrame(gb.Interrupts)

	gb.Bus.BootROMLock = bootROMLock
	gb.Bus.APU = &gb.APU
	gb.Bus.PPU = &gb.PPU
	gb.Bus.Cartridge = &gb.Cartridge
	gb.Bus.Joypad = &gb.Joypad
	gb.Bus.Interrupts = interrupts
	gb.Bus.Timer = &timer
	gb.Bus.Config = gb.Config

	gb.Mem = mem
	gb.FrameSync = fs
	gb.Interrupts = interrupts
	gb.Audio = audio
	gb.Timer = &timer

	go gb.CLK.run(gb)

	gb.CPU.Reset()

	gb.CLK.Onpanic = gb.CPU.Dump
}
