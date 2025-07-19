package main

import (
	"fmt"
	"os"

	"github.com/jonathangjertsen/toyboy/model"
)

func main() {
	rom := "assets/cartridges/unbricked.gb"
	bytes, err := os.ReadFile(rom)
	if err != nil {
		panic(err)
	}
	dis := model.NewDisassembler(bytes)
	dis.SetPC(0x0100)
	out := dis.Disassembly()
	for _, section := range out.Code {
		fmt.Printf("\nCode section at 0x%04x\n", section.Address())
		for _, inst := range section.Instructions {
			fmt.Printf("%04x | %v\n", inst.Address, inst.Opcode)
		}
	}
	data := splitSections(out.Data)
	prevEndAddr := uint16(0xffff)
	for _, section := range data {
		if prevEndAddr != section.Address {
			fmt.Printf("\nData section at 0x%04x\n", section.Address)
		}
		prevEndAddr = section.Address + uint16(len(section.Raw))
		allEqual := true
		testByte := section.Raw[0]
		for _, b := range section.Raw {
			if b != testByte {
				allEqual = false
				break
			}
		}
		if allEqual {
			fmt.Printf("0x%0x bytes of 0x%02x\n", len(section.Raw), testByte)
			continue
		}

		i := 0
		for line := range (len(section.Raw) + 15) / 16 {
			fmt.Printf("%04x | ", int(section.Address)+line*16)
			for range 16 {
				if i >= len(section.Raw) {
					break
				}
				fmt.Printf("%02x ", section.Raw[i])
				i++
			}
			fmt.Printf("\n")
		}
	}
}
func splitSections(sections []model.DataSection) []model.DataSection {
	var result []model.DataSection
	for _, section := range sections {
		raw := section.Raw
		start := 0

		for start < len(raw) {
			runStart := start
			runByte := raw[start]
			runLen := 1

			// Find run of same byte
			for i := start + 1; i < len(raw) && raw[i] == runByte; i++ {
				runLen++
			}

			if runLen >= 32 {
				// Add preceding non-uniform part
				if runStart > 0 {
					result = append(result, model.DataSection{
						Address: section.Address,
						Raw:     raw[0:runStart],
					})
				}
				// Add the uniform run
				result = append(result, model.DataSection{
					Address: section.Address + uint16(runStart),
					Raw:     raw[runStart : runStart+runLen],
				})
				// Continue after the run
				raw = raw[runStart+runLen:]
				section.Address += uint16(runStart + runLen)
				start = 0
			} else {
				start++
			}
		}

		// Add any remaining non-uniform tail
		if len(raw) > 0 {
			result = append(result, model.DataSection{
				Address: section.Address,
				Raw:     raw,
			})
		}
	}
	return result
}
