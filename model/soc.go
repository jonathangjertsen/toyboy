package model

import (
	"context"
	"log/slog"
)

type SOC struct {
	SystemClock *Clock
}

func NewSOC(ctx context.Context, logger *slog.Logger, config HWConfig) *SOC {
	systemClock := NewClock(ctx, logger, config.SystemClockFrequency)

	soc := &SOC{}
	soc.SystemClock = systemClock

	return soc
}
