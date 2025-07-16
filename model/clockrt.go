package model

import (
	"fmt"
	"runtime"
	"time"
)

type ClockRT struct {
	ticker        *time.Ticker
	tickInterval  time.Duration
	cyclesPerTick uint64
	resume        chan struct{}
	pause         chan struct{}
	jobs          chan func()
	uiDevices     []func()
	divided       []clockRTDivided
}

// Executes the function in the clocks' goroutine
func (r *ClockRT) Sync(f func()) {
	done := make(chan struct{})
	r.jobs <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}

func (r *ClockRT) SetFrequency(f float64) {
	r.Sync(func() {
		r.setFreq(f)
	})
}

func NewRealtimeClock(config ClockConfig) *ClockRT {
	clockRT := ClockRT{
		resume: make(chan struct{}),
		pause:  make(chan struct{}),
		jobs:   make(chan func()),
	}
	go clockRT.run(config.Frequency)
	return &clockRT
}

func (clockRT *ClockRT) AttachUIDevice(dev func()) {
	clockRT.uiDevices = append(clockRT.uiDevices, dev)
}

type clockRTDivided struct {
	clock   *Clock
	top     uint64
	counter uint64
	cycle   uint64
}

func (clockRT *ClockRT) Divide(top uint64) *Clock {
	clock := NewClock()
	clockRT.divided = append(clockRT.divided, clockRTDivided{clock, top, 0, 0})
	return clock
}

func (clockRT *ClockRT) wait() {
	for {
		resumed := false
		select {
		case <-clockRT.pause:
			fmt.Printf("Ignored pause\n")
		case <-clockRT.resume:
			resumed = true
		case job := <-clockRT.jobs:
			job()
		}
		if resumed {
			break
		}
	}
}

func (clockRT *ClockRT) setFreq(f float64) {
	cycleInterval := time.Duration(float64(time.Second) / f)
	clockRT.tickInterval = time.Millisecond * 2
	if clockRT.tickInterval < cycleInterval {
		clockRT.tickInterval = cycleInterval
	}
	clockRT.cyclesPerTick = uint64(clockRT.tickInterval / cycleInterval)
	if clockRT.ticker != nil {
		clockRT.ticker.Reset(clockRT.tickInterval)
	}
}

func (clockRT *ClockRT) run(initFreq float64) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	clockRT.setFreq(initFreq)

	var count uint64
	clockRT.wait()
	clockRT.ticker = time.NewTicker(clockRT.tickInterval)
	for {
		select {
		case <-clockRT.ticker.C:
			count = clockRT.Cycles(count, clockRT.cyclesPerTick)
			clockRT.uiCycle()
		case <-clockRT.resume:
			fmt.Printf("Ignored resume\n")
		case <-clockRT.pause:
			clockRT.wait()
		case job := <-clockRT.jobs:
			job()
		}
	}
}

func (clockRT *ClockRT) uiCycle() {
	for _, dev := range clockRT.uiDevices {
		dev()
	}
}

// Start the clock
// When this function returns, the clock has started
func (clockRT *ClockRT) Start() {
	clockRT.resume <- struct{}{}
}

// Stop the clock
// When this function returns, the clock has stopped
func (clockRT *ClockRT) Stop() {
	clockRT.pause <- struct{}{}
}

func (clockRT *ClockRT) Cycle(currCycle uint64) {
	for i := range clockRT.divided {
		div := &clockRT.divided[i]

		if div.counter == 0 {
			// TODO: should just be Rising, and the Falling part is commented out below.
			// This breaks the CPU somehow, though.
			div.clock.Cycle(div.cycle)
			div.cycle++
			div.counter = div.top
		} /*else if div.counter == (div.top >> 1) {
			div.clock.Falling(div.cycle)
			div.cycle++
		}*/
		div.counter--
	}
}

func (clockRT *ClockRT) Cycles(currCycle uint64, n uint64) uint64 {
	for offs := range n {
		clockRT.Cycle(currCycle + offs)
	}
	return currCycle + n
}
