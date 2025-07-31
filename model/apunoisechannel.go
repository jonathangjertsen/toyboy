package model

type NoiseChannel struct {
	PeriodCounter PeriodCounter
	LengthTimer   LengthTimer
	Envelope      Envelope

	RegLengthTimer    Data8
	RegVolumeEnvelope Data8
	RegRNG            Data8
	RegCtl            Data8

	Activated       bool
	DacEnabled      bool
	ClockDivider    Data8
	ClockShift      Data8
	LFSRWidth       Data8
	LFSR            Data16
	Output          AudioSample
	ClockCounter    int
	ClockCounterTop int
}

func (nc *NoiseChannel) tickLengthTimer() {
	if disable := nc.LengthTimer.clock(64); disable {
		nc.Activated = false
	}
}
func (nc *NoiseChannel) clock() {
	if !nc.Activated {
		return
	}

	nc.ClockCounter--
	if nc.ClockCounter > 0 {
		return
	}
	nc.ClockCounter = nc.ClockCounterTop

	lfsr := nc.LFSR

	// Compute (LFSR) XOR (LFSR >> 1)
	// This puts LFSR0 XOR LFSR1 in bit position 0
	bit := ^(lfsr ^ (lfsr >> 1)) & 1

	// Place that bit where appropriate
	if nc.LFSRWidth == 7 {
		lfsr &= ^Data16(1<<15) & ^Data16(1<<7)
		lfsr |= (bit << 15) | (bit << 7)
	} else {
		lfsr &= ^Data16(1 << 15)
		lfsr |= (bit << 15)
	}

	// Shift out the bit
	lfsr >>= 1

	nc.LFSR = lfsr

	if bit == 1 {
		nc.Output = AudioSample(nc.Envelope.Volume)
	} else {
		nc.Output = 0
	}
}

func (nc *NoiseChannel) SetLengthTimer(v Data8) {
	nc.RegLengthTimer = v
	nc.LengthTimer.SetResetValue(v & 0x3f)
}

func (nc *NoiseChannel) SetVolumeEnvelope(v Data8) {
	nc.RegVolumeEnvelope = v
	nc.DacEnabled = nc.Envelope.SetVolumeEnvelope(v)
	if !nc.DacEnabled {
		nc.Activated = false
	}
}

func (nc *NoiseChannel) SetRNG(v Data8) {
	nc.RegRNG = v
	nc.ClockDivider = v & 0x7
	if v&Bit3 != 0 {
		nc.LFSRWidth = 7
	} else {
		nc.LFSRWidth = 15
	}
	nc.ClockShift = v >> 4

	// F = 262144 / (div * 2 ^ shift) Hz
	// T = (div * (2 ^ shift)) / 262144
	//   = (div * (2 ^ shift)) / (MCLK / 4)
	// So if we call clock() every MCLK / 4 cycles,
	//   we should clock the lfsr every div * 2 ^ shift times
	// But, div=0 is treated as div=0.5, and if both
	//   div== and shift==0 then 0.5*2^0 = 0.5 and we can't do half a clock
	// So we need to clock() every MCLK / 2 cycles,
	//   let div2 = div2 == 0 ? 1 : div * 2,
	//   and the nclock the lfsr every div' * 2 ^ shift times
	div2 := int(nc.ClockDivider) * 2
	if div2 == 0 {
		div2++
	}
	div2 <<= nc.ClockShift
	nc.ClockCounterTop = div2
}

func (nc *NoiseChannel) SetCtl(v Data8) {
	nc.RegCtl = v
	nc.LengthTimer.Enable = v&Bit6 != 0
	if v&Bit7 != 0 {
		nc.trigger()
	}
}

func (nc *NoiseChannel) trigger() {
	// Ch4 is enabled.
	if nc.DacEnabled {
		nc.Activated = true
	}

	// If length timer expired it is reset.
	if nc.LengthTimer.Counter == 64 {
		nc.LengthTimer.Counter = Data16(nc.LengthTimer.Reset)
	}

	// The period divider is set to the contents of NR33 and NR34?
	nc.PeriodCounter.Counter = nc.PeriodCounter.Reset

	// Envelope timer is reset.
	nc.Envelope.EnvTimer = 0

	// Volume is set to contents of NR42 initial volume.
	nc.Envelope.Volume = nc.Envelope.VolumeReset

	// LFSR is reset
	nc.LFSR = 0
}

func (nc *NoiseChannel) Sample() AudioSample {
	if !nc.Activated {
		return 0
	}

	return nc.Output
}
