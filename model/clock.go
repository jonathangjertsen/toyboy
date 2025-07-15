package model

import (
	"fmt"
	"time"
)

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

type RealtimeClock struct {
	Clock
	resume chan struct{}
	pause  chan struct{}
	freq   chan float64
}

func (r *RealtimeClock) SetFrequency(f float64) {
	r.freq <- f
}

func NewRealtimeClock(config ClockConfig) *RealtimeClock {
	rtClock := RealtimeClock{
		Clock:  *NewClock(),
		resume: make(chan struct{}),
		pause:  make(chan struct{}),
		freq:   make(chan float64, 1),
	}
	tickInterval := time.Millisecond
	cycleInterval := time.Duration(float64(time.Second) / config.Frequency)
	cyclesPerTick := uint64(tickInterval / cycleInterval)
	go func() {
		var count uint64
		<-rtClock.resume
		ticker := time.NewTicker(tickInterval)
		for {
			select {
			case <-ticker.C:
				count = rtClock.Cycles(count, cyclesPerTick)
			case <-rtClock.resume:
				fmt.Printf("Ignored resume\n")
			case <-rtClock.pause:
				for {
					resumed := false
					select {
					case <-rtClock.pause:
						fmt.Printf("Ignored pause\n")
					case <-rtClock.resume:
						resumed = true
					}
					if resumed {
						break
					}
				}
			}
		}
	}()
	go func() {
		for f := range rtClock.freq {
			cycleInterval = time.Duration(float64(time.Second) / f)
			cyclesPerTick = uint64(tickInterval / cycleInterval)
		}
	}()
	return &rtClock
}

// Start the clock
// When this function returns, the clock has started
func (rtClock *RealtimeClock) Start() {
	rtClock.resume <- struct{}{}
}

// Stop the clock
// When this function returns, the clock has stopped
func (rtClock *RealtimeClock) Stop() {
	rtClock.pause <- struct{}{}
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

func (c *Clock) Cycles(currCycle uint64, n uint64) uint64 {
	for offs := range n {
		c.Cycle(currCycle + offs)
	}
	return currCycle + n
}

func (c *Clock) AttachDevice(dev func(c Cycle)) {
	c.devices = append(c.devices, dev)
}

func (c *Clock) Divide(div uint64) *Clock {
	if div == 0 {
		panic("divide clock by 0?")
	}
	if div == 1 {
		return c
	}
	child := NewClock()
	c.AttachDevice(func(cyc Cycle) {
		d, m := cyc.C/div, cyc.C%div
		if !cyc.Falling && m == 0 {
			child.Rising(d)
		}
		if cyc.Falling && m == (div/2) {
			child.Falling(d)
		}
	})
	return child
}
