package model

import (
	"fmt"
	"slices"
)

var audioDebugEvents = []string{
	"SetPulse1Sweep",
	"SetPulse1LengthDuty",
	"SetPulse1VolumeEnvelope",
	"SetPulse1PeriodLow",
	"SetPulse1PeriodHighCtl",
	"SetPulse2LengthDuty",
	"SetPulse2VolumeEnvelope",
	"SetPulse2PeriodLow",
	"SetPulse2PeriodHighCtl",
	"SetWaveDACEn",
	"SetWaveLengthTimer",
	"SetWaveOutputLevel",
	"SetWavePeriodLow",
	"SetWavePeriodHighCtl",
	"SetNoiseLengthTimer",
	"SetNoiseVolumeEnvelope",
	"SetNoiseRNG",
	"SetNoiseCtl",
	"SetMasterVolumePan",
	"SetChannelPan",
	"SetMasterCtl",
}

type AudioCtl struct {
	MasterCtl uint8

	Pulse1 PulseChannelWithSweep
	Pulse2 PulseChannel
	Wave   WaveChannel
	Noise  NoiseChannel

	Mixer Mixer

	canWriteLengthTimersWithAPUOff bool
}

type Mixer struct {
	RegChannelPan         uint8
	RegMasterVolumeVINPan uint8

	leftVIN        bool
	rightVIN       bool
	leftVol        uint8
	rightVol       uint8
	leftEnChannel  [4]bool
	rightEnChannel [4]bool
}

func (mixer *Mixer) SetChannelPan(v uint8) {
	mixer.RegChannelPan = v
	mixer.leftEnChannel[0] = v&(1<<0) != 0
	mixer.leftEnChannel[1] = v&(1<<1) != 0
	mixer.leftEnChannel[2] = v&(1<<2) != 0
	mixer.leftEnChannel[3] = v&(1<<3) != 0
	mixer.rightEnChannel[0] = v&(1<<4) != 0
	mixer.rightEnChannel[1] = v&(1<<5) != 0
	mixer.rightEnChannel[2] = v&(1<<6) != 0
	mixer.rightEnChannel[3] = v&(1<<7) != 0
}

func (mixer *Mixer) SetMasterVolumeVINPan(v uint8) {
	mixer.RegMasterVolumeVINPan = v
	mixer.rightVol = v & 0x7
	mixer.rightVIN = (v>>3)&0x1 != 0
	mixer.leftVol = (v >> 4) & 0x7
	mixer.leftVIN = (v>>7)&0x1 != 0
}

type PulseChannel struct {
	RegLengthDuty     uint8
	RegVolumeEnvelope uint8
	RegPeriodLow      uint8
	RegPeriodHighCtl  uint8

	waveDuty     uint8
	lengthTimer  uint8
	initVolume   uint8
	envDir       uint8
	envSweepPace uint8
}

func (pc *PulseChannel) SetLengthDuty(v uint8) {
	pc.RegLengthDuty = v

	pc.lengthTimer = v & 0x3f
	pc.waveDuty = (v >> 5) & 0x3
}

func (pc *PulseChannel) SetVolumeEnvelope(v uint8) {
	pc.RegVolumeEnvelope = v

	pc.envSweepPace = v & 0x7
	pc.envDir = v >> 3 & 0x1
	pc.initVolume = (v >> 4) & 0xf

	// TODO "Setting bits 3-7 of this register all to 0 (initial volume = 0, envelope = decreasing) turns the DAC off (and thus, the channel as well), which may cause an audio pop."
}

func (pc *PulseChannel) SetPeriodLow(v uint8) {
	pc.RegPeriodLow = v
	panic("not implemented: SetPulse1PeriodLow")
}

func (pc *PulseChannel) SetPeriodHighCtl(v uint8) {
	pc.RegPeriodHighCtl = v
	panic("not implemented: SetPeriodHighCtl")
}

type PulseChannelWithSweep struct {
	PulseChannel
	RegSweep uint8
}

func (pc *PulseChannelWithSweep) SetSweep(v uint8) {
	pc.RegSweep = v
}

type WaveChannel struct {
	RegDACEn         uint8
	RegLengthTimer   uint8
	RegOutputLevel   uint8
	RegPeriodLow     uint8
	RegPeriodHighCtl uint8
}

func (wc *WaveChannel) SetDACEn(v uint8) {
	wc.RegDACEn = v
	panic("not implemented: SetDACEn")
}

func (wc *WaveChannel) SetLengthTimer(v uint8) {
	wc.RegLengthTimer = v
	panic("not implemented: SetLengthTimer")
}

func (wc *WaveChannel) SetOutputLevel(v uint8) {
	wc.RegOutputLevel = v
	panic("not implemented: SetOutputLevel")
}

func (wc *WaveChannel) SetPeriodLow(v uint8) {
	wc.RegPeriodLow = v
	panic("not implemented: SetPeriodLow")
}

func (wc *WaveChannel) SetPeriodHighCtl(v uint8) {
	wc.RegPeriodHighCtl = v
	panic("not implemented: SetPeriodHighCtl")
}

type NoiseChannel struct {
	RegLengthTimer    uint8
	RegVolumeEnvelope uint8
	RegRNG            uint8
	RegCtl            uint8
}

func (nc *NoiseChannel) SetLengthTimer(v uint8) {
	nc.RegLengthTimer = v
	panic("not implemented: SetLengthTimer")
}

func (nc *NoiseChannel) SetVolumeEnvelope(v uint8) {
	nc.RegVolumeEnvelope = v
	panic("not implemented: SetVolumeEnvelope")
}

func (nc *NoiseChannel) SetRNG(v uint8) {
	nc.RegRNG = v
	panic("not implemented: SetRNG")
}

func (nc *NoiseChannel) SetCtl(v uint8) {
	nc.RegCtl = v
	panic("not implemented: SetCtl")
}

func NewAudioCtl() *AudioCtl {
	return &AudioCtl{
		canWriteLengthTimersWithAPUOff: true, // on monochrome models
	}
}

func (audio *AudioCtl) Debug(event string, f string, v ...any) {
	if !slices.Contains(audioDebugEvents, event) {
		return
	}
	fmt.Printf("AUDIO | %s | ", event)
	fmt.Printf(f, v...)
}

func (audio *AudioCtl) APUEnabled() bool {
	return audio.MasterCtl&0x80 != 0
}

func (audio *AudioCtl) Name() string {
	return "AUDIO"
}

func (audio *AudioCtl) Range() (uint16, uint16) {
	return 0xff10, 0x0017
}

func (audio *AudioCtl) Read(addr uint16) uint8 {
	switch addr {
	case 0xff10:
		return audio.Pulse1.RegSweep
	case 0xff11:
		return audio.Pulse1.RegLengthDuty
	case 0xff12:
		return audio.Pulse1.RegVolumeEnvelope
	case 0xff13:
		return audio.Pulse1.RegPeriodLow // WO
	case 0xff14:
		return audio.Pulse1.RegPeriodHighCtl
	case 0xff15:
		return 0
	case 0xff16:
		return audio.Pulse2.RegLengthDuty
	case 0xff17:
		return audio.Pulse2.RegVolumeEnvelope
	case 0xff18:
		return audio.Pulse2.RegPeriodLow
	case 0xff19:
		return audio.Pulse2.RegPeriodHighCtl
	case 0xff1a:
		return audio.Wave.RegDACEn
	case 0xff1b:
		return audio.Wave.RegLengthTimer
	case 0xff1c:
		return audio.Wave.RegOutputLevel
	case 0xff1d:
		return audio.Wave.RegPeriodLow
	case 0xff1e:
		return audio.Wave.RegPeriodHighCtl
	case 0xff1f:
		return 0
	case 0xff20:
		return audio.Noise.RegLengthTimer
	case 0xff21:
		return audio.Noise.RegVolumeEnvelope
	case 0xff22:
		return audio.Noise.RegRNG
	case 0xff23:
		return audio.Noise.RegCtl
	case 0xff24:
		return audio.Mixer.RegMasterVolumeVINPan
	case 0xff25:
		return audio.Mixer.RegChannelPan
	case 0xff26:
		return audio.MasterCtl
	}
	panicf("Read from unknown audio register %#v", addr)
	return 0
}

func (audio *AudioCtl) Write(addr uint16, v uint8) {
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
	default:
		panicf("Write to unknown audio register %#v", addr)
	}
}

func (audio *AudioCtl) SetPulse1Sweep(v uint8) {
	audio.Debug("SetPulse1Sweep", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1.SetSweep(v)
}

func (audio *AudioCtl) SetPulse1LengthDuty(v uint8) {
	audio.Debug("SetPulse1LengthDuty", "0x%02x\n", v)
	if !audio.canWriteLengthTimersWithAPUOff && !audio.APUEnabled() {
		return
	}
	audio.Pulse1.SetLengthDuty(v)
}

func (audio *AudioCtl) SetPulse1VolumeEnvelope(v uint8) {
	audio.Debug("SetPulse1VolumeEnvelope", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1.SetVolumeEnvelope(v)
}

func (audio *AudioCtl) SetPulse1PeriodLow(v uint8) {
	audio.Debug("SetPulse1PeriodLow", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1.SetPeriodLow(v)
}

func (audio *AudioCtl) SetPulse1PeriodHighCtl(v uint8) {
	audio.Debug("SetPulse1PeriodHighCtl", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse1.SetPeriodHighCtl(v)
}

func (audio *AudioCtl) SetPulse2LengthDuty(v uint8) {
	audio.Debug("SetPulse2LengthDuty", "0x%02x\n", v)
	if !audio.canWriteLengthTimersWithAPUOff && !audio.APUEnabled() {
		return
	}
	audio.Pulse2.SetLengthDuty(v)
}

func (audio *AudioCtl) SetPulse2VolumeEnvelope(v uint8) {
	audio.Debug("SetPulse2VolumeEnvelope", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse2.SetVolumeEnvelope(v)
}

func (audio *AudioCtl) SetPulse2PeriodLow(v uint8) {
	audio.Debug("SetPulse2PeriodLow", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse2.SetPeriodLow(v)
}

func (audio *AudioCtl) SetPulse2PeriodHighCtl(v uint8) {
	audio.Debug("SetPulse2PeriodHighCtl", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Pulse2.SetPeriodHighCtl(v)
}

func (audio *AudioCtl) SetWaveDACEn(v uint8) {
	audio.Debug("SetWaveDACEn", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Wave.SetDACEn(v)
}

func (audio *AudioCtl) SetWaveLengthTimer(v uint8) {
	audio.Debug("SetWaveLengthTimer", "0x%02x\n", v)
	if !audio.canWriteLengthTimersWithAPUOff && !audio.APUEnabled() {
		return
	}
	audio.Wave.SetLengthTimer(v)
}

func (audio *AudioCtl) SetWaveOutputLevel(v uint8) {
	audio.Debug("SetWaveOutputLevel", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Wave.SetOutputLevel(v)
}

func (audio *AudioCtl) SetWavePeriodLow(v uint8) {
	audio.Debug("SetWavePeriodLow", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Wave.SetPeriodLow(v)
}

func (audio *AudioCtl) SetWavePeriodHighCtl(v uint8) {
	audio.Debug("SetWavePeriodHighCtl", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Wave.SetPeriodHighCtl(v)
}

func (audio *AudioCtl) SetNoiseLengthTimer(v uint8) {
	audio.Debug("SetNoiseLengthTimer", "0x%02x\n", v)
	if !audio.canWriteLengthTimersWithAPUOff && !audio.APUEnabled() {
		return
	}
	audio.Noise.SetLengthTimer(v)
}

func (audio *AudioCtl) SetNoiseVolumeEnvelope(v uint8) {
	audio.Debug("SetNoiseVolumeEnvelope", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Noise.SetVolumeEnvelope(v)
}

func (audio *AudioCtl) SetNoiseRNG(v uint8) {
	audio.Debug("SetNoiseRNG", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Noise.SetRNG(v)
}

func (audio *AudioCtl) SetNoiseCtl(v uint8) {
	audio.Debug("SetNoiseCtl", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Noise.SetCtl(v)
}

func (audio *AudioCtl) SetMasterVolumePan(v uint8) {
	audio.Debug("SetMasterVolumePan", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Mixer.SetMasterVolumeVINPan(v)
}

func (audio *AudioCtl) SetChannelPan(v uint8) {
	audio.Debug("SetChannelPan", "0x%02x\n", v)
	if !audio.APUEnabled() {
		return
	}
	audio.Mixer.SetChannelPan(v)
}

func (audio *AudioCtl) SetMasterCtl(v uint8) {
	audio.Debug("SetMasterCtl", "0x%02x\n", v)

	// only bit 7 is R/W
	audio.MasterCtl = maskedWrite(audio.MasterCtl, v, 0x80)

	// Turning the APU off clears all APU registers
	if !audio.APUEnabled() {
		audio.Pulse1.SetSweep(0)
		audio.Pulse1.SetLengthDuty(0)
		audio.Pulse1.SetVolumeEnvelope(0)
		audio.Pulse1.SetPeriodLow(0)
		audio.Pulse1.SetPeriodHighCtl(0)
		audio.Pulse2.SetLengthDuty(0)
		audio.Pulse2.SetVolumeEnvelope(0)
		audio.Pulse2.SetPeriodLow(0)
		audio.Pulse2.SetPeriodHighCtl(0)
		audio.Wave.SetDACEn(0)
		audio.Wave.SetLengthTimer(0)
		audio.Wave.SetOutputLevel(0)
		audio.Wave.SetPeriodLow(0)
		audio.Wave.SetPeriodHighCtl(0)
		audio.Noise.SetLengthTimer(0)
		audio.Noise.SetVolumeEnvelope(0)
		audio.Noise.SetRNG(0)
		audio.Noise.SetCtl(0)
		audio.Mixer.SetMasterVolumeVINPan(0)
		audio.Mixer.SetChannelPan(0)
		audio.MasterCtl = 0
	}
}

func maskedWrite(prev, v, mask uint8) uint8 {
	return (prev & (^mask)) | (v & mask)
}
