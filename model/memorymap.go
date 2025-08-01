package model

//go:generate go-enum --flag --no-iota --nocomments

// ENUM(
// Zero                 = 0x0000
// BootROMEnd           = 0x00ff
// CartridgeEntryPoint  = 0x0100
// NintendoLogoBegin    = 0x0104
// NintendoLogoEnd      = 0x0133
// TitleBegin           = 0x0134
// TitleEnd             = 0x0143
// NewLicenseeCodeBegin = 0x0144
// NewLicenseeCodeEnd   = 0x0145
// SGBFlag              = 0x0146
// CartridgeType        = 0x0147
// ROMSize              = 0x0148
// RAMSize              = 0x0149
// DestCode             = 0x014a
// OldLicenseeCode      = 0x014b
// MaskROMVersionNo     = 0x014c
// HeaderChecksum       = 0x014d
// GlobalChecksumBegin  = 0x014e
// GlobalChecksumEnd    = 0x014f
// CartridgeBank0End    = 0x3fff
// CartridgeBankNBegin  = 0x4000
// CartridgeBankNEnd    = 0x7fff
// TileDataBegin        = 0x8000
// TileDataEnd          = 0x97ff
// TileMap0Begin        = 0x9800
// TileMap0End          = 0x9bff
// TileMap1Begin        = 0x9c00
// TileMap1End          = 0x9fff
// CartridgeRAMBegin    = 0xa000
// CartridgeRAMEnd      = 0xbfff
// WRAMBegin            = 0xc000
// WRAMEnd              = 0xdfff
// EchoRAMBegin         = 0xe000
// EchoRAMEnd           = 0xfdff
// OAMBegin             = 0xfe00
// OAMEnd               = 0xfe9f
// ProhibitedBegin      = 0xfea0
// ProhibitedEnd        = 0xfeff
// P1                   = 0xff00
// SB                   = 0xff01
// SC                   = 0xff02
// DIV                  = 0xff04
// TIMA                 = 0xff05
// TMA                  = 0xff06
// TAC                  = 0xff07
// IF                   = 0xff0f
// LCDC                 = 0xff40
// STAT                 = 0xff41
// SCY                  = 0xff42
// SCX                  = 0xff43
// LY                   = 0xff44
// LYC                  = 0xff45
// DMA                  = 0xff46
// BGP                  = 0xff47
// OBP0                 = 0xff48
// OBP1                 = 0xff49
// WY                   = 0xff4a
// WX                   = 0xff4b
// NR10                 = 0xff10
// NR11                 = 0xff11
// NR12                 = 0xff12
// NR13                 = 0xff13
// NR14                 = 0xff14
// NR21                 = 0xff16
// NR22                 = 0xff17
// NR23                 = 0xff18
// NR24                 = 0xff19
// NR30                 = 0xff1a
// NR31                 = 0xff1b
// NR32                 = 0xff1c
// NR33                 = 0xff1d
// NR34                 = 0xff1e
// NR41                 = 0xff20
// NR42                 = 0xff21
// NR43                 = 0xff22
// NR44                 = 0xff23
// NR50                 = 0xff24
// NR51                 = 0xff25
// NR52                 = 0xff26
// WaveRAMBegin         = 0xff30
// WaveRAMEnd           = 0xff3f
// BootROMLock          = 0xff50
// HRAMBegin            = 0xff80
// HRAMEnd              = 0xfffe
// IE                   = 0xffff
// )
type Addr uint16

const (
	AddrCartridgeHeaderBegin = AddrCartridgeEntryPoint
	AddrCartridgeHeaderEnd   = AddrGlobalChecksumEnd
	AddrVRAMBegin            = AddrTileDataBegin
	AddrVRAMEnd              = AddrTileMap1End
	AddrAPUBegin             = AddrNR10
	AddrAPUEnd               = AddrNR52
	AddrPPUBegin             = AddrLCDC
	AddrPPUEnd               = AddrWX
	AddrTimerBegin           = AddrDIV
	AddrTimerEnd             = AddrTAC
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

func U8Slice(in []Data8) []uint8 {
	out := make([]uint8, len(in))
	for i := range in {
		out[i] = uint8(in[i])
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

type Size16 uint16

const (
	SizeCartridgeHeader Size16 = 0x014f
	SizeBootROM         Size16 = 0x0100
	SizeCartridgeBank   Size16 = 0x4000
	SizeVRAM            Size16 = 0x2000
	SizeTileData        Size16 = 0x0400
	SizeCartridgeRAM    Size16 = 0x2000
	SizeWRAM            Size16 = 0x2000
	SizeEchoRAM         Size16 = 0x1f00
	SizeAPU             Size16 = 0x0030
	SizePPU             Size16 = 0x000c
	SizeHRAM            Size16 = 0x007f
	SizeOAM             Size16 = 0x00a0
	SizeProhibited      Size16 = 0x0060
	SizeTimer           Size16 = 0x0004
	SizeWaveRAM         Size16 = 0x0010
)

type Offset8 int8
