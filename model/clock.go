package model

type Clock struct {
	devices []func(c Cycle)
}

type Cycle struct {
	C       uint64
	Falling bool
}

func NewClock() *Clock {
	clock := &Clock{}
	return clock
}

func (c *Clock) Cycle(currCycle uint64) {
	c.Rising(currCycle)
	c.Falling(currCycle)
}

func (c *Clock) Rising(currCycle uint64) {
	for _, dev := range c.devices {
		dev(Cycle{currCycle, false})
	}
}

func (c *Clock) Falling(currCycle uint64) {
	for _, dev := range c.devices {
		dev(Cycle{currCycle, true})
	}
}

func (c *Clock) AttachDevice(dev func(c Cycle)) {
	c.devices = append(c.devices, dev)
}
