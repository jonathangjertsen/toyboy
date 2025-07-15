package model

import (
	"fmt"
	"io"
)

type CoreDump struct {
	Regs         RegisterFile
	ProgramStart uint16
	ProgramEnd   uint16
	Program      []uint8
	HRAM         []uint8
	OAM          []uint8
	VRAM         []uint8
}

func (cd *CoreDump) PrintRegs(f io.Writer) {
	fmt.Fprintf(f, "PC = 0x%04x\n", cd.Regs.PC)
	fmt.Fprintf(f, "SP = 0x%04x\n", cd.Regs.SP)
	fmt.Fprintf(f, "A  =   0x%02x\n", cd.Regs.A)
	fmt.Fprintf(f, "F  =   0x%02x\n", cd.Regs.F)
	fmt.Fprintf(f, "B  =   0x%02x\n", cd.Regs.B)
	fmt.Fprintf(f, "C  =   0x%02x\n", cd.Regs.C)
	fmt.Fprintf(f, "D  =   0x%02x\n", cd.Regs.D)
	fmt.Fprintf(f, "E  =   0x%02x\n", cd.Regs.E)
	fmt.Fprintf(f, "H  =   0x%02x\n", cd.Regs.H)
	fmt.Fprintf(f, "L  =   0x%02x\n", cd.Regs.L)
	fmt.Fprintf(f, "W  =   0x%02x\n", cd.Regs.TempW)
	fmt.Fprintf(f, "Z  =   0x%02x\n", cd.Regs.TempZ)
	fmt.Fprintf(f, "IR =   0x%02x\n", uint8(cd.Regs.IR))
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
	memdump(f, cd.Program, cd.ProgramStart, cd.ProgramEnd, cd.Regs.PC-1)
}

func (cd *CoreDump) PrintHRAM(f io.Writer) {
	memdump(f, cd.HRAM, 0xff80, 0xfffe, cd.Regs.SP)
}

func (cd *CoreDump) PrintOAM(f io.Writer) {
	memdump(f, cd.OAM, 0xfe00, 0xfe99, 0)
}

func (cd *CoreDump) PrintVRAM(f io.Writer) {
	memdump(f, cd.VRAM, 0x8000, 0x9fff, 0)
}

func (cd *CoreDump) Print(f io.Writer) {
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
	fmt.Fprintf(f, "VRAM:\n")
	cd.PrintVRAM(f)
	fmt.Fprintf(f, "--------\n")
}

func memdump(f io.Writer, mem []uint8, start, end, highlight uint16) {
	alignedStart := (start / 0x10) * 0x10
	for addr := alignedStart; addr < start; addr++ {
		if addr%0x10 == 0 {
			fmt.Fprintf(f, "\n %04x |", addr)
		}
		fmt.Fprintf(f, " .. ")
	}

	for addr := start; addr <= end; addr++ {
		if addr%0x10 == 0 {
			fmt.Fprintf(f, "\n %04x |", addr)
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
