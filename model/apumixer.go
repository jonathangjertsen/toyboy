package model

type Mixer struct {
	RegChannelPan         Data8
	RegMasterVolumeVINPan Data8

	leftVIN        bool
	rightVIN       bool
	leftVol        Data8
	rightVol       Data8
	leftEnChannel  [4]bool
	rightEnChannel [4]bool
}

func (mixer *Mixer) SetChannelPan(v Data8) {
	mixer.RegChannelPan = v
	mixer.leftEnChannel[0] = v&(1<<0) != 0
	mixer.leftEnChannel[1] = v&(1<<1) != 0
	mixer.leftEnChannel[2] = v&(1<<2) != 0
	mixer.leftEnChannel[3] = v&(1<<3) != 0
	mixer.rightEnChannel[0] = v&(1<<4) != 0
	mixer.rightEnChannel[1] = v&(1<<5) != 0
	mixer.rightEnChannel[2] = v&(1<<6) != 0
	mixer.rightEnChannel[3] = v&(1<<7) != 0
}

func (mixer *Mixer) SetMasterVolumeVINPan(v Data8) {
	mixer.RegMasterVolumeVINPan = v
	mixer.rightVol = v & 0x7
	mixer.rightVIN = (v>>3)&0x1 != 0
	mixer.leftVol = (v >> 4) & 0x7
	mixer.leftVIN = (v>>7)&0x1 != 0
}
