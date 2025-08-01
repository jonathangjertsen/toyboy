package model

import (
	"fmt"
)

type Cartridge struct {
	ROM             [512][ROMBankSize]Data8
	RAM             [4][RAMBankSize]Data8
	MBCFeatures     MBCFeatures
	ExtRAMEnabled   bool
	BankNo1         Data8
	BankNo2         Data8
	BankModeSel     Data8
	SelectedRAMBank Data8
	SelectedROMBank Data8

	mem []Data8
}

func (cart *Cartridge) LoadSave(save *SaveState, mem []Data8) {
	*cart = save.Cartridge
	cart.mem = mem
}

func (cart *Cartridge) Save(save *SaveState) {
	save.Cartridge = *cart
}

func (cart *Cartridge) SetROMBank(which Data8) {
	copy(cart.mem[AddrCartridgeBankNBegin:AddrCartridgeBankNEnd], cart.ROM[which][:])
	cart.SelectedROMBank = which
}

func (cart *Cartridge) SetRAMBank(which Data8) {
	fmt.Printf("SetRAMBank %d\n", which)

	// Store current RAM contents to bank
	copy(cart.RAM[cart.SelectedRAMBank][:], cart.mem[AddrCartridgeRAMBegin:AddrCartridgeRAMEnd])

	// Load from bank to RAM
	copy(cart.mem[AddrCartridgeRAMBegin:AddrCartridgeRAMEnd], cart.RAM[which][:])
	cart.SelectedRAMBank = which
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
