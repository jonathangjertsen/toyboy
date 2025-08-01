package model

// Bit 16: BG priority
// Bit 8-9: Color idx
// Bit 0-7: Palette
type Pixel = Data16

const (
	PxMaskColorIDX  = 0x300
	PxShiftColorIDX = 8
	PxMaskPriority  = 0x8000
	PxMaskPalette   = 0x00ff
)

func (p Pixel) Color() Color {
	shift := ((p & PxMaskColorIDX) >> (PxShiftColorIDX - 1))
	return Color((p & (0x3 << shift)) >> shift)
}
