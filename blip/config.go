package blip

var DefaultBlipConfig = Config{
	InputRange:        32,
	Volume:            0.5, // up to 1.0, which would map -InputRange:InputRange to -32768:32767
	Quality:           12,
	MaxBufferSize:     100000,
	BufferAccuracy:    16,
	PhaseBits:         6,
	SampleBits:        30,
	WidestImpulse:     16,
	TrebleBoostdB:     -8.0,
	LPFRolloffFreq:    0,
	LPFCutoffFreq:     0,
	HPFFrequency:      0,
	InitialSampleRate: 44100,
	InitialBufferSize: 44100,
	InitialClockRate:  44100,

	//baseUnit :=44800.0 - 128 * 18; // allows treble up to +0 dB
	//baseUnit :=37888.0; // allows treble to +5 dB
	KernelBaseUnit: 32768.0,
}

type Config struct {
	// Expected input range for the samples
	InputRange int

	// From 0.0 (silent) to 1.0 (maps input range to max range of 16-bit PCM)
	// Not adjustable after setup; the BL
	Volume float64

	Quality        int
	MaxBufferSize  int
	BufferAccuracy int

	// Sets the subsample resolution (subsample = 2^PhaseBits)
	// http://slack.net/~ant/bl-synth/11.implementation.html
	PhaseBits int

	SampleBits    int
	WidestImpulse int

	// High shelf EQ
	// Boost amount (or more likely, cut)
	TrebleBoostdB  float64
	LPFRolloffFreq float64
	LPFCutoffFreq  float64

	InitialSampleRate int
	InitialBufferSize int
	InitialClockRate  int
	HPFFrequency      int
	KernelBaseUnit    float64
}

func (c *Config) resolution() int {
	return 1 << c.PhaseBits
}

func (c *Config) bufferExta() int {
	return c.WidestImpulse + 2
}
