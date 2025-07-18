package model

import "sync/atomic"

type Debugger struct {
	BreakX  atomic.Int64
	BreakY  atomic.Int64
	BreakPC atomic.Int64
	y       uint8
	CLK     *ClockRT
}

func NewDebugger(clk *ClockRT) *Debugger {
	dbg := &Debugger{
		CLK: clk,
	}
	dbg.BreakX.Store(-1)
	dbg.BreakY.Store(-1)
	dbg.BreakPC.Store(-1)
	return dbg
}

func (dbg *Debugger) Break() {
	dbg.CLK.pauseAfterCycle.Add(1)
}

func (dbg *Debugger) SetY(y uint8) {
	dbg.y = y
}

func (dbg *Debugger) SetX(x uint8) {
	bx, by := dbg.BreakX.Load(), dbg.BreakY.Load()
	if int64(x) == bx && int64(dbg.y) == by {
		dbg.Break()
	}
}

func (dbg *Debugger) SetPC(pc uint16) {
	bpc := dbg.BreakPC.Load()
	if bpc == int64(pc) {
		dbg.Break()
	}
}
