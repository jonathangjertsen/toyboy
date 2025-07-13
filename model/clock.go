package model

import (
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
	C      uint64
	Rising bool
}

func NewClock() *Clock {
	clock := &Clock{
		m: &sync.Mutex{},
	}
	return clock
}

type RealtimeClock struct {
	Clock
	begin chan struct{}
	freq  chan float64
}

func (r *RealtimeClock) SetFrequency(f float64) {
	r.freq <- f
}

func NewRealtimeClock(config ClockConfig) *RealtimeClock {
	rtClock := RealtimeClock{
		Clock: *NewClock(),
		begin: make(chan struct{}),
		freq:  make(chan float64, 1),
	}
	tickInterval := time.Millisecond
	cycleInterval := time.Duration(float64(time.Second) / config.Frequency)
	cyclesPerTick := uint64(tickInterval / cycleInterval)
	go func() {
		var count uint64
		<-rtClock.begin
		ticker := time.NewTicker(tickInterval)
		for range ticker.C {
			rtClock.m.Lock()
			count = rtClock.Cycles(count, cyclesPerTick)
			rtClock.m.Unlock()
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

func (rtClock *RealtimeClock) Start() {
	rtClock.begin <- struct{}{}
}

func (c *Clock) Cycle(currCycle uint64) {
	for _, cb := range c.rising {
		cb(Cycle{currCycle, true})
	}
	for _, cb := range c.falling {
		cb(Cycle{currCycle, false})
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
