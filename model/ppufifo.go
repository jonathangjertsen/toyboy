package model

type EightPixels uint64

func (ep *EightPixels) Get(i int) Pixel {
	return Pixel((*ep >> (i * 8)) & 0xFF)
}

func (ep *EightPixels) Set(i int, p Pixel) {
	*ep &= ^(EightPixels(0xFF) << (i * 8))
	*ep |= EightPixels(p) << (i * 8)
}

type FIFO struct {
	Slots    EightPixels // 8 pixels as 8 bytes
	Level    int
	ShiftPos int
	PushPos  int
}

func (fifo *FIFO) Slot(i int) Pixel {
	return fifo.Slots.Get(i)
}

func (fifo *FIFO) SetSlot(i int, p Pixel) {
	fifo.Slots.Set(i, p)
}

func (fifo *FIFO) Clear() {
	fifo.Slots = 0
	fifo.Level = 0
	fifo.ShiftPos = 0
	fifo.PushPos = 0
}

func (fifo *FIFO) ShiftOut() (Pixel, bool) {
	if fifo.Level == 0 {
		return 0, false
	}
	p := Pixel((fifo.Slots >> (fifo.ShiftPos * 8)) & 0xFF)
	fifo.ShiftPos = (fifo.ShiftPos + 1) & 7
	fifo.Level--
	return p, true
}

func (fifo *FIFO) Write8(pixels [8]Pixel) {
	var packed EightPixels
	for i := 0; i < 8; i++ {
		packed |= EightPixels(pixels[i]) << (i * 8)
	}
	if fifo.ShiftPos == 0 {
		fifo.Slots = packed
	} else {
		var newSlots EightPixels
		remaining := 8 - fifo.ShiftPos
		for i := 0; i < remaining; i++ {
			newSlots |= EightPixels(pixels[i]) << ((fifo.ShiftPos + i) * 8)
		}
		for i := 0; i < fifo.ShiftPos; i++ {
			newSlots |= EightPixels(pixels[remaining+i]) << (i * 8)
		}
		fifo.Slots = newSlots
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
	for i := 0; i < fifo.Level; i++ {
		shift := (pos * 8) & 63
		dump.Slots[i] = Pixel((fifo.Slots >> shift) & 0xFF)
		pos = (pos + 1) & 7
	}
	return dump
}
