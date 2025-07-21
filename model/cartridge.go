package model

import "fmt"

type Cartridge struct {
	Bank0 MemoryRegion
	BankN MemoryRegion
}

func NewCartridge(clk *ClockRT) *Cartridge {
	return &Cartridge{
		Bank0: NewMemoryRegion(clk, AddrZero, SizeCartridgeBank),
		BankN: NewMemoryRegion(clk, AddrCartridgeBankNBegin, SizeCartridgeBank),
	}
}

func (cart *Cartridge) LoadROM(data []uint8) {
	// No MCB impl yet
	if len(data) != 0x8000 {
		panic(fmt.Sprintf("len(ROM)=%04x, want 0x8000", len(data)))
	}
	cart.Bank0.Data = Data8Slice(data[:0x4000])
	cart.BankN.Data = Data8Slice(data[0x4000:0x8000])
}

func (cart *Cartridge) Read(addr Addr) Data8 {
	if addr <= AddrCartridgeBank0End {
		return cart.Bank0.Read(addr)
	}
	return cart.BankN.Read(addr)
}

func (cart *Cartridge) Write(addr Addr, v Data8) {
	if addr <= AddrCartridgeBank0End {
		// TODO: implement MBCs
		return
	}
	if addr <= AddrCartridgeBankNEnd {
		// TODO: implement MBCs
		return
	}
	panicf("write to bank other than 0 not implemented, addr=%s v=%s\n", addr.Hex(), v.Hex())
}

func (cart *Cartridge) GetCounters(addr Addr) (uint64, uint64) {
	if addr <= AddrCartridgeBank0End {
		return cart.Bank0.GetCounters(addr)
	}
	if addr <= AddrCartridgeBankNEnd {
		// TODO
		return cart.BankN.GetCounters(addr)
	}
	return 0, 0
}
