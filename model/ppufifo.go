package model

func GetPixel(ep uint64, i int) Pixel {
	return Pixel((ep >> (i * 8)) & 0xFF)
}

func SetPixel(ep *uint64, i int, p Pixel) {
	*ep &= ^(uint64(0xFF) << (i * 8))
	*ep |= uint64(p) << (i * 8)
}

type FIFO struct {
	Slots uint64 // 8 pixels as 8 bytes
	Level int    // 0 to 8
}

func (fifo *FIFO) Clear() {
	fifo.Slots = 0
	fifo.Level = 0
}

func (fifo *FIFO) ShiftOut() (Pixel, bool) {
	if fifo.Level == 0 {
		return 0, false
	}
	p := Pixel(fifo.Slots)
	fifo.Slots >>= 8
	fifo.Level--
	return p, true
}
