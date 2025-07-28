package model

type PulseChannel struct {
	PeriodCounter PeriodCounter
	LengthTimer   LengthTimer
	DutyGenerator DutyGenerator
	Envelope      Envelope

	RegLengthDuty     Data8
	RegVolumeEnvelope Data8
	RegPeriodLow      Data8
	RegPeriodHighCtl  Data8

	dacEnabled bool
	activated  bool

	test1 bool
}

type PulseChannelWithSweep struct {
	Sweep Sweep
	PulseChannel
}

func (pc *PulseChannel) tickLengthTimer() {
	if disable := pc.LengthTimer.clock(64); disable {
		pc.activated = false
	}
}

func (pc *PulseChannel) SetLengthDuty(v Data8) {
	pc.RegLengthDuty = v
	pc.LengthTimer.SetResetValue(v & 0x3f)
	pc.DutyGenerator.SetDuty(v)
}

func (pc *PulseChannel) SetVolumeEnvelope(v Data8) {
	pc.RegVolumeEnvelope = v
	pc.dacEnabled = pc.Envelope.SetVolumeEnvelope(v)
	if !pc.dacEnabled {
		pc.activated = false
	}
}

func (pc *PulseChannel) SetPeriodLow(v Data8) {
	pc.RegPeriodLow = v
	pc.PeriodCounter.SetPeriodLow(v)
}

func (pc *PulseChannel) SetPeriodHighCtl(v Data8) {
	pc.RegPeriodHighCtl = v
	pc.PeriodCounter.SetPeriodHigh(v)
	pc.LengthTimer.lengthEnable = v&Bit6 != 0
	if v&Bit7 != 0 {
		pc.trigger()
	}
}

func (pc *PulseChannel) trigger() {
	// Ch1 is enabled.
	if pc.dacEnabled {
		pc.activated = true
	}

	// If length timer expired it is reset.
	if pc.LengthTimer.lengthTimer == 64 {
		pc.LengthTimer.lengthTimer = Data16(pc.LengthTimer.lengthTimerReset)
	}

	// The period divider is set to the contents of NR13 and NR14.
	pc.PeriodCounter.periodDivider = pc.PeriodCounter.periodDividerReset

	// Envelope timer is reset.
	pc.Envelope.envTimer = 0

	// Volume is set to contents of NR12 initial volume.
	pc.Envelope.volume = pc.Envelope.volumeReset
}

func (pc *PulseChannel) clock() {
	if !pc.activated {
		return
	}
	if pc.PeriodCounter.clock() {
		pc.DutyGenerator.clock()
	}
}

func (pc *PulseChannel) Sample() int8 {
	if !pc.activated {
		return 0
	}

	out := pc.DutyGenerator.output
	out = pc.Envelope.scale(out)
	return out
}
