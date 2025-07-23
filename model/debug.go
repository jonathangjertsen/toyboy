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
	d.Disassembler.SetPC(addr)
	d.Debugger.SetPC(addr)
}

func (d *Debug) SetWarning(key string, message string) {
	d.Warnings[key] = UserMessage{
		Time:    time.Now(),
		Warn:    true,
		Message: message,
	}
	fmt.Printf("WARNING (%s): %s\n", key, message)
}
