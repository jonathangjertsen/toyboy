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
	FrameSync  FrameSync
	Interrupts Interrupts
	Timer      Timer
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
	gb.Mem = NewAddressSpace()
	gb.Debug = Debug{
		Debugger:     NewDebugger(),
		Disassembler: NewDisassembler(&gb.Config.Debug.Disassembler),
		Warnings:     map[string]UserMessage{},
	}
	gb.Debug.Init()
	gb.FrameSync.ch = make(chan func(*ViewPort), 1)
	gb.CLK = ClockRT{
		resume:  make(chan struct{}),
		pause:   make(chan struct{}),
		stop:    make(chan struct{}),
		jobs:    make(chan func()),
		Onpanic: func(mem []Data8) {},
	}

	if gb.Config.BootROM.Variant == "DMGBoot" {
		bootROM := Data8Slice(assets.DMGBoot)
		copy(gb.Mem[:SizeBootROM], bootROM)
		gb.Debug.SetProgram(bootROM)
		gb.Debug.SetPC(0, &gb.CLK)
	} else {
		panic("unknown boot ROM")
	}
	gb.Debug.HRAM.Source = gb.Mem[AddrHRAMBegin : AddrHRAMEnd+1]
	gb.Debug.WRAM.Source = gb.Mem[AddrWRAMBegin : AddrWRAMEnd+1]

	gb.Cartridge = Cartridge{
		BankNo1:         1,
		SelectedROMBank: 1,
	}
	bootROMLock := NewBootROMLock(gb.Mem, &gb.Cartridge, &gb.Debug)
	gb.Joypad.Action = 0xf
	gb.Joypad.Direction = 0xf
	gb.Mem[AddrP1] = 0x1f

	gb.CPU = CPU{
		Config:     gb.Config,
		Bus:        &gb.Bus,
		Debug:      &gb.Debug,
		Interrupts: &gb.Interrupts,
		rewind:     NewRewind(8192),
	}
	gb.CPU.handlers = handlers(&gb.CPU)

	gb.PPU.SpriteFetcher.Suspended = true
	gb.PPU.SpriteFetcher.DoneX = 0xff
	gb.PPU.beginFrame(gb.Mem, &gb.Interrupts)

	gb.Bus.BootROMLock = bootROMLock
	gb.Bus.APU = &gb.APU
	gb.Bus.PPU = &gb.PPU
	gb.Bus.Cartridge = &gb.Cartridge
	gb.Bus.Joypad = &gb.Joypad
	gb.Bus.Interrupts = &gb.Interrupts
	gb.Bus.Timer = &gb.Timer
	gb.Bus.Config = gb.Config

	gb.CLK.Onpanic = gb.CPU.Dump
}
