package model

import (
	"fmt"
	"io"
	"os"
)

type CoreDump struct {
	Cycle           Cycle
	Regs            RegisterFile
	ProgramStart    Addr
	ProgramEnd      Addr
	Program         MemDump
	HRAM            MemDump
	OAM             MemDump
	VRAM            MemDump
	APU             MemDump
	PPU             PPUDump
	RewindBuffer    []ExecLogEntry
	RewindBufferIdx int
	Disassembly     *Disassembly
}

type PPUDump struct {
	Registers                 MemDump
	BGFIFO                    PixelFIFODump
	SpriteFIFO                PixelFIFODump
	LastShifted               Color
	OAMScanCycle              uint64
	PixelDrawCycle            uint64
	HBlankRemainingCycles     uint64
	VBlankLineRemainingCycles uint64
	PixelShifter              PixelShifter
	BackgroundFetcher         BackgroundFetcher
	SpriteFetcher             SpriteFetcher
	OAMBuffer                 OAMBuffer
}

type MemDump []MemInfo

func (md MemDump) Bytes() []Data8 {
	out := make([]Data8, len(md))
	for i := range out {
		out[i] = md[i].Value
	}
	return out
}

func (cd *CoreDump) PrintRegs(f io.Writer) {
	fmt.Fprintf(f, "PC = %s\n", cd.Regs.PC.Hex())
	fmt.Fprintf(f, "SP = %s\n", cd.Regs.SP.Hex())
	fmt.Fprintf(f, "A  = 0x%02x\n", cd.Regs.A)
	fmt.Fprintf(f, "F  = 0x%02x\n", cd.Regs.F)
	fmt.Fprintf(f, "B  = 0x%02x\n", cd.Regs.B)
	fmt.Fprintf(f, "C  = 0x%02x\n", cd.Regs.C)
	fmt.Fprintf(f, "D  = 0x%02x\n", cd.Regs.D)
	fmt.Fprintf(f, "E  = 0x%02x\n", cd.Regs.E)
	fmt.Fprintf(f, "H  = 0x%02x\n", cd.Regs.H)
	fmt.Fprintf(f, "L  = 0x%02x\n", cd.Regs.L)
	fmt.Fprintf(f, "W  = 0x%02x\n", cd.Regs.TempW)
	fmt.Fprintf(f, "Z  = 0x%02x\n", cd.Regs.TempZ)
	fmt.Fprintf(f, "IR = 0x%02x\n", uint8(cd.Regs.IR))
	z, h, n, c := 0, 0, 0, 0
	if cd.Regs.GetFlagZ() {
		z = 1
	}
	if cd.Regs.GetFlagH() {
		h = 1
	}
	if cd.Regs.GetFlagN() {
		n = 1
	}
	if cd.Regs.GetFlagC() {
		c = 1
	}
	fmt.Fprintf(f, "Z=%v C=%v\n", z, c)
	fmt.Fprintf(f, "N=%v H=%v\n", n, h)
}

func (cd *CoreDump) PrintProgram(f io.Writer) {
	if cd.ProgramEnd >= 0x8000 {
		fmt.Printf("out of bounds\n")
		return
	}
	memdump(f, cd.Program[cd.ProgramStart:cd.ProgramEnd+1], cd.ProgramStart, cd.ProgramEnd, cd.Regs.PC-1)
}

func (cd *CoreDump) PrintHRAM(f io.Writer) {
	memdump(f, cd.HRAM, AddrHRAMBegin, AddrHRAMEnd, cd.Regs.SP)
}

func (cd *CoreDump) PrintOAM(f io.Writer) {
	memdump(f, cd.OAM, AddrOAMBegin, AddrOAMEnd, 0)
}

func (cd *CoreDump) PrintOAMAttrs(f io.Writer) {
	oam := cd.OAM.Bytes()
	for idx := range 40 {
		obj := DecodeSprite(oam[idx*4 : (idx+1)*4])
		fmt.Fprintf(f, "%02d T=%03d X=%03d Y=%03d Attr=%x\n", idx, obj.TileIndex, obj.X, obj.Y, obj.Attributes.Hex())
	}
}

func (cd *CoreDump) PrintVRAM(f io.Writer) {
	memdump(f, cd.VRAM, AddrVRAMBegin, AddrVRAMEnd, 0)
}

func (cd *CoreDump) PrintWRAM(f io.Writer) {
	memdump(f, cd.VRAM, AddrWRAMBegin, AddrWRAMEnd, 0)
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (cd *CoreDump) PrintPPU(f io.Writer) {
	ppu := cd.PPU
	regdump(f, ppu.Registers, AddrPPUBegin, AddrPPUEnd)
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
	fmt.Fprintf(f, "BFetch.Susp:   %d\n", bool2int(ppu.BackgroundFetcher.Suspended))
	fmt.Fprintf(f, "BFetch.WYRch:  %d\n", bool2int(ppu.BackgroundFetcher.WindowYReached))
	fmt.Fprintf(f, "BFetch.WFetch: %d\n", bool2int(ppu.BackgroundFetcher.WindowFetching))
	fmt.Fprintf(f, "BFetch.WLC:    %d\n", ppu.BackgroundFetcher.WindowLineCounter)
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "BGFIFO: /[")
	for i := range ppu.BGFIFO.Level {
		fmt.Fprintf(f, "%d", ppu.BGFIFO.Slots[i].Color)
	}
	for range 8 - ppu.BGFIFO.Level {
		fmt.Fprintf(f, " ")
	}
	fmt.Fprintf(f, "]\n")
	fmt.Fprintf(f, "       [%d]\n", ppu.LastShifted)
	fmt.Fprintf(f, " SFIFO: \\[")
	for i := range ppu.SpriteFIFO.Level {
		fmt.Fprintf(f, "%d", ppu.SpriteFIFO.Slots[i].Color)
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
	fmt.Fprintf(f, "SFetch.Susp:   %d\n", bool2int(ppu.SpriteFetcher.Suspended))
	fmt.Fprintf(f, "Shift.Discard: %d\n", ppu.PixelShifter.Discard)
	fmt.Fprintf(f, "Shift.X:       %d\n", ppu.PixelShifter.X)
	fmt.Fprintf(f, "Shift.Susp:    %d\n", bool2int(ppu.PixelShifter.Suspended))
	fmt.Fprintf(f, "OAMBuffer.LV:  %d\n", ppu.OAMBuffer.Level)
}

func (cd *CoreDump) PrintAPU(f io.Writer) {
	regdump(f, cd.APU, AddrAPUBegin, AddrAPUEnd)
}

func regdump(f io.Writer, mem []MemInfo, start, end Addr) {
	for addr := start; addr <= end; addr++ {
		a := Addr(addr)
		if !a.IsValid() {
			continue
		}
		fmt.Fprintf(f, "%5s = %02x\n", Addr(addr), mem[addr-start].Value)
	}
}

func memdump(f io.Writer, mem []MemInfo, start, end, highlight Addr) {
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
			fmt.Fprintf(f, "[%02x]", mem[addr-start].Value)
		} else {
			fmt.Fprintf(f, " %02x ", mem[addr-start].Value)
		}
	}

	alignedEnd := (end/0x10)*0x10 + 0x10 - 1
	for addr := end; addr < alignedEnd; addr++ {
		fmt.Fprintf(f, " .. ")
	}
	fmt.Fprintf(f, "\n")
}

func (cpu *CPU) GetCoreDump() CoreDump {
	cpu.Bus.CoreDumpBegin()
	defer func() {
		cpu.Bus.CoreDumpEnd()
	}()

	var cd CoreDump
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
	cd.Program = cpu.Bus.getmem(0x0000, AddrCartridgeBank0End)
	cd.Disassembly = cpu.Disassembler.Disassembly()
	cd.HRAM = cpu.Bus.getmem(AddrHRAMBegin, AddrHRAMEnd)
	cd.OAM = cpu.Bus.getmem(AddrOAMBegin, AddrOAMEnd)
	cd.VRAM = cpu.Bus.getmem(AddrVRAMBegin, AddrVRAMEnd)
	cd.APU = cpu.Bus.getmem(AddrAPUBegin, AddrAPUEnd)
	cd.PPU.Registers = cpu.Bus.getmem(AddrPPUBegin, AddrPPUEnd)
	cd.PPU.BGFIFO = cpu.Bus.PPU.BackgroundFIFO.Dump()
	cd.PPU.SpriteFIFO = cpu.Bus.PPU.SpriteFIFO.Dump()
	cd.PPU.LastShifted = cpu.Bus.PPU.PixelShifter.LastShifted
	cd.PPU.OAMScanCycle = cpu.Bus.PPU.OAMScanCycle
	cd.PPU.PixelDrawCycle = cpu.Bus.PPU.PixelDrawCycle
	cd.PPU.HBlankRemainingCycles = cpu.Bus.PPU.HBlankRemainingCycles
	cd.PPU.VBlankLineRemainingCycles = cpu.Bus.PPU.VBlankLineRemainingCycles
	cd.PPU.PixelShifter = cpu.Bus.PPU.PixelShifter
	cd.PPU.BackgroundFetcher = cpu.Bus.PPU.BackgroundFetcher
	cd.PPU.SpriteFetcher = cpu.Bus.PPU.SpriteFetcher
	cd.PPU.OAMBuffer = cpu.Bus.PPU.OAMBuffer
	cd.RewindBuffer = cpu.rewindBuffer
	cd.RewindBufferIdx = cpu.rewindBufferIdx
	return cd
}

type MemInfo struct {
	Value        Data8
	ReadCounter  uint64
	WriteCounter uint64
}

func (bus *Bus) getmem(start, end Addr) []MemInfo {
	address := bus.Address
	data := bus.Data
	defer func() {
		bus.Address = address
		bus.Data = data
	}()

	out := make([]MemInfo, end-start+1)
	for addr := start; addr <= end; addr++ {
		memInfo := &out[addr-start]
		memInfo.ReadCounter, memInfo.WriteCounter = bus.GetCounters(addr)
		bus.WriteAddress(addr)
		memInfo.Value = bus.Data
	}

	return out
}

func (cpu *CPU) Dump() {
	cd := cpu.GetCoreDump()
	f := os.Stdout
	fmt.Fprintf(f, "\n--------\nCore dump:\n")
	cd.PrintRegs(f)
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
	for i := (cd.RewindBufferIdx + 1) % len(cd.RewindBuffer); i != cd.RewindBufferIdx; i = (i + 1) % len(cd.RewindBuffer) {
		entry := cd.RewindBuffer[i]
		extra := ""
		if entry.BranchResult == +1 {
			extra = "(taken)"
		} else if entry.BranchResult == -1 {
			extra = "(not taken)"
		}
		fmt.Fprintf(f, "[PC=%s] %s %s\n", entry.PC.Hex(), entry.Opcode, extra)
	}
}

func (cd *CoreDump) PrintDisassembly(f io.Writer) {
	cd.Disassembly.Print(f)
}
