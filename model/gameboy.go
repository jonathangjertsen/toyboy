package model

import (
	"sync/atomic"

	"github.com/jonathangjertsen/toyboy/assets"
)

type Gameboy struct {
	Mem         []Data8
	Bus         Bus
	Debug       Debug
	CPU         CPU
	PPU         PPU
	APU         APU
	Cartridge   Cartridge
	Joypad      Joypad
	FrameSync   FrameSync
	Interrupts  Interrupts
	Timer       Timer
	BootROMLock BootROMLock
}

func Start(clk *ClockRT, runFlag *atomic.Bool) {
	clk.Start()
	if runFlag != nil {
		runFlag.Store(true)
	}
}

func Pause(clk *ClockRT, runFlag *atomic.Bool) {
	clk.Pause()
	if runFlag != nil {
		runFlag.Store(false)
	}
}

func NewGameboy(config *Config, clk *ClockRT) *Gameboy {
	var gb Gameboy
	gb.Mem = NewAddressSpace()
	gb.Debug = Debug{
		Debugger:     NewDebugger(),
		Disassembler: NewDisassembler(&config.Debug.Disassembler),
		Warnings:     map[string]UserMessage{},
	}
	gb.Debug.Init()
	gb.FrameSync.ch = make(chan func(*ViewPort), 1)

	if config.BootROM.Variant == "DMGBoot" {
		bootROM := Data8Slice(assets.DMGBoot)
		copy(gb.Mem[:SizeBootROM], bootROM)
		gb.Debug.SetProgram(bootROM)
		gb.Debug.SetPC(0, clk)
	} else {
		panic("unknown boot ROM")
	}
	gb.Debug.HRAM.Source = gb.Mem[AddrHRAMBegin : AddrHRAMEnd+1]
	gb.Debug.WRAM.Source = gb.Mem[AddrWRAMBegin : AddrWRAMEnd+1]

	gb.Cartridge = Cartridge{
		BankNo1:         1,
		SelectedROMBank: 1,
	}
	gb.Joypad.Action = 0xf
	gb.Joypad.Direction = 0xf
	gb.Mem[AddrP1] = 0x1f

	gb.CPU = CPU{
		GB:     &gb,
		rewind: NewRewind(8192),
	}
	gb.CPU.handlers = handlers(&gb.CPU)

	gb.PPU.SpriteFetcher.Suspended = true
	gb.PPU.SpriteFetcher.DoneX = 0xff
	gb.PPU.beginFrame(gb.Mem, &gb.Interrupts)

	gb.Bus.GB = &gb
	gb.Bus.BootROMLock = &gb.BootROMLock
	gb.Bus.APU = &gb.APU
	gb.Bus.PPU = &gb.PPU
	gb.Bus.Cartridge = &gb.Cartridge
	gb.Bus.Joypad = &gb.Joypad
	gb.Bus.Interrupts = &gb.Interrupts
	gb.Bus.Timer = &gb.Timer

	clk.Onpanic = gb.CPU.Dump

	//mustTriviallySerialize(&gb)

	return &gb
}
