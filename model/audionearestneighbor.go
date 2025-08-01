package model

import "time"

const (
	FracBits = 8
	Scale    = 1 << FracBits
	HPFactor = int32(255) // most relaxed hp filter possible
	LPFactor = int32(128) // most relaxed lp filter possible
)

type AudioNN struct {
	SampleBuffers  SampleBuffers
	SampleInterval time.Duration
	SampleDivider  int
	SubSampling    int
	MCounter       int
	Capacitor      AudioSample
	Out            chan []AudioSample
}

type SampleBuffers struct {
	Left  []AudioSample
	Right []AudioSample
	Size  int
	Idx   int
}

func NewSampleBuffers(size int) SampleBuffers {
	return SampleBuffers{
		Left:  make([]AudioSample, size),
		Right: make([]AudioSample, size),
		Size:  size,
	}
}

func (ab *SampleBuffers) Add(l, r AudioSample) bool {
	ab.Left[ab.Idx] = l
	ab.Right[ab.Idx] = r

	ab.Idx++
	if ab.Idx == ab.Size {
		ab.Idx = 0
		return true
	}
	return false
}

func (audio *AudioNN) SetMPeriod(mPeriod time.Duration) {
	if mPeriod > 0 {
		audio.SampleDivider = int(audio.SampleInterval / mPeriod)
		audio.MCounter = audio.SampleDivider * audio.SubSampling
	} else {
		audio.SampleDivider = 0
	}
}

func (audio *AudioNN) Clock(apu *APU) {
	if audio.SampleDivider <= 0 {
		return
	}
	audio.MCounter -= audio.SubSampling
	if audio.MCounter > 0 {
		return
	}
	audio.MCounter = audio.SampleDivider*audio.SubSampling - audio.MCounter
	if !audio.SampleBuffers.Add(
		apu.Mixer.MixStereoSimple(
			apu.Pulse1.Sample(),
			apu.Pulse2.Sample(),
			apu.Noise.Sample(),
			apu.Wave.Sample(),
		),
	) {
		return
	}

	mono := make([]AudioSample, len(audio.SampleBuffers.Left))
	for i := range audio.SampleBuffers.Left {
		mono[i] = (audio.SampleBuffers.Left[i] + audio.SampleBuffers.Right[i])
	}
	highpass(mono, &audio.Capacitor)
	audio.Out <- mono
}

func highpass(audio []AudioSample, capacitor *AudioSample) {
	for i := range audio {
		in := audio[i] << FracBits
		out := in - *capacitor
		*capacitor = in - AudioSample((int32(out)*HPFactor)>>FracBits)
		audio[i] = out >> FracBits
	}
}
