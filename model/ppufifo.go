package model

type FIFO struct {
	Slots    [8]Pixel
	Level    int
	ShiftPos int
	PushPos  int
}

func (fifo *FIFO) Clear() {
	clear(fifo.Slots[:])
	fifo.Level = 0
	fifo.ShiftPos = 0
	fifo.PushPos = 0
}

func (fifo *FIFO) ShiftOut() (Pixel, bool) {
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

func (fifo *FIFO) Write8(pixels [8]Pixel) {
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

type PixelFIFODump struct {
	Slots [8]Pixel
	Level int
}

func (fifo *FIFO) Dump() PixelFIFODump {
	var dump PixelFIFODump
	dump.Level = fifo.Level
	pos := fifo.ShiftPos
	for i := range fifo.Level {
		dump.Slots[i] = fifo.Slots[pos]
		i++
		if i == 8 {
			i = 0
		}
	}
	return dump
}
