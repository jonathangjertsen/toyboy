package model

type Envelope struct {
	volume       Data8
	volumeReset  Data8
	envTimer     Data8
	envDir       bool
	envSweepPace Data8
}

func (env *Envelope) SetVolumeEnvelope(v Data8) bool {
	env.envSweepPace = v & 0x7
	env.envDir = v&Bit3 != 0
	env.volumeReset = (v >> 4) & 0xf
	return v&0xf8 != 0
}

func (env *Envelope) clock(envtick Data8) {
	if env.envSweepPace == 0 {
		return
	}
	if envtick%env.envSweepPace != 0 {
		return
	}
	if env.envDir {
		if env.volume < 16 {
			env.volume++
		}
	} else {
		if env.volume > 0 {
			env.volume--
		}
	}
}

func (env *Envelope) scale(sample int8) int8 {
	return sample * int8(env.volume)
}
