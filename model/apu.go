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

type APU struct {
	MemoryRegion
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

	period uint16

	lengthEnable bool
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

	// TODO "Setting bits 3-7 of this register all to 0 (initial volume = 0, envelope = decreasing) turns the DAC off (and thus, the channel as well), which may cause an apu pop."
}

func (pc *PulseChannel) SetPeriodLow(v uint8) {
	pc.RegPeriodLow = v

	// keep upper 3 bits, overwrite lower 8
	pc.period &= uint16(0x0700)
	pc.period |= uint16(v)
}

func (pc *PulseChannel) SetPeriodHighCtl(v uint8) {
	pc.RegPeriodHighCtl = v

	// keep lower 8 bits, overwrite upper 3
	pc.period &= uint16(0x00ff)
	pc.period |= uint16(v&0x7) << 8

	pc.lengthEnable = v&0x40 != 0
	if v&0x80 != 0 {
		pc.trigger()
	}
}

func (pc *PulseChannel) trigger() {
	fmt.Printf("Trigger not implemented\n")
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
	fmt.Printf("not implemented: SetDACEn\n")
}

func (wc *WaveChannel) SetLengthTimer(v uint8) {
	wc.RegLengthTimer = v
	fmt.Printf("not implemented: SetLengthTimer\n")
}

func (wc *WaveChannel) SetOutputLevel(v uint8) {
	wc.RegOutputLevel = v
	fmt.Printf("not implemented: SetOutputLevel\n")
}

func (wc *WaveChannel) SetPeriodLow(v uint8) {
	wc.RegPeriodLow = v
	fmt.Printf("not implemented: SetPeriodLow\n")
}

func (wc *WaveChannel) SetPeriodHighCtl(v uint8) {
	wc.RegPeriodHighCtl = v
	fmt.Printf("not implemented: SetPeriodHighCtl\n")
}

type NoiseChannel struct {
	RegLengthTimer    uint8
	RegVolumeEnvelope uint8
	RegRNG            uint8
	RegCtl            uint8
}

func (nc *NoiseChannel) SetLengthTimer(v uint8) {
	nc.RegLengthTimer = v
	fmt.Printf("not implemented: SetLengthTimer\n")
}

func (nc *NoiseChannel) SetVolumeEnvelope(v uint8) {
	nc.RegVolumeEnvelope = v
	fmt.Printf("not implemented: SetVolumeEnvelope\n")
}

func (nc *NoiseChannel) SetRNG(v uint8) {
	nc.RegRNG = v
	fmt.Printf("not implemented: SetRNG\n")
}

func (nc *NoiseChannel) SetCtl(v uint8) {
	nc.RegCtl = v
	fmt.Printf("not implemented: SetCtl\n")
}

func NewAPU(clock *ClockRT) *APU {
	return &APU{
		MemoryRegion:                   NewMemoryRegion(clock, AddrAPUBegin, AddrAPUEnd),
		canWriteLengthTimersWithAPUOff: true, // on monochrome models
	}
}

func (apu *APU) Debug(event string, f string, v ...any) {
	if !slices.Contains(audioDebugEvents, event) {
		return
	}
	fmt.Printf("AUDIO | %s | ", event)
	fmt.Printf(f, v...)
}

func (apu *APU) Enabled() bool {
	return apu.MasterCtl&0x80 != 0
}

func (apu *APU) Read(addr uint16) uint8 {
	_ = apu.MemoryRegion.Read(addr)

	switch Addr(addr) {
	case AddrNR10:
		return apu.Pulse1.RegSweep
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
		return apu.MasterCtl
	}
	panicf("Read from unknown apu register %#v", addr)
	return 0
}

func (apu *APU) Write(addr uint16, v uint8) {
	apu.MemoryRegion.Write(addr, v)

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

func (apu *APU) SetPulse1Sweep(v uint8) {
	apu.Debug("SetPulse1Sweep", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Pulse1.SetSweep(v)
}

func (apu *APU) SetPulse1LengthDuty(v uint8) {
	apu.Debug("SetPulse1LengthDuty", "0x%02x\n", v)
	if !apu.canWriteLengthTimersWithAPUOff && !apu.Enabled() {
		return
	}
	apu.Pulse1.SetLengthDuty(v)
}

func (apu *APU) SetPulse1VolumeEnvelope(v uint8) {
	apu.Debug("SetPulse1VolumeEnvelope", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Pulse1.SetVolumeEnvelope(v)
}

func (apu *APU) SetPulse1PeriodLow(v uint8) {
	apu.Debug("SetPulse1PeriodLow", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Pulse1.SetPeriodLow(v)
}

func (apu *APU) SetPulse1PeriodHighCtl(v uint8) {
	apu.Debug("SetPulse1PeriodHighCtl", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Pulse1.SetPeriodHighCtl(v)
}

func (apu *APU) SetPulse2LengthDuty(v uint8) {
	apu.Debug("SetPulse2LengthDuty", "0x%02x\n", v)
	if !apu.canWriteLengthTimersWithAPUOff && !apu.Enabled() {
		return
	}
	apu.Pulse2.SetLengthDuty(v)
}

func (apu *APU) SetPulse2VolumeEnvelope(v uint8) {
	apu.Debug("SetPulse2VolumeEnvelope", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Pulse2.SetVolumeEnvelope(v)
}

func (apu *APU) SetPulse2PeriodLow(v uint8) {
	apu.Debug("SetPulse2PeriodLow", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Pulse2.SetPeriodLow(v)
}

func (apu *APU) SetPulse2PeriodHighCtl(v uint8) {
	apu.Debug("SetPulse2PeriodHighCtl", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Pulse2.SetPeriodHighCtl(v)
}

func (apu *APU) SetWaveDACEn(v uint8) {
	apu.Debug("SetWaveDACEn", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Wave.SetDACEn(v)
}

func (apu *APU) SetWaveLengthTimer(v uint8) {
	apu.Debug("SetWaveLengthTimer", "0x%02x\n", v)
	if !apu.canWriteLengthTimersWithAPUOff && !apu.Enabled() {
		return
	}
	apu.Wave.SetLengthTimer(v)
}

func (apu *APU) SetWaveOutputLevel(v uint8) {
	apu.Debug("SetWaveOutputLevel", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Wave.SetOutputLevel(v)
}

func (apu *APU) SetWavePeriodLow(v uint8) {
	apu.Debug("SetWavePeriodLow", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Wave.SetPeriodLow(v)
}

func (apu *APU) SetWavePeriodHighCtl(v uint8) {
	apu.Debug("SetWavePeriodHighCtl", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Wave.SetPeriodHighCtl(v)
}

func (apu *APU) SetNoiseLengthTimer(v uint8) {
	apu.Debug("SetNoiseLengthTimer", "0x%02x\n", v)
	if !apu.canWriteLengthTimersWithAPUOff && !apu.Enabled() {
		return
	}
	apu.Noise.SetLengthTimer(v)
}

func (apu *APU) SetNoiseVolumeEnvelope(v uint8) {
	apu.Debug("SetNoiseVolumeEnvelope", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Noise.SetVolumeEnvelope(v)
}

func (apu *APU) SetNoiseRNG(v uint8) {
	apu.Debug("SetNoiseRNG", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Noise.SetRNG(v)
}

func (apu *APU) SetNoiseCtl(v uint8) {
	apu.Debug("SetNoiseCtl", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Noise.SetCtl(v)
}

func (apu *APU) SetMasterVolumePan(v uint8) {
	apu.Debug("SetMasterVolumePan", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Mixer.SetMasterVolumeVINPan(v)
}

func (apu *APU) SetChannelPan(v uint8) {
	apu.Debug("SetChannelPan", "0x%02x\n", v)
	if !apu.Enabled() {
		return
	}
	apu.Mixer.SetChannelPan(v)
}

func (apu *APU) SetMasterCtl(v uint8) {
	apu.Debug("SetMasterCtl", "0x%02x\n", v)

	// only bit 7 is R/W
	apu.MasterCtl = maskedWrite(apu.MasterCtl, v, 0x80)

	// Turning the APU off clears all APU registers
	if !apu.Enabled() {
		apu.Pulse1.SetSweep(0)
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

func maskedWrite(prev, v, mask uint8) uint8 {
	return (prev & (^mask)) | (v & mask)
}
