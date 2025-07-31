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
	Onpanic         func(mem []Data8)
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

func (r *ClockRT) SetSpeedPercent(pct float64, audio *Audio) {
	r.Sync(func() {
		r.setSpeedPercent(pct, audio)
	})
}

func NewRealtimeClock() *ClockRT {
	clockRT := ClockRT{
		resume:  make(chan struct{}),
		pause:   make(chan struct{}),
		stop:    make(chan struct{}),
		jobs:    make(chan func()),
		Onpanic: func(mem []Data8) {},
	}
	return &clockRT
}

func (clockRT *ClockRT) AttachUIDevice(dev func()) {
	clockRT.uiDevices = append(clockRT.uiDevices, dev)
}

func (clockRT *ClockRT) wait() bool {
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

func (clockRT *ClockRT) setSpeedPercent(pct float64, audio *Audio) {
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

func (clockRT *ClockRT) run(initSpeedPercent float64, ints *Interrupts, debug *Debug, mem []Data8, fs *FrameSync, audio *Audio, apu *APU, ppu *PPU, cpu *CPU, timer *Timer) {
	defer func() {
		if e := recover(); e != nil {
			clockRT.Onpanic(mem)
			panic(e)
		}
	}()
	clockRT.setSpeedPercent(initSpeedPercent, audio)

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
			clockRT.MCycle(clockRT.mCyclesPerTick, ints, debug, mem, fs, audio, apu, ppu, cpu, timer)
		case <-uiTicker.C:
			clockRT.uiCycle()
		case <-clockRT.resume:
			fmt.Printf("Ignored resume\n")
		case <-clockRT.pause:
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
// When this function returns, the clock has paused
func (clockRT *ClockRT) Pause() {
	clockRT.pause <- struct{}{}
}

// Stop the clock
// When this function returns, the clock has stopped
func (clockRT *ClockRT) Stop() {
	clockRT.stop <- struct{}{}
}

func (clockRT *ClockRT) MCycle(
	n int,
	ints *Interrupts,
	debug *Debug,
	mem []Data8,
	fs *FrameSync,
	audio *Audio,
	apu *APU,
	ppu *PPU,
	cpu *CPU,
	timer *Timer,
) {
	for range n {
		// Breakpoints will stop here, right before executing next M-cycle
		if clockRT.pauseAfterCycle.Load() > 0 {
			exit := clockRT.wait()
			if exit {
				return
			}
			clockRT.pauseAfterCycle.Add(-1)
		}

		audio.Clock(apu)

		// Clock the CPU. This is the only place where the enabled-state of APU/PPU can change.
		cpu.fsm(clockRT, mem)

		m := clockRT.Cycle >> 2
		clockRT.Cycle += 4
		if m&0x3f == 0 {
			timer.tickDIV(mem, ints, apu)
		}

		// Clock the peripherals.
		// 99.99% of the time, both PPU and APU are on, so we clock everything
		if ppu.RegLCDC&apu.MasterCtl&Bit7 != 0 {
			// T0
			apu.Wave.clock(mem)
			if m&0x1 == 0 {
				apu.Pulse1.clock()
				apu.Pulse2.clock()
			}
			if clockRT.Cycle&0xf == 0 {
				apu.Noise.clock()
			}
			ppu.fsm(ints, debug, clockRT, mem, fs)

			// T1

			// T2
			ppu.fsm(ints, debug, clockRT, mem, fs)

			// T3
		} else {
			clockRT.mCycleSlowPath(
				m,
				ints,
				debug,
				mem,
				fs,
				apu,
				ppu,
			)
		}
	}
}

func (clockRT *ClockRT) mCycleSlowPath(m uint, ints *Interrupts, debug *Debug, mem []Data8, fs *FrameSync, apu *APU, ppu *PPU) {
	// T0
	if ppu.RegLCDC&Bit7 != 0 {
		ppu.fsm(ints, debug, clockRT, mem, fs)
	}
	if apu.MasterCtl&Bit7 != 0 {
		apu.Wave.clock(mem)
		if m&0x1 == 0 {
			apu.Pulse1.clock()
			apu.Pulse2.clock()
		}
		if clockRT.Cycle&0xf == 0 {
			apu.Noise.clock()
		}
	}

	// T1

	// T2
	if ppu.RegLCDC&Bit7 != 0 {
		ppu.fsm(ints, debug, clockRT, mem, fs)
	}

	// T3
}
