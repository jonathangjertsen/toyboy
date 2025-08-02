package model

import (
	"fmt"
	"io"
	"os"
)

type CoreDump struct {
	Cycle        uint
	Regs         RegisterFile
	ProgramStart Addr
	ProgramEnd   Addr
	Mem          []Data8
	PPU          PPUDump
	Rewind       Rewind
	Disassembly  *Disassembly
}

type PPUDump struct {
	LastShifted               Color
	OAMScanCycle              uint64
	PixelDrawCycle            uint64
	HBlankRemainingCycles     uint64
	VBlankLineRemainingCycles uint64
	PixelShifter              Shifter
	BackgroundFetcher         BackgroundFetcher
	SpriteFetcher             SpriteFetcher
	OAMBuffer                 OAMBuffer
}

func PrintRegs(f io.Writer, regs RegisterFile) {
	fmt.Fprintf(f, "PC = %s\n", regs.PC.Hex())
	fmt.Fprintf(f, "SP = %s\n", regs.SP.Hex())
	fmt.Fprintf(f, "A  = 0x%02x\n", regs.A)
	fmt.Fprintf(f, "F  = 0x%02x\n", regs.F)
	fmt.Fprintf(f, "B  = 0x%02x\n", regs.B)
	fmt.Fprintf(f, "C  = 0x%02x\n", regs.C)
	fmt.Fprintf(f, "D  = 0x%02x\n", regs.D)
	fmt.Fprintf(f, "E  = 0x%02x\n", regs.E)
	fmt.Fprintf(f, "H  = 0x%02x\n", regs.H)
	fmt.Fprintf(f, "L  = 0x%02x\n", regs.L)
	fmt.Fprintf(f, "W  = 0x%02x\n", regs.TempW)
	fmt.Fprintf(f, "Z  = 0x%02x\n", regs.TempZ)
	fmt.Fprintf(f, "IR = 0x%02x\n", uint8(regs.IR))
	z, h, n, c := 0, 0, 0, 0
	if regs.GetFlagZ() {
		z = 1
	}
	if regs.GetFlagH() {
		h = 1
	}
	if regs.GetFlagN() {
		n = 1
	}
	if regs.GetFlagC() {
		c = 1
	}
	fmt.Fprintf(f, "Z=%v C=%v\n", z, c)
	fmt.Fprintf(f, "N=%v H=%v\n", n, h)
}

func (cd *CoreDump) PrintProgram(f io.Writer) {
	if cd.ProgramEnd >= 0x8000 {
		return
	}
	MemDump(f, cd.Mem, cd.ProgramStart, cd.ProgramEnd, cd.Regs.PC-1)
}

func (cd *CoreDump) PrintHRAM(f io.Writer) {
	MemDump(f, cd.Mem, AddrHRAMBegin, AddrHRAMEnd, cd.Regs.SP)
}

func (cd *CoreDump) PrintOAM(f io.Writer) {
	MemDump(f, cd.Mem, AddrOAMBegin, AddrOAMEnd, 0)
}

func (cd *CoreDump) PrintOAMAttrs(f io.Writer) {
	oam := cd.Mem[AddrOAMBegin : AddrOAMEnd+1]
	for idx := range 40 {
		obj := DecodeObject(oam[idx*4 : (idx+1)*4])
		fmt.Fprintf(f, "%02d T=%03d X=%03d Y=%03d Attr=%x\n", idx, obj.TileIndex, obj.X, obj.Y, obj.Attributes.Hex())
	}
}

func (cd *CoreDump) PrintVRAM(f io.Writer) {
	MemDump(f, cd.Mem, AddrVRAMBegin, AddrVRAMEnd, 0)
}

func (cd *CoreDump) PrintWRAM(f io.Writer) {
	MemDump(f, cd.Mem, AddrWRAMBegin, AddrWRAMEnd, cd.Regs.SP)
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func PrintAPU(f io.Writer, mem []Data8, apu *APU) {
	RegDump(f, mem, AddrAPUBegin, AddrAPUEnd)
	fmt.Fprintf(f, "DivAPU:    %d\n", apu.DIVAPU)
	fmt.Fprintf(f, "Pulse1 on=%d dac=%d\n", b2i(apu.Pulse1.Activated), b2i(apu.Pulse1.DacEnabled))
	printPeriodCounter(f, &apu.Pulse1.PeriodCounter)
	printLengthTimer(f, &apu.Pulse1.LengthTimer)
	printDutyGenerator(f, &apu.Pulse1.PulseChannel)
	printEnvelope(f, &apu.Pulse1.Envelope)
	fmt.Fprintf(f, "Pulse2 on=%d dac=%d\n", b2i(apu.Pulse2.Activated), b2i(apu.Pulse1.DacEnabled))
	printPeriodCounter(f, &apu.Pulse2.PeriodCounter)
	printLengthTimer(f, &apu.Pulse2.LengthTimer)
	printDutyGenerator(f, &apu.Pulse2)
	printEnvelope(f, &apu.Pulse2.Envelope)
	fmt.Fprintf(f, "Wave on=%d dac=%d\n", b2i(apu.Wave.Activated), b2i(apu.Wave.DacEnabled))
	printPeriodCounter(f, &apu.Wave.PeriodCounter)
	printLengthTimer(f, &apu.Wave.LengthTimer)
	fmt.Fprintf(f, "Noise on=%d dac=%d\n", b2i(apu.Noise.Activated), b2i(apu.Noise.DacEnabled))
	printPeriodCounter(f, &apu.Noise.PeriodCounter)
	printLengthTimer(f, &apu.Noise.LengthTimer)
	printEnvelope(f, &apu.Noise.Envelope)
	fmt.Fprintf(f, "                               ")
}

func printPeriodCounter(f io.Writer, pc *PeriodCounter) {
	fmt.Fprintf(f, "  PC RST=%d V=%d\n", pc.Reset, pc.Counter)
}

func printLengthTimer(f io.Writer, lt *LengthTimer) {
	fmt.Fprintf(f, "  LT EN=%d RST=%d V=%d\n", b2i(lt.Enable), lt.Reset, lt.Counter)
}

func printDutyGenerator(f io.Writer, pc *PulseChannel) {
	fmt.Fprintf(f, "  DG WF=%d V=%d\n", pc.Waveform, pc.Output)
}

func printEnvelope(f io.Writer, env *Envelope) {
	fmt.Fprintf(f, "  ENV SP=%d T=%d D=%d R=%d V=%d\n", env.EnvSweepPace, env.EnvTimer, b2i(env.EnvDir), env.VolumeReset, env.Volume)
}

func PrintPPU(f io.Writer, ppu PPUDump, mem []Data8) {
	RegDump(f, mem, AddrPPUBegin, AddrPPUEnd)
	fmt.Fprintf(f, "OAMCt:         %d\n", ppu.OAMScanCycle)
	fmt.Fprintf(f, " PDCt:         %d\n", ppu.PixelDrawCycle)
	fmt.Fprintf(f, " HBCt:         %d\n", ppu.HBlankRemainingCycles)
	fmt.Fprintf(f, " VBCt:         %d\n", ppu.VBlankLineRemainingCycles)
	fmt.Fprintf(f, "BFetch.C:      %d\n", ppu.BackgroundFetcher.Cycle)
	fmt.Fprintf(f, "BFetch.State:  %d\n", int(ppu.BackgroundFetcher.State))
	fmt.Fprintf(f, "BFetch.X:      %d\n", ppu.BackgroundFetcher.X)
	fmt.Fprintf(f, "BFetch.TOffX:  %d\n", ppu.BackgroundFetcher.TileOffsetX)
	fmt.Fprintf(f, "BFetch.TOffY:  %d\n", ppu.BackgroundFetcher.TileOffsetY)
	fmt.Fprintf(f, "BFetch.&TIdx:  %s\n", ppu.BackgroundFetcher.TileIndexAddr.Hex())
	fmt.Fprintf(f, "BFetch.TIdx:   %d\n", ppu.BackgroundFetcher.TileIndex)
	fmt.Fprintf(f, "BFetch.TAddr:  %s\n", ppu.BackgroundFetcher.TileLSBAddr.Hex())
	fmt.Fprintf(f, "BFetch.Susp:   %d\n", b2i(ppu.BackgroundFetcher.Suspended))
	fmt.Fprintf(f, "BFetch.WYRch:  %d\n", b2i(ppu.BackgroundFetcher.WindowYReached))
	fmt.Fprintf(f, "BFetch.WFetch: %d\n", b2i(ppu.BackgroundFetcher.WindowFetching))
	fmt.Fprintf(f, "BFetch.WLC:    %d\n", ppu.BackgroundFetcher.WindowLineCounter)
	fmt.Fprintf(f, "SFetch.C:      %d\n", ppu.SpriteFetcher.Cycle)
	fmt.Fprintf(f, "SFetch.State:  %d\n", int(ppu.SpriteFetcher.State))
	fmt.Fprintf(f, "SFetch.X:      %d\n", ppu.SpriteFetcher.X)
	fmt.Fprintf(f, "SFetch.SIdx:   %d\n", ppu.SpriteFetcher.SpriteIDX)
	fmt.Fprintf(f, "SFetch.TIdx:   %d\n", ppu.SpriteFetcher.TileIndex)
	fmt.Fprintf(f, "SFetch.TAddr:  %s\n", ppu.SpriteFetcher.TileLSBAddr.Hex())
	fmt.Fprintf(f, "SFetch.Susp:   %d\n", b2i(ppu.SpriteFetcher.Suspended))
	fmt.Fprintf(f, "Shift.Discard: %d\n", ppu.PixelShifter.Discard)
	fmt.Fprintf(f, "Shift.X:       %d\n", ppu.PixelShifter.X)
	fmt.Fprintf(f, "Shift.Susp:    %d\n", b2i(ppu.PixelShifter.Suspended))
	fmt.Fprintf(f, "OAMBuffer.LV:  %d\n", ppu.OAMBuffer.Level)
}

func RegDump(f io.Writer, mem []Data8, start, end Addr) {
	for addr := start; addr <= end; addr++ {
		a := Addr(addr)
		if !a.IsValid() {
			continue
		}
		fmt.Fprintf(f, "%5s = %02x\n", Addr(addr), mem[addr])
	}
}

func MemDump(f io.Writer, mem []Data8, start, end, highlight Addr) {
	alignedStart := (start / 0x10) * 0x10
	for addr := alignedStart; addr < start; addr++ {
		if addr%0x10 == 0 {
			fmt.Fprintf(f, "\n %s |", addr.Hex())
		}
		fmt.Fprintf(f, " .. ")
	}

	for addr := start; addr <= end; addr++ {
		if addr%0x10 == 0 {
			fmt.Fprintf(f, "\n%s |", addr.Hex())
		}
		if highlight == addr {
			fmt.Fprintf(f, "[%02x]", mem[addr])
		} else {
			fmt.Fprintf(f, " %02x ", mem[addr])
		}
	}

	alignedEnd := (end/0x10)*0x10 + 0x10 - 1
	for addr := end; addr < alignedEnd; addr++ {
		fmt.Fprintf(f, " .. ")
	}
	fmt.Fprintf(f, "\n")
}

func (cpu *CPU) GetCoreDump(gb *Gameboy) CoreDump {
	var cd CoreDump
	cd.Mem = gb.Mem
	cd.Regs = cpu.Regs
	cd.ProgramStart = 0
	if cpu.Regs.PC > 0x40 {
		cd.ProgramStart = cpu.Regs.PC - 0x40
	}
	cd.ProgramStart = (cd.ProgramStart / 0x10) * 0x10

	cd.ProgramEnd = 0xffff
	if cpu.Regs.PC < 0xffff-0x40 {
		cd.ProgramEnd = cpu.Regs.PC + 0x40
	}
	cd.ProgramEnd = (cd.ProgramEnd/0x10)*0x10 + 0x10 - 1
	cd.Disassembly = gb.Debug.Disassembly(0, 0xffff)
	cd.Rewind = cpu.Rewind
	return cd
}

func (cpu *CPU) Dump(gb *Gameboy) {
	f := os.Stdout
	/*
		fmt.Fprintf(f, "\n--------\nCore dump:\n")
		PrintRegs(f, cd.Regs)
		fmt.Fprintf(f, "--------\n")
		fmt.Fprintf(f, "Code (PC highlighted)\n")
		cd.PrintProgram(f)
		fmt.Fprintf(f, "--------\n")
		fmt.Fprintf(f, "HRAM (SP highlighted):\n")
		cd.PrintHRAM(f)
		fmt.Fprintf(f, "--------\n")
		fmt.Fprintf(f, "OAM:\n")
		cd.PrintOAM(f)
		fmt.Fprintf(f, "--------\n")
	*/
	fmt.Printf("Last executed instructions:\n")
	gb.PrintRewindBuffer(f, false)
	fmt.Fprintf(f, "--------\n")
}

func (gb *Gameboy) PrintRewindBuffer(f io.Writer, reverse bool) {
	rw := &gb.CPU.Rewind
	curr := rw.Curr()
	currTxt := fmt.Sprintf("Current: [PC=%s] %s (%d)\n", curr.Instruction.Address.Hex(), curr.Instruction.Asm(), gb.CPU.UOpCycle)

	if reverse {
		fmt.Fprint(f, currTxt)
		for i := rw.Prev(rw.End()); i != rw.Prev(rw.Start()); i = rw.Prev(i) {
			gb.printRewindEntry(f, i)
		}
	} else {
		for i := rw.Start(); i != rw.End(); i = rw.Next(i) {
			gb.printRewindEntry(f, i)
		}
		fmt.Fprint(f, currTxt)
	}
	fmt.Fprintf(f, "                                \n")
}

func (gb *Gameboy) printRewindEntry(f io.Writer, i int) {
	entry := gb.CPU.Rewind.At(i)
	extra := ""
	if entry.BranchResult == +1 {
		extra += " (taken)"
	} else if entry.BranchResult == -1 {
		extra += " (not taken)"
	}
	if entry.Instruction.InISR {
		extra += " (in ISR)"
	}
	if entry.Instruction.NopCount > 0 {
		extra += fmt.Sprintf(" (x%d)", entry.Instruction.NopCount+1)
	}
	fmt.Fprintf(f, "[PC=%s] %s%s\n", entry.Instruction.Address.Hex(), entry.Instruction.Asm(), extra)
}

func (cd *CoreDump) PrintDisassembly(f io.Writer) {
	cd.Disassembly.Print(f)
}
