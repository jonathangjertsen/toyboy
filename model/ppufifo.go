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
	if fifo.Level == 0 {
		return Pixel{}, false
	}
	fifo.Level--
	p := fifo.Slots[fifo.ShiftPos]
	fifo.ShiftPos = (fifo.ShiftPos + 1) & 7
	return p, true
}

func (fifo *FIFO) Write8(pixels [8]Pixel) {
	if fifo.ShiftPos == 0 {
		fifo.Slots = pixels
	} else {
		remaining := 8 - fifo.ShiftPos
		copy(fifo.Slots[fifo.ShiftPos:], pixels[:remaining])
		copy(fifo.Slots[:fifo.ShiftPos], pixels[remaining:])
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
