package model

import (
	"fmt"
	"sync"
	"time"
)

type Callback func(c Cycle)

type Clock struct {
	m       *sync.Mutex
	rising  []Callback
	falling []Callback
}

type Cycle struct {
	C       uint64
	Falling bool
}

func NewClock() *Clock {
	clock := &Clock{
		m: &sync.Mutex{},
	}
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
				rtClock.m.Lock()
				count = rtClock.Cycles(count, cyclesPerTick)
				rtClock.m.Unlock()
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
			rtClock.m.Lock()
			cycleInterval = time.Duration(float64(time.Second) / f)
			cyclesPerTick = uint64(tickInterval / cycleInterval)
			rtClock.m.Unlock()
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
	for _, cb := range c.rising {
		cb(Cycle{currCycle, false})
	}
	for _, cb := range c.falling {
		cb(Cycle{currCycle, true})
	}
}

func (c *Clock) Cycles(currCycle uint64, n uint64) uint64 {
	for offs := range n {
		c.Cycle(currCycle + offs)
	}
	return currCycle + n
}

func (c *Clock) AddRiseCallback(cb Callback) {
	c.m.Lock()
	c.rising = append(c.rising, cb)
	c.m.Unlock()
}

func (c *Clock) AddFallCallback(cb Callback) {
	c.m.Lock()
	c.falling = append(c.falling, cb)
	c.m.Unlock()
}

func (c *Clock) Divide(div uint64) *Clock {
	child := NewClock()
	c.AddRiseCallback(func(cyc Cycle) {
		d, m := cyc.C/div, cyc.C%div
		if m == 0 {
			child.m.Lock()
			child.Cycle(d)
			child.m.Unlock()
		}
	})
	return child
}

func (c *Clock) Invert(div uint64) *Clock {
	child := NewClock()
	c.AddFallCallback(func(cyc Cycle) {
		child.m.Lock()
		child.Cycle(cyc.C)
		child.m.Unlock()
	})
	return child
}
