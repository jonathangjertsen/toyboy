package model

type NoiseChannel struct {
	PeriodCounter PeriodCounter
	LengthTimer   LengthTimer
	Envelope      Envelope

	RegLengthTimer    Data8
	RegVolumeEnvelope Data8
	RegRNG            Data8
	RegCtl            Data8

	activated  bool
	dacEnabled bool
}

func (nc *NoiseChannel) clock() {

}

func (nc *NoiseChannel) SetLengthTimer(v Data8) {
	nc.RegLengthTimer = v
}

func (nc *NoiseChannel) SetVolumeEnvelope(v Data8) {
	nc.RegVolumeEnvelope = v

	nc.dacEnabled = nc.Envelope.SetVolumeEnvelope(v)
	if !nc.dacEnabled {
		nc.activated = false
	}
}

func (nc *NoiseChannel) SetRNG(v Data8) {
	nc.RegRNG = v
}

func (nc *NoiseChannel) SetCtl(v Data8) {
	nc.RegCtl = v
}

func (nc *NoiseChannel) tickLengthTimer() {
	nc.LengthTimer.clock(64)
}
