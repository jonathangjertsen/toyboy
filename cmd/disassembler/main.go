package main

import (
	"os"

	"github.com/jonathangjertsen/toyboy/model"
)

func main() {
	rom := "assets/cartridges/01-special.gb"
	bytes, err := os.ReadFile(rom)
	if err != nil {
		panic(err)
	}
	dis := model.NewDisassembler(&model.DefaultConfig.Debug.Disassembler)
	dis.SetProgram(bytes)
	dis.ExploreFrom(0x0100)
	out := dis.Disassembly()
	out.Print(os.Stdout)
}
