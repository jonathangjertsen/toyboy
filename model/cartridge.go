package model

import "fmt"

type Cartridge struct {
	ROM             [512][ROMBankSize]uint8
	RAM             [4][RAMBankSize]uint8
	CurrROMBank0    MemoryRegion
	CurrROMBankN    MemoryRegion
	CurrRAMBank     MemoryRegion
	MBCFeatures     MBCFeatures
	ExtRAMEnabled   bool
	BankNo1         Data8
	BankNo2         Data8
	BankModeSel     Data8
	SelectedRAMBank Data8
	SelectedROMBank Data8
}

func NewCartridge(clk *ClockRT) *Cartridge {
	return &Cartridge{
		CurrROMBank0:    NewMemoryRegion(clk, AddrZero, SizeCartridgeBank),
		CurrROMBankN:    NewMemoryRegion(clk, AddrCartridgeBankNBegin, SizeCartridgeBank),
		CurrRAMBank:     NewMemoryRegion(clk, AddrCartridgeRAMBegin, SizeCartridgeRAM),
		BankNo1:         1,
		SelectedROMBank: 1,
	}
}

func (cart *Cartridge) LoadROM(data []uint8) {
	if len(data)%ROMBankSize != 0 {
		panicf("ROM size %d is not a multiple of %d", len(data), ROMBankSize)
	}
	for i := range len(data) / ROMBankSize {
		copy(cart.ROM[i][:], data[i*ROMBankSize:(i+1)*ROMBankSize])
	}
	cart.MBCFeatures = GetMBCFeatures(data[AddrCartridgeType], data[AddrROMSize], data[AddrRAMSize])
	cart.CurrROMBank0.Data = Data8Slice(cart.ROM[0][:])
	cart.CurrROMBankN.Data = Data8Slice(cart.ROM[1][:])
}

func (cart *Cartridge) SetROMBank(which Data8) {
	fmt.Printf("SetROMBank %d\n", which)

	copy(cart.CurrROMBankN.Data, Data8Slice(cart.ROM[which][:]))
	cart.SelectedROMBank = which
}

func (cart *Cartridge) SetRAMBank(which Data8) {
	fmt.Printf("SetRAMBank %d\n", which)

	// Store current RAM contents to bank
	copy(cart.RAM[cart.SelectedRAMBank][:], U8Slice(cart.CurrRAMBank.Data))

	// Load from bank to RAM
	copy(cart.CurrRAMBank.Data, Data8Slice(cart.RAM[which][:]))

	cart.SelectedRAMBank = which
}

func (cart *Cartridge) Read(addr Addr) Data8 {
	if addr <= AddrCartridgeBank0End {
		return cart.CurrROMBank0.Read(addr)
	}
	return cart.CurrROMBankN.Read(addr)
}

func (cart *Cartridge) Write(addr Addr, v Data8) {
	switch cart.MBCFeatures.ID {
	case MBCIDNone:
		return
	case MBCID1:
		if addr <= 0x1fff {
			cart.ExtRAMEnabled = v&0x0f == 0x0a
		} else if addr <= 0x3fff {
			if v == 0x00 {
				v = 0x01
			}
			cart.BankNo1 = v & 0x1f
			cart.updateBank()
		} else if addr <= 0x5fff {
			cart.BankNo2 = v & 0x03
			cart.updateBank()
		} else if addr <= 0x7fff {
			cart.BankModeSel = v & 0x01
			cart.updateBank()
		}
	default:
		panic("not implemented MBC")
	}
}

func (cart *Cartridge) updateBank() {
	if cart.BankModeSel != 0x00 {
		panic("advanced banking mode not implemented")
	}
	newRAMBank := cart.SelectedRAMBank
	newROMBank := cart.SelectedROMBank

	if cart.MBCFeatures.NRAMBanks == 4 {
		newRAMBank = cart.BankNo2
	} else if cart.MBCFeatures.NROMBanks >= 64 {
		newROMBank = (cart.BankNo2 << 5) | cart.BankNo1
	} else {
		newROMBank = cart.BankNo1
	}

	if newRAMBank != cart.SelectedRAMBank {
		cart.SetRAMBank(newRAMBank)
	}
	if newROMBank != cart.SelectedROMBank {
		cart.SetROMBank(newROMBank)
	}
}

func (cart *Cartridge) GetCounters(addr Addr) (uint64, uint64) {
	if addr <= AddrCartridgeBank0End {
		return cart.CurrROMBank0.GetCounters(addr)
	}
	if addr <= AddrCartridgeBankNEnd {
		// TODO
		return cart.CurrROMBankN.GetCounters(addr)
	}
	return 0, 0
}
