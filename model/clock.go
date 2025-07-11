package model

import (
	"context"
	"log/slog"
	"time"
)

type Clock struct {
	SetInterval chan time.Duration
	AddCallback chan func(Cycle)
	callbacks   []func(c Cycle)
}

type Cycle struct {
	T time.Time
	C uint64
}

func NewClock(ctx context.Context, logger *slog.Logger, config ClockConfig) *Clock {
	cycleInterval := time.Duration(float64(time.Second) / config.Frequency)
	tickInterval := time.Millisecond * 16
	ticker := time.NewTicker(tickInterval)
	clock := &Clock{
		SetInterval: make(chan time.Duration, 1),
		AddCallback: make(chan func(Cycle), 1),
	}
	go func() {
		var cycle uint64
		cyclesPerTick := tickInterval / cycleInterval
		for {
			select {
			case t := <-ticker.C:
				for range cyclesPerTick {
					clock.cycle(t, cycle)
					cycle++
				}
			case d := <-clock.SetInterval:
				cyclesPerTick = tickInterval / d
			case cb := <-clock.AddCallback:
				clock.callbacks = append(clock.callbacks, cb)
			}
		}
	}()
	return clock
}

func (c *Clock) cycle(t time.Time, count uint64) {
	for _, cb := range c.callbacks {
		cb(Cycle{t, count})
	}
}

func (c *Clock) Divide(ctx context.Context, logger *slog.Logger, div uint64, buffer int) *Clock {
	child := &Clock{
		AddCallback: make(chan func(Cycle), 1),
	}
	c.AddCallback <- func(cyc Cycle) {
		d, m := cyc.C/div, cyc.C%div
		if m == 0 {
			child.cycle(cyc.T, d)
		}
	}
	go func() {
		for cb := range child.AddCallback {
			child.callbacks = append(child.callbacks, cb)
		}
	}()
	return child
}
