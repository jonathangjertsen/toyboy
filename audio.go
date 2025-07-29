package main

import (
	"io"

	"github.com/ebitengine/oto/v3"
	"github.com/jonathangjertsen/toyboy/model"
)

type AudioInterface struct {
	In       chan []model.AudioSample
	buffered []model.AudioSample
	player   *oto.Player
}

func NewAudio() *AudioInterface {
	aif := &AudioInterface{}
	aif.In = make(chan []model.AudioSample, 200)
	opts := &oto.NewContextOptions{}
	opts.SampleRate = 44100
	opts.ChannelCount = 1
	opts.Format = oto.FormatSignedInt16LE
	otoCtx, readyChan, err := oto.NewContext(opts)
	if err != nil {
		panic("oto.NewContext failed with " + err.Error())
	}
	<-readyChan
	aif.player = otoCtx.NewPlayer(aif)
	go aif.player.Play()
	return aif
}

func emitSample(sample model.AudioSample, buf []uint8, offset *int) {
	data16 := uint16(sample)
	lsb, msb := uint8(data16), uint8(data16>>8)
	buf[*offset] = lsb
	*offset++
	buf[*offset] = msb
	*offset++
}

func emitSamples(
	samples []model.AudioSample,
	buf []uint8,
	n int,
	offsetOut *int,
) {
	for i := range n {
		emitSample(samples[i], buf, offsetOut)
	}
}

func (aif *AudioInterface) Read(p []byte) (int, error) {
	nBytesRequested := len(p)
	nBytesOut := 0

	// Start by emitting samples buffered from the last round
	nSamplesBuffered := len(aif.buffered)
	if nSamplesBuffered > 0 {
		nSamplesRequested := nBytesRequested / 2
		if nSamplesBuffered > nSamplesRequested {
			// Emit only the requested number of samples
			// then adjust for next round and we are done
			emitSamples(aif.buffered, p, nSamplesRequested, &nBytesOut)
			aif.buffered = aif.buffered[nSamplesRequested:]
			return nBytesRequested, nil
		} else {
			// Emit all buffered samples
			// then continue receiving buffers
			emitSamples(aif.buffered, p, nSamplesBuffered, &nBytesOut)
			aif.buffered = nil
		}
	}

	// We can't block here, so only provide samples from the channel
	// if there's something there
	nAvailable := len(aif.In)
	for range nAvailable {
		// This should never block, but the channel might be closed in which case we dip
		data, ok := <-aif.In
		if !ok {
			aif.player.Close()
			return 0, io.EOF
		}

		// Amplify samples to PCM level
		// TODO: this really should be done on the sender side, but it gets distorted somehow
		for i := range data {
			data[i] *= 256
		}

		// If we can't emit everything now, buffer the remainder for next round
		nSamplesToEmit := len(data)
		nBytesAfter := nBytesOut + nSamplesToEmit*2
		if nBytesAfter > nBytesRequested {
			nSamplesToEmit -= (nBytesAfter - nBytesRequested) / 2
			aif.buffered = append(aif.buffered, data[nSamplesToEmit:]...)
			break
		}

		// Emit as much as possible
		emitSamples(data, p, nSamplesToEmit, &nBytesOut)
	}

	return nBytesOut, nil
}
