package model

import "fmt"

//go:generate go-enum --marshal --flag --nocomments

// ENUM(
// Zero              = 0x0000
// BootROMEnd        = 0x00ff
// CartridgeBank0End = 0x3fff
// TileDataBegin     = 0x8000
// TileDataEnd       = 0x97ff
// TileMap0Begin     = 0x9800
// TileMap0End       = 0x9bff
// TileMap1Begin     = 0x9c00
// TileMap1End       = 0x9fff
// RAMBegin          = 0xa000
// RAMEnd            = 0xbfff
// WRAMBegin         = 0xc000
// WRAMEnd           = 0xdfff
// EchoRAMBegin      = 0xe000
// EchoRAMEnd        = 0xfdff
// OAMBegin          = 0xfe00
// OAMEnd = 0xfe9f
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
// BootROMLock = 0xff50
// HRAMBegin = 0xff80
// HRAMEnd = 0xfffe
// IE   = 0xffff
// )
type Addr uint16

const (
	AddrVRAMBegin = AddrTileDataBegin
	AddrVRAMEnd   = AddrTileMap1End
	AddrAPUBegin  = AddrNR10
	AddrAPUEnd    = AddrNR52
	AddrPPUBegin  = AddrLCDC
	AddrPPUEnd    = AddrWX
)

func ByteSlice(in []Data8) []byte {
	out := make([]byte, len(in))
	for i := range in {
		out[i] = byte(in[i])
	}
	return out
}

func Data8Slice(in []byte) []Data8 {
	out := make([]Data8, len(in))
	for i := range in {
		out[i] = Data8(in[i])
	}
	return out
}

func (a Addr) Hex() string {
	return Data16(a).Hex()
}

func (a Addr) Dec() string {
	return Data16(a).Dec()
}

func (a Addr) MSB() Data8 {
	return Data16(a).MSB()
}

func (a Addr) LSB() Data8 {
	return Data16(a).LSB()
}

func (a Addr) Split() (Data8, Data8) {
	return Data16(a).Split()
}

type Data16 uint16

func (a Data16) MSB() Data8 {
	return Data8(a >> 8)
}

func (a Data16) LSB() Data8 {
	return Data8(a)
}

func (a Data16) Hex() string {
	return Hex16(uint16(a))
}

func (a Data16) Dec() string {
	return fmt.Sprintf("%dd", a)
}

func (a Data16) Split() (Data8, Data8) {
	return a.MSB(), a.LSB()
}

type Data8 uint8

func (a Data8) Bit(i int) bool {
	if i < 0 || i > 8 {
		panic(i)
	}
	return a&(1<<i) != 0
}

func (a Data8) SignBit() bool {
	return a.Bit(7)
}

func (a Data8) SignedOffset() Offset8 {
	return Offset8(int8(uint8(a)))
}
func (a Data8) SignedAbs() Data8 {
	if a.SignBit() {
		return Data8(-int8(a))
	}
	return a
}

func (a Data8) Hex() string {
	return Hex8(uint8(a))
}

func (a Data8) Dec() string {
	return fmt.Sprintf("%dd", a)
}

func Hex16(x uint16) string {
	return fmt.Sprintf("%04xh", x)
}

func Hex8(x uint8) string {
	return fmt.Sprintf("%02xh", x)
}

type Size16 uint16

const (
	SizeBootROM        Size16 = 0x0100
	SizeCartridgeBank0 Size16 = 0x4000
	SizeVRAM           Size16 = 0x2000
	SizeTileData       Size16 = 0x0400
	SizeCartridgeRAM   Size16 = 0x2000
	SizeWRAM           Size16 = 0x2000
	SizeEchoRAM        Size16 = 0x1f00
	SizeAPU            Size16 = 0x0017
	SizePPU            Size16 = 0x000c
	SizeHRAM           Size16 = 0x007f
	SizeOAM            Size16 = 0x00a0
)

type Offset8 int8
