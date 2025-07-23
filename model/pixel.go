package model

type Pixel struct {
	ColorIDX           Data8
	Palette            [4]Color
	SpritePriority     bool // CGB only
	BackgroundPriority bool
}

func (p Pixel) Color() Color {
	return p.Palette[p.ColorIDX]
}
