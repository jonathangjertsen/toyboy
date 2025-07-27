package model

type ViewPort [144][160]Color

func (vp *ViewPort) Flatten() [144 * 160]Color {
	var out [144 * 160]Color
	for i := range 144 {
		for j := range 160 {
			out[i*160+j] = vp[i][j]
		}
	}
	return out
}

func (vp *ViewPort) Grayscale() [144 * 160]uint8 {
	var out [144 * 160]uint8
	for i := range 144 {
		for j := range 160 {
			out[i*160+j] = vp[i][j].Grayscale()
		}
	}
	return out
}
