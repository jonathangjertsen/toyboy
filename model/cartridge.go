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
}

func (gb *Gameboy) SetROMBank(which Data8) {
	copy(gb.Mem[AddrCartridgeBankNBegin:AddrCartridgeBankNEnd], gb.Cartridge.ROM[which][:])
	gb.Cartridge.SelectedROMBank = which
}

func (gb *Gameboy) SetRAMBank(which Data8) {
	fmt.Printf("SetRAMBank %d\n", which)

	// Store current RAM contents to bank
	copy(gb.Cartridge.RAM[gb.Cartridge.SelectedRAMBank][:], gb.Mem[AddrCartridgeRAMBegin:AddrCartridgeRAMEnd])

	// Load from bank to RAM
	copy(gb.Mem[AddrCartridgeRAMBegin:AddrCartridgeRAMEnd], gb.Cartridge.RAM[which][:])
	gb.Cartridge.SelectedRAMBank = which
}

func (gb *Gameboy) WriteCartridge(addr Addr, v Data8) {
	switch gb.Cartridge.MBCFeatures.ID {
	case MBCIDNone:
		return
	case MBCID1:
		if addr <= 0x1fff {
			gb.Cartridge.ExtRAMEnabled = v&0x0f == 0x0a
		} else if addr <= 0x3fff {
			if v == 0x00 {
				v = 0x01
			}
			gb.Cartridge.BankNo1 = v & 0x1f
			gb.updateBank()
		} else if addr <= 0x5fff {
			gb.Cartridge.BankNo2 = v & 0x03
			gb.updateBank()
		} else if addr <= 0x7fff {
			gb.Cartridge.BankModeSel = v & 0x01
			gb.updateBank()
		}
	default:
		panic("not implemented MBC")
	}
}

func (gb *Gameboy) updateBank() {
	if gb.Cartridge.BankModeSel != 0x00 {
		panic("advanced banking mode not implemented")
	}
	newRAMBank := gb.Cartridge.SelectedRAMBank
	newROMBank := gb.Cartridge.SelectedROMBank

	if gb.Cartridge.MBCFeatures.NRAMBanks == 4 {
		newRAMBank = gb.Cartridge.BankNo2
	} else if gb.Cartridge.MBCFeatures.NROMBanks >= 64 {
		newROMBank = (gb.Cartridge.BankNo2 << 5) | gb.Cartridge.BankNo1
	} else {
		newROMBank = gb.Cartridge.BankNo1
	}

	if newRAMBank != gb.Cartridge.SelectedRAMBank {
		gb.SetRAMBank(newRAMBank)
	}
	if newROMBank != gb.Cartridge.SelectedROMBank {
		gb.SetROMBank(newROMBank)
	}
}
