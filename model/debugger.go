package model

import (
	"fmt"
)

type Debugger struct {
	BreakX  int64
	BreakY  int64
	BreakPC int64
	BreakIR int64
	CurrY   Data8
}

func NewDebugger() Debugger {
	return Debugger{}
}

func (dbg *Debugger) Init() {
	if dbg == nil {
		return
	}
	dbg.BreakX = -1
	dbg.BreakY = -1
	dbg.BreakPC = -1
	dbg.BreakIR = -1
}

func (dbg *Debugger) Break(clk *ClockRT) {
	if dbg == nil {
		return
	}
	clk.PauseAfterCycle.Add(1)
}

func (dbg *Debugger) SetY(y Data8) {
	if dbg == nil {
		return
	}
	dbg.CurrY = y
}

func (dbg *Debugger) SetX(x Data8, clk *ClockRT) {
	if dbg == nil {
		return
	}
	bx, by := dbg.BreakX, dbg.BreakY
	if int64(x) == bx && int64(dbg.CurrY) == by {
		dbg.Break(clk)
		fmt.Printf("PPU breakpoint\n")
	}
}

func (dbg *Debugger) SetIR(ir Opcode, clk *ClockRT) {
	if dbg == nil {
		return
	}
	bpc := dbg.BreakIR
	if bpc == int64(ir) {
		dbg.Break(clk)
		fmt.Printf("IR breakpoint\n")
	}
}

func (dbg *Debugger) SetPC(pc Addr, clk *ClockRT) {
	if dbg == nil {
		return
	}
	bpc := dbg.BreakPC
	if bpc == int64(pc) {
		dbg.Break(clk)
		fmt.Printf("PC breakpoint\n")
	}
}
