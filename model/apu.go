package model

type APU struct {
	MemoryRegion
	MasterCtl Data8

	Pulse1 PulseChannelWithSweep
	Pulse2 PulseChannel
	Wave   WaveChannel
	Noise  NoiseChannel

	Mixer Mixer

	DIVAPU Data8

	canWriteLengthTimersWithAPUOff bool
}

func NewAPU(clock *ClockRT, config *Config) *APU {
	apu := &APU{
		MemoryRegion:                   NewMemoryRegion(clock, AddrAPUBegin, SizeAPU),
		canWriteLengthTimersWithAPUOff: true, // on monochrome models
	}
	if config.BootROM.Skip {
		apu.Reset()
	}
	clock.apu = apu
	return apu
}

func (apu *APU) Reset() {
	apu.SetMasterCtl(0x80)
	apu.Pulse1.SetLengthDuty(0x80)
	apu.Pulse1.SetVolumeEnvelope(0xf3)
	apu.SetChannelPan(0xf3)
	apu.SetMasterVolumePan(0x77)
}

func (apu *APU) incDIVAPU() {
	apu.DIVAPU++
	divapu := apu.DIVAPU

	if apu.MasterCtl&Bit7 == 0 {
		return
	}

	if divapu&0x1 == 0x0 {
		apu.Pulse1.tickLengthTimer()
		apu.Pulse2.tickLengthTimer()
		apu.Wave.tickLengthTimer()
		apu.Noise.tickLengthTimer()
	}

	if divapu&0x3 == 0x0 {
		// Every 4 ticks: CH1 Freq Sweep
	}

	if divapu&0x7 == 0x0 {
		divapu >>= 3
		apu.Pulse1.Envelope.clock(divapu)
		apu.Pulse2.Envelope.clock(divapu)
		apu.Noise.Envelope.clock(divapu)
	}
}

func (apu *APU) Read(addr Addr) Data8 {
	switch Addr(addr) {
	case AddrNR10:
		return apu.Pulse1.Sweep.RegSweep
	case AddrNR11:
		return apu.Pulse1.RegLengthDuty
	case AddrNR12:
		return apu.Pulse1.RegVolumeEnvelope
	case AddrNR13:
		return apu.Pulse1.RegPeriodLow // WO
	case AddrNR14:
		return apu.Pulse1.RegPeriodHighCtl
	case 0xff15:
		return 0
	case AddrNR21:
		return apu.Pulse2.RegLengthDuty
	case AddrNR22:
		return apu.Pulse2.RegVolumeEnvelope
	case AddrNR23:
		return apu.Pulse2.RegPeriodLow
	case AddrNR24:
		return apu.Pulse2.RegPeriodHighCtl
	case AddrNR30:
		return apu.Wave.RegDACEn
	case AddrNR31:
		return apu.Wave.RegLengthTimer
	case AddrNR32:
		return apu.Wave.RegOutputLevel
	case AddrNR33:
		return apu.Wave.RegPeriodLow
	case AddrNR34:
		return apu.Wave.RegPeriodHighCtl
	case 0xff1f:
		return 0
	case AddrNR41:
		return apu.Noise.RegLengthTimer
	case AddrNR42:
		return apu.Noise.RegVolumeEnvelope
	case AddrNR43:
		return apu.Noise.RegRNG
	case AddrNR44:
		return apu.Noise.RegCtl
	case AddrNR50:
		return apu.Mixer.RegMasterVolumeVINPan
	case AddrNR51:
		return apu.Mixer.RegChannelPan
	case AddrNR52:
		return apu.ReadMasterCtl()
	}
	panicf("Read from unknown apu register %#v", addr)
	return 0
}

func (apu *APU) Write(addr Addr, v Data8) {
	apu.MemoryRegion.Data[addr-AddrAPUBegin] = v

	switch Addr(addr) {
	case AddrNR10:
		apu.SetPulse1Sweep(v)
	case AddrNR11:
		apu.SetPulse1LengthDuty(v)
	case AddrNR12:
		apu.SetPulse1VolumeEnvelope(v)
	case AddrNR13:
		apu.SetPulse1PeriodLow(v) // WO
	case AddrNR14:
		apu.SetPulse1PeriodHighCtl(v)
	case 0xff15:
	case AddrNR21:
		apu.SetPulse2LengthDuty(v)
	case AddrNR22:
		apu.SetPulse2VolumeEnvelope(v)
	case AddrNR23:
		apu.SetPulse2PeriodLow(v)
	case AddrNR24:
		apu.SetPulse2PeriodHighCtl(v)
	case AddrNR30:
		apu.SetWaveDACEn(v)
	case AddrNR31:
		apu.SetWaveLengthTimer(v)
	case AddrNR32:
		apu.SetWaveOutputLevel(v)
	case AddrNR33:
		apu.SetWavePeriodLow(v)
	case AddrNR34:
		apu.SetWavePeriodHighCtl(v)
	case 0xff1f:
	case AddrNR41:
		apu.SetNoiseLengthTimer(v)
	case AddrNR42:
		apu.SetNoiseVolumeEnvelope(v)
	case AddrNR43:
		apu.SetNoiseRNG(v)
	case AddrNR44:
		apu.SetNoiseCtl(v)
	case AddrNR50:
		apu.SetMasterVolumePan(v)
	case AddrNR51:
		apu.SetChannelPan(v)
	case AddrNR52:
		apu.SetMasterCtl(v)
	default:
		panicf("Write to unknown apu register %#v", addr)
	}
}

func (apu *APU) SetPulse1Sweep(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Pulse1.Sweep.SetSweep(v)
}

func (apu *APU) SetPulse1LengthDuty(v Data8) {
	if !apu.canWriteLengthTimersWithAPUOff && (apu.MasterCtl&Bit7 == 0) {
		return
	}
	apu.Pulse1.SetLengthDuty(v)
}

func (apu *APU) SetPulse1VolumeEnvelope(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Pulse1.SetVolumeEnvelope(v)
}

func (apu *APU) SetPulse1PeriodLow(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Pulse1.SetPeriodLow(v)
}

func (apu *APU) SetPulse1PeriodHighCtl(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Pulse1.SetPeriodHighCtl(v)
}

func (apu *APU) SetPulse2LengthDuty(v Data8) {
	if !apu.canWriteLengthTimersWithAPUOff && (apu.MasterCtl&Bit7 == 0) {
		return
	}
	apu.Pulse2.SetLengthDuty(v)
}

func (apu *APU) SetPulse2VolumeEnvelope(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Pulse2.SetVolumeEnvelope(v)
}

func (apu *APU) SetPulse2PeriodLow(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Pulse2.SetPeriodLow(v)
}

func (apu *APU) SetPulse2PeriodHighCtl(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Pulse2.SetPeriodHighCtl(v)
}

func (apu *APU) SetWaveDACEn(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Wave.SetDACEn(v)
}

func (apu *APU) SetWaveLengthTimer(v Data8) {
	if !apu.canWriteLengthTimersWithAPUOff && (apu.MasterCtl&Bit7 == 0) {
		return
	}
	apu.Wave.SetLengthTimer(v)
}

func (apu *APU) SetWaveOutputLevel(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Wave.SetOutputLevel(v)
}

func (apu *APU) SetWavePeriodLow(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Wave.SetPeriodLow(v)
}

func (apu *APU) SetWavePeriodHighCtl(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Wave.SetPeriodHighCtl(v)
}

func (apu *APU) SetNoiseLengthTimer(v Data8) {
	if !apu.canWriteLengthTimersWithAPUOff && (apu.MasterCtl&Bit7 == 0) {
		return
	}
	apu.Noise.SetLengthTimer(v)
}

func (apu *APU) SetNoiseVolumeEnvelope(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Noise.SetVolumeEnvelope(v)
}

func (apu *APU) SetNoiseRNG(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Noise.SetRNG(v)
}

func (apu *APU) SetNoiseCtl(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Noise.SetCtl(v)
}

func (apu *APU) SetMasterVolumePan(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Mixer.SetMasterVolumeVINPan(v)
}

func (apu *APU) SetChannelPan(v Data8) {
	if apu.MasterCtl&Bit7 == 0 {
		return
	}
	apu.Mixer.SetChannelPan(v)
}

func (apu *APU) ReadMasterCtl() Data8 {
	v := apu.MasterCtl
	v &= 0x80
	if apu.Pulse1.activated {
		v |= 1
	}
	if apu.Pulse2.activated {
		v |= 2
	}
	if apu.Wave.activated {
		v |= 4
	}
	if apu.Noise.activated {
		v |= 8
	}
	return v
}

func (apu *APU) SetMasterCtl(v Data8) {
	// Write the whole register - we will calculate bits 0:3 on read
	apu.MasterCtl = v

	// Turning the APU off clears all APU registers
	if apu.MasterCtl&Bit7 == 0 {
		apu.Pulse1.Sweep.SetSweep(0)
		apu.Pulse1.SetLengthDuty(0)
		apu.Pulse1.SetVolumeEnvelope(0)
		apu.Pulse1.SetPeriodLow(0)
		apu.Pulse1.SetPeriodHighCtl(0)
		apu.Pulse2.SetLengthDuty(0)
		apu.Pulse2.SetVolumeEnvelope(0)
		apu.Pulse2.SetPeriodLow(0)
		apu.Pulse2.SetPeriodHighCtl(0)
		apu.Wave.SetDACEn(0)
		apu.Wave.SetLengthTimer(0)
		apu.Wave.SetOutputLevel(0)
		apu.Wave.SetPeriodLow(0)
		apu.Wave.SetPeriodHighCtl(0)
		apu.Noise.SetLengthTimer(0)
		apu.Noise.SetVolumeEnvelope(0)
		apu.Noise.SetRNG(0)
		apu.Noise.SetCtl(0)
		apu.Mixer.SetMasterVolumeVINPan(0)
		apu.Mixer.SetChannelPan(0)
		apu.MasterCtl = 0
	}
}

func maskedWrite(prev, v, mask Data8) Data8 {
	return (prev & (^mask)) | (v & mask)
}
