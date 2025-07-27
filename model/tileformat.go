package model

const DefaultPalette = (0 << 0) | (1 << 2) | (2 << 4) | (3 << 6)

type TileLine [8]Pixel

func DecodeTileLine(msb, lsb Data8) TileLine {
	var pixels TileLine
	for i := range 8 {
		pixels[i].ColorIDXBGPriority = ((lsb >> (7 - i)) & 1) | (((msb >> (7 - i)) & 1) << 1)
		pixels[i].Palette = DefaultPalette
	}
	return pixels
}

type Tile [8]TileLine

func DecodeTile(b []Data8) Tile {
	var tile Tile
	for i := range 8 {
		tile[i] = DecodeTileLine(b[2*i], b[2*i+1])
	}
	return tile
}
