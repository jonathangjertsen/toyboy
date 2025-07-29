package model

type Mixer struct {
	RegChannelPan         Data8
	RegMasterVolumeVINPan Data8
}

func (mixer *Mixer) SetChannelPan(v Data8) {
	mixer.RegChannelPan = v
}

func (mixer *Mixer) SetMasterVolumeVINPan(v Data8) {
	mixer.RegMasterVolumeVINPan = v
}
func (mixer *Mixer) MixStereoSimple(p1, p2, w, n AudioSample) (AudioSample, AudioSample) {
	return (p1 + p2 + w + n), (p1 + p2 + w + n)
}

func (mixer *Mixer) MixStereo(p1, p2, w, n AudioSample) (AudioSample, AudioSample) {
	left := AudioSample(0)
	right := AudioSample(0)

	if mixer.RegChannelPan&Bit0 != 0 {
		left += p1
	}
	if mixer.RegChannelPan&Bit1 != 0 {
		left += p2
	}
	if mixer.RegChannelPan&Bit2 != 0 {
		left += w
	}
	if mixer.RegChannelPan&Bit3 != 0 {
		left += n
	}
	if mixer.RegChannelPan&Bit4 != 0 {
		right += p1
	}
	if mixer.RegChannelPan&Bit5 != 0 {
		right += p2
	}
	if mixer.RegChannelPan&Bit6 != 0 {
		right += w
	}
	if mixer.RegChannelPan&Bit7 != 0 {
		right += n
	}

	rightVol := mixer.RegMasterVolumeVINPan & 0x07
	leftVol := (mixer.RegMasterVolumeVINPan >> 4) & 0x07
	left *= (AudioSample(leftVol) + 1)
	right *= (AudioSample(rightVol) + 1)

	// p1, p2, w and n can be at most 16
	// left and right pre mul: at most 16+16+16+16=64
	// left and right post mul: at most 64*8=512
	// 16 bit PCM is +/-32768
	// so multiply each by 64
	return left, right
}
