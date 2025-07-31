package model

type Envelope struct {
	Volume       Data8
	VolumeReset  Data8
	EnvTimer     Data8
	EnvDir       bool
	EnvSweepPace Data8
}

func (env *Envelope) SetVolumeEnvelope(v Data8) bool {
	env.EnvSweepPace = v & 0x7
	env.EnvDir = v&Bit3 != 0
	env.VolumeReset = (v >> 4) & 0xf
	return v&0xf8 != 0
}

func (env *Envelope) clock(envtick Data8) {
	if env.EnvSweepPace == 0 {
		return
	}
	if envtick%env.EnvSweepPace != 0 {
		return
	}
	if env.EnvDir {
		if env.Volume < 16 {
			env.Volume++
		}
	} else {
		if env.Volume > 0 {
			env.Volume--
		}
	}
}

func (env *Envelope) scale(sample AudioSample) AudioSample {
	return sample * AudioSample(env.Volume)
}
