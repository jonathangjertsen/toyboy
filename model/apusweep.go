package model

type Sweep struct {
	RegSweep Data8
}

func (sw *Sweep) SetSweep(v Data8) {
	sw.RegSweep = v
}
