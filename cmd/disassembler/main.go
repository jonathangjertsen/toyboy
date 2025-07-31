package main

import (
	"os"

	"github.com/jonathangjertsen/toyboy/model"
)

func main() {
	bytes, err := os.ReadFile("assets/cartridges/02-interrupts.gb")
	if err != nil {
		panic(err)
	}
	dis := model.NewDisassembler(&model.DefaultConfig.Debug.Disassembler)
	dis.SetProgram(bytes)
	dis.ExploreFrom(model.AddrCartridgeEntryPoint)
	out := dis.Disassembly(0, model.AddrCartridgeBankNEnd)
	out.Print(os.Stdout)
}
