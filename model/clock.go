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
	T      time.Time
	C      uint64
	Rising bool
}

func newClock() *Clock {
	clock := &Clock{
		m: &sync.Mutex{},
	}
	return clock
}

type RootClock struct {
	Clock
	freq chan float64
}

func (r *RootClock) SetFrequency(f float64) {
	r.freq <- f
}

func NewRootClock(config ClockConfig) *RootClock {
	rootClock := RootClock{
		Clock: *newClock(),
		freq:  make(chan float64, 1),
	}
	tickInterval := time.Millisecond
	ticker := time.NewTicker(tickInterval)
	cycleInterval := time.Duration(float64(time.Second) / config.Frequency)
	cyclesPerTick := tickInterval / cycleInterval
	go func() {
		var count uint64
		for t := range ticker.C {
			rootClock.m.Lock()
			for range cyclesPerTick {
				rootClock.cycle(t, count)
				count++
			}
			rootClock.m.Unlock()
		}
	}()

	go func() {
		for f := range rootClock.freq {
			rootClock.m.Lock()
			cycleInterval = time.Duration(float64(time.Second) / f)
			cyclesPerTick = tickInterval / cycleInterval
			rootClock.m.Unlock()
		}
	}()
	return &rootClock
}

func (c *Clock) cycle(t time.Time, count uint64) {
	for _, cb := range c.rising {
		cb(Cycle{t, count, true})
	}
	for _, cb := range c.falling {
		cb(Cycle{t, count, false})
	}
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
	child := newClock()
	c.AddRiseCallback(func(cyc Cycle) {
		d, m := cyc.C/div, cyc.C%div
		if m == 0 {
			child.m.Lock()
			child.cycle(cyc.T, d)
			child.m.Unlock()
		}
	})
	return child
}

func (c *Clock) Invert(div uint64) *Clock {
	child := newClock()
	c.AddFallCallback(func(cyc Cycle) {
		child.m.Lock()
		child.cycle(cyc.T, cyc.C)
		child.m.Unlock()
	})
	return child
}
