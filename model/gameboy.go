package model

import (
	"sync/atomic"

	"github.com/jonathangjertsen/toyboy/assets"
	"github.com/jonathangjertsen/toyboy/save"
)

type Gameboy struct {
	Mem         []Data8
	Address     Addr
	Data        Data8
	Debug       Debug
	CPU         CPU
	PPU         PPU
	APU         APU
	Cartridge   Cartridge
	Joypad      Joypad
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
		Rewind: NewRewind(8192),
	}

	gb.PPU.SpriteFetcher.Suspended = true
	gb.PPU.SpriteFetcher.DoneX = 0xff
	gb.PPU.beginFrame(gb.Mem, &gb.Interrupts)

	clk.Onpanic = gb.CPU.Dump

	save.MustTriviallySerialize(&gb)

	return &gb
}

func (gb *Gameboy) WriteAddress(addr Addr) {
	gb.Address = addr
	gb.Data = gb.ProbeAddress(addr)
}

func (gb *Gameboy) ProbeAddress(addr Addr) Data8 {
	if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		return gb.APU.Read(addr)
	}
	if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		return gb.PPU.Read(addr)
	}
	if addr == AddrP1 {
		return gb.Joypad.Read(gb.Mem[AddrP1], addr)
	}
	return gb.Mem[addr]
}

func (gb *Gameboy) WriteData(v Data8) {
	gb.Data = v
	addr := gb.Address

	if addr <= AddrBootROMEnd {
		if gb.BootROMLock.BootOff {
			gb.Cartridge.Write(gb.Mem, addr, v)
		}
		return
	}
	if addr <= AddrCartridgeBankNEnd {
		gb.Cartridge.Write(gb.Mem, addr, v)
		return
	}
	gb.Mem[addr] = v

	if addr == AddrBootROMLock {
		gb.BootROMLock.Write(gb.Mem, &gb.Debug, &gb.Cartridge, v)
	} else if addr == AddrP1 {
		gb.Joypad.Write(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		gb.Interrupts.IRQCheck(gb.Mem)
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		gb.APU.Write(addr, v)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		gb.PPU.Write(gb.Mem, addr, v, &gb.Interrupts)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		gb.Timer.Write(addr, v)
	}
}
