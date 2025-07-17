package model

type TileLine [8]Pixel

func DecodeTileLine(msb, lsb uint8) TileLine {
	var pixels TileLine
	for i := range 8 {
		lsbMask := lsb & (1 << (7 - i))
		lsbMask >>= (7 - i)
		msbMask := msb & (1 << (7 - i))
		msbMask >>= (7 - i)
		msbMask <<= 1
		pixels[i].Color = Color(lsbMask | msbMask)
	}
	return pixels
}

type Tile [8]TileLine

func DecodeTile(b []byte) Tile {
	var tile Tile
	for i := range 8 {
		tile[i] = DecodeTileLine(b[2*i], b[2*i+1])
	}
	return tile
}
