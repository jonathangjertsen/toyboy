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
	Time    string
	Message string
	Warn    bool
}

func (d *Debug) SetPC(gb *Gameboy, addr Addr, clk *ClockRT) {
	if d == nil {
		return
	}
	if addr != 0 {
		addr--
	}
	d.Disassembler.SetPC(addr)
}

func (d *Debug) SetIR(gb *Gameboy, op Opcode, clk *ClockRT) {
	if d == nil {
		return
	}
	d.Debugger.SetIR(gb, op, clk)
}

func (d *Debug) SetWarning(key string, message string) {
	if d == nil {
		return
	}
	d.Warnings[key] = UserMessage{
		Time:    time.Now().Format(time.RFC3339),
		Warn:    true,
		Message: message,
	}
	fmt.Printf("WARNING (%s): %s\n", key, message)
}
