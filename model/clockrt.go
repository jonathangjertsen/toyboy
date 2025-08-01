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
	pause           chan chan bool
	stop            chan struct{}
	jobs            chan func()
	uiDevices       []func()
	Onpanic         func(gb *Gameboy)
	PauseAfterCycle atomic.Int32
	Running         atomic.Bool
}

func NewClock() *ClockRT {
	return &ClockRT{
		resume:  make(chan struct{}),
		pause:   make(chan chan bool),
		stop:    make(chan struct{}),
		jobs:    make(chan func()),
		Onpanic: func(gb *Gameboy) {},
	}
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

func (r *ClockRT) SetSpeedPercent(pct float64, audio Audio) {
	r.Sync(func() {
		r.setSpeedPercent(pct, audio)
	})
}

func (clockRT *ClockRT) AttachUIDevice(dev func()) {
	clockRT.uiDevices = append(clockRT.uiDevices, dev)
}

func (clockRT *ClockRT) wait() bool {
	clockRT.Running.Store(false)
	for {
		resumed := false
		select {
		case res := <-clockRT.pause:
			res <- false
		case <-clockRT.resume:
			resumed = true
		case job := <-clockRT.jobs:
			job()
		case <-clockRT.stop:
			return true
		}
		if resumed {
			break
		}
	}
	clockRT.Running.Store(true)
	return false
}

func (clockRT *ClockRT) setSpeedPercent(pct float64, audio Audio) {
	// Target frequency
	tFreq := 4194304.0 * pct / 100
	mFreq := tFreq / 4

	// Convert to interval
	mCycleInterval := time.Duration(float64(time.Second) / mFreq)
	mCycleInterval /= 2

	// Update audio
	audio.SetMPeriod(mCycleInterval)

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

func (clockRT *ClockRT) Run(gb *Gameboy, config *Config, audio Audio, fs *FrameSync) {
	defer func() {
		if e := recover(); e != nil {
			clockRT.Onpanic(gb)
			panic(e)
		}
	}()
	clockRT.setSpeedPercent(config.Clock.SpeedPercent, audio)

	exit := clockRT.wait()
	if exit {
		return
	}
	clockRT.ticker = time.NewTicker(clockRT.tickInterval)
	uiTicker := time.NewTicker(time.Second / 60)
	for {
		var exit bool
		select {
		case <-clockRT.ticker.C:
			clockRT.MCycle(clockRT.mCyclesPerTick, gb, audio, fs)
		case <-uiTicker.C:
			clockRT.uiCycle()
		case <-clockRT.resume:
			fmt.Printf("Ignored resume\n")
		case res := <-clockRT.pause:
			res <- true
			exit = clockRT.wait()
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
// When this function returns, the clock has paused.
// Returns the previous run state of the clock
func (clockRT *ClockRT) Pause() bool {
	res := make(chan bool)
	clockRT.pause <- res
	return <-res
}

// Stop the clock
// When this function returns, the clock has stopped
func (clockRT *ClockRT) Stop() {
	clockRT.stop <- struct{}{}
}

func (clockRT *ClockRT) MCycle(
	n int,
	gb *Gameboy,
	audio Audio,
	fs *FrameSync,
) {
	for i := range n {
		_ = i

		gb.TCycle = clockRT.Cycle
		gb.MCycle = clockRT.Cycle >> 2

		// Breakpoints will stop here, right before executing next M-cycle
		if clockRT.PauseAfterCycle.Load() > 0 {
			exit := clockRT.wait()
			if exit {
				return
			}
			clockRT.PauseAfterCycle.Add(-1)
		}

		audio.Clock(&gb.APU)

		gb.applyPendingIME()

		// Clock the CPU. This is the only place where the enabled-state of APU/PPU can change.
		gb.CPU.fsm(clockRT, gb)

		m := clockRT.Cycle >> 2
		clockRT.Cycle += 4

		// Clock the peripherals.
		// 99.99% of the time, both PPU and APU are on, so we clock everything
		if gb.PPU.RegLCDC&gb.APU.MasterCtl&Bit7 != 0 {
			// T0
			gb.tickDIV()
			gb.APU.Wave.clock(gb.Mem)
			if m&0x1 == 0 {
				gb.APU.Pulse1.clock()
				gb.APU.Pulse2.clock()
			}
			if clockRT.Cycle&0xf == 0 {
				gb.APU.Noise.clock()
			}
			gb.PPU.fsm(gb, clockRT, fs)

			// T1
			gb.tickDIV()

			// T2
			gb.PPU.fsm(gb, clockRT, fs)
			gb.tickDIV()

			// T3
			gb.tickDIV()
		} else {
			clockRT.mCycleSlowPath(m, gb, fs)
		}
	}
}

func (clockRT *ClockRT) mCycleSlowPath(m uint, gb *Gameboy, fs *FrameSync) {
	// T0
	gb.tickDIV()
	if gb.PPU.RegLCDC&Bit7 != 0 {
		gb.PPU.fsm(gb, clockRT, fs)
	}
	if gb.APU.MasterCtl&Bit7 != 0 {
		gb.APU.Wave.clock(gb.Mem)
		if m&0x1 == 0 {
			gb.APU.Pulse1.clock()
			gb.APU.Pulse2.clock()
		}
		if clockRT.Cycle&0xf == 0 {
			gb.APU.Noise.clock()
		}
	}

	// T1
	gb.tickDIV()

	// T2
	gb.tickDIV()
	if gb.PPU.RegLCDC&Bit7 != 0 {
		gb.PPU.fsm(gb, clockRT, fs)
	}

	// T3
	gb.tickDIV()
}
