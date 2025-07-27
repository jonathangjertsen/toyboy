package model

type Pixel struct {
	ColorIDXBGPriority Data8 // Bit 7: BG priority. Bit 0-1: color idx
	Palette            Data8
}

func (p Pixel) ColorIDX() Data8 {
	return p.ColorIDXBGPriority & 0x3
}

func (p Pixel) BGPriority() bool {
	return p.ColorIDXBGPriority&Bit7 != 0
}

func (p *Pixel) SetBGPriority() {
	p.ColorIDXBGPriority |= 0x80
}

func (p Pixel) Color() Color {
	shift := p.ColorIDX() * 2
	return Color((p.Palette & (0x3 << shift)) >> shift)
}
