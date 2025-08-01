package model

type BootROMLock struct {
	BootOff bool
}

func (brl *BootROMLock) Write(gb *Gameboy, v Data8) {
	if brl.BootOff {
		return
	}
	if v&1 == 1 {
		brl.Lock(gb)
	}
}

func (brl *BootROMLock) Lock(gb *Gameboy) {
	brl.BootOff = true
	copy(gb.Mem[:SizeBootROM], gb.Cartridge.ROM[0][:SizeBootROM])

	// Update debug
	gb.Debug.SetProgram(gb.Mem[:AddrCartridgeBankNEnd])

	// Explore from known entry points (Cartridge entrypoint and interrupt vector)
	gb.Debug.Disassembler.SetPC(0x100)
	gb.Debug.Disassembler.SetPC(0x40)
	gb.Debug.Disassembler.SetPC(0x48)
	gb.Debug.Disassembler.SetPC(0x50)
	gb.Debug.Disassembler.SetPC(0x58)
	gb.Debug.Disassembler.SetPC(0x60)
}
