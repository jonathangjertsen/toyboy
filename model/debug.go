package model

import (
	"fmt"
	"time"
)

type Debug struct {
	Disassembler
	Debugger
	Warnings map[string]UserMessage
}

type UserMessage struct {
	Time    time.Time
	Message string
	Warn    bool
}

func (d *Debug) SetPC(addr Addr, clk *ClockRT) {
	if d == nil {
		return
	}
	if addr != 0 {
		addr--
	}
	d.Disassembler.SetPC(addr)
	d.Debugger.SetPC(addr, clk)
}

func (d *Debug) SetIR(op Opcode, clk *ClockRT) {
	if d == nil {
		return
	}
	d.Debugger.SetIR(op, clk)
}

func (d *Debug) SetWarning(key string, message string) {
	if d == nil {
		return
	}
	d.Warnings[key] = UserMessage{
		Time:    time.Now(),
		Warn:    true,
		Message: message,
	}
	fmt.Printf("WARNING (%s): %s\n", key, message)
}
