package model

type PulseChannel struct {
	PeriodCounter PeriodCounter
	LengthTimer   LengthTimer
	Envelope      Envelope

	Waveform Data8
	Output   AudioSample

	RegLengthDuty     Data8
	RegVolumeEnvelope Data8
	RegPeriodLow      Data8
	RegPeriodHighCtl  Data8

	DacEnabled bool
	Activated  bool
}

type PulseChannelWithSweep struct {
	Sweep Sweep
	PulseChannel
}

func (pc *PulseChannel) tickLengthTimer() {
	if disable := pc.LengthTimer.clock(64); disable {
		pc.Activated = false
	}
}

func (pc *PulseChannel) SetLengthDuty(v Data8) {
	pc.RegLengthDuty = v
	pc.LengthTimer.SetResetValue(v & 0x3f)
	switch (v >> 6) & 0x3 {
	case 0:
		pc.Waveform = 0b1111_1110 // 12.5%
	case 1:
		pc.Waveform = 0b0111_1110 // 25.0%
	case 2:
		pc.Waveform = 0b0111_1000 // 50.0%
	case 3:
		pc.Waveform = 0b1000_0001 // 75.0%
	}
}

func (pc *PulseChannel) SetVolumeEnvelope(v Data8) {
	pc.RegVolumeEnvelope = v
	pc.DacEnabled = pc.Envelope.SetVolumeEnvelope(v)
	if !pc.DacEnabled {
		pc.Activated = false
	}
}

func (pc *PulseChannel) SetPeriodLow(v Data8) {
	pc.RegPeriodLow = v
	pc.PeriodCounter.SetPeriodLow(v)
}

func (pc *PulseChannel) SetPeriodHighCtl(v Data8) {
	pc.RegPeriodHighCtl = v
	pc.PeriodCounter.SetPeriodHigh(v)
	pc.LengthTimer.Enable = v&Bit6 != 0
	if v&Bit7 != 0 {
		pc.trigger()
	}
}

func (pc *PulseChannel) trigger() {
	// Ch1 is enabled.
	if pc.DacEnabled {
		pc.Activated = true
	}

	// If length timer expired it is reset.
	if pc.LengthTimer.Counter == 64 {
		pc.LengthTimer.Counter = Data16(pc.LengthTimer.Reset)
	}

	// The period divider is set to the contents of NR13 and NR14.
	pc.PeriodCounter.Counter = pc.PeriodCounter.Reset

	// Envelope timer is reset.
	pc.Envelope.EnvTimer = 0

	// Volume is set to contents of NR12 initial volume.
	pc.Envelope.Volume = pc.Envelope.VolumeReset
}

func (pc *PulseChannel) clock() {
	if !pc.Activated {
		return
	}
	if pc.PeriodCounter.clock() {
		if pc.Waveform&Bit0 != 0 {
			pc.Output = 1
			pc.Waveform = (pc.Waveform >> 1) | 0x80
		} else {
			pc.Output = 0
			pc.Waveform >>= 1
		}
	}
}

func (pc *PulseChannel) Sample() AudioSample {
	if !pc.Activated {
		return 0
	}
	return pc.Envelope.scale(pc.Output)
}
