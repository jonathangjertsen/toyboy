package model

import (
	"context"
	"log/slog"
	"time"

	"github.com/jonathangjertsen/gameboy/util"
)

type Clock struct {
	broker      *util.Broker[Tick]
	ticker      *time.Ticker
	SetDuration chan time.Duration
}

type Tick struct {
}

func NewClock(ctx context.Context, logger *slog.Logger, freq float64) *Clock {
	interval := time.Duration(float64(time.Second) / freq)
	ticker := time.NewTicker(interval)
	clock := &Clock{
		broker:      util.NewBroker[Tick](ctx, logger, interval),
		ticker:      ticker,
		SetDuration: make(chan time.Duration, 1),
	}
	go clock.worker(ctx)
	return clock
}

func (c *Clock) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.ticker.C:
			c.broker.In <- Tick{}
		case d := <-c.SetDuration:
			c.ticker.Reset(d)
		}
	}
}
