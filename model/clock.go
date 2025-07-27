package model

type Clock struct {
	Curr    Cycle
	Devices []func(c Cycle)
}

type Cycle struct {
	C       uint64
	Falling bool
}

func NewClock() *Clock {
	clock := &Clock{}
	return clock
}

func (c *Clock) Cycle() {
	c.Rising()
	c.Falling()
	c.Curr.C++
}

func (c *Clock) Rising() {
	for _, dev := range c.Devices {
		dev(Cycle{c.Curr.C, false})
	}
}

func (c *Clock) Falling() {
	for _, dev := range c.Devices {
		dev(Cycle{c.Curr.C, true})
	}
}

func (c *Clock) AttachDevice(dev func(c Cycle)) {
	c.Devices = append(c.Devices, dev)
}
