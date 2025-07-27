package model

type WaveChannel struct {
	PeriodCounter PeriodCounter
	LengthTimer   LengthTimer

	RegDACEn         Data8
	RegLengthTimer   Data8
	RegOutputLevel   Data8
	RegPeriodLow     Data8
	RegPeriodHighCtl Data8

	dacEnabled bool
	activated  bool
}

func (wc *WaveChannel) clock() {

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
}

func (wc *WaveChannel) SetPeriodLow(v Data8) {
	wc.RegPeriodLow = v
	wc.PeriodCounter.SetPeriodLow(v)
}

func (wc *WaveChannel) SetPeriodHighCtl(v Data8) {
	wc.RegPeriodHighCtl = v
	wc.PeriodCounter.SetPeriodHigh(v)
}

func (wc *WaveChannel) tickLengthTimer() {
	if disable := wc.LengthTimer.clock(256); disable {
		wc.activated = false
	}
}
