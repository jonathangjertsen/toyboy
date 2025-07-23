package model

type TileLine [8]Pixel

func DecodeTileLine(msb, lsb Data8) TileLine {
	var pixels TileLine
	for i := range 8 {
		lsbMask := lsb & (1 << (7 - i))
		lsbMask >>= (7 - i)
		msbMask := msb & (1 << (7 - i))
		msbMask >>= (7 - i)
		msbMask <<= 1
		pixels[i].ColorIDX = lsbMask | msbMask
		pixels[i].Palette[0] = 0
		pixels[i].Palette[1] = 1
		pixels[i].Palette[2] = 2
		pixels[i].Palette[3] = 3
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
