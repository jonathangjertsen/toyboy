package model

import (
	"fmt"
	"runtime"
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
	go rtClock.run(config.Frequency)
	return &rtClock
}

func (rtClock *RealtimeClock) run(initFreq float64) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tickInterval := time.Millisecond * 2
	cycleInterval := time.Duration(float64(time.Second) / initFreq)
	cyclesPerTick := uint64(tickInterval / cycleInterval)

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
		case f := <-rtClock.freq:
			cycleInterval = time.Duration(float64(time.Second) / f)
			cyclesPerTick = uint64(tickInterval / cycleInterval)
		}
	}
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

func (c *Clock) Divide(pow uint64) *Clock {
	if pow == 0 {
		return c
	}
	if pow > 63 {
		panic("too big division")
	}
	mask := uint64((1 << pow) - 1)
	fallV := uint64(1 << (pow - 1))
	child := NewClock()
	c.AttachDevice(func(cyc Cycle) {
		d, m := cyc.C>>pow, cyc.C&mask
		if !cyc.Falling && m == 0 {
			child.Rising(d)
		}
		if cyc.Falling && m == fallV {
			child.Falling(d)
		}
	})
	return child
}
