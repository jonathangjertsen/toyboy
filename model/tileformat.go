package model

func decodeTileLine(lsb, msb uint8) [8]Pixel {
	var pixels [8]Pixel
	for i := 0; i < 8; i++ {
		lsbMask := lsb & (1 << (8 - i))
		lsbMask >>= (8 - i)
		msbMask := msb & (1 << (8 - i))
		msbMask >>= (8 - i)
		msbMask <<= 1
		pixels[i].Color = Color(lsbMask | msbMask)
	}
	return pixels
}
