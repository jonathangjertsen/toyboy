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
	return Debugger{
		BreakX:  -1,
		BreakY:  -1,
		BreakPC: -1,
		BreakIR: -1,
	}
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

func (dbg *Debugger) SetIR(gb *Gameboy, ir Opcode, clk *ClockRT) {
	if dbg == nil {
		return
	}
	if dbg.BreakIR == int64(ir) {
		dbg.Break(clk)
		fmt.Printf("IR breakpoint\n")
	}
	if dbg.BreakPC >= int64(gb.CPU.Regs.PC) && dbg.BreakPC < int64(gb.CPU.Regs.PC+Addr(instSize[ir])) {
		dbg.Break(clk)
		fmt.Printf("PC breakpoint\n")
	}
}
