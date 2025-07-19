package main

import (
	"os"

	"github.com/jonathangjertsen/toyboy/model"
)

func main() {
	rom := "assets/cartridges/tetris.gb"
	bytes, err := os.ReadFile(rom)
	if err != nil {
		panic(err)
	}
	dis := model.NewDisassembler()
	dis.Trace = true
	dis.SetProgram(bytes)
	dis.SetPC(0x0100)
	out := dis.Disassembly()
	out.Print(os.Stdout)
}
