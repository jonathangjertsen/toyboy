package model

import (
	"context"
	"log/slog"
)

type SOC struct {
	CLK *Clock
	PHI *Clock
}

func NewSOC(ctx context.Context, logger *slog.Logger, config HWConfig) *SOC {
	clk := NewClock(ctx, logger, config.SystemClock)
	phi := clk.Divide(ctx, logger, 4, 100)

	soc := &SOC{}
	soc.CLK = clk
	soc.PHI = phi

	return soc
}
