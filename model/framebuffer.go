package model

type FrameBuffer [256][256]Color

func (vp *FrameBuffer) Flatten() [256 * 256]Color {
	var out [256 * 256]Color
	for i := range 256 {
		for j := range 256 {
			out[i*256+j] = vp[i][j]
		}
	}
	return out
}
