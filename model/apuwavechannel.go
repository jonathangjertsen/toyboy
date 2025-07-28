package model

type WaveChannel struct {
	PeriodCounter PeriodCounter
	LengthTimer   LengthTimer
	WaveGenerator WaveGenerator

	RegDACEn         Data8
	RegLengthTimer   Data8
	RegOutputLevel   Data8
	RegPeriodLow     Data8
	RegPeriodHighCtl Data8

	dacEnabled bool
	activated  bool
}

func (wc *WaveChannel) SetDACEn(v Data8) {
	wc.RegDACEn = v
	wc.dacEnabled = v&Bit7 != 0
	if !wc.dacEnabled {
		wc.activated = false
	}
}

func (wc *WaveChannel) SetLengthTimer(v Data8) {
	wc.RegLengthTimer = v
	wc.LengthTimer.SetResetValue(v)
}

func (wc *WaveChannel) SetOutputLevel(v Data8) {
	wc.RegOutputLevel = v
	wc.WaveGenerator.OutLevel = (v >> 5) & 0x3
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
		wc.activated = false
	}
}

func (wc *WaveChannel) trigger() {
	// Ch1 is enabled.
	if wc.dacEnabled {
		wc.activated = true
	}

	// If length timer expired it is reset.
	if wc.LengthTimer.lengthTimer == 256 {
		wc.LengthTimer.lengthTimer = Data16(wc.LengthTimer.lengthTimerReset)
	}

	// The period divider is set to the contents of NR13 and NR14.
	wc.PeriodCounter.periodDivider = wc.PeriodCounter.periodDividerReset

	wc.WaveGenerator.Index = 1
}

func (wc *WaveChannel) clock() {
	if !wc.activated {
		return
	}
	if wc.PeriodCounter.clock() {
		wc.WaveGenerator.clock()
	}
}

func (wc *WaveChannel) Sample() AudioSample {
	if !wc.activated {
		return 0
	}

	out := wc.WaveGenerator.output
	return out
}
