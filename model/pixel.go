package model

type Pixel struct {
	Color              Color
	Palette            Data8
	SpritePriority     bool // CGB only
	BackgroundPriority bool
}
