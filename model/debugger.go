package model

import (
	"fmt"
	"sync/atomic"
)

type Debugger struct {
	BreakX  atomic.Int64
	BreakY  atomic.Int64
	BreakPC atomic.Int64
	BreakIR atomic.Int64
	y       Data8
	CLK     *ClockRT
}

func NewDebugger(clk *ClockRT) Debugger {
	return Debugger{
		CLK: clk,
	}
}

func (dbg *Debugger) Init() {
	dbg.BreakX.Store(-1)
	dbg.BreakY.Store(-1)
	dbg.BreakPC.Store(-1)
	dbg.BreakIR.Store(-1)
}

func (dbg *Debugger) Break() {
	dbg.CLK.pauseAfterCycle.Add(1)
}

func (dbg *Debugger) SetY(y Data8) {
	dbg.y = y
}

func (dbg *Debugger) SetX(x Data8) {
	bx, by := dbg.BreakX.Load(), dbg.BreakY.Load()
	if int64(x) == bx && int64(dbg.y) == by {
		dbg.Break()
		fmt.Printf("PPU breakpoint\n")
	}
}

func (dbg *Debugger) SetIR(ir Opcode) {
	bpc := dbg.BreakIR.Load()
	if bpc == int64(ir) {
		dbg.Break()
		fmt.Printf("IR breakpoint\n")
	}
}

func (dbg *Debugger) SetPC(pc Addr) {
	bpc := dbg.BreakPC.Load()
	if bpc == int64(pc) {
		dbg.Break()
		fmt.Printf("PC breakpoint\n")
	}
}
