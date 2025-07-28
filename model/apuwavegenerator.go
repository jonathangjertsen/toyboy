package model

type WaveGenerator struct {
	WaveRAM  [16]Data8
	Index    int
	OutLevel Data8
	output   AudioSample
}

func (wg *WaveGenerator) clock() {
	data := wg.WaveRAM[wg.Index>>1]
	if wg.Index&1 == 0 {
		// upper nibble on even index
		data >>= 4
	} else {
		// lower nibble on odd index
		data &= 0x0f
	}
	switch wg.OutLevel {
	case 0:
		data = 0
	case 1:
	case 2:
		data >>= 1
	case 3:
		data >>= 2
	}
	wg.output = AudioSample(data)
	wg.Index++
	wg.Index &= 0x1f
}
