package model

type BootROMLock struct {
	BootOff bool
}

func (brl *BootROMLock) Write(mem []Data8, debug *Debug, cart *Cartridge, v Data8) {
	if brl.BootOff {
		return
	}
	if v&1 == 1 {
		brl.Lock(mem, debug, cart)
	}
}

func (brl *BootROMLock) Lock(mem []Data8, debug *Debug, cart *Cartridge) {
	brl.BootOff = true
	copy(mem[:SizeBootROM], cart.ROM[0][:SizeBootROM])

	if debug != nil {
		// Update debug
		debug.SetProgram(mem[:AddrCartridgeBankNEnd])

		// Explore from known entry points (Cartridge entrypoint and interrupt vector)
		debug.Disassembler.SetPC(0x100)
		debug.Disassembler.SetPC(0x40)
		debug.Disassembler.SetPC(0x48)
		debug.Disassembler.SetPC(0x50)
		debug.Disassembler.SetPC(0x58)
		debug.Disassembler.SetPC(0x60)
	}
}
