package model

import (
	"fmt"
	"time"

	"github.com/jonathangjertsen/toyboy/blip"
)

type AudioBLIP struct {
	buffer         *blip.Buffer
	synthPC1       *blip.Synth
	synthPC2       *blip.Synth
	synthW         *blip.Synth
	synthN         *blip.Synth
	config         *blip.Config
	clock          int
	bufferSize     int
	n              int
	Out            chan []AudioSample
	prevPC1        AudioSample
	prevPC2        AudioSample
	prevWave       AudioSample
	prevNoise      AudioSample
	SampleInterval time.Duration
	SampleDivider  int
	MCounter       int
	SCounter       int
	cps            int
}

func NewAudioBLIP(size int, out chan []AudioSample) *AudioBLIP {
	config := blip.DefaultBlipConfig
	buffer := blip.NewBuffer(config)
	buffer.SetSamplingParams(config.InitialSampleRate, 1024*1000/44100)
	return &AudioBLIP{
		buffer:         buffer,
		synthPC1:       blip.NewSynth(buffer),
		synthPC2:       blip.NewSynth(buffer),
		synthW:         blip.NewSynth(buffer),
		synthN:         blip.NewSynth(buffer),
		config:         &config,
		clock:          0,
		bufferSize:     size,
		n:              0,
		Out:            out,
		SampleInterval: time.Second / time.Duration(config.InitialSampleRate),
	}
}

func (audio *AudioBLIP) SetMPeriod(mPeriod time.Duration) {
	if mPeriod > 0 {
		audio.cps = int(time.Second + (mPeriod/2)/mPeriod)
		audio.buffer.SetClockRate(audio.cps)
		audio.SampleDivider = int(audio.SampleInterval / mPeriod)
		audio.MCounter = audio.SampleDivider
		audio.SCounter = audio.bufferSize
	} else {
		audio.buffer.SetClockRate(0)
	}
}

func (audio *AudioBLIP) Clock(apu *APU) {
	newPulse1 := apu.Pulse1.Sample()
	newPulse2 := apu.Pulse2.Sample()
	newNoise := apu.Noise.Sample()
	newWave := apu.Wave.Sample()

	if newPulse1 != audio.prevPC1 {
		audio.synthPC1.Update(audio.clock, int(newPulse1))
		audio.prevPC1 = newPulse1
	}
	if newPulse2 != audio.prevPC2 {
		audio.synthPC2.Update(audio.clock, int(newPulse2))
		audio.prevPC2 = newPulse2
	}
	if newNoise != audio.prevNoise {
		audio.synthN.Update(audio.clock, int(newNoise))
		audio.prevNoise = newNoise
	}
	if newWave != audio.prevWave {
		audio.synthW.Update(audio.clock, int(newWave))
		audio.prevWave = newWave
	}

	audio.clock++
	audio.MCounter--
	if audio.MCounter > 0 {
		return
	}
	audio.MCounter = audio.SampleDivider

	audio.SCounter--
	if audio.SCounter > 0 {
		return
	}
	audio.SCounter = audio.bufferSize
	dt := (int(audio.SampleInterval) * audio.bufferSize) / int(time.Millisecond)

	audio.buffer.EndFrame(audio.bufferSize * audio.SampleDivider)
	out := make([]AudioSample, audio.bufferSize)
	n := audio.buffer.Read(out)
	fmt.Printf("read %d of %d sd=%d dt=%d avail=%d\n", n, audio.bufferSize, audio.SampleDivider, dt, audio.buffer.SamplesAvailable())
	audio.Out <- out
}
