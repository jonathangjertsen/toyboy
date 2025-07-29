package model

import "time"

// Port of Blargg's Blip_Buffer library
// Here is the original license for that library:

/* Copyright (C) 2003-2006 Shay Green. This module is free software; you
can redistribute it and/or modify it under the terms of the GNU Lesser
General Public License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version. This
module is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for
more details. You should have received a copy of the GNU Lesser General
Public License along with this module; if not, write to the Free Software
Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA 02111-1307 USA */

const (
	BlipBufferAccuracy = 16
	BlipPhaseBits      = 6
	BlipRes            = 1 << BlipPhaseBits
	BlipSampleBits     = 30
	BlipGoodQuality    = 12
	BlipWidestImpulse  = 16
)

type Time uint
type ResampledTime uint

type BlipBuffer struct {
	SampleRate   int
	ClockRate    int
	Buffer       []AudioSample
	HighPassFreq int

	Factor uint          // ?
	Offset ResampledTime // ?
}

func NewBlipBuffer() *BlipBuffer {
	return &BlipBuffer{}
}

func (bb *BlipBuffer) SetSampleRate(
	samplesPerSec int,
	bufferLength time.Duration,
) {
	bb.SampleRate = samplesPerSec
}

func (bb *BlipBuffer) SetClockRate(
	clocksPerSecond int,
) {
	bb.ClockRate = clocksPerSecond
}

func (bb *BlipBuffer) EndFrame() {
}

func (bb *BlipBuffer) SetHighPassFreq(freq int) {
	bb.HighPassFreq = freq
}

func (bb *BlipBuffer) OutputLatency() int {
	return 0
}

func (bb *BlipBuffer) ClearWaiting() {
}

func (bb *BlipBuffer) ClearAll() {
}

func (bb *BlipBuffer) SamplesAvailable() int {
	return 0
}

func (bb *BlipBuffer) RemoveSamples(n int) {
}

func (bb *BlipBuffer) CountSamples(dt Time) {
}

func (bb *BlipBuffer) MixSamples(buf []AudioSample) {
	n := len(buf)
	_ = n
}

func (bb *BlipBuffer) CountClocks(n int) int {
	return 0
}

func (bb *BlipBuffer) RemoveSilence(n int) {
}

func (bb *BlipBuffer) ResampledDuration(n int) ResampledTime {
	return ResampledTime(n) * ResampledTime(bb.Factor)
}

func (bb *BlipBuffer) ResampledTime(t Time) ResampledTime {
	return ResampledTime(t)*ResampledTime(bb.Factor) + bb.Offset
}

func (bb *BlipBuffer) ClockRateFactor(clockRate uint) ResampledTime {
	return 0
}

func (bb *BlipBuffer) ReadMono(buf []AudioSample) {
	n := len(buf)
	_ = n
}

type BlipSynth struct {
}
