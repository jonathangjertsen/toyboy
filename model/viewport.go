package model

type ViewPort [144][160]Color

func (vp *ViewPort) Flatten() [144 * 160]Color {
	var out [144 * 160]Color
	for i := range 144 {
		for j := range 160 {
			out[i*144+j] = vp[i][j]
		}
	}
	return out
}
