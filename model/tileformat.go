package model

const DefaultPalette = (0 << 0) | (1 << 2) | (2 << 4) | (3 << 6)

type TileLine [8]Pixel

func DecodeTileLine(msb, lsb Data8, palette Data16) TileLine {
	var pixels TileLine
	for i := range 8 {
		pixels[i] = (Data16(((lsb>>(7-i))&1)|(((msb>>(7-i))&1)<<1)) << 8) | palette
	}
	return pixels
}
