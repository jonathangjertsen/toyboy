package model

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ClockRT struct {
	ticker          *time.Ticker
	tickInterval    time.Duration
	Cycle           uint
	mCyclesPerTick  int
	resume          chan struct{}
	pause           chan struct{}
	stop            chan struct{}
	jobs            chan func()
	uiDevices       []func()
	Onpanic         func()
	pauseAfterCycle atomic.Int32
	Running         atomic.Bool
	Audio           *Audio

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

func NewRealtimeClock(config ConfigClock, audio *Audio) *ClockRT {
	clockRT := ClockRT{
		resume:  make(chan struct{}),
		pause:   make(chan struct{}),
		stop:    make(chan struct{}),
		jobs:    make(chan func()),
		Onpanic: func() {},
		Audio:   audio,
	}
	go clockRT.run(config.SpeedPercent)
	return &clockRT
}

func (clockRT *ClockRT) AttachUIDevice(dev func()) {
	clockRT.uiDevices = append(clockRT.uiDevices, dev)
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
	// Target frequency
	tFreq := 4194304.0 * pct / 100
	mFreq := tFreq / 4

	// Convert to interval
	mCycleInterval := time.Duration(float64(time.Second) / mFreq)

	// Update audio
	clockRT.Audio.SetMPeriod(mCycleInterval)

	// How often we run the real ticker
	minTickInterval := time.Millisecond * 2
	if mCycleInterval > minTickInterval {
		clockRT.tickInterval = mCycleInterval
		clockRT.mCyclesPerTick = 1
	} else {
		clockRT.tickInterval = minTickInterval
		clockRT.mCyclesPerTick = int(clockRT.tickInterval / mCycleInterval)
	}

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

		clockRT.Audio.Clock()

		m := clockRT.Cycle >> 2
		clockRT.Cycle += 4

		// Clock the CPU. This is the only place where the enabled-state of APU/PPU can change.
		clockRT.cpu.fsm(m)

		// Clock the peripherals.
		// 99.99% of the time, both PPU and APU are on, so we clock everything
		if clockRT.ppu.RegLCDC&clockRT.apu.MasterCtl&Bit7 != 0 {
			if m&0x3f == 0 {
				clockRT.timer.tickDIV()
			}

			// T0
			clockRT.apu.Wave.clock()
			if m&0x1 == 0 {
				clockRT.apu.Pulse1.clock()
				clockRT.apu.Pulse2.clock()
				clockRT.apu.Noise.clock()
			}
			clockRT.ppu.fsm()

			// T1

			// T2
			clockRT.ppu.fsm()

			// T3
		} else {
			clockRT.mCycleSlowPath(
				m,
				clockRT.ppu.RegLCDC&Bit7 != 0,
				clockRT.apu.MasterCtl&Bit7 != 0,
			)
		}
	}
}

func (clockRT *ClockRT) mCycleSlowPath(m uint, ppu, apu bool) {
	// T0
	if apu {
		clockRT.apu.Wave.clock()

		if m&0x1 == 0 {
			clockRT.apu.Pulse1.clock()
			clockRT.apu.Pulse2.clock()
		}
		if m&0x3 == 0 {
			clockRT.apu.Noise.clock()
		}
	}
	if ppu {
		clockRT.ppu.fsm()
	}

	// T1

	// T2
	if ppu {
		clockRT.ppu.fsm()
	}
	if apu {
		clockRT.apu.Wave.clock()
	}

	// T3
}
