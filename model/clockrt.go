package model

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ClockRT struct {
	ticker          *time.Ticker
	tickInterval    time.Duration
	cycle           Cycle
	cyclesPerTick   uint64
	resume          chan struct{}
	pause           chan struct{}
	stop            chan struct{}
	jobs            chan func()
	uiDevices       []func()
	divided         []clockRTDivided
	Onpanic         func()
	pauseAfterCycle atomic.Int32
	Running         atomic.Bool
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

func (r *ClockRT) SetSpeedPercent(pct float64) {
	r.Sync(func() {
		r.setSpeedPercent(pct)
	})
}

func NewRealtimeClock(config ConfigClock) *ClockRT {
	clockRT := ClockRT{
		resume:  make(chan struct{}),
		pause:   make(chan struct{}),
		stop:    make(chan struct{}),
		jobs:    make(chan func()),
		Onpanic: func() {},
	}
	go clockRT.run(config.SpeedPercent)
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
	clockRT.Running.Store(false)
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
	clockRT.Running.Store(true)
}

func (clockRT *ClockRT) setSpeedPercent(pct float64) {
	f := 4194304.0 * pct / 100
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

func (clockRT *ClockRT) run(initSpeedPercent float64) {
	defer func() {
		if e := recover(); e != nil {
			clockRT.Onpanic()
			panic(e)
		}
	}()
	clockRT.setSpeedPercent(initSpeedPercent)

	clockRT.wait()
	clockRT.ticker = time.NewTicker(clockRT.tickInterval)
	for {
		var exit bool
		select {
		case <-clockRT.ticker.C:
			clockRT.Cycles(clockRT.cyclesPerTick)
			clockRT.uiCycle()
		case <-clockRT.resume:
			fmt.Printf("Ignored resume\n")
		case <-clockRT.pause:
			clockRT.wait()
		case job := <-clockRT.jobs:
			job()
		case <-clockRT.stop:
			clockRT.Running.Store(false)
			exit = true
		}
		if exit {
			break
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

// Pause the clock
// When this function returns, the clock has stopped
func (clockRT *ClockRT) Pause() {
	clockRT.pause <- struct{}{}
}

// Stop the clock
// When this function returns, the clock has stopped
func (clockRT *ClockRT) Stop() {
	clockRT.stop <- struct{}{}
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

func (clockRT *ClockRT) Cycles(n uint64) {
	for range n {
		clockRT.Cycle(clockRT.cycle.C)
		if clockRT.pauseAfterCycle.Load() > 0 {
			clockRT.wait()
			clockRT.pauseAfterCycle.Add(-1)
		}
		clockRT.cycle.C++
	}
}
