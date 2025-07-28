package model

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

type Audio struct {
	APU            *APU
	SampleBuffers  SampleBuffers
	SampleInterval time.Duration
	SampleDivider  int
	MCounter       int
	Output         AudioOutput
}

type AudioOutput interface {
	io.Writer
	Start()
	Stop()
}

type SampleBuffers struct {
	Pulse1 []uint8
	Pulse2 []uint8
	Wave   []uint8
	Size   int
	Idx    int
}

func (sb *SampleBuffers) Mix() []uint8 {
	mixed := make([]uint8, len(sb.Pulse1))
	for i := range len(mixed) {
		mixed[i] = sb.Pulse1[i] + sb.Pulse2[i] + sb.Wave[i]
	}
	return mixed
}

type AudioTestOutput struct {
	File string
	buf  bytes.Buffer
	mu   sync.Mutex
	w    io.WriteCloser
}

func NewAudioTestOutput(file string) *AudioTestOutput {
	return &AudioTestOutput{
		File: file,
	}
}

func (ato *AudioTestOutput) Start() {
	w, err := os.Create(ato.File)
	if err != nil {
		return
	}
	ato.w = w
}

func (ato *AudioTestOutput) Stop() {
	ato.mu.Lock()
	ato.w.Write(ato.buf.Bytes())
	ato.w = nil
	ato.buf = bytes.Buffer{}
	ato.mu.Unlock()
}

func (ato *AudioTestOutput) Write(p []uint8) (int, error) {
	ato.mu.Lock()
	defer ato.mu.Unlock()

	if ato.w != nil {
		return ato.buf.Write(p)
	}
	return 0, nil
}

func NewSampleBuffers(size int) SampleBuffers {
	return SampleBuffers{
		Pulse1: make([]uint8, size),
		Pulse2: make([]uint8, size),
		Wave:   make([]uint8, size),
		Size:   size,
	}
}

func (ab *SampleBuffers) Add(pulse1, pulse2, wave int8) bool {
	ab.Pulse1[ab.Idx] = uint8(pulse1)
	ab.Pulse2[ab.Idx] = uint8(pulse2)
	ab.Wave[ab.Idx] = uint8(wave)

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
	if audio.SampleInterval > 0 {
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
	) {
		return
	}
	audio.Output.Write(audio.SampleBuffers.Mix())
}

func (audio *Audio) Start() {
	audio.Output.Start()
}

func (audio *Audio) Stop() {
	audio.Output.Stop()
}
