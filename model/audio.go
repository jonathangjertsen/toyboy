package model

import (
	"time"
)

type AudioSample int16

const (
	FracBits = 8
	Scale    = 1 << FracBits
	HPFactor = int32(255) // most relaxed hp filter possible
	LPFactor = int32(128) // most relaxed lp filter possible
)

type Audio struct {
	APU            *APU
	SampleBuffers  SampleBuffers
	SampleInterval time.Duration
	SampleDivider  int
	MCounter       int
	Capacitor      AudioSample
	Out            chan []AudioSample
}

type SampleBuffers struct {
	Pulse1 []AudioSample
	Pulse2 []AudioSample
	Wave   []AudioSample
	Noise  []AudioSample
	Size   int
	Idx    int
}

func (sb *SampleBuffers) Mix() []AudioSample {
	mixed := make([]AudioSample, len(sb.Pulse1))
	for i := range len(mixed) {
		mixed[i] = sb.Pulse1[i] + sb.Pulse2[i] + sb.Wave[i] + sb.Noise[i]
	}
	return mixed
}

func NewSampleBuffers(size int) SampleBuffers {
	return SampleBuffers{
		Pulse1: make([]AudioSample, size),
		Pulse2: make([]AudioSample, size),
		Wave:   make([]AudioSample, size),
		Noise:  make([]AudioSample, size),
		Size:   size,
	}
}

func (ab *SampleBuffers) Add(pulse1, pulse2, wave, noise AudioSample) bool {
	ab.Pulse1[ab.Idx] = pulse1
	ab.Pulse2[ab.Idx] = pulse2
	ab.Wave[ab.Idx] = wave
	ab.Noise[ab.Idx] = noise

	ab.Idx++
	if ab.Idx == ab.Size {
		ab.Idx = 0
		return true
	}
	return false
}

func (audio *Audio) Enabled() bool {
	return audio.APU != nil && audio.SampleDivider > 0
}

func (audio *Audio) SetMPeriod(mPeriod time.Duration) {
	if mPeriod > 0 {
		audio.SampleDivider = int(audio.SampleInterval / mPeriod)
		audio.MCounter = audio.SampleDivider
	} else {
		audio.SampleDivider = 0
	}
}

func (audio *Audio) Clock() {
	if !audio.Enabled() {
		return
	}
	audio.MCounter--
	if audio.MCounter > 0 {
		return
	}
	audio.MCounter = audio.SampleDivider
	if !audio.SampleBuffers.Add(
		audio.APU.Pulse1.Sample(),
		audio.APU.Pulse2.Sample(),
		audio.APU.Wave.Sample(),
		audio.APU.Noise.Sample(),
	) {
		return
	}

	mix := audio.SampleBuffers.Mix()
	highpass(mix, &audio.Capacitor)
	/*select {
	case audio.Out <- mix:
	default:
		fmt.Printf("MISSED AUDIO BUFFER\n")
	}
	*/
}

func highpass(audio []AudioSample, capacitor *AudioSample) {
	for i := range audio {
		in := audio[i] << FracBits
		out := in - *capacitor
		*capacitor = in - AudioSample((int32(out)*HPFactor)>>FracBits)
		audio[i] = out >> FracBits
	}
}
