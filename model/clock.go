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
	tickInterval  time.Duration
	cyclesPerTick uint64
	resume        chan struct{}
	pause         chan struct{}
	jobs          chan func()
}

// Executes the function in the clocks' goroutine
func (r *RealtimeClock) Sync(f func()) {
	done := make(chan struct{})
	r.jobs <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}

func (r *RealtimeClock) SetFrequency(f float64) {
	r.Sync(func() {
		r.setFreq(f)
	})
}

func NewRealtimeClock(config ClockConfig) *RealtimeClock {
	rtClock := RealtimeClock{
		Clock:  *NewClock(),
		resume: make(chan struct{}),
		pause:  make(chan struct{}),
		jobs:   make(chan func()),
	}
	go rtClock.run(config.Frequency)
	return &rtClock
}

func (rtClock *RealtimeClock) wait() {
	for {
		resumed := false
		select {
		case <-rtClock.pause:
			fmt.Printf("Ignored pause\n")
		case <-rtClock.resume:
			resumed = true
		case job := <-rtClock.jobs:
			job()
		}
		if resumed {
			break
		}
	}
}
func (rtClock *RealtimeClock) setFreq(f float64) {
	rtClock.tickInterval = time.Millisecond * 2
	cycleInterval := time.Duration(float64(time.Second) / f)
	rtClock.cyclesPerTick = uint64(rtClock.tickInterval / cycleInterval)
}

func (rtClock *RealtimeClock) run(initFreq float64) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	rtClock.setFreq(initFreq)

	var count uint64
	rtClock.wait()
	ticker := time.NewTicker(rtClock.tickInterval)
	for {
		select {
		case <-ticker.C:
			count = rtClock.Cycles(count, rtClock.cyclesPerTick)
		case <-rtClock.resume:
			fmt.Printf("Ignored resume\n")
		case <-rtClock.pause:
			rtClock.wait()
		case job := <-rtClock.jobs:
			job()
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
