package model_test

import (
	"testing"

	"github.com/jonathangjertsen/gameboy/model"
)

func TestOpcodeLDAE(t *testing.T) {
	coreTest(func(core *model.CPUCore, mr *model.MemoryRegion, clock *model.Clock) {
		mr.Write(0x0000, uint8(model.OpcodeLDAE))
		core.Regs.A = 0xff
		core.Regs.E = 0x77
		clock.Cycles(0, 2)
		if core.Regs.A != 0x77 {
			t.Fatalf("LDAE failed: expected A=0x77, got A=0x%02x", core.Regs.A)
		}
		if core.Regs.E != 0x77 {
			t.Fatalf("LDAE failed: expected E=0x77, got E=0x%02x", core.Regs.A)
		}
	})
}

func coreTest(f func(core *model.CPUCore, mr *model.MemoryRegion, clock *model.Clock)) {
	phi := model.NewClock()
	core := model.NewCPUCore(phi)
	mr := model.NewMemoryRegion("PROGRAM", 0x0000, 0x10)
	core.AttachPeripheral(&mr)
	f(core, &mr, phi)
}
