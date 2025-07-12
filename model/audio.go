package model

import (
	"fmt"
	"slices"
)

var audioDebugEvents = []string{}

type Audio struct {
	Pulse1Sweep          uint8
	Pulse1LengthDuty     uint8
	Pulse1VolumeEnvelope uint8
	Pulse1PeriodLow      uint8
	Pulse1PeriodHighCtl  uint8
	Pulse2LengthDuty     uint8
	Pulse2VolumeEnvelope uint8
	Pulse2PeriodLow      uint8
	Pulse2PeriodHighCtl  uint8
	WaveDACEn            uint8
	WaveLengthTimer      uint8
	WaveOutputLevel      uint8
	WavePeriodLow        uint8
	WavePeriodHighCtl    uint8
	NoiseLengthTimer     uint8
	NoiseVolumeEnvelope  uint8
	NoiseRNG             uint8
	NoiseCtl             uint8
	MasterVolumePan      uint8
	ChannelPan           uint8
	MasterCtl            uint8

	clockCycle Cycle

	canWriteLengthTimersWithAPUOff bool
}

func NewAudio() *Audio {
	return &Audio{
		canWriteLengthTimersWithAPUOff: true, // on monochrome models
	}
}

func (audio *Audio) Debug(event string, f string, v ...any) {
	if !slices.Contains(audioDebugEvents, event) {
		return
	}
	dir := "v"
	if audio.clockCycle.Rising {
		dir = "^"
	}
	fmt.Printf("%d %s | %s | ", audio.clockCycle.C, dir, event)
	fmt.Printf(f, v...)
}

func (audio *Audio) APUEnabled() bool {
	return audio.MasterCtl&0x80 != 0
}

func (audio *Audio) Name() string {
	return "AUDIO"
}

func (audio *Audio) Range() (uint16, uint16) {
	return 0xff10, 0x0017
}

func (audio *Audio) Read(addr uint16) uint8 {
	switch addr {
	case 0xff10:
		return audio.Pulse1Sweep
	case 0xff11:
		return audio.Pulse1LengthDuty
	case 0xff12:
		return audio.Pulse1VolumeEnvelope
	case 0xff13:
		return audio.Pulse1PeriodLow // WO
	case 0xff14:
		return audio.Pulse1PeriodHighCtl
	case 0xff15:
		return 0
	case 0xff16:
		return audio.Pulse2LengthDuty
	case 0xff17:
		return audio.Pulse2VolumeEnvelope
	case 0xff18:
		return audio.Pulse2PeriodLow
	case 0xff19:
		return audio.Pulse2PeriodHighCtl
	case 0xff1a:
		return audio.WaveDACEn
	case 0xff1b:
		return audio.WaveLengthTimer
	case 0xff1c:
		return audio.WaveOutputLevel
	case 0xff1d:
		return audio.WavePeriodLow
	case 0xff1e:
		return audio.WavePeriodHighCtl
	case 0xff1f:
		return 0
	case 0xff20:
		return audio.NoiseLengthTimer
	case 0xff21:
		return audio.NoiseVolumeEnvelope
	case 0xff22:
		return audio.NoiseRNG
	case 0xff23:
		return audio.NoiseCtl
	case 0xff24:
		return audio.MasterVolumePan
	case 0xff25:
		return audio.ChannelPan
	case 0xff26:
		return audio.MasterCtl
	}
	panicv(addr)
	return 0
}

func (audio *Audio) Write(addr uint16, v uint8) {
	switch addr {
	case 0xff10:
		audio.SetPulse1Sweep(v)
	case 0xff11:
		audio.SetPulse1LengthDuty(v)
	case 0xff12:
		audio.SetPulse1VolumeEnvelope(v)
	case 0xff13:
		audio.SetPulse1PeriodLow(v) // WO
	case 0xff14:
		audio.SetPulse1PeriodHighCtl(v)
	case 0xff15:
	case 0xff16:
		audio.SetPulse2LengthDuty(v)
	case 0xff17:
		audio.SetPulse2VolumeEnvelope(v)
	case 0xff18:
		audio.SetPulse2PeriodLow(v)
	case 0xff19:
		audio.SetPulse2PeriodHighCtl(v)
	case 0xff1a:
		audio.SetWaveDACEn(v)
	case 0xff1b:
		audio.SetWaveLengthTimer(v)
	case 0xff1c:
		audio.SetWaveOutputLevel(v)
	case 0xff1d:
		audio.SetWavePeriodLow(v)
	case 0xff1e:
		audio.SetWavePeriodHighCtl(v)
	case 0xff1f:
	case 0xff20:
		audio.SetNoiseLengthTimer(v)
	case 0xff21:
		audio.SetNoiseVolumeEnvelope(v)
	case 0xff22:
		audio.SetNoiseRNG(v)
	case 0xff23:
		audio.SetNoiseCtl(v)
	case 0xff24:
		audio.SetMasterVolumePan(v)
	case 0xff25:
		audio.SetChannelPan(v)
	case 0xff26:
		audio.SetMasterCtl(v)
	}
}

func (audio *Audio) SetPulse1Sweep(v uint8) {
	audio.Debug("SetPulse1Sweep", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1Sweep = v
	panic("not implemented: SetPulse1Sweep")
}

func (audio *Audio) SetPulse1LengthDuty(v uint8) {
	audio.Debug("SetPulse1LengthDuty", "%v\n", v)
	if audio.canWriteLengthTimersWithAPUOff || !audio.APUEnabled() {
		return
	}
	audio.Pulse1LengthDuty = v
	panic("not implemented: SetPulse1LengthDuty")
}

func (audio *Audio) SetPulse1VolumeEnvelope(v uint8) {
	audio.Debug("SetPulse1VolumeEnvelope", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1VolumeEnvelope = v
	panic("not implemented: SetPulse1VolumeEnvelope")
}

func (audio *Audio) SetPulse1PeriodLow(v uint8) {
	audio.Debug("SetPulse1PeriodLow", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1PeriodLow = v
	panic("not implemented: SetPulse1PeriodLow")
}

func (audio *Audio) SetPulse1PeriodHighCtl(v uint8) {
	audio.Debug("SetPulse1PeriodHighCtl", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1PeriodHighCtl = v
	panic("not implemented: SetPulse1PeriodHighCtl")
}

func (audio *Audio) SetPulse2LengthDuty(v uint8) {
	audio.Debug("SetPulse2LengthDuty", "%v\n", v)
	if audio.canWriteLengthTimersWithAPUOff || !audio.APUEnabled() {
		return
	}
	audio.Pulse2LengthDuty = v
	panic("not implemented: SetPulse2LengthDuty")
}

func (audio *Audio) SetPulse2VolumeEnvelope(v uint8) {
	audio.Debug("SetPulse2VolumeEnvelope", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse2VolumeEnvelope = v
	panic("not implemented: SetPulse2VolumeEnvelope")
}

func (audio *Audio) SetPulse2PeriodLow(v uint8) {
	audio.Debug("SetPulse2PeriodLow", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse2PeriodLow = v
	panic("not implemented: SetPulse2PeriodLow")
}

func (audio *Audio) SetPulse2PeriodHighCtl(v uint8) {
	audio.Debug("SetPulse2PeriodHighCtl", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse2PeriodHighCtl = v
	panic("not implemented: SetPulse2PeriodHighCtl")
}

func (audio *Audio) SetWaveDACEn(v uint8) {
	audio.Debug("SetWaveDACEn", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.WaveDACEn = v
	panic("not implemented: SetWaveDACEn")
}

func (audio *Audio) SetWaveLengthTimer(v uint8) {
	audio.Debug("SetWaveLengthTimer", "%v\n", v)
	if audio.canWriteLengthTimersWithAPUOff || !audio.APUEnabled() {
		return
	}
	audio.WaveLengthTimer = v
	panic("not implemented: SetWaveLengthTimer")
}

func (audio *Audio) SetWaveOutputLevel(v uint8) {
	audio.Debug("SetWaveOutputLevel", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.WaveOutputLevel = v
	panic("not implemented: SetWaveOutputLevel")
}

func (audio *Audio) SetWavePeriodLow(v uint8) {
	audio.Debug("SetWavePeriodLow", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.WavePeriodLow = v
	panic("not implemented: SetWavePeriodLow")
}

func (audio *Audio) SetWavePeriodHighCtl(v uint8) {
	audio.Debug("SetWavePeriodHighCtl", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.WavePeriodHighCtl = v
	panic("not implemented: SetWavePeriodHighCtl")
}

func (audio *Audio) SetNoiseLengthTimer(v uint8) {
	audio.Debug("SetNoiseLengthTimer", "%v\n", v)
	if audio.canWriteLengthTimersWithAPUOff || !audio.APUEnabled() {
		return
	}
	audio.NoiseLengthTimer = v
	panic("not implemented: SetNoiseLengthTimer")
}

func (audio *Audio) SetNoiseVolumeEnvelope(v uint8) {
	audio.Debug("SetNoiseVolumeEnvelope", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.NoiseVolumeEnvelope = v
	panic("not implemented: SetNoiseVolumeEnvelope")
}

func (audio *Audio) SetNoiseRNG(v uint8) {
	audio.Debug("SetNoiseRNG", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.NoiseRNG = v
	panic("not implemented: SetNoiseRNG")
}

func (audio *Audio) SetNoiseCtl(v uint8) {
	audio.Debug("SetNoiseCtl", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.NoiseCtl = v
	panic("not implemented: SetNoiseCtl")
}

func (audio *Audio) SetMasterVolumePan(v uint8) {
	audio.Debug("SetMasterVolumePan", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.MasterVolumePan = v
	panic("not implemented: SetMasterVolumePan")
}

func (audio *Audio) SetChannelPan(v uint8) {
	audio.Debug("SetChannelPan", "%v\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.ChannelPan = v
	panic("not implemented: SetChannelPan")
}

func (audio *Audio) SetMasterCtl(v uint8) {
	audio.Debug("SetMasterCtl", "%v\n", v)

	// only bit 7 is R/W
	audio.MasterCtl = maskedWrite(audio.MasterCtl, v, 0x80)

	// Turning the APU off clears all APU registers
	if !audio.APUEnabled() {
		audio.Pulse1Sweep = 0
		audio.Pulse1LengthDuty = 0
		audio.Pulse1VolumeEnvelope = 0
		audio.Pulse1PeriodLow = 0
		audio.Pulse1PeriodHighCtl = 0
		audio.Pulse2LengthDuty = 0
		audio.Pulse2VolumeEnvelope = 0
		audio.Pulse2PeriodLow = 0
		audio.Pulse2PeriodHighCtl = 0
		audio.WaveDACEn = 0
		audio.WaveLengthTimer = 0
		audio.WaveOutputLevel = 0
		audio.WavePeriodLow = 0
		audio.WavePeriodHighCtl = 0
		audio.NoiseLengthTimer = 0
		audio.NoiseVolumeEnvelope = 0
		audio.NoiseRNG = 0
		audio.NoiseCtl = 0
		audio.MasterVolumePan = 0
		audio.ChannelPan = 0
		audio.MasterCtl = 0
	}
}

func maskedWrite(prev, v, mask uint8) uint8 {
	return (prev & (^mask)) | (v & mask)
}
