package tests_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jonathangjertsen/toyboy/model"
	"github.com/jonathangjertsen/toyboy/tests"
)

func Test(t *testing.T) {
	for i := range 256 {
		if i == 0xcb {
			for j := range 256 {
				testCB(t, uint8(j))
				if t.Failed() {
					return
				}
			}
		} else {
			testOpcode(t, uint8(i))
			if t.Failed() {
				return
			}
		}
	}
}

func testOpcode(t *testing.T, opcodeRaw uint8) {
	t.Helper()
	opcode := model.Opcode(opcodeRaw)

	if slices.Contains([]model.Opcode{
		model.OpcodeHALT,
		model.OpcodeSTOP,
		model.OpcodeUndefD3,
		model.OpcodeUndefDB,
		model.OpcodeUndefDD,
		model.OpcodeUndefE3,
		model.OpcodeUndefE4,
		model.OpcodeUndefEB,
		model.OpcodeUndefEC,
		model.OpcodeUndefED,
		model.OpcodeUndefF4,
		model.OpcodeUndefFC,
		model.OpcodeUndefFD,
	}, opcode) {
		// not implemented yet
		return
	}

	if !t.Run(opcode.String(), func(t *testing.T) {
		tcs := tests.MustReadFile(fmt.Sprintf("%02x", uint8(opcode)))
		tests.Run(t, tcs, opcode)
	}) {
		return
	}
}

func testCB(t *testing.T, cb uint8) {
	t.Helper()
	opcode := model.NewCBOp(model.Data8(cb))
	if !t.Run("CB"+opcode.String(), func(t *testing.T) {
		tcs := tests.MustReadFile(fmt.Sprintf("cb %02x", uint8(cb)))
		tests.Run(t, tcs, model.OpcodeCB)
	}) {
		return
	}
}
