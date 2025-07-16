package model

type PixelFIFO struct {
	Slots    [8]Pixel
	Level    int
	ShiftPos int
	PushPos  int
}

func (fifo *PixelFIFO) Clear() {
	clear(fifo.Slots[:])
	fifo.Level = 0
	fifo.ShiftPos = 0
	fifo.PushPos = 0
}

func (fifo *PixelFIFO) ShiftOut() (Pixel, bool) {
	var p Pixel
	if fifo.Level == 0 {
		return p, false
	}
	fifo.Level--
	p = fifo.Slots[fifo.ShiftPos]
	fifo.ShiftPos++
	if fifo.ShiftPos == 8 {
		fifo.ShiftPos = 0
	}
	return p, true
}

func (fifo *PixelFIFO) Write8(pixels [8]Pixel) {
	pos := fifo.ShiftPos
	for i := range 8 {
		fifo.Slots[pos] = pixels[i]
		pos++
		if pos == 8 {
			pos = 0
		}
	}
	fifo.Level = 8
}
