package model

type SaveState struct {
	Mem       []Data8
	APU       APU
	Cartridge Cartridge

	BusData    Data8
	BusAddress Addr

	CPURegisterFile               RegisterFile
	CPUMachineCycle               int
	CPUClockCycle                 uint
	CPUWroteToAddressBusThisCycle bool
	CPULastBranchResult           int
	CPUHalted                     bool
	CPURewindBuffer               Rewind
	CPUCBOp                       CBOp
}

func Save(
	mem []Data8,
	apu *APU,
	bus *Bus,
	cart *Cartridge,
	cpu *CPU,
) *SaveState {
	save := &SaveState{}

	save.Mem = make([]Data8, len(mem))
	copy(save.Mem, mem)

	apu.Save(save)
	bus.Save(save)
	cart.Save(save)
	cpu.Save(save)

	return save
}

func LoadSave(
	save *SaveState,
	mem []Data8,
	apu *APU,
	bus *Bus,
	cart *Cartridge,
	cpu *CPU,
) {
	copy(mem, save.Mem)

	apu.LoadSave(save)
	bus.LoadSave(save)
	cart.LoadSave(save, mem)
	cpu.LoadSave(save, mem)
}
