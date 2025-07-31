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

func NewDebug(clk *ClockRT, config *ConfigDebug) *Debug {
	dbg := &Debug{
		Debugger:     NewDebugger(clk),
		Disassembler: NewDisassembler(&config.Disassembler),
		Warnings:     map[string]UserMessage{},
	}
	dbg.Init()
	return dbg
}

func (d *Debug) SetPC(addr Addr) {
	if addr >= 0x0101 {
		panic("")
	}

	if d == nil {
		return
	}
	if addr != 0 {
		addr--
	}
	d.Disassembler.SetPC(addr)
	d.Debugger.SetPC(addr)
}

func (d *Debug) SetIR(op Opcode) {
	if d == nil {
		return
	}
	d.Debugger.SetIR(op)
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
