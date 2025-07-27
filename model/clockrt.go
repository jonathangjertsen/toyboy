package model

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ClockRT struct {
	ticker          *time.Ticker
	tickInterval    time.Duration
	Cycle           Cycle
	mCyclesPerTick  int
	resume          chan struct{}
	pause           chan struct{}
	stop            chan struct{}
	jobs            chan func()
	uiDevices       []func()
	Onpanic         func()
	pauseAfterCycle atomic.Int32
	Running         atomic.Bool

	cpu   *CPU
	ppu   *PPU
	apu   *APU
	timer *Timer
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
	clockRT.mCyclesPerTick = int(clockRT.tickInterval / (4 * cycleInterval))
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
	uiTicker := time.NewTicker(time.Second / 60)
	for {
		var exit bool
		select {
		case <-clockRT.ticker.C:
			clockRT.MCycle(clockRT.mCyclesPerTick)
		case <-uiTicker.C:
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

func (clockRT *ClockRT) MCycle(n int) {
	for range n {
		// Breakpoints will stop here, right before executing next M-cycle
		if clockRT.pauseAfterCycle.Load() > 0 {
			clockRT.wait()
			clockRT.pauseAfterCycle.Add(-1)
		}

		m := clockRT.Cycle.C >> 2
		clockRT.Cycle.C += 4

		// Clock the CPU. This is the only place where the enabled-state of APU/PPU can change.
		clockRT.cpu.fsm(Cycle{m, false})
		clockRT.cpu.fsm(Cycle{m, true})

		// Clock the peripherals.
		// 99.99% of the time, both PPU and APU are on, so we clock everything
		if clockRT.ppu.RegLCDC&clockRT.apu.MasterCtl&Bit7 != 0 {
			// T0
			clockRT.timer.tickDIVTimer()
			clockRT.apu.Wave.clock()
			clockRT.apu.Pulse1.clock()
			clockRT.apu.Pulse2.clock()
			if clockRT.Cycle.C&0xf == 0 {
				clockRT.apu.Noise.clock()
			}
			clockRT.ppu.fsm()

			// T1
			clockRT.timer.tickDIVTimer()

			// T2
			clockRT.timer.tickDIVTimer()
			clockRT.apu.Wave.clock()
			clockRT.ppu.fsm()

			// T3
			clockRT.timer.tickDIVTimer()
		} else {
			clockRT.mCycleSlowPath(
				clockRT.ppu.RegLCDC&Bit7 != 0,
				clockRT.apu.MasterCtl&Bit7 != 0,
			)
		}
	}
}

func (clockRT *ClockRT) mCycleSlowPath(ppu, apu bool) {
	// T0
	clockRT.timer.tickDIVTimer()
	if ppu {
		clockRT.ppu.fsm()
	}
	if apu {
		clockRT.apu.Wave.clock()
		clockRT.apu.Pulse1.clock()
		clockRT.apu.Pulse2.clock()
		if clockRT.Cycle.C&0xf == 0 {
			clockRT.apu.Noise.clock()
		}
	}

	// T1
	clockRT.timer.tickDIVTimer()

	// T2
	clockRT.timer.tickDIVTimer()
	if ppu {
		clockRT.ppu.fsm()
	}
	if apu {
		clockRT.apu.Wave.clock()
	}

	// T3
	clockRT.timer.tickDIVTimer()
}
