package model

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(White, LightGray, DarkGray, Black)
type Color uint8

type PPU struct {
	RegLCDC uint8
	RegSTAT uint8
	RegSCY  uint8
	RegSCX  uint8
	RegWY   uint8
	RegWX   uint8
	RegLY   uint8
	RegLYC  uint8
	RegBGP  uint8
	RegOBP0 uint8
	RegOBP1 uint8

	Palette [4]Color

	FBBackground [256][256]uint8
	FBWindow     [256][256]uint8
	FBViewport   [144][160]uint8
}

func (ppu *PPU) Enabled() bool {
	return ppu.RegLCDC&0x80 != 0
}

func (ppu *PPU) WindowTilemapArea() uint16 {
	if ppu.RegLCDC&0x40 != 0 {
		return 0x9800
	}
	return 0x9c00
}

func (ppu *PPU) WindowEnable() bool {
	if ppu.RegLCDC&0x01 == 0 {
		return false // DMG only
	}
	return ppu.RegLCDC&0x20 != 0
}

func (ppu *PPU) BGWindowTileDataArea() uint8 {
	bitSet := ppu.RegLCDC&0x10 != 0
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) BGTilemapArea() uint16 {
	if ppu.RegLCDC&0x08 != 0 {
		return 0x9800
	}
	return 0x9c00
}

func (ppu *PPU) OBJSize() uint8 {
	bitSet := ppu.RegLCDC&0x04 != 0
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) OBJEnable() uint8 {
	bitSet := ppu.RegLCDC&0x02 != 0
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) BGWindowEnablePriority() uint8 {
	bitSet := ppu.RegLCDC&0x01 != 0
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) SetLCDC(v uint8) {
	ppu.RegLCDC = v
}

func (ppu *PPU) SetSTAT(v uint8) {
	ppu.RegSTAT = v
	panic("not implemented: SetSTAT")
}

func (ppu *PPU) SetSCY(v uint8) {
	ppu.RegSCY = v
}

func (ppu *PPU) SetSCX(v uint8) {
	ppu.RegSCX = v
	panic("not implemented: SetSCX")
}

func (ppu *PPU) SetWY(v uint8) {
	ppu.RegWY = v
	panic("not implemented: SetWY")
}

func (ppu *PPU) SetWX(v uint8) {
	ppu.RegWX = v
	panic("not implemented: SetWX")
}

func (ppu *PPU) SetLY(v uint8) {
	ppu.RegLY = v
	panic("not implemented: SetLY")
}

func (ppu *PPU) SetLYC(v uint8) {
	ppu.RegLYC = v
	panic("not implemented: SetLYC")
}

func (ppu *PPU) SetBGP(v uint8) {
	ppu.RegBGP = v

	ppu.Palette[0] = Color((v >> 0) & 0x3)
	ppu.Palette[1] = Color((v >> 2) & 0x3)
	ppu.Palette[2] = Color((v >> 4) & 0x3)
	ppu.Palette[3] = Color((v >> 6) & 0x3)
}

func (ppu *PPU) SetOBP0(v uint8) {
	ppu.RegOBP0 = v
	panic("not implemented: SetOBP0")
}

func (ppu *PPU) SetOBP1(v uint8) {
	ppu.RegOBP1 = v
	panic("not implemented: SetOBP1")
}

func NewPPU() *PPU {
	return &PPU{}
}

func (ppu *PPU) Name() string {
	return "LCD"
}

func (ppu *PPU) Range() (uint16, uint16) {
	return 0xff40, 0x000c
}

func (ppu *PPU) Read(addr uint16) uint8 {
	switch addr {
	case 0xff40:
		return ppu.RegLCDC
	case 0xff41:
		return ppu.RegSTAT
	case 0xff42:
		return ppu.RegSCY
	case 0xff43:
		return ppu.RegSCX
	case 0xff44:
		return ppu.RegLY
	case 0xff45:
		return ppu.RegLYC
	case 0xff47:
		return ppu.RegBGP
	case 0xff48:
		return ppu.RegOBP0
	case 0xff49:
		return ppu.RegOBP1
	case 0xff4a:
		return ppu.RegWY
	case 0xff4b:
		return ppu.RegWX
	}
	panicf("Read from unknown LCD register %#v", addr)
	return 0
}

func (ppu *PPU) Write(addr uint16, v uint8) {
	switch addr {
	case 0xff40:
		ppu.SetLCDC(v)
	case 0xff41:
		ppu.SetSTAT(v)
	case 0xff42:
		ppu.SetSCY(v)
	case 0xff43:
		ppu.SetSCX(v)
	case 0xff44:
		ppu.SetLY(v)
	case 0xff45:
		ppu.SetLYC(v)
	case 0xff47:
		ppu.SetBGP(v)
	case 0xff48:
		ppu.SetOBP0(v)
	case 0xff49:
		ppu.SetOBP1(v)
	case 0xff4a:
		ppu.SetWY(v)
	case 0xff4b:
		ppu.SetWX(v)
	default:
		panicf("Write to unknown LCD register %#v", addr)
	}
}
