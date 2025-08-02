package model

import "time"

func AudioStub() (Audio, chan []AudioSample) {
	devnull := make(chan []AudioSample, 1024)
	go func() {
		for range devnull {
		}
	}()
	audio := &AudioNN{
		SampleInterval: time.Second / 44100,
		SampleBuffers:  NewSampleBuffers(512),
		SubSampling:    1024,
		Out:            devnull,
	}
	return audio, devnull
}
