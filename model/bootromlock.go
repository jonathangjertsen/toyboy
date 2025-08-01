package model

type BootROMLock struct {
	BootOff bool

	mem       []Data8 `json:"-"`
	cartridge *Cartridge
	debug     *Debug
}

func (brl *BootROMLock) Write(addr Addr, v Data8) {
	if brl.BootOff {
		return
	}
	if v&1 == 1 {
		brl.Lock()
	}
}

func (brl *BootROMLock) Lock() {
	brl.BootOff = true
	copy(brl.mem[:SizeBootROM], brl.cartridge.ROM[0][:SizeBootROM])

	if brl.debug != nil {
		// Update debug
		brl.debug.SetProgram(brl.mem[:AddrCartridgeBankNEnd])

		// Explore from known entry points (Cartridge entrypoint and interrupt vector)
		brl.debug.Disassembler.SetPC(0x100)
		brl.debug.Disassembler.SetPC(0x40)
		brl.debug.Disassembler.SetPC(0x48)
		brl.debug.Disassembler.SetPC(0x50)
		brl.debug.Disassembler.SetPC(0x58)
		brl.debug.Disassembler.SetPC(0x60)
	}
}
