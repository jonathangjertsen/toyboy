package model

// Bit 16: BG priority
// Bit 8-9: Color idx
// Bit 0-7: Palette
type Pixel = Data8

const (
	PxMaskPriority  = 0x80
	PxMaskPriority8 = 0x8080808080808080
	PXMaskColor     = 0x03
	PXMaskColor8    = 0x0303030303030303
)

func (p Pixel) Color() Color {
	return Color(p & 0x3)
}

const DefaultPalette = (0 << 0) | (1 << 2) | (2 << 4) | (3 << 6)

func DecodeLineBG(msb, lsb Data8, palette Data8) uint64 {
	var line uint64
	for i := 0; i < 8; i++ {
		shift := (((lsb>>(7-i))&1)<<1 | ((msb>>(7-i))&1)<<2)
		color := (palette >> shift) & 0x3
		line |= uint64(color) << (i * 8)
	}
	return line
}

func DecodeLineSprite(msb, lsb Data8, palette Data8) uint64 {
	var line uint64
	for i := 0; i < 8; i++ {
		shift := (((lsb>>(7-i))&1)<<1 | ((msb>>(7-i))&1)<<2)
		color := (palette >> shift) & 0x3
		if shift == 0 {
			color = 0
		}
		line |= uint64(color) << (i * 8)
	}
	return line
}
