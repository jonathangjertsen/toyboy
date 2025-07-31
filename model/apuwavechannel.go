package model

type WaveChannel struct {
	PeriodCounter PeriodCounter
	LengthTimer   LengthTimer

	Index    Addr
	OutLevel Data8
	Output   AudioSample

	RegDACEn         Data8
	RegLengthTimer   Data8
	RegOutputLevel   Data8
	RegPeriodLow     Data8
	RegPeriodHighCtl Data8

	DacEnabled bool
	Activated  bool
}

func (wc *WaveChannel) SetDACEn(v Data8) {
	wc.RegDACEn = v
	wc.DacEnabled = v&Bit7 != 0
	if !wc.DacEnabled {
		wc.Activated = false
	}
}

func (wc *WaveChannel) SetLengthTimer(v Data8) {
	wc.RegLengthTimer = v
	wc.LengthTimer.SetResetValue(v)
}

func (wc *WaveChannel) SetOutputLevel(v Data8) {
	wc.RegOutputLevel = v
	wc.OutLevel = (v >> 5) & 0x3
}

func (wc *WaveChannel) SetPeriodLow(v Data8) {
	wc.RegPeriodLow = v
	wc.PeriodCounter.SetPeriodLow(v)
}

func (wc *WaveChannel) SetPeriodHighCtl(v Data8) {
	wc.RegPeriodHighCtl = v
	wc.PeriodCounter.SetPeriodHigh(v)
	if v&Bit7 != 0 {
		wc.trigger()
	}
}

func (wc *WaveChannel) tickLengthTimer() {
	if disable := wc.LengthTimer.clock(256); disable {
		wc.Activated = false
	}
}

func (wc *WaveChannel) trigger() {
	// Ch1 is enabled.
	if wc.DacEnabled {
		wc.Activated = true
	}

	// If length timer expired it is reset.
	if wc.LengthTimer.Counter == 256 {
		wc.LengthTimer.Counter = Data16(wc.LengthTimer.Reset)
	}

	// The period divider is set to the contents of NR13 and NR14.
	wc.PeriodCounter.Counter = wc.PeriodCounter.Reset

	wc.Index = 1
}

func (wc *WaveChannel) clock(mem []Data8) {
	if !wc.Activated {
		return
	}
	if !wc.PeriodCounter.clock() {
		return
	}

	data := mem[AddrWaveRAMBegin+wc.Index>>1]
	if wc.Index&1 == 0 {
		// upper nibble on even index
		data >>= 4
	} else {
		// lower nibble on odd index
		data &= 0x0f
	}
	switch wc.OutLevel {
	case 0:
		data = 0
	case 1:
	case 2:
		data >>= 1
	case 3:
		data >>= 2
	}
	wc.Output = AudioSample(data)
	wc.Index++
	wc.Index &= 0x1f
}

func (wc *WaveChannel) Sample() AudioSample {
	if !wc.Activated {
		return 0
	}

	return wc.Output
}
