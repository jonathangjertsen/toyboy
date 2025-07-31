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
	AddressSpace *AddressSpace
	Program      []Data8
	APU          []Data8
	PPU          PPUDump
	Rewind       *Rewind
	Disassembly  *Disassembly
}

type PPUDump struct {
	Registers                 []Data8
	BGFIFO                    PixelFIFODump
	SpriteFIFO                PixelFIFODump
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
	fmt.Fprintf(f, "B  = 0x%02x\n", regs.BC.MSR)
	fmt.Fprintf(f, "C  = 0x%02x\n", regs.BC.LSR)
	fmt.Fprintf(f, "D  = 0x%02x\n", regs.DE.MSR)
	fmt.Fprintf(f, "E  = 0x%02x\n", regs.DE.LSR)
	fmt.Fprintf(f, "H  = 0x%02x\n", regs.HL.MSR)
	fmt.Fprintf(f, "L  = 0x%02x\n", regs.HL.LSR)
	fmt.Fprintf(f, "W  = 0x%02x\n", regs.Temp.MSR)
	fmt.Fprintf(f, "Z  = 0x%02x\n", regs.Temp.LSR)
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
	if cd.ProgramEnd >= min(Addr(len(cd.Program)), 0x8000) {
		return
	}
	MemDump(f, cd.Program[cd.ProgramStart:cd.ProgramEnd+1], cd.ProgramStart, cd.ProgramEnd, cd.Regs.PC-1)
}

func (cd *CoreDump) PrintHRAM(f io.Writer) {
	MemDump(f, cd.AddressSpace[:], 0, Addr(SizeHRAM), cd.Regs.SP)
}

func (cd *CoreDump) PrintOAM(f io.Writer) {
	MemDump(f, cd.AddressSpace[:], 0, Addr(SizeOAM), 0)
}

func (cd *CoreDump) PrintOAMAttrs(f io.Writer) {
	oam := cd.AddressSpace[AddrOAMBegin : AddrOAMEnd+1]
	for idx := range 40 {
		obj := DecodeObject(oam[idx*4 : (idx+1)*4])
		fmt.Fprintf(f, "%02d T=%03d X=%03d Y=%03d Attr=%x\n", idx, obj.TileIndex, obj.X, obj.Y, obj.Attributes.Hex())
	}
}

func (cd *CoreDump) PrintVRAM(f io.Writer) {
	MemDump(f, cd.AddressSpace[:], 0, Addr(SizeVRAM), 0)
}

func (cd *CoreDump) PrintWRAM(f io.Writer) {
	MemDump(f, cd.AddressSpace[:], 0, Addr(SizeWRAM), cd.Regs.SP)
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func PrintAPU(f io.Writer, apu *APU) {
	RegDump(f, apu.MemoryRegion.Data, AddrAPUBegin, AddrAPUEnd)
	fmt.Fprintf(f, "DivAPU:    %d\n", apu.DIVAPU)
	fmt.Fprintf(f, "Pulse1 on=%d dac=%d\n", b2i(apu.Pulse1.activated), b2i(apu.Pulse1.dacEnabled))
	printPeriodCounter(f, &apu.Pulse1.PeriodCounter)
	printLengthTimer(f, &apu.Pulse1.LengthTimer)
	printDutyGenerator(f, &apu.Pulse1.DutyGenerator)
	printEnvelope(f, &apu.Pulse1.Envelope)
	fmt.Fprintf(f, "Pulse2 on=%d dac=%d\n", b2i(apu.Pulse2.activated), b2i(apu.Pulse1.dacEnabled))
	printPeriodCounter(f, &apu.Pulse2.PeriodCounter)
	printLengthTimer(f, &apu.Pulse2.LengthTimer)
	printDutyGenerator(f, &apu.Pulse2.DutyGenerator)
	printEnvelope(f, &apu.Pulse2.Envelope)
	fmt.Fprintf(f, "Wave on=%d dac=%d\n", b2i(apu.Wave.activated), b2i(apu.Wave.dacEnabled))
	printPeriodCounter(f, &apu.Wave.PeriodCounter)
	printLengthTimer(f, &apu.Wave.LengthTimer)
	fmt.Fprintf(f, "Noise on=%d dac=%d\n", b2i(apu.Noise.activated), b2i(apu.Noise.dacEnabled))
	printPeriodCounter(f, &apu.Noise.PeriodCounter)
	printLengthTimer(f, &apu.Noise.LengthTimer)
	printEnvelope(f, &apu.Noise.Envelope)
	fmt.Fprintf(f, "                               ")
}

func printPeriodCounter(f io.Writer, pc *PeriodCounter) {
	fmt.Fprintf(f, "  PC RST=%d V=%d\n", pc.periodDividerReset, pc.periodDivider)
}

func printLengthTimer(f io.Writer, lt *LengthTimer) {
	fmt.Fprintf(f, "  LT EN=%d RST=%d V=%d\n", b2i(lt.lengthEnable), lt.lengthTimerReset, lt.lengthTimer)
}

func printDutyGenerator(f io.Writer, dg *DutyGenerator) {
	fmt.Fprintf(f, "  DG WF=%d V=%d\n", dg.waveform, dg.output)
}

func printEnvelope(f io.Writer, env *Envelope) {
	fmt.Fprintf(f, "  ENV SP=%d T=%d D=%d R=%d V=%d\n", env.envSweepPace, env.envTimer, b2i(env.envDir), env.volumeReset, env.volume)
}

func PrintPPU(f io.Writer, ppu PPUDump) {
	RegDump(f, ppu.Registers, AddrPPUBegin, AddrPPUEnd)
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
	fmt.Fprintf(f, "BFetch.Tile:   %s\n", join16(ppu.BackgroundFetcher.TileMSB, ppu.BackgroundFetcher.TileLSB).Hex())
	fmt.Fprintf(f, "BFetch.Susp:   %d\n", b2i(ppu.BackgroundFetcher.Suspended))
	fmt.Fprintf(f, "BFetch.WYRch:  %d\n", b2i(ppu.BackgroundFetcher.WindowYReached))
	fmt.Fprintf(f, "BFetch.WFetch: %d\n", b2i(ppu.BackgroundFetcher.WindowFetching))
	fmt.Fprintf(f, "BFetch.WLC:    %d\n", ppu.BackgroundFetcher.WindowLineCounter)
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "BGFIFO: /[")
	for i := range ppu.BGFIFO.Level {
		fmt.Fprintf(f, "%d", ppu.BGFIFO.Slots[i].ColorIDX())
	}
	for range 8 - ppu.BGFIFO.Level {
		fmt.Fprintf(f, " ")
	}
	fmt.Fprintf(f, "]\n")
	fmt.Fprintf(f, "       [%d]\n", ppu.LastShifted)
	fmt.Fprintf(f, " SFIFO: \\[")
	for i := range ppu.SpriteFIFO.Level {
		fmt.Fprintf(f, "%d", ppu.SpriteFIFO.Slots[i].ColorIDX())
	}
	for range 8 - ppu.SpriteFIFO.Level {
		fmt.Fprintf(f, " ")
	}
	fmt.Fprintf(f, "]\n")
	fmt.Fprintf(f, "SFetch.C:      %d\n", ppu.SpriteFetcher.Cycle)
	fmt.Fprintf(f, "SFetch.State:  %d\n", int(ppu.SpriteFetcher.State))
	fmt.Fprintf(f, "SFetch.X:      %d\n", ppu.SpriteFetcher.X)
	fmt.Fprintf(f, "SFetch.SIdx:   %d\n", ppu.SpriteFetcher.SpriteIDX)
	fmt.Fprintf(f, "SFetch.TIdx:   %d\n", ppu.SpriteFetcher.TileIndex)
	fmt.Fprintf(f, "SFetch.TAddr:  %s\n", ppu.SpriteFetcher.TileLSBAddr.Hex())
	fmt.Fprintf(f, "SFetch.Tile:   %s\n", join16(ppu.SpriteFetcher.TileMSB, ppu.SpriteFetcher.TileLSB).Hex())
	fmt.Fprintf(f, "SFetch.Susp:   %d\n", b2i(ppu.SpriteFetcher.Suspended))
	fmt.Fprintf(f, "Shift.Discard: %d\n", ppu.PixelShifter.Discard)
	fmt.Fprintf(f, "Shift.X:       %d\n", ppu.PixelShifter.X)
	fmt.Fprintf(f, "Shift.Susp:    %d\n", b2i(ppu.PixelShifter.Suspended))
	fmt.Fprintf(f, "OAMBuffer.LV:  %d\n", ppu.OAMBuffer.Level)
}

func (cd *CoreDump) PrintAPU(f io.Writer) {
	RegDump(f, cd.APU, AddrAPUBegin, AddrAPUEnd)
}

func RegDump(f io.Writer, mem []Data8, start, end Addr) {
	for addr := start; addr <= end; addr++ {
		a := Addr(addr)
		if !a.IsValid() {
			continue
		}
		fmt.Fprintf(f, "%5s = %02x\n", Addr(addr), mem[addr-start])
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
			fmt.Fprintf(f, "[%02x]", mem[addr-start])
		} else {
			fmt.Fprintf(f, " %02x ", mem[addr-start])
		}
	}

	alignedEnd := (end/0x10)*0x10 + 0x10 - 1
	for addr := end; addr < alignedEnd; addr++ {
		fmt.Fprintf(f, " .. ")
	}
	fmt.Fprintf(f, "\n")
}

func (cpu *CPU) GetCoreDump() CoreDump {
	bus, ok := cpu.Bus.(*Bus)
	if !ok {
		return CoreDump{}
	}

	end := cpu.Bus.BeginCoreDump()
	defer end()

	var cd CoreDump
	cd.AddressSpace = bus.AddressSpace
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
	cd.Program = bus.ProbeRange(0x0000, AddrCartridgeBank0End)
	cd.Disassembly = cpu.Debug.Disassembly(0, 0xffff)
	cd.APU = bus.APU.Data
	var ppu *PPU
	bus.GetPeripheral(&ppu)
	cd.PPU = ppu.GetDump()
	cd.Rewind = cpu.rewind
	return cd
}

func (cpu *CPU) Dump() {
	cd := cpu.GetCoreDump()
	f := os.Stdout
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
	fmt.Printf("Last executed instructions:\n")
	cd.PrintRewindBuffer(f)
	fmt.Fprintf(f, "--------\n")
}

func (cd *CoreDump) PrintRewindBuffer(f io.Writer) {
	for i := cd.Rewind.Start(); i != cd.Rewind.End(); i = cd.Rewind.Next(i) {
		entry := cd.Rewind.At(i)
		extra := ""
		if entry.BranchResult == +1 {
			extra = "(taken)"
		} else if entry.BranchResult == -1 {
			extra = "(not taken)"
		}
		fmt.Fprintf(f, "[PC=%s] %s %s\n", entry.Instruction.Address.Hex(), entry.Instruction.Asm(), extra)
	}
}

func (cd *CoreDump) PrintDisassembly(f io.Writer) {
	cd.Disassembly.Print(f)
}
