package model

import (
	"fmt"
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
	TCycle      uint
	MCycle      uint
	PureRAM     bool
}

func (gb *Gameboy) Print(format string, args ...any) {
	fmt.Printf("M=%8d+%d | PC=%s | ", gb.MCycle, gb.TCycle-(gb.MCycle<<2), gb.CPU.Regs.PC.Hex())
	fmt.Printf(format, args...)
	fmt.Printf("\n")
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

func (gb *Gameboy) AllocMem() {
	gb.Mem = make([]Data8, 65536)
	gb.Cartridge.ROM = make([][ROMBankSize]Data8, 512)
	gb.Cartridge.RAM = make([][RAMBankSize]Data8, 4)
}

func (gb *Gameboy) Init(config *Config, clk *ClockRT) {
	save.MustTriviallySerialize(gb)

	gb.initMemory()
	gb.initDebug(config)
	gb.loadBootROM(config)
	gb.initCartridge()
	gb.initJoypad()
	gb.initCPU(config, clk)
	gb.initPPU()
}

func (gb *Gameboy) initMemory() {
	clear(gb.Mem)
}

func (gb *Gameboy) loadBootROM(config *Config) {
	if config.BootROM.Variant == "DMGBoot" {
		bootROM := Data8Slice(assets.DMGBoot)
		copy(gb.Mem[:SizeBootROM], bootROM)
		gb.Debug.SetProgram(bootROM)
	} else if config.BootROM.Variant == "None" {
	} else {
		panicf("unknown bootROM variant '%s'", config.BootROM.Variant)
	}
}
func (gb *Gameboy) initDebug(config *Config) {
	gb.Debug = Debug{
		Debugger:     NewDebugger(),
		Disassembler: NewDisassembler(&config.Debug.Disassembler),
		Warnings:     map[string]UserMessage{},
	}
	gb.Debug.HRAM.Source = gb.Mem[AddrHRAMBegin : AddrHRAMEnd+1]
	gb.Debug.WRAM.Source = gb.Mem[AddrWRAMBegin : AddrWRAMEnd+1]
}

func (gb *Gameboy) initCartridge() {
	gb.Cartridge.RegLow = 1
	gb.Cartridge.SelectedROMBank0 = 1
}

func (gb *Gameboy) initJoypad() {
	gb.Joypad.Action = 0xf
	gb.Joypad.Direction = 0xf
	gb.Mem[AddrP1] = 0x1f
}

func (gb *Gameboy) initCPU(config *Config, clk *ClockRT) {
	gb.CPU = CPU{
		Rewind: NewRewind(config.Debug.RewindSize),
	}
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.instructionFetch(clk)
	gb.CPU.IncPC()
	gb.CPU.UOpCycle++
	clk.Onpanic = gb.CPU.Dump
}

func (gb *Gameboy) initPPU() {
	gb.PPU.SpriteFetcher.Suspended = true
	gb.PPU.SpriteFetcher.DoneX = 0xff
	gb.beginFrame()
}

func (gb *Gameboy) WriteAddress(addr Addr) {
	gb.Address = addr
	gb.Data = gb.ProbeAddress(addr)
}

const LowestSpecialAddress = AddrP1

func (gb *Gameboy) ProbeAddress(addr Addr) Data8 {
	if addr < LowestSpecialAddress || gb.PureRAM {
		return gb.Mem[addr]
	}

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
	if gb.PureRAM {
		gb.Mem[gb.Address] = v
		return
	}

	gb.Data = v
	addr := gb.Address

	if addr <= AddrBootROMEnd {
		if gb.BootROMLock.BootOff {
			gb.WriteCartridge(addr, v)
		}
		return
	}
	if addr <= AddrCartridgeBankNEnd {
		gb.WriteCartridge(addr, v)
		return
	}
	gb.Mem[addr] = v

	if addr == AddrBootROMLock {
		gb.WriteBootROMLock(v)
	} else if addr == AddrP1 {
		gb.WriteJoypad(addr, v)
	} else if addr == AddrIF || addr == AddrIE {
		gb.IRQCheck()
	} else if addr >= AddrAPUBegin && addr <= AddrAPUEnd {
		gb.APU.Write(addr, v)
	} else if addr >= AddrPPUBegin && addr <= AddrPPUEnd {
		gb.WritePPU(addr, v)
	} else if addr >= AddrTimerBegin && addr <= AddrTimerEnd {
		gb.Timer.Write(addr, v)
	}
}
