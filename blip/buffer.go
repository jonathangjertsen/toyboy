package blip

import "fmt"

// Public API

func NewBuffer(config Config) *Buffer {
	bb := &Buffer{}

	// The original implementation sets this to nullptr initially, and reallocs in SetSamplingParams
	// That's not really
	// We just allocate a reasonable amount up front and never change the allocation
	bb.mem = make([]int, config.MaxBufferSize)
	bb.buf = bb.mem[:config.InitialBufferSize]

	bb.config = &config
	bb.SetSamplingParams(config.InitialSampleRate, config.InitialBufferSize)
	bb.SetClockRate(config.InitialClockRate)
	bb.updateHPFShift()

	return bb
}

func (bb *Buffer) SetSamplingParams(sampleRate, bufferSize int) {
	if bufferSize > cap(bb.mem) {
		bufferSize = cap(bb.mem)
	}
	bb.buf = bb.mem[:bufferSize]

	bb.sampleRate = sampleRate
	bb.updateHPFShift() // hpfShift depends on sampleRate, so need to update that as well
	if bb.clockRate != 0 {
		bb.SetClockRate(bb.clockRate)
	}
	bb.offset = 0
	bb.accumulator = 0
	clear(bb.buf[:bb.SamplesAvailable()+bb.config.bufferExta()])
}

func (bb *Buffer) SetClockRate(cps int) bool {
	r := float64(bb.sampleRate) / float64(cps)
	rFixedPoint := Float64ToFixed(r, bb.config.BufferAccuracy)
	ok := rFixedPoint > 0 || (bb.sampleRate == 0)
	if ok {
		bb.clockRate = cps
		bb.factor = uint(rFixedPoint)
	}
	return ok
}

func (bb *Buffer) EndFrame(t int) {
	bb.offset += uint(t) * bb.factor
	fmt.Printf("offset=%d factor=%d\n", bb.offset, bb.factor)
}

func (bb *Buffer) SamplesAvailable() int {
	return int(bb.offset >> bb.config.BufferAccuracy)
}

func (bb *Buffer) CountSamples(t int) int {
	lastSample := bb.clockToRealTime(t) >> bb.config.BufferAccuracy
	firstSample := bb.offset >> bb.config.BufferAccuracy
	return int(lastSample - firstSample)
}

func (bb *Buffer) MixSamples(samples []int16) {
	baseOffset := int((bb.offset >> bb.config.BufferAccuracy) + uint(bb.config.WidestImpulse)/2)
	sampleShift := bb.config.SampleBits - 16
	prev := 0
	for i, x := range samples {
		s := int(x) << sampleShift
		bb.buf[baseOffset+i] += s - prev
		prev = s
	}
	bb.buf[baseOffset+len(samples)-1] -= prev
}

func (bb *Buffer) Read(out []int16) int {
	count := bb.SamplesAvailable()

	if count > len(out) {
		count = len(out)
	}
	if count == 0 {
		return count
	}
	sampleShift := bb.config.SampleBits - 16
	hpfShift := bb.hpfShift
	accumulator := bb.accumulator

	accumulator -= accumulator >> hpfShift
	accumulator += bb.buf[0]

	for i := range count {
		s := accumulator >> sampleShift
		accumulator -= accumulator >> hpfShift
		bufferValue := bb.buf[i]
		accumulator += bufferValue

		if int(int16(s)) == s {
			out[i] = int16(s)
		} else {
			out[i] = int16(0x7fff - (s >> 24))
		}
	}
	bb.accumulator = accumulator

	if count > 0 {
		bb.offset -= uint(count) << bb.config.BufferAccuracy
		remain := bb.SamplesAvailable() + bb.config.bufferExta()
		copy(bb.buf[:remain], bb.buf[count:count+remain])
		clear(bb.buf[remain : remain+count])
	}
	return count
}

// Implementation

type Buffer struct {
	mem         []int
	buf         []int // sub-slice of mem
	sampleRate  int
	clockRate   int
	hpfShift    int
	accumulator int
	factor      uint // Fixed point with BufferAccuracy bits
	offset      uint // Fixed point with BufferAccuracy bits
	config      *Config
}

func (bb *Buffer) updateHPFShift() {
	shift := int(31)
	if bb.config.HPFFrequency > 0 {
		shift = 13
		f := (bb.config.HPFFrequency << 16) / bb.sampleRate
		for {
			f >>= 1
			if f == 0 {
				break
			}
			shift--
			if shift == 0 {
				break
			}
		}
	}
	bb.hpfShift = shift
}

func (bb *Buffer) clockToRealTime(t int) uint {
	return uint(t)*bb.factor + bb.offset
}
