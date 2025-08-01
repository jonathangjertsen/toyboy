package model

// Bit 16: BG priority
// Bit 8-9: Color idx
// Bit 0-7: Palette
type Pixel = Data8

type TileLine [8]Pixel

const (
	PxMaskPriority = 0x80
)

func (p Pixel) Color() Color {
	return Color(p & 0x3)
}

const DefaultPalette = (0 << 0) | (1 << 2) | (2 << 4) | (3 << 6)

func DecodeTileLineBG(msb, lsb Data8, palette Data8) TileLine {
	var pixels TileLine
	for i := range 8 {
		shift := (((lsb>>(7-i))&1)<<1 | (((msb >> (7 - i)) & 1) << 2))
		color := (palette & (0x3 << shift)) >> shift
		pixels[i] = color
	}
	return pixels
}

func DecodeTileLineSprite(msb, lsb Data8, palette Data8) TileLine {
	var pixels TileLine
	for i := range 8 {
		shift := (((lsb>>(7-i))&1)<<1 | (((msb >> (7 - i)) & 1) << 2))
		color := (palette & (0x3 << shift)) >> shift
		if shift == 0 {
			color = 0
		}
		pixels[i] = color
	}
	return pixels
}
