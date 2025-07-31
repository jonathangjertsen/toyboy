package blip

import "math"

// Public API

func NewSynth(buf *Buffer) *Synth {
	bs := &Synth{
		intIR: make(
			[]int16,
			buf.config.resolution()*(buf.config.Quality/2)+1,
		),
		normalizationUnit: 0,
		output:            buf,
		lastAmplitude:     0,
		deltaFactor:       0,
		config:            buf.config,
	}
	bs.computeIntegratedImpulseResponse()
	bs.updateVolume()
	return bs
}

// Called when the emulated audio device output changes
func (bs *Synth) Update(clockT int, amplitude int) {
	bs.update(bs.output.clockToRealTime(clockT), amplitude)
}

/// Implementation

type Synth struct {
	output *Buffer

	// Integral of impulse response kernel
	// There's (1 << PhaseBits) subsamples for every output sample
	// TODO: should be int?
	intIR             []int16
	normalizationUnit int
	config            *Config

	// synth state
	lastAmplitude int
	deltaFactor   int
}

func (bs *Synth) computeIntegratedImpulseResponse() {
	// Create temporary buffer to construct the impulse response kernel
	// This is only done once up front, so we do this using floating point
	cfg := bs.config
	res := cfg.resolution()
	impulse := make([]float32, res/2*(bs.config.WidestImpulse-1)+res*2)
	halfSize := res / 2 * (bs.config.Quality - 1)

	// Determine required amount of oversampling
	oversample := float64(res)*2.25/float64(halfSize) + 0.85
	nyquist := float64(bs.output.sampleRate) * 0.5
	if cfg.LPFCutoffFreq > 0 {
		oversample = nyquist / float64(cfg.LPFCutoffFreq)
	}

	// Set up parameters for impulse response
	cutoff := min(cfg.LPFRolloffFreq*oversample/nyquist, 0.999)
	treble := min(5.0, max(-300.0, cfg.TrebleBoostdB))
	maxHarmonic := float64(4096.0)
	rolloffFactor := math.Pow(10.0, 1.0/(maxHarmonic*20)*treble/(1.0-cutoff))
	powAN := math.Pow(rolloffFactor, maxHarmonic-maxHarmonic*cutoff)
	angleIncrement := float64(math.Pi / 2 / maxHarmonic / (float64(res) * oversample))

	// Construct pre-impulse part of response
	// This generates half a sinc
	for i := range halfSize {
		angle := ((float64(i-halfSize)*2 + 1) * angleIncrement)
		c := rolloffFactor*math.Cos((maxHarmonic-1.0)*angle) - math.Cos(maxHarmonic*angle)
		cosNCAngle := math.Cos(maxHarmonic * cutoff * angle)
		cosNC1Angle := math.Cos((maxHarmonic*cutoff - 1.0) * angle)
		cosAngle := math.Cos(angle)
		c = c*powAN - rolloffFactor*cosNC1Angle*cosNCAngle
		d := 1.0 + rolloffFactor*(rolloffFactor-cosAngle-cosAngle)
		b := 2.0 - cosAngle - cosAngle
		a := 1.0 - cosAngle - cosNCAngle + cosNC1Angle
		impulse[i] = float32((a*d + c*b) / (b * d)) // a / b + c / d
	}

	// Apply (half of) a Hamming window to the half generated so far
	toFraction := float64(math.Pi) / (float64(halfSize - 1))
	for i := halfSize - 1; i >= 0; i-- {
		impulse[i] *= float32(0.54 - 0.46*math.Cos(float64(i)*toFraction))
	}

	// Construct post-impulse part of response by mirroring the pre-impulse part
	for i := res - 1; i >= 0; i-- {
		impulse[res+halfSize+i] = impulse[res+halfSize-1-i]
	}

	// Set first sample (so first res subsamples) to 0
	clear(impulse[:res])

	// Normalize filter
	total := 0.0
	for i := range halfSize {
		total += float64(impulse[res+i])
	}
	rescale := bs.config.KernelBaseUnit / 2 / total
	bs.normalizationUnit = int(bs.config.KernelBaseUnit)

	// Integrate and rescale
	// The integrated impulse response can be applied to deltas to reconstruct the actual signal
	sum := 0.0
	next := 0.0
	impulsesSize := len(bs.intIR)
	for i := range impulsesSize {
		bs.intIR[i] = int16(math.Floor((next-sum)*rescale + 0.5))
		sum += float64(impulse[i])
		next += float64(impulse[i+res])
	}

	bs.doErrorCorrection()
	bs.updateVolume()
}

func (bs *Synth) doErrorCorrection() {
	res := bs.config.resolution()

	// lSubsample goes from left to center
	// rSubsample goes from right to center
	for rSubsample := res - 1; rSubsample >= res/2; rSubsample-- {
		lSubsample := res - 2 - rSubsample

		// Start with theoretical sum of phase pairs
		correction := bs.normalizationUnit
		for i := 1; i < len(bs.intIR); i += res {
			correction -= int(bs.intIR[i+rSubsample])
			correction -= int(bs.intIR[i+lSubsample])
		}

		// Center impulse uses same half for both sides, gets double-counted above
		if rSubsample == lSubsample {
			correction /= 2
		}

		// Apply correction to end of first half
		bs.intIR[len(bs.intIR)-res+rSubsample] += int16(correction)
	}
}

func (bs *Synth) updateVolume() {
	volumeUnit := bs.config.Volume / float64(bs.config.InputRange)
	deltaFactor := volumeUnit * float64(FixedOne(bs.config.SampleBits)) / float64(bs.normalizationUnit)

	// If volume is really low, increase deltaFactor and attentuate the kernel
	// Probably something to do with precision.
	if deltaFactor > 0.0 {
		shift := 0

		for deltaFactor < 2.0 {
			shift++
			deltaFactor *= 2.0
		}

		if shift > 0 {
			bs.normalizationUnit >>= shift
			assert(bs.normalizationUnit > 0, "volume unit is too low")

			half := (1 << (shift - 1))
			rounding := 0x8000 + half
			for i, original := range bs.intIR {
				bs.intIR[i] = int16(((int(original) + rounding) >> shift) - (rounding >> shift))
			}
			bs.doErrorCorrection()
		}
	}
	bs.deltaFactor = int(deltaFactor + 0.5)
}

func (bs *Synth) update(t uint, amplitude int) {
	delta := (amplitude - bs.lastAmplitude) * bs.deltaFactor
	bs.lastAmplitude = amplitude

	subsamples := bs.config.resolution()
	offset := int(t >> bs.config.BufferAccuracy)
	if offset >= len(bs.output.buf) {
		// time is beyond end of buf
		return
	}

	phase := int(t>>(bs.config.BufferAccuracy-bs.config.PhaseBits)) & (subsamples - 1)

	imp := bs.intIR[subsamples-phase:]
	buf := bs.output.buf[offset:]

	currImpulseSample := int(imp[0])

	// Convolve first part
	fwd := (bs.config.WidestImpulse - bs.config.Quality) / 2
	{
		i := 0

		buf[fwd+i] += currImpulseSample * delta
		buf[fwd+1+i] += int(imp[subsamples*(i+1)]) * delta
		currImpulseSample = int(imp[subsamples*(i+2)])
	}
	if bs.config.Quality > 8 {
		i := 2

		buf[fwd+i] += currImpulseSample * delta
		buf[fwd+1+i] += int(imp[subsamples*(i+1)]) * delta
		currImpulseSample = int(imp[subsamples*(i+2)])
	}
	if bs.config.Quality > 12 {
		i := 4

		buf[fwd+i] += currImpulseSample * delta
		buf[fwd+1+i] += int(imp[subsamples*(i+1)]) * delta
		currImpulseSample = int(imp[subsamples*(i+2)])
	}

	// Convolve center
	mid := bs.config.Quality/2 - 1
	buf[fwd+mid-1] += currImpulseSample * delta
	currImpulseSample = int(imp[subsamples*mid])
	buf[fwd+mid] += currImpulseSample * delta
	imp = imp[phase:]

	// Convolve tail
	rev := fwd + bs.config.Quality - 2
	if bs.config.Quality > 12 {
		r := 6

		buf[rev-r] += currImpulseSample * delta
		buf[rev+1-r] += int(imp[subsamples*r]) * delta
		currImpulseSample = int(imp[subsamples*(r-1)])
	}
	if bs.config.Quality > 8 {
		r := 4

		buf[rev-r] += currImpulseSample * delta
		buf[rev+1-r] += int(imp[subsamples*r]) * delta
		currImpulseSample = int(imp[subsamples*(r-1)])
	}
	{
		r := 2

		buf[rev-r] += currImpulseSample * delta
		buf[rev+1-r] += int(imp[subsamples*r]) * delta
		currImpulseSample = int(imp[subsamples*(r-1)])
	}

	buf[rev] += currImpulseSample * delta
	buf[rev+1] += int(imp[0]) * delta
}
