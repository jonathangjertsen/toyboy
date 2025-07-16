package model

const (
	AddrBootROMBegin = 0x0000
	AddrBootROMSize  = 0x0100
	AddrBootROMEnd   = AddrBootROMBegin + AddrBootROMSize - 1

	AddrCartridgeBank0Begin = 0x0000
	AddrCartridgeBank0Size  = 0x4000
	AddrCartridgeBank0End   = AddrCartridgeBank0Begin + AddrCartridgeBank0Size - 1

	AddrVRAMBegin = 0x8000
	SizeVRAM      = 0x2000
	AddrVRAMEnd   = AddrVRAMBegin + SizeVRAM - 1

	AddrAPUBegin = 0xff10
	SizeAPU      = 0x0017
	AddrAPUEnd   = AddrAPUBegin + SizeAPU - 1

	AddrPPUBegin = 0xff40
	SizePPU      = 0x000c
	AddrPPUEnd   = AddrPPUBegin + SizePPU - 1

	AddrBootROMLock = 0xff50

	AddrHRAMBegin = 0xff80
	SizeHRAM      = 0x007f
	AddrHRAMEnd   = AddrHRAMBegin + SizeHRAM - 1

	AddrOAMBegin = 0xfe00
	SizeOAM      = 0x00a0
	AddrOAMEnd   = AddrOAMBegin + SizeOAM - 1
)

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(
// P1   = 0xff00
// SB   = 0xff01
// SC   = 0xff02
// DIV  = 0xff04
// TIMA = 0xff05
// TMA  = 0xff06
// TAC  = 0xff07
// IF   = 0xff0f
// LCDC = 0xff40
// STAT = 0xff41
// SCY  = 0xff42
// SCX  = 0xff43
// LY   = 0xff44
// LYC  = 0xff45
// DMA  = 0xff46
// BGP  = 0xff47
// OBP0 = 0xff48
// OBP1 = 0xff49
// WY   = 0xff4a
// WX   = 0xff4b
// NR10 = 0xff10
// NR11 = 0xff11
// NR12 = 0xff12
// NR13 = 0xff13
// NR14 = 0xff14
// NR21 = 0xff16
// NR22 = 0xff17
// NR23 = 0xff18
// NR24 = 0xff19
// NR30 = 0xff1a
// NR31 = 0xff1b
// NR32 = 0xff1c
// NR33 = 0xff1d
// NR34 = 0xff1e
// NR41 = 0xff20
// NR42 = 0xff21
// NR43 = 0xff22
// NR44 = 0xff23
// NR50 = 0xff24
// NR51 = 0xff25
// NR52 = 0xff26
// )
type Addr uint16
